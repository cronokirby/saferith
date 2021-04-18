package safenum

import (
	"crypto/subtle"
	"math/bits"
)

// Constant Time Utilities

// ctEq compares x and y for equality, returning 1 if equal, and 0 otherwise
//
// This doesn't leak any information about either of them
func ctEq(x, y Word) Word {
	zero := uint64(x ^ y)
	// The usual trick in Go's subtle library doesn't work for the case where
	// x and y differ in every single bit. Instead, we do the same testing mechanism,
	// but over each "half" of the number
	//
	// I'm not sure if this is optimal.
	hiZero := ((zero >> 32) - 1) >> 63
	loZero := ((zero & 0xFF_FF_FF_FF) - 1) >> 63
	return Word(hiZero & loZero)
}

// ctGt checks x > y, returning 1 or 0
//
// This doesn't leak any information about either of them
func ctGt(x, y Word) Word {
	z := y - x
	return (z ^ ((x ^ y) & (x ^ z))) >> (_W - 1)
}

// ctIfElse selects x if v = 1, and y otherwise
//
// This doesn't leak the value of any of its inputs
func ctIfElse(v, x, y Word) Word {
	// mask should be all 1s if v is 1, otherwise all 0s
	mask := -v
	return y ^ (mask & (y ^ x))
}

// ctCondCopy copies y into x, if v == 1, otherwise does nothing
//
// Both slices must have the same length.
//
// LEAK: the length of the slices
//
// Otherwise, which branch was taken isn't leaked
func ctCondCopy(v Word, x, y []Word) {
	// see ctMux
	mask := -v
	for i := 0; i < len(x); i++ {
		x[i] = x[i] ^ (mask & (x[i] ^ y[i]))
	}
}

// ctCondSwap swaps the contents of a and b, when v == 1, otherwise does nothing
//
// Both slices must have the same length.
//
// LEAK: the length of the slices
//
// Whether or not a swap happened isn't leaked
func ctCondSwap(v Word, a, b []Word) {
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
		sel := ctEq(w, d) | ctGt(w, d) | (hi >> i)
		hi2 := (w - d) >> j
		lo2 := lo - (d << i)
		hi = ctIfElse(sel, hi2, hi)
		lo = ctIfElse(sel, lo2, lo)
		quo |= sel
		quo <<= 1
	}
	sel := ctEq(lo, d) | ctGt(lo, d) | hi
	quo |= sel
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
	limbs []Word
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

