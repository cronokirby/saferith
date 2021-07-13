package safenum

import (
	"crypto/subtle"
	"math/big"
	"math/bits"
)

// Constant Time Utilities

// choice represents a constant-time boolean.
//
// The value of choice is always either 1 or 0.
//
// We use a separate type instead of bool, in order to be able to make decisions without leaking
// which decision was made.
type choice Word

// ctEq compares x and y for equality, returning 1 if equal, and 0 otherwise
//
// This doesn't leak any information about either of them
func ctEq(x, y Word) choice {
	// If x == y, then x ^ y should be all zero bits.
	q := uint64(x ^ y)
	// For any q != 0, either the MSB of q, or the MSB of -q is 1.
	// We can thus or those together, and check the top bit. When q is zero,
	// that means that x and y are equal, so we negate that top bit.
	return 1 ^ choice((q|-q)>>63)
}

// ctGt checks x > y, returning 1 or 0
//
// This doesn't leak any information about either of them
func ctGt(x, y Word) choice {
	z := y - x
	return choice((z ^ ((x ^ y) & (x ^ z))) >> (_W - 1))
}

// ctIfElse selects x if v = 1, and y otherwise
//
// This doesn't leak the value of any of its inputs
func ctIfElse(v choice, x, y Word) Word {
	// mask should be all 1s if v is 1, otherwise all 0s
	mask := -Word(v)
	return y ^ (mask & (y ^ x))
}

// ctCondCopy copies y into x, if v == 1, otherwise does nothing
//
// Both slices must have the same length.
//
// LEAK: the length of the slices
//
// Otherwise, which branch was taken isn't leaked
func ctCondCopy(v choice, x, y []Word) {
	for i := 0; i < len(x); i++ {
		x[i] = ctIfElse(v, y[i], x[i])
	}
}

// ctCondSwap swaps the contents of a and b, when v == 1, otherwise does nothing
//
// Both slices must have the same length.
//
// LEAK: the length of the slices
//
// Whether or not a swap happened isn't leaked
func ctCondSwap(v choice, a, b []Word) {
	for i := 0; i < len(a) && i < len(b); i++ {
		ai := a[i]
		a[i] = ctIfElse(v, b[i], ai)
		b[i] = ctIfElse(v, ai, b[i])
	}
}

// "Missing" Functions
// These are routines that could in theory be implemented in assembly,
// but aren't already present in Go's big number routines

// div calculates the quotient and remainder of hi:lo / d
//
// Unlike bits.Div, this doesn't leak anything about the inputs
func div(hi, lo, d Word) (Word, Word) {
	var quo Word
	hi = ctIfElse(ctEq(hi, d), 0, hi)
	for i := _W - 1; i > 0; i-- {
		j := _W - i
		w := (hi << j) | (lo >> i)
		sel := ctEq(w, d) | ctGt(w, d) | choice(hi>>i)
		hi2 := (w - d) >> j
		lo2 := lo - (d << i)
		hi = ctIfElse(sel, hi2, hi)
		lo = ctIfElse(sel, lo2, lo)
		quo |= Word(sel)
		quo <<= 1
	}
	sel := ctEq(lo, d) | ctGt(lo, d) | choice(hi)
	quo |= Word(sel)
	rem := ctIfElse(sel, lo-d, lo)
	return quo, rem
}

// mulSubVVW calculates z -= y * x
//
// This also results in a carry.
func mulSubVVW(z, x []Word, y Word) (c Word) {
	for i := 0; i < len(z) && i < len(x); i++ {
		hi, lo := mulAddWWW_g(x[i], y, c)
		sub, cc := bits.Sub(uint(z[i]), uint(lo), 0)
		c, z[i] = Word(cc), Word(sub)
		c += hi
	}
	return
}

// Nat represents an arbitrary sized natural number.
//
// Different methods on Nats will talk about a "capacity". The capacity represents
// the announced size of some number. Operations may vary in time *only* relative
// to this capacity, and not to the actual value of the number.
//
// The capacity of a number is usually inherited through whatever method was used to
// create the number in the first place.
type Nat struct {
	// Two natural numbers are not allowed to share the same slice. This allows
	// us to use pointer comparison to check that Nats don't alias eachother
	limbs   []Word
	reduced *Modulus
}

// ensureLimbCapacity makes sure that a Nat has capacity for a certain number of limbs
//
// This will modify the slice contained inside the natural, but won't change the size of
// the slice, so it doesn't affect the value of the natural.
//
// LEAK: Probably the current number of limbs, and size
// OK: both of these should be public
func (z *Nat) ensureLimbCapacity(size int) {
	if cap(z.limbs) < size {
		newLimbs := make([]Word, len(z.limbs), size)
		copy(newLimbs, z.limbs)
		z.limbs = newLimbs
	}
}

// resizedLimbs returns a slice of limbs with size lengths
//
// LEAK: the current number of limbs, and size
// OK: both are public
func (z *Nat) resizedLimbs(size int) []Word {
	z.ensureLimbCapacity(size)
	res := z.limbs[:size]
	// Make sure that the expansion (if any) is cleared
	for i := len(z.limbs); i < size; i++ {
		res[i] = 0
	}
	return res
}

// unaliasedLimbs returns a set of limbs for z, such that they do not alias those of x
//
// This will create a copy of the limbs, if necessary.
//
// LEAK: the size of z, whether or not z and x are the same Nat
func (z *Nat) unaliasedLimbs(x *Nat) []Word {
	res := z.limbs
	if z == x {
		res = make([]Word, len(z.limbs))
		copy(res, z.limbs)
	}
	return res
}

// trueSize calculates the actual size necessary for representing these limbs
//
// This is the size with leading zeros removed. This leaks the number
// of such zeros, but nothing else.
func trueSize(limbs []Word) uint {
	var size uint
	// Instead of checking == 0 directly, which may leak the value, we instead
	// compare with zero in constant time, and check if that succeeded in a leaky way.
	for size = uint(len(limbs)); size > 0 && ctEq(limbs[size-1], 0) == 1; size-- {
	}
	return size
}