// SetUint64 sets z to x, and returns z
//
// This will have the exact same capacity as a 64 bit number
func (z *Nat) SetUint64(x uint64) *Nat {
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

// Modulus represents a natural number used for modular reduction
//
// Unlike with natural numbers, the number of bits need to contain the modulus
// is assumed to be public. Operations are allowed to leak this size, and creating
// a modulus will remove unnecessary zeros.
type Modulus struct {
	nat Nat
	// the number of leading zero bits
	leading uint
	// The inverse of the least significant limb, modulo W
	m0inv Word
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
	if m.nat.limbs[0]&1 == 0 {
		panic("Modulus is even")
	}
	m.leading = uint(bits.LeadingZeros(uint(m.nat.limbs[len(m.nat.limbs)-1])))
	m.m0inv = invertModW(m.nat.limbs[0])
	m.m0inv = -m.m0inv
}

// SetUint64 sets the modulus according to an integer
func ModulusFromUint64(x uint64) Modulus {
	var m Modulus
	m.nat.SetUint64(x)
	// edge case for 32 bit limb size
	if _W < 64 && len(m.nat.limbs) > 1 && m.nat.limbs[1] == 0 {
		m.nat.limbs = m.nat.limbs[:1]
	}
	m.precomputeValues()
	return m
}

// trueSize calculates the actual size necessary for representing these limbs
//
// This is the size with leading zeros removed. This leaks the number
// of such zeros
func trueSize(limbs []Word) int {
	var size int
	for size = len(limbs); size > 0 && limbs[size-1] == 0; size-- {
	}
	return size
}

// FromBytes creates a new Modulus, converting from big endian bytes
//
// This function will remove leading zeros, thus leaking the true size of the modulus.
// See the documentation for the Modulus type, for more information about this contract.
func ModulusFromBytes(bytes []byte) Modulus {
	var m Modulus
	// TODO: You could allocate a smaller buffer to begin with, versus using the Nat method
	m.nat.SetBytes(bytes)

	m.nat.limbs = m.nat.limbs[:trueSize(m.nat.limbs)]
	m.precomputeValues()
	return m
}

// FromNat creates a new Modulus, using the value of a Nat
//
// This will leak the true size of this natural number. Because of this,
// the true size of the number should not be sensitive information. This is
// a stronger requirement than we usually have for Nat.
func ModulusFromNat(nat Nat) Modulus {
	var m Modulus
	// We make a copy here, to avoid any aliasing between buffers
	size := trueSize(nat.limbs)
	m.nat.limbs = m.nat.resizedLimbs(size)
	copy(m.nat.limbs, nat.limbs)
	m.precomputeValues()
	return m
}

// Bytes returns the big endian bytes making up the modulus
func (m *Modulus) Bytes() []byte {
	return m.nat.Bytes()
}

// shiftAddIn calculates z = z << _W + x mod m
//
// The length of z and scratch should be len(m) + 1
func shiftAddIn(z, scratch []Word, x Word, m *Modulus) {
	// Making tests on the exact bit length of m is ok,
	// since that's part of the contract for moduli
	size := len(m.nat.limbs)
	if size == 0 {
		return
	}
	if size == 1 {
		// In this case, z:x % m is exactly what we need to calculate
		_, r := div(z[0], x, m.nat.limbs[0])
		z[0] = r
		return
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
	q := ctIfElse(ctEq(a1, b0), ^Word(0), ctIfElse(ctEq(rawQ, 0), 0, rawQ-1))

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
	subVV(scratch, z, m.nat.limbs)
	ctCondCopy(over, z, scratch)
}

// Mod calculates z <- x mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) Mod(x *Nat, m *Modulus) *Nat {
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
	size := len(m.nat.limbs)
	var yModM Nat
	yModM.Mod(y, m)
	xLimbs := x.unaliasedLimbs(z)
	// The idea is to position x after size zeros, making modular reduction
	// also calculate the montgomery representation
	z.limbs = z.resizedLimbs(size + len(xLimbs))
	copy(z.limbs[size:], xLimbs)
	for i := 0; i < size; i++ {
		z.limbs[i] = 0
	}
	z.Mod(z, m)
	z.limbs = z.resizedLimbs(2 * size)

	zLimbs := z.limbs[:size]
	scratch := z.limbs[size:]
	montgomeryMul(zLimbs, yModM.limbs, zLimbs, scratch, m)

	z.limbs = zLimbs
	return z
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
	return z
}

// Exp calculates z <- x^y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) Exp(x *Nat, y *Nat, m *Modulus) *Nat {
	size := len(m.nat.limbs)

	// We create a new nat for this operation, so we only worry about aliasing y
	var xsquared Nat
	xsquared.Mod(x, m)

	yLimbs := y.unaliasedLimbs(z)
	scratch := z.resizedLimbs(3 * size)
	z.limbs = scratch[:size]
	z.limbs[0] = 1
	scratchA := scratch[size : 2*size]
	scratchB := scratch[2*size:]

	montgomeryRepresentation(xsquared.limbs, scratchA, m)

	// LEAK: y's length
	// OK: this should be public
	for i := 0; i < len(yLimbs); i++ {
		yi := yLimbs[i]
		for j := 0; j < _W; j++ {
			montgomeryMul(z.limbs, xsquared.limbs, scratchA, scratchB, m)
			selectMultiply := yi & 1
			ctCondCopy(selectMultiply, z.limbs, scratchA)
			montgomeryMul(xsquared.limbs, xsquared.limbs, xsquared.limbs, scratchB, m)
			yi >>= 1
		}
	}
	return z
}

// cmpGeq compares two limbs (same size) returning 1 if x >= y, and 0 otherwise
func cmpGeq(x []Word, y []Word) Word {
	res := Word(1)
	// LEAK: length of x, y
	// OK: this should be public
	for i := 0; i < len(x) && i < len(y); i++ {
		res = ctIfElse(ctEq(x[i], y[i]), res, ctGt(x[i], y[i]))
	}
	return res
}

// cmpEq compares two slices of the same length for equality, in constant time
//
// LEAK: the length of a, and b
func cmpEq(a, b []Word) Word {
	var v Word
	for i := 0; i < len(a) && i < len(b); i++ {
		v |= a[i] ^ b[i]
	}
	return ctEq(v, 0)
}

// cmpZero checks if a slice is equal to zero, in constant time
//
// LEAK: the length of a
func cmpZero(a []Word) Word {
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

	eq := Word(1)
	geq := Word(0)
	for i := 0; i < len(zLimbs) && i < len(xLimbs); i++ {
		eq_at_i := ctEq(zLimbs[i], xLimbs[i])
		eq &= eq_at_i
		geq = ctIfElse(eq_at_i, geq, ctGt(zLimbs[i], xLimbs[i]))
	}
	return subtle.ConstantTimeSelect(int(eq), 0, subtle.ConstantTimeSelect(int(geq), 1, -1))
}

// EqZero compares z to 0, returning true if they're equal
func (z *Nat) EqZero() bool {
	return cmpZero(z.limbs) == 1
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
	copy(a, xLimbs)
	// v = 0
	// u = 1
	for i := 0; i < size; i++ {
		u[i] = 0
		v[i] = 0
	}
	u[0] = 1

	// halfm = (m + 1) / 2
	halfm[size] = addVW(halfm, mLimbs, 1)
	shrVU(halfm, halfm, 1)
	halfm = halfm[:size]

	copy(b, m.limbs)

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
		aOdd := shrVU(a1, a, 1) >> (_W - 1)
		aEven := 1 ^ aOdd
		uOdd := shrVU(u2, u, 1) >> (_W - 1)
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
		subCarry := subVV(u, u, v)
		addVV(u1, u, m.limbs)
		ctCondCopy(subCarry, u, u1)
		uOdd = shrVU(u, u, 1) >> (_W - 1)
		addVV(u1, u, halfm)
		ctCondCopy(uOdd, u, u1)

		// If a was indeed even, we use the results we produced earlier
		ctCondCopy(aEven, a, a1)
		ctCondCopy(aEven, u, u2)
	}

	z.limbs = v
	return z
}

// ModInverse calculates z <- x^-1 mod m
//
// This will produce nonsense if the modulus is even.
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModInverse(x *Nat, m *Modulus) *Nat {
	z.Mod(x, m)
	return z.modInverse(z, &m.nat)
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
				out[i] = ((out[i] << 1) | sel)
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
func (z *Nat) ModInverseEven(x *Nat, m *Nat) *Nat {
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
	size := len(m.limbs)
	// We want to invert m modulo x, so we first calculate the reduced version, before inverting
	var newZ Nat
	newZ.limbs = divDouble(m.limbs, x.limbs, []Word{})
	newZ.modInverse(&newZ, x)
	inverseZero := cmpZero(newZ.limbs)
	newZ.Mul(&newZ, m, uint(2*size*_W))
	newZ.limbs = newZ.resizedLimbs(2 * size)
	subVW(newZ.limbs, newZ.limbs, 1)
	divDouble(newZ.limbs, x.limbs, newZ.limbs)
	// The result fits on a single half of newZ, but we need to subtract it from m.
	// We can use the other half of newZ, and then copy it back over if we need to keep it
	subVV(newZ.limbs[size:], m.limbs, newZ.limbs[:size])
	// If the inverse was zero, then x was 1, and so we should return 1.
	// We go ahead and prepare this result, but expect to copy over the subtraction
	// we just calculated soon over, in the usual case.
	newZ.limbs[0] = 1
	for i := 1; i < size; i++ {
		newZ.limbs[i] = 0
	}
	ctCondCopy(1^inverseZero, newZ.limbs[:size], newZ.limbs[size:])

	z.limbs = newZ.limbs[:size]
	return z
}