// AnnouncedLen returns the number of bits this number is publicly known to have
func (z *Nat) AnnouncedLen() uint {
	return uint(len(z.limbs)) * _W
}

// TrueLen calculates the exact number of bits needed to represent z
//
// This function will leak this value, and violates the standard contract
// around Nats and announced length. For most purposes, `AnnouncedLen` should
// be used instead.
func (z *Nat) TrueLen() uint {
	limbSize := trueSize(z.limbs)
	size := limbSize * _W
	if limbSize > 0 {
		size -= uint(bits.LeadingZeros(uint(z.limbs[limbSize-1])))
	}
	return size
}

// FillBytes writes out the big endian bytes of a natural number.
//
// This will always write out the full capacity of the number, without
// any kind trimming.
//
// This will panic if the buffer's length cannot accomodate the capacity of the number
func (z *Nat) FillBytes(buf []byte) []byte {
	length := len(z.limbs) * _S
	i := length
	// LEAK: Number of limbs
	// OK: The number of limbs is public
	// LEAK: The addresses touched in the out array
	// OK: Every member of out is touched
	for _, x := range z.limbs {
		y := x
		for j := 0; j < _S; j++ {
			i--
			buf[i] = byte(y)
			y >>= 8
		}
	}
	return buf
}

// extendFront pads the front of a slice to a certain size
//
// LEAK: the length of the buffer, size
func extendFront(buf []byte, size int) []byte {
	// LEAK: the length of the buffer
	if len(buf) >= size {
		return buf
	}

	shift := size - len(buf)
	// LEAK: the capacity of the buffer
	// OK: assuming the capacity of the buffer is related to the length,
	// and the length is ok to leak
	if cap(buf) < size {
		newBuf := make([]byte, size)
		copy(newBuf[shift:], buf)
		return newBuf
	}

	newBuf := buf[:size]
	copy(newBuf[shift:], buf)
	for i := 0; i < shift; i++ {
		newBuf[i] = 0
	}
	return newBuf
}

// SetBytes interprets a number in big-endian format, stores it in z, and returns z.
//
// The exact length of the buffer must be public information! This length also dictates
// the capacity of the number returned, and thus the resulting timings for operations
// involving that number.
func (z *Nat) SetBytes(buf []byte) *Nat {
	z.reduced = nil
	// We pad the front so that we have a multiple of _S
	// Padding the front is adding extra zeros to the BE representation
	necessary := (len(buf) + _S - 1) &^ (_S - 1)
	// LEAK: the size of buf
	// OK: this is public information
	buf = extendFront(buf, necessary)
	limbCount := necessary / _S
	// LEAK: limbCount
	// OK: this is derived from the length of buf, which is public
	z.limbs = z.resizedLimbs(limbCount)
	j := necessary
	// LEAK: The number of limbs
	// OK: This is public information
	for i := 0; i < limbCount; i++ {
		z.limbs[i] = 0
		j -= _S
		for k := 0; k < _S; k++ {
			z.limbs[i] <<= 8
			z.limbs[i] |= Word(buf[j+k])
		}
	}
	return z
}

// Bytes creates a slice containing the contents of this Nat, in big endian
//
// This will always fill the output byte slice based on the announced length of this Nat.
func (z *Nat) Bytes() []byte {
	length := len(z.limbs) * _S
	out := make([]byte, length)
	return z.FillBytes(out)
}

// Byte will access the ith byte in this nat, with 0 being the least significant byte.
//
// This will leak the value of i, and panic if i is < 0.
func (z *Nat) Byte(i int) byte {
	if i < 0 {
		panic("negative byte")
	}
	limbCount := len(z.limbs)
	bytesPerLimb := _W / 8
	if i >= bytesPerLimb*limbCount {
		return 0
	}
	return byte(z.limbs[i/bytesPerLimb] >> (8 * (i % bytesPerLimb)))
}

// Big converts a Nat into a big.Int
//
// This will leak information about the true size of z, so caution
// should be exercised when using this method with sensitive values.
func (z *Nat) Big() *big.Int {
	res := new(big.Int)
	// Unfortunate that there's no good way to handle this
	bigLimbs := make([]big.Word, len(z.limbs))
	for i := 0; i < len(bigLimbs) && i < len(z.limbs); i++ {
		bigLimbs[i] = big.Word(z.limbs[i])
	}
	res.SetBits(bigLimbs)
	return res
}

// SetBig modifies z to contain the value of x
//
// The size parameter is used to pad or truncate z to a certain number of bits.
func (z *Nat) SetBig(x *big.Int, size uint) *Nat {
	limbCount := int((size + _W - 1) / _W)
	z.limbs = z.resizedLimbs(limbCount)
	bigLimbs := x.Bits()
	for i := 0; i < len(z.limbs) && i < len(bigLimbs); i++ {
		z.limbs[i] = Word(bigLimbs[i])
	}
	return z
}

// SetUint64 sets z to x, and returns z
//
// This will have the exact same capacity as a 64 bit number
func (z *Nat) SetUint64(x uint64) *Nat {
	z.reduced = nil
	// LEAK: Whether or not _W == 64
	// OK: This is known in advance based on the architecture
	if _W == 64 {
		z.limbs = z.resizedLimbs(1)
		z.limbs[0] = Word(x)
	} else {
		// This works since _W is a power of 2
		limbCount := 64 / _W
		z.limbs = z.resizedLimbs(limbCount)
		// LEAK: limbCount
		// OK: This is public
		for i := 0; i < limbCount; i++ {
			z.limbs[i] = Word(x)
			x >>= _W
		}
	}
	return z
}

// Uint64 represents this number as uint64
//
// The behavior of this function is undefined if the announced length of z is > 64.
func (z *Nat) Uint64() uint64 {
	if len(z.limbs) < 1 {
		return 0
	}
	if _W == 64 {
		return uint64(z.limbs[0])
	} else {
		var ret uint64
		for i := len(z.limbs) - 1; i >= 0; i-- {
			ret = (ret << _W) | uint64(z.limbs[i])
		}
		return ret
	}
}

// SetNat copies the value of x into z
//
// z will have the same announced length as x.
func (z *Nat) SetNat(x *Nat) *Nat {
	z.limbs = z.resizedLimbs(len(x.limbs))
	copy(z.limbs, x.limbs)
	z.reduced = x.reduced
	return z
}

// Modulus represents a natural number used for modular reduction
//
// Unlike with natural numbers, the number of bits need to contain the modulus
// is assumed to be public. Operations are allowed to leak this size, and creating
// a modulus will remove unnecessary zeros.
//
// Operations on a Modulus may leak whether or not a Modulus is even.
type Modulus struct {
	nat Nat
	// the number of leading zero bits
	leading uint
	// The inverse of the least significant limb, modulo W
	m0inv Word
	// If true, then this modulus is even
	even bool
}

// invertModW calculates x^-1 mod _W
func invertModW(x Word) Word {
	y := x
	// This is enough for 64 bits, and the extra iteration is not that costly for 32
	for i := 0; i < 5; i++ {
		y = y * (2 - x*y)
	}
	return y
}

// precomputeValues calculates the desirable modulus fields in advance
//
// This sets the leading number of bits, leaking the true bit size of m,
// as well as the inverse of the least significant limb (without leaking it).
//
// This will also do integrity checks, namely that the modulus isn't empty or even
func (m *Modulus) precomputeValues() {
	if len(m.nat.limbs) < 1 {
		panic("Modulus is empty")
	}
	m.leading = uint(bits.LeadingZeros(uint(m.nat.limbs[len(m.nat.limbs)-1])))
	// I think checking the bit directly might leak more data than we'd like
	m.even = ctEq(m.nat.limbs[0]&1, 0) == 1
	// There's no point calculating this if m isn't even, and we can leak evenness
	if !m.even {
		m.m0inv = invertModW(m.nat.limbs[0])
		m.m0inv = -m.m0inv
	}
}

// SetUint64 sets the modulus according to an integer
func ModulusFromUint64(x uint64) *Modulus {
	var m Modulus
	m.nat.SetUint64(x)
	// edge case for 32 bit limb size
	if _W < 64 && len(m.nat.limbs) > 1 && m.nat.limbs[1] == 0 {
		m.nat.limbs = m.nat.limbs[:1]
	}
	m.precomputeValues()
	return &m
}

// FromBytes creates a new Modulus, converting from big endian bytes
//
// This function will remove leading zeros, thus leaking the true size of the modulus.
// See the documentation for the Modulus type, for more information about this contract.
func ModulusFromBytes(bytes []byte) *Modulus {
	var m Modulus
	// TODO: You could allocate a smaller buffer to begin with, versus using the Nat method
	m.nat.SetBytes(bytes)

	m.nat.limbs = m.nat.limbs[:trueSize(m.nat.limbs)]
	m.precomputeValues()
	return &m
}

// FromNat creates a new Modulus, using the value of a Nat
//
// This will leak the true size of this natural number. Because of this,
// the true size of the number should not be sensitive information. This is
// a stronger requirement than we usually have for Nat.
func ModulusFromNat(nat *Nat) *Modulus {
	var m Modulus
	// We make a copy here, to avoid any aliasing between buffers
	size := trueSize(nat.limbs)
	m.nat.limbs = m.nat.resizedLimbs(int(size))
	copy(m.nat.limbs, nat.limbs)
	m.precomputeValues()
	return &m
}

// Bytes returns the big endian bytes making up the modulus
func (m *Modulus) Bytes() []byte {
	return m.nat.Bytes()
}

// BitLen returns the exact number of bits used to store this Modulus
//
// Moduli are allowed to leak this value.
func (m *Modulus) BitLen() uint {
	return uint(len(m.nat.limbs))*_W - m.leading
}

// Cmp compares two moduli, returning -1 if m < n, 0 if m == n, and 1 if m > n
//
// This will not leak information about the value of this comparison.
func (m *Modulus) Cmp(n *Modulus) int {
	return m.nat.Cmp(&n.nat)
}

// shiftAddIn calculates z = z << _W + x mod m
//
// The length of z and scratch should be len(m) + 1
func shiftAddIn(z, scratch []Word, x Word, m *Modulus) (q Word) {
	// Making tests on the exact bit length of m is ok,
	// since that's part of the contract for moduli
	size := len(m.nat.limbs)
	if size == 0 {
		return
	}
	if size == 1 {
		// In this case, z:x (/, %) m is exactly what we need to calculate
		q, r := div(z[0], x, m.nat.limbs[0])
		z[0] = r
		return q
	}

	// The idea is as follows:
	//
	// We want to shift x into z, and then divide by m. Instead of dividing by
	// m, we can get a good estimate, using the top two 2 * _W bits of z, and the
	// top _W bits of m. These are stored in a1:a0, and b0 respectively.

	// We need to keep around the top word of z, pre-shifting
	hi := z[size-1]

	a1 := (z[size-1] << m.leading) | (z[size-2] >> (_W - m.leading))
	// The actual shift can be performed by moving the limbs of z up, then inserting x
	for i := size - 1; i > 0; i-- {
		z[i] = z[i-1]
	}
	z[0] = x
	a0 := (z[size-1] << m.leading) | (z[size-2] >> (_W - m.leading))
	b0 := (m.nat.limbs[size-1] << m.leading) | (m.nat.limbs[size-2] >> (_W - m.leading))

	// We want to use a1:a0 / b0 - 1 as our estimate. If rawQ is 0, we should
	// use 0 as our estimate. Another edge case when an overflow happens in the quotient.
	// It can be shown that this happens when a1 == b0. In this case, we want
	// to use the maximum value for q
	rawQ, _ := div(a1, a0, b0)
	q = ctIfElse(ctEq(a1, b0), ^Word(0), ctIfElse(ctEq(rawQ, 0), 0, rawQ-1))

	// This estimate is off by +- 1, so we subtract q * m, and then either add
	// or subtract m, based on the result.
	c := mulSubVVW(z, m.nat.limbs, q)
	// If the carry from subtraction is greater than the limb of z we've shifted out,
	// then we've underflowed, and need to add in m
	under := ctGt(c, hi)
	// For us to be too large, we first need to not be too low, as per the previous flag.
	// Then, if the lower limbs of z are still larger, or the top limb of z is equal to the carry,
	// we can conclude that we're too large, and need to subtract m
	stillBigger := cmpGeq(z, m.nat.limbs)
	over := (1 ^ under) & (stillBigger | (1 ^ ctEq(c, hi)))
	addVV(scratch, z, m.nat.limbs)
	ctCondCopy(under, z, scratch)
	q = ctIfElse(under, q-1, q)
	subVV(scratch, z, m.nat.limbs)
	ctCondCopy(over, z, scratch)
	q = ctIfElse(over, q+1, q)
	return
}

// Mod calculates z <- x mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) Mod(x *Nat, m *Modulus) *Nat {
	if x.reduced == m {
		xLimbs := x.unaliasedLimbs(z)
		z.limbs = z.resizedLimbs(len(m.nat.limbs))
		z.reduced = m
		copy(z.limbs, xLimbs)
		return z
	}
	size := len(m.nat.limbs)
	xLimbs := x.unaliasedLimbs(z)
	z.limbs = z.resizedLimbs(2 * size)
	for i := 0; i < len(z.limbs); i++ {
		z.limbs[i] = 0
	}
	// Multiple times in this section:
	// LEAK: the length of x
	// OK: this is public information
	i := len(xLimbs) - 1
	// We can inject at least size - 1 limbs while staying under m
	// Thus, we start injecting from index size - 2
	start := size - 2
	// That is, if there are at least that many limbs to choose from
	if i < start {
		start = i
	}
	for j := start; j >= 0; j-- {
		z.limbs[j] = xLimbs[i]
		i--
	}
	// We shift in the remaining limbs, making sure to reduce modulo M each time
	for ; i >= 0; i-- {
		shiftAddIn(z.limbs[:size], z.limbs[size:], xLimbs[i], m)
	}
	z.limbs = z.limbs[:size]
	z.reduced = m
	return z
}

// Div calculates z <- x / m, with m a Modulus.
//
// This might seem like an odd signature, but by using a Modulus,
// we can achieve the same speed as the Mod method. This wouldn't be the case for
// an arbitrary Nat.
//
// cap determines the number of bits to keep in the result.
func (z *Nat) Div(x *Nat, m *Modulus, cap uint) *Nat {
	if len(x.limbs) < len(m.nat.limbs) || x.reduced == m {
		z.limbs = z.resizedLimbs(int((cap + _W - 1) / _W))
		for i := 0; i < len(z.limbs); i++ {
			z.limbs[i] = 0
		}
		return z
	}

	size := len(m.nat.limbs)

	xLimbs := x.unaliasedLimbs(z)

	// Enough for 2 buffers the size of m, and to store the full quotient
	z.limbs = z.resizedLimbs(2*size + len(xLimbs))

	remainder := z.limbs[:size]
	scratch := z.limbs[size : 2*size]
	// Our full quotient, in big endian order.
	quotientBE := z.limbs[2*size:]
	// We use this to append without actually reallocating. We fill our quotient
	// in from 0 upwards.
	qI := 0

	i := len(xLimbs) - 1
	// We can inject at least size - 1 limbs while staying under m
	// Thus, we start injecting from index size - 2
	start := size - 2
	// That is, if there are at least that many limbs to choose from
	if i < start {
		start = i
	}
	for j := start; j >= 0; j-- {
		remainder[j] = xLimbs[i]
		i--
		quotientBE[qI] = 0
		qI++
	}

	for ; i >= 0; i-- {
		q := shiftAddIn(remainder, scratch, xLimbs[i], m)
		quotientBE[qI] = q
		qI++
	}

	limbCount := int((cap + _W - 1) / _W)

	z.limbs = z.resizedLimbs(limbCount)
	// First, reverse all the limbs we want, from the last part of the buffer we used.
	for i := 0; i < len(z.limbs) && i < len(quotientBE); i++ {
		z.limbs[i] = quotientBE[qI-i-1]
	}
	// Now, we need to truncate the last limb
	extraBits := uint(_W*limbCount) - cap
	bitsToKeep := _W - extraBits
	mask := ^(^Word(0) << bitsToKeep)
	// LEAK: the size of z (since we're making an extra access at the end)
	// OK: this is public information, since cap is public
	z.limbs[len(z.limbs)-1] &= mask

	z.reduced = nil
	return z
}

// ModAdd calculates z <- x + y mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) ModAdd(x *Nat, y *Nat, m *Modulus) *Nat {
	var xModM, yModM Nat
	// This is necessary for the correctness of the algorithm, since
	// we don't assume that x and y are in range.
	// Furthermore, we can now assume that x and y have the same number
	// of limbs as m
	xModM.Mod(x, m)
	yModM.Mod(y, m)

	// The only thing we have to resize is z, everything else has m's length
	size := len(m.nat.limbs)
	scratch := z.resizedLimbs(2 * size)
	z.limbs = scratch[:size]
	subResult := scratch[size:]

	addCarry := addVV(z.limbs, xModM.limbs, yModM.limbs)
	subCarry := subVV(subResult, z.limbs, m.nat.limbs)
	// Three cases are possible:
	//
	// addCarry, subCarry = 0 -> subResult
	// 	 we didn't overflow our buffer, but our result was big
	//   enough to subtract m without underflow, so it was larger than m
	// addCarry, subCarry = 1 -> subResult
	//   we overflowed the buffer, and the subtraction of m is correct,
	//   because our result only looks too small because of the missing carry bit
	// addCarry = 0, subCarry = 1 -> addResult
	// 	 we didn't overflow our buffer, and the subtraction of m is wrong,
	//   because our result was already smaller than m
	// The other case is impossible, because it would mean we have a result big
	// enough to both overflow the addition by at least m. But, we made sure that
	// x and y are at most m - 1, so this isn't possible.
	selectSub := ctEq(addCarry, subCarry)
	ctCondCopy(selectSub, z.limbs[:size], subResult)
	z.reduced = m
	return z
}

func (z *Nat) ModSub(x *Nat, y *Nat, m *Modulus) *Nat {
	var xModM, yModM Nat
	// First reduce x and y mod m
	xModM.Mod(x, m)
	yModM.Mod(y, m)

	size := len(m.nat.limbs)
	scratch := z.resizedLimbs(2 * size)
	z.limbs = scratch[:size]
	addResult := scratch[size:]

	subCarry := subVV(z.limbs, xModM.limbs, yModM.limbs)
	underflow := ctEq(subCarry, 1)
	addVV(addResult, z.limbs, m.nat.limbs)
	ctCondCopy(underflow, z.limbs, addResult)
	z.reduced = m
	return z
}

// Add calculates z <- x + y, modulo 2^cap
//
// The capacity is given in bits, and also controls the size of the result.
func (z *Nat) Add(x *Nat, y *Nat, cap uint) *Nat {
	limbCount := int((cap + _W - 1) / _W)
	xLimbs := x.resizedLimbs(limbCount)
	yLimbs := y.resizedLimbs(limbCount)
	z.limbs = z.resizedLimbs(limbCount)
	addVV(z.limbs, xLimbs, yLimbs)
	// Now, we need to truncate the last limb
	bitsToKeep := cap % _W
	// Checking a function of the cap is ok
	if bitsToKeep == 0 {
		return z
	}
	mask := ^(^Word(0) << bitsToKeep)
	// LEAK: the size of z (since we're making an extra access at the end)
	// OK: this is public information, since cap is public
	z.limbs[len(z.limbs)-1] &= mask
	z.reduced = nil
	return z
}

// Sub calculates z <- x - y, modulo 2^cap
//
// The capacity is given in bits, and also controls the size of the result.
func (z *Nat) Sub(x *Nat, y *Nat, cap uint) *Nat {
	limbCount := int((cap + _W - 1) / _W)
	xLimbs := x.resizedLimbs(limbCount)
	yLimbs := y.resizedLimbs(limbCount)
	z.limbs = z.resizedLimbs(limbCount)
	subVV(z.limbs, xLimbs, yLimbs)
	// Now, we need to truncate the last limb
	bitsToKeep := cap % _W
	// Checking a function of the cap is ok
	if bitsToKeep == 0 {
		return z
	}
	mask := ^(^Word(0) << bitsToKeep)
	// LEAK: the size of z (since we're making an extra access at the end)
	// OK: this is public information, since cap is public
	z.limbs[len(z.limbs)-1] &= mask
	z.reduced = nil
	return z
}

// montgomeryRepresentation calculates zR mod m
func montgomeryRepresentation(z []Word, scratch []Word, m *Modulus) {
	// Our strategy is to shift by W, n times, each time reducing modulo m
	size := len(m.nat.limbs)
	// LEAK: the size of the modulus
	// OK: this is public
	for i := 0; i < size; i++ {
		shiftAddIn(z, scratch, 0, m)
	}
}

// You might have the urge to replace this with []Word, and use the routines
// that already exist for doing operations. This would be a mistake.
// Go doesn't seem to be able to optimize and inline slice operations nearly as
// well as it can for this little type. Attempts to replace this struct with a
// slice were an order of magnitude slower (as per the exponentiation operation)
type triple struct {
	w0 Word
	w1 Word
	w2 Word
}

func (a *triple) add(b triple) {
	w0, c0 := bits.Add(uint(a.w0), uint(b.w0), 0)
	w1, c1 := bits.Add(uint(a.w1), uint(b.w1), c0)
	w2, _ := bits.Add(uint(a.w2), uint(b.w2), c1)
	a.w0 = Word(w0)
	a.w1 = Word(w1)
	a.w2 = Word(w2)
}

func tripleFromMul(a Word, b Word) triple {
	// You might be tempted to use mulWW here, but for some reason, Go cannot
	// figure out how to inline that assembly routine, but using bits.Mul directly
	// gets inlined by the compiler into effectively the same assembly.
	//
	// Beats me.
	w1, w0 := bits.Mul(uint(a), uint(b))
	return triple{w0: Word(w0), w1: Word(w1), w2: 0}
}

// montgomeryMul performs z <- xy / R mod m
//
// LEAK: the size of the modulus
//
// out, x, y must have the same length as the modulus, and be reduced already.
//
// out can alias x and y, but not scratch
func montgomeryMul(x []Word, y []Word, out []Word, scratch []Word, m *Modulus) {
	size := len(m.nat.limbs)

	for i := 0; i < size; i++ {
		scratch[i] = 0
	}
	dh := Word(0)
	for i := 0; i < size; i++ {
		f := (scratch[0] + x[i]*y[0]) * m.m0inv
		var c triple
		for j := 0; j < size; j++ {
			z := triple{w0: scratch[j], w1: 0, w2: 0}
			z.add(tripleFromMul(x[i], y[j]))
			z.add(tripleFromMul(f, m.nat.limbs[j]))
			z.add(c)
			if j > 0 {
				scratch[j-1] = z.w0
			}
			c.w0 = z.w1
			c.w1 = z.w2
		}
		z := triple{w0: dh, w1: 0, w2: 0}
		z.add(c)
		scratch[size-1] = z.w0
		dh = z.w1
	}
	c := subVV(out, scratch, m.nat.limbs)
	ctCondCopy(1^ctEq(dh, c), out, scratch)
}

// ModMul calculates z <- x * y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModMul(x *Nat, y *Nat, m *Modulus) *Nat {
	xModM := new(Nat).Mod(x, m)
	yModM := new(Nat).Mod(y, m)
	bitLen := m.BitLen()
	z.Mul(xModM, yModM, 2*bitLen)
	return z.Mod(z, m)
}

// Mul calculates z <- x * y, modulo 2^cap
//
// The capacity is given in bits, and also controls the size of the result.
func (z *Nat) Mul(x *Nat, y *Nat, cap uint) *Nat {
	limbCount := int((cap + _W - 1) / _W)
	// Since we neex to set z to zero, we have no choice to use a new buffer,
	// because we allow z to alias either of the arguments
	zLimbs := make([]Word, limbCount)
	xLimbs := x.resizedLimbs(limbCount)
	yLimbs := y.resizedLimbs(limbCount)
	// LEAK: limbCount
	// OK: the capacity is public, or should be
	for i := 0; i < limbCount; i++ {
		addMulVVW(zLimbs[i:], xLimbs, yLimbs[i])
	}
	// Now, we need to truncate the last limb
	extraBits := uint(_W*limbCount) - cap
	bitsToKeep := _W - extraBits
	mask := ^(^Word(0) << bitsToKeep)
	// LEAK: the size of z (since we're making an extra access at the end)
	// OK: this is public information, since cap is public
	zLimbs[len(zLimbs)-1] &= mask
	// Now we can write over
	z.limbs = zLimbs
	z.reduced = nil
	return z
}

func (z *Nat) expOdd(x *Nat, y *Nat, m *Modulus) *Nat {
	size := len(m.nat.limbs)

	xModM := new(Nat).Mod(x, m)
	yLimbs := y.unaliasedLimbs(z)

	scratch := z.resizedLimbs(18 * size)
	scratch1 := scratch[16*size : 17*size]
	scratch2 := scratch[17*size:]

	z.limbs = scratch[:size]
	for i := 0; i < size; i++ {
		z.limbs[i] = 0
	}
	z.limbs[0] = 1
	montgomeryRepresentation(z.limbs, scratch1, m)

	x1 := scratch[size : 2*size]
	copy(x1, xModM.limbs)
	montgomeryRepresentation(scratch[size:2*size], scratch1, m)
	for i := 2; i < 16; i++ {
		ximinus1 := scratch[(i-1)*size : i*size]
		xi := scratch[i*size : (i+1)*size]
		montgomeryMul(ximinus1, x1, xi, scratch1, m)
	}

	// LEAK: y's length
	// OK: this should be public
	for i := len(yLimbs) - 1; i >= 0; i-- {
		yi := yLimbs[i]
		for j := _W - 4; j >= 0; j -= 4 {
			montgomeryMul(z.limbs, z.limbs, z.limbs, scratch1, m)
			montgomeryMul(z.limbs, z.limbs, z.limbs, scratch1, m)
			montgomeryMul(z.limbs, z.limbs, z.limbs, scratch1, m)
			montgomeryMul(z.limbs, z.limbs, z.limbs, scratch1, m)

			window := (yi >> j) & 0b1111
			for i := 1; i < 16; i++ {
				xToI := scratch[i*size : (i+1)*size]
				ctCondCopy(ctEq(window, Word(i)), scratch1, xToI)
			}
			montgomeryMul(z.limbs, scratch1, scratch1, scratch2, m)
			ctCondCopy(1^ctEq(window, 0), z.limbs, scratch1)
		}
	}
	for i := 0; i < size; i++ {
		scratch2[i] = 0
	}
	scratch2[0] = 1
	montgomeryMul(z.limbs, scratch2, z.limbs, scratch1, m)
	z.reduced = m
	return z
}

func (z *Nat) expEven(x *Nat, y *Nat, m *Modulus) *Nat {
	xModM := new(Nat).Mod(x, m)
	yLimbs := y.unaliasedLimbs(z)

	scratch := new(Nat)

	// LEAK: y's length
	// OK: this should be public
	for i := len(yLimbs) - 1; i >= 0; i-- {
		yi := yLimbs[i]
		for j := _W; j >= 0; j-- {
			z.ModMul(z, z, m)

			sel := choice((yi >> j) & 1)
			scratch.ModMul(z, xModM, m)
			ctCondCopy(sel, z.limbs, scratch.limbs)
		}
	}
	z.reduced = m
	return z
}

// Exp calculates z <- x^y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) Exp(x *Nat, y *Nat, m *Modulus) *Nat {
	if m.even {
		return z.expEven(x, y, m)
	} else {
		return z.expOdd(x, y, m)
	}
}

// cmpEq compares two limbs (same size) returning 1 if x >= y, and 0 otherwise
func cmpEq(x []Word, y []Word) choice {
	res := choice(1)
	for i := 0; i < len(x) && i < len(y); i++ {
		res &= ctEq(x[i], y[i])
	}
	return res
}

// cmpGeq compares two limbs (same size) returning 1 if x >= y, and 0 otherwise
func cmpGeq(x []Word, y []Word) choice {
	var c uint
	for i := 0; i < len(x) && i < len(y); i++ {
		_, c = bits.Sub(uint(x[i]), uint(y[i]), c)
	}
	return 1 ^ choice(c)
}

// cmpZero checks if a slice is equal to zero, in constant time
//
// LEAK: the length of a
func cmpZero(a []Word) choice {
	var v Word
	for i := 0; i < len(a); i++ {
		v |= a[i]
	}
	return ctEq(v, 0)
}

// Cmp compares two natural numbers, returning 0 if equal, -1 if z < x, and 1 if z > x
//
// This leaks nothing about the values of the numbers, except for their rspective
// announced lengths.
//
// Be mindful that even though the comparison is done without leaks, branching on the result
// can introduce a leak of this result. Because of this, be careful with how you use
// this function.
func (z *Nat) Cmp(x *Nat) int {
	// Rough Idea: Resize both slices to the maximum length, then compare
	// using that length

	// LEAK: z's length, x's length, the maximum
	// OK: These should be public information
	size := len(x.limbs)
	zLen := len(z.limbs)
	if zLen > size {
		size = zLen
	}
	zLimbs := z.resizedLimbs(size)
	xLimbs := x.resizedLimbs(size)

	eq := choice(1)
	geq := choice(0)
	for i := 0; i < len(zLimbs) && i < len(xLimbs); i++ {
		eq_at_i := ctEq(zLimbs[i], xLimbs[i])
		eq &= eq_at_i
		geq = (eq_at_i & geq) | ((1 ^ eq_at_i) & ctGt(zLimbs[i], xLimbs[i]))
	}
	return subtle.ConstantTimeSelect(int(eq), 0, subtle.ConstantTimeSelect(int(geq), 1, -1))
}

// CmpMod compares this natural number with a modulus, returning 0 if z == m, -1 if z < m, and 1 if z > m
//
// This doesn't leak anything about the values of the numbers, but the same warnings as for Cmp apply.
func (z *Nat) CmpMod(m *Modulus) int {
	return z.Cmp(&m.nat)
}

// EqZero compares z to 0, returning true if they're equal
//
// While calling this function will not leak whether or not this number is equal to zero,
// an if condition using its result can leak this information. Because of this,
// be mindful to not misuse this function.
func (z *Nat) EqZero() bool {
	return cmpZero(z.limbs) == 1
}

// eGCD calculates and returns d := gcd(x, m) and v s.t. vx = d mod m
//
// This function assumes that m is and odd number, but doesn't assume
// that m is truncated to its full size.
//
// The slices returned should be copied into the result, and not used
// directly, for aliasing reasons.
//
// The recipient Nat is used only for scratch space.
func (z *Nat) eGCD(x []Word, m []Word) ([]Word, []Word) {
	size := len(m)

	scratch := z.resizedLimbs(8 * size)
	v := scratch[:size]
	u := scratch[size : 2*size]
	b := scratch[2*size : 3*size]
	a := scratch[3*size : 4*size]
	halfm := scratch[4*size : 5*size+1]
	a1 := scratch[5*size : 6*size]
	u1 := scratch[6*size : 7*size]
	u2 := scratch[7*size:]

	// a = x
	copy(a, x)
	// v = 0
	// u = 1
	for i := 0; i < size; i++ {
		u[i] = 0
		v[i] = 0
	}
	u[0] = 1

	// halfm = (m + 1) / 2
	halfm[size] = addVW(halfm, m, 1)
	shrVU(halfm, halfm, 1)
	halfm = halfm[:size]

	copy(b, m)

	// Idea:
	//
	// while a != 0:
	//   if a is even:
	//	   a = a / 2
	//     u = (u / 2) mod m
	//   else:
	//     if a < b:
	//       swap(a, b)
	//       swap(u, v)
	//     a = (a - b) / 2
	//     u = (u - v) / 2 mod m
	//
	// We run for 2 * k - 1 iterations, with k the number of bits of the modulus
	for i := 0; i < 2*_W*size-1; i++ {
		// a1 and u2 will hold the results to use if a is even
		aOdd := choice(shrVU(a1, a, 1) >> (_W - 1))
		aEven := 1 ^ aOdd
		uOdd := choice(shrVU(u2, u, 1) >> (_W - 1))
		addVV(u1, u2, halfm)
		ctCondCopy(uOdd, u2, u1)

		// Now we calculate the results if a is not even, which may get overwritten later
		aSmaller := 1 ^ cmpGeq(a, b)
		swap := aOdd & aSmaller
		ctCondSwap(swap, a, b)
		ctCondSwap(swap, u, v)

		subVV(a, a, b)
		shrVU(a, a, 1)
		// u = (u - v) / 2 mod m
		subCarry := choice(subVV(u, u, v))
		addVV(u1, u, m)
		ctCondCopy(subCarry, u, u1)
		uOdd = choice(shrVU(u, u, 1) >> (_W - 1))
		addVV(u1, u, halfm)
		ctCondCopy(uOdd, u, u1)

		// If a was indeed even, we use the results we produced earlier
		ctCondCopy(aEven, a, a1)
		ctCondCopy(aEven, u, u2)
	}

	return b, v
}

// Coprime returns 1 if gcd(x, y) == 1, and 0 otherwise
func (x *Nat) Coprime(y *Nat) int {
	size := len(x.limbs)
	if len(y.limbs) > size {
		size = len(y.limbs)
	}
	a := x.resizedLimbs(size)
	b := y.resizedLimbs(size)

	// Our gcd(a, b) routine requires b to be odd, and will return garbage otherwise.
	aOdd := choice(a[0] & 1)
	ctCondSwap(aOdd, a, b)

	scratch := new(Nat)
	d, _ := scratch.eGCD(a, b)

	scratch.SetUint64(1)
	one := scratch.resizedLimbs(len(d))

	bOdd := choice(b[0] & 1)
	// If at least one of a or b is odd, then our GCD calculation will have been correct,
	// otherwise, both are even, so we want to return false anyways.
	return int((aOdd | bOdd) & cmpEq(d, one))
}

// modInverse calculates the inverse of a reduced x modulo m
//
// This assumes that m is an odd number, but not that it's truncated
// to its true size. This routine will only leak the announced sizes of
// x and m.
//
// We also assume that x is already reduced modulo m
func (z *Nat) modInverse(x *Nat, m *Nat) *Nat {
	size := len(m.limbs)
	// Make sure that z doesn't alias either of m or x
	xLimbs := x.unaliasedLimbs(z)
	mLimbs := m.unaliasedLimbs(z)
	_, v := z.eGCD(xLimbs, mLimbs)
	z.limbs = z.resizedLimbs(size)
	copy(z.limbs, v)
	return z
}

// ModInverse calculates z <- x^-1 mod m
//
// This will produce nonsense if the modulus is even.
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModInverse(x *Nat, m *Modulus) *Nat {
	z.Mod(x, m)
	if m.even {
		return z.modInverseEven(x, m)
	} else {
		return z.modInverse(z, &m.nat)
	}
}

// divDouble divides x by d, outputtting the quotient in out, and a remainder
//
// This routine assumes nothing about the padding of either of its inputs, and
// leaks nothing beyond their announced length.
//
// If out is not empty, it's assumed that x has at most twice the bit length of d,
// and the quotient can thus fit in a slice the length of d, which out is assumed to be.
//
// If out is empty, no quotient is produced, but the remainder is still calculated.
// This remainder will be correct regardless of the size difference between x and d.
func divDouble(x []Word, d []Word, out []Word) []Word {
	size := len(d)
	r := make([]Word, size)
	scratch := make([]Word, size)

	// We use free injection, like in Mod
	i := len(x) - 1
	// We can inject at least size - 1 limbs while staying under m
	// Thus, we start injecting from index size - 2
	start := size - 2
	// That is, if there are at least that many limbs to choose from
	if i < start {
		start = i
	}
	for j := start; j >= 0; j-- {
		r[j] = x[i]
		i--
	}

	for ; i >= 0; i-- {
		// out can alias x, because we recover x[i] early
		xi := x[i]
		// Hopefully the branch predictor can make these checks not too expensive,
		// otherwise we'll have to duplicate the routine
		if len(out) > 0 {
			out[i] = 0
		}
		for j := _W - 1; j >= 0; j-- {
			xij := (xi >> j) & 1
			shiftCarry := shlVU(r, r, 1)
			r[0] |= xij
			subCarry := subVV(scratch, r, d)
			sel := ctEq(shiftCarry, subCarry)
			ctCondCopy(sel, r, scratch)
			if len(out) > 0 {
				out[i] = ((out[i] << 1) | Word(sel))
			}
		}
	}
	return r
}

// ModInverseEven calculates the modular inverse of x, mod m
//
// This routine will work even if m is an even number, unlike ModInverse.
// Furthermore, it doesn't require the modulus to be truncated to its true size, and
// will only leak information about the public sizes of its inputs. It is slower
// than the standard routine though.
//
// This function assumes that x has an inverse modulo m, naturally
func (z *Nat) modInverseEven(x *Nat, m *Modulus) *Nat {
	// Idea:
	//
	// You want to find Z such that ZX = 1 mod M. The problem is
	// that the usual routine assumes that m is odd. In this case m is even.
	// For X to be invertible, we need it to be odd. We can thus invert M mod X,
	// finding an A satisfying AM = 1 mod X. This means that AM = 1 + KX, for some
	// positive integer K. Modulo M, this entails that KX = -1 mod M, so -K provides
	// us with an inverse for X.
	//
	// To find K, we can calculate (AM - 1) / X, and then subtract this from M, to get our inverse.
	size := len(m.nat.limbs)
	// We want to invert m modulo x, so we first calculate the reduced version, before inverting
	var newZ Nat
	newZ.limbs = divDouble(m.nat.limbs, x.limbs, []Word{})
	newZ.modInverse(&newZ, x)
	inverseZero := cmpZero(newZ.limbs)
	newZ.Mul(&newZ, &m.nat, uint(2*size*_W))
	newZ.limbs = newZ.resizedLimbs(2 * size)
	subVW(newZ.limbs, newZ.limbs, 1)
	divDouble(newZ.limbs, x.limbs, newZ.limbs)
	// The result fits on a single half of newZ, but we need to subtract it from m.
	// We can use the other half of newZ, and then copy it back over if we need to keep it
	subVV(newZ.limbs[size:], m.nat.limbs, newZ.limbs[:size])
	// If the inverse was zero, then x was 1, and so we should return 1.
	// We go ahead and prepare this result, but expect to copy over the subtraction
	// we just calculated soon over, in the usual case.
	newZ.limbs[0] = 1
	for i := 1; i < size; i++ {
		newZ.limbs[i] = 0
	}
	ctCondCopy(1^inverseZero, newZ.limbs[:size], newZ.limbs[size:])

	z.limbs = newZ.limbs[:size]
	z.reduced = nil
	return z
}

// modSqrt3Mod4 sets z <- sqrt(x) mod p, when p is a prime with p = 3 mod 4
func (z *Nat) modSqrt3Mod4(x *Nat, p *Modulus) *Nat {
	// In this case, we can do x^(p + 1) / 4
	e := new(Nat).SetNat(&p.nat)
	carry := addVW(e.limbs, e.limbs, 1)
	shrVU(e.limbs, e.limbs, 2)
	e.limbs[len(e.limbs)-1] |= (carry << (_W - 2))
	return z.Exp(x, e, p)
}

// tonelliShanks sets z <- sqrt(x) mod p, for any prime modulus
func (z *Nat) tonelliShanks(x *Nat, p *Modulus) *Nat {
	// c.f. https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-hash-to-curve-09#appendix-G.4
	scratch := new(Nat)
	x = new(Nat).SetNat(x)

	one := new(Nat).SetUint64(1)
	trailingZeros := 1
	reducedPminusOne := new(Nat).Sub(&p.nat, one, p.BitLen())
	shrVU(reducedPminusOne.limbs, reducedPminusOne.limbs, 1)

	nonSquare := new(Nat).SetUint64(2)
	for scratch.Exp(nonSquare, reducedPminusOne, p).Cmp(one) == 0 {
		nonSquare.Add(nonSquare, one, p.BitLen())
	}

	for reducedPminusOne.limbs[0]&1 == 0 {
		trailingZeros += 1
		shrVU(reducedPminusOne.limbs, reducedPminusOne.limbs, 1)
	}

	reducedQminusOne := new(Nat).Sub(reducedPminusOne, one, p.BitLen())
	shrVU(reducedQminusOne.limbs, reducedQminusOne.limbs, 1)

	c := new(Nat).Exp(nonSquare, reducedPminusOne, p)

	z.Exp(x, reducedQminusOne, p)
	t := new(Nat).ModMul(z, z, p)
	t.ModMul(t, x, p)
	z.ModMul(z, x, p)
	b := new(Nat).SetNat(t)
	one.limbs = one.resizedLimbs(len(b.limbs))
	for i := trailingZeros; i > 1; i-- {
		for j := 1; j < i-1; j++ {
			b.ModMul(b, b, p)
		}
		sel := 1 ^ cmpEq(b.limbs, one.limbs)
		scratch.ModMul(z, c, p)
		ctCondCopy(sel, z.limbs, scratch.limbs)
		c.ModMul(c, c, p)
		scratch.ModMul(t, c, p)
		ctCondCopy(sel, t.limbs, scratch.limbs)
		b.SetNat(t)
	}
	z.reduced = p
	return z
}

// ModSqrt calculates the square root of x modulo p
//
// p must be a prime number, and x must actually have a square root
// modulo p. The result is undefined if these conditions aren't satisfied
//
// This function will leak information about the value of p. This isn't intended
// to be used in situations where the modulus isn't publicly known.
func (z *Nat) ModSqrt(x *Nat, p *Modulus) *Nat {
	if len(p.nat.limbs) == 0 {
		panic("Can't take square root mod 0")
	}
	if p.nat.limbs[0]&0b11 == 3 {
		return z.modSqrt3Mod4(x, p)
	}
	return z.tonelliShanks(x, p)
}
