package safenum

import (
	"crypto/subtle"
	"math/big"
	"math/bits"
)

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
	limbs []Word
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
	return z.limbs[:size]
}

func fromInt(i *big.Int) Nat {
	var n Nat
	n.SetBytes(i.Bytes())
	return n
}

func (z Nat) toInt() *big.Int {
	var ret big.Int
	ret.SetBytes(z.Bytes())
	return &ret
}

func div(hi, lo, y Word) (Word, Word) {
	// TODO: Use constant time division
	quo, rem := bits.Div(uint(hi), uint(lo), uint(y))
	return Word(quo), Word(rem)
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
		_, r := div(z[0], x, m.nat.limbs[0])
		z[0] = r
		return
	}

	scratch[0] = x
	var a1, a0, b0 Word
	if m.leading == 0 {
		a1 = z[size-1]
		copy(scratch[1:], z)
		a0 = scratch[size-1]
		b0 = m.nat.limbs[size-1]
	} else {
		a1 = (z[size-1] << m.leading) | (z[size-2] >> (_W - m.leading))
		copy(scratch[1:], z)
		a0 = (scratch[size-1] << m.leading) | (scratch[size-2] >> (_W - m.leading))
		b0 = (m.nat.limbs[size-1] << m.leading) | (m.nat.limbs[size-2] >> (_W - m.leading))
	}

	var q Word
	// TODO: After division is constant time, MUX here instead
	if a0 >= b0 {
		q = ^Word(0)
	} else {
		q, _ = div(a1, a0, b0)
		if q != 0 {
			q -= 1
		}
	}
	c1 := mulSubVVW(scratch, m.nat.limbs, q)
	if c1 > a1 {
		addVV(scratch, scratch, m.nat.limbs)
		copy(z, scratch)
		return
	}
	c := subVV(z, scratch, m.nat.limbs)
	sel := 1
	if c == a1-c1 {
		sel = 0
	}
	constantTimeWordCopy(sel, z, scratch)
}

// Mod calculates z <- x mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) Mod(x *Nat, m *Modulus) *Nat {
	/*
		size := len(m.nat.limbs)
		xLimbs := x.limbs
		z.limbs = make([]Word, 2*size)
		for i := len(xLimbs) - 1; i >= 0; i-- {
			shiftAddIn(z.limbs[:size], z.limbs[size:], xLimbs[i], m)
		}
		z.limbs = z.limbs[:size]
		return z
	*/
	limbCount := len(m.nat.limbs)
	// We need two buffers, because of aliasing
	subScratch := make([]Word, limbCount)
	rLimbs := make([]Word, limbCount)
	// LEAK: the length of x
	// OK: this should be public
	for i := len(x.limbs) - 1; i >= 0; i-- {
		limb := x.limbs[i]
		for j := _W - 1; j >= 0; j-- {
			xi := (limb >> j) & 1
			shiftCarry := shlVU(rLimbs, rLimbs, 1)
			rLimbs[0] |= xi
			subCarry := subVV(subScratch, rLimbs, m.nat.limbs)
			selectSub := constantTimeWordEq(shiftCarry, subCarry)
			constantTimeWordCopy(selectSub, rLimbs, subScratch)
		}
	}
	// Now we can safely swap things out
	z.limbs = rLimbs
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
	limbCount := len(m.nat.limbs)
	z.limbs = z.resizedLimbs(limbCount)

	// LEAK: limbCount
	// OK: the size of the modulus should be public information
	addCarry := addVV(z.limbs, xModM.limbs, yModM.limbs)
	// I don't think we can avoid using an extra scratch buffer
	subResult := make([]Word, limbCount)
	// LEAK: limbCount
	// OK: see above
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
	selectSub := constantTimeWordEq(addCarry, subCarry)
	constantTimeWordCopy(selectSub, z.limbs, subResult)
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
	mask := ^(^Word(0) << bitsToKeep)
	// LEAK: the size of z (since we're making an extra access at the end)
	// OK: this is public information, since cap is public
	z.limbs[len(z.limbs)-1] &= mask
	return z
}

// ModMul calculates z <- x * y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModMul(x *Nat, y *Nat, m *Modulus) *Nat {
	limbCount := len(x.limbs) + len(y.limbs)
	cap := _W * limbCount
	z.Mul(x, y, uint(cap))
	z.Mod(z, m)
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

// ModInverse calculates z <- x^-1 mod m
//
// This will produce nonsense if the modulus is even.
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModInverse(x *Nat, m *Modulus) *Nat {
	limbCount := len(m.nat.limbs)

	// aHalf <- a / 2
	// aMinusBHalf <- (a - b) / 2
	var a, aHalf, aMinusBHalf Nat
	a.Mod(x, m)
	aHalf.limbs = make([]Word, limbCount)
	aMinusBHalf.limbs = make([]Word, limbCount)

	// bHalf <- b / 2
	// bMinusAHalf <- (b - a) / 2
	var b, bHalf, bMinusAHalf Nat
	b.limbs = make([]Word, limbCount)
	copy(b.limbs, m.nat.limbs)
	bHalf.limbs = make([]Word, limbCount)
	bMinusAHalf.limbs = make([]Word, limbCount)

	// uHalf <- u / 2
	// uHalfAdjust <- u / 2 + adjust (if u wasn't even)
	// uMinusVHalf <- (u - v) / 2
	// uMinusVHalfUnder <- (u - v) + m (when the subtraction overflows)
	// uMinusVHalfUnder <- (u - v) / 2 + adjust (if this wasn't even)
	var u, uHalf, uHalfAdjust, uMinusVHalf, uMinusVHalfUnder, uMinusVHalfAdjust Nat
	u.limbs = make([]Word, limbCount)
	u.limbs[0] = 1
	uHalf.limbs = make([]Word, limbCount)
	uHalfAdjust.limbs = make([]Word, limbCount)
	uMinusVHalf.limbs = make([]Word, limbCount)
	uMinusVHalfUnder.limbs = make([]Word, limbCount)
	uMinusVHalfAdjust.limbs = make([]Word, limbCount)

	// vHalf <- v / 2
	// vHalfAdjust <- v / 2 + adjust (if v wasn't even)
	// vMinusUHalf <- (v - u) / 2
	// vMinusUHalfUnder <- (v - u) + m (when the subtraction overflows)
	// vMinusUHalfUnder <- (v - u) / 2 + adjust (if this wasn't even)
	var v, vHalf, vHalfAdjust, vMinusUHalf, vMinusUHalfUnder, vMinusUHalfAdjust Nat
	v.limbs = make([]Word, limbCount)
	vHalf.limbs = make([]Word, limbCount)
	vHalfAdjust.limbs = make([]Word, limbCount)
	vMinusUHalf.limbs = make([]Word, limbCount)
	vMinusUHalfUnder.limbs = make([]Word, limbCount)
	vMinusUHalfAdjust.limbs = make([]Word, limbCount)

	// In order to implement a / 2 mod m, if a might not be even,
	// we shift right by 2, and the conditionally add in (m + 1) / 2.
	// Adjust contains (m + 1) / 2
	var adjust Nat
	// We just want to add 1 to m, and then shift down, so we need to have an extra
	// bit of capacity in case adding 1 to m needs an extra limb. I guess this is necessary
	// e.g. you're using a mersenne prime as a modulus?
	adjust.Add(&u, &m.nat, _W*uint(limbCount)+1)
	shrVU(adjust.limbs, adjust.limbs, 1)
	adjust.limbs = adjust.limbs[:limbCount]

	for i := 1; i < _W*limbCount; i++ {
		aOdd := shrVU(aHalf.limbs, a.limbs, 1) >> (_W - 1)
		bLarger := subVV(aMinusBHalf.limbs, a.limbs, b.limbs)
		shrVU(aMinusBHalf.limbs, aMinusBHalf.limbs, 1)

		bOdd := shrVU(bHalf.limbs, b.limbs, 1) >> (_W - 1)
		aLarger := subVV(bMinusAHalf.limbs, b.limbs, a.limbs)
		shrVU(bMinusAHalf.limbs, bMinusAHalf.limbs, 1)

		uOdd := shrVU(uHalf.limbs, u.limbs, 1) >> (_W - 1)
		addVV(uHalfAdjust.limbs, uHalf.limbs, adjust.limbs)
		constantTimeWordCopy(int(uOdd), uHalf.limbs, uHalfAdjust.limbs)
		uUnder := subVV(uMinusVHalf.limbs, u.limbs, v.limbs)
		addVV(uMinusVHalfUnder.limbs, uMinusVHalf.limbs, m.nat.limbs)
		constantTimeWordCopy(int(uUnder), uMinusVHalf.limbs, uMinusVHalfUnder.limbs)
		uAdjust := shrVU(uMinusVHalf.limbs, uMinusVHalf.limbs, 1) >> (_W - 1)
		addVV(uMinusVHalfAdjust.limbs, uMinusVHalf.limbs, adjust.limbs)
		constantTimeWordCopy(int(uAdjust), uMinusVHalf.limbs, uMinusVHalfAdjust.limbs)

		vOdd := shrVU(vHalf.limbs, v.limbs, 1) >> (_W - 1)
		addVV(vHalfAdjust.limbs, vHalf.limbs, adjust.limbs)
		constantTimeWordCopy(int(vOdd), vHalf.limbs, vHalfAdjust.limbs)
		vUnder := subVV(vMinusUHalf.limbs, v.limbs, u.limbs)
		addVV(vMinusUHalfUnder.limbs, vMinusUHalf.limbs, m.nat.limbs)
		constantTimeWordCopy(int(vUnder), vMinusUHalf.limbs, vMinusUHalfUnder.limbs)
		vAdjust := shrVU(vMinusUHalf.limbs, vMinusUHalf.limbs, 1) >> (_W - 1)
		addVV(vMinusUHalfAdjust.limbs, vMinusUHalf.limbs, adjust.limbs)
		constantTimeWordCopy(int(vAdjust), vMinusUHalf.limbs, vMinusUHalfAdjust.limbs)

		// Here's the big idea:
		//
		// if a == b:
		//	 pass
		// else if even(a):
		//	 a <- a / 2
		//   u <- u / 2 mod m
		// else if even(b):
		//   b <- b / 2
		//   v <- v / 2 mod m
		// else if a > b:
		//   a <- (a - b) / 2
		//   u <- (u - v) / 2 mod m
		// else if b > a:
		//   b <- (b - a) / 2
		//   v <- (v - u) / 2 mod m

		// TODO: Is this the best way of making the selection matrix?
		// Exactly one of these is going to be true, in theory
		select1 := 1 - int(aOdd)
		select2 := (1 - select1) & (1 - int(bOdd))
		select3 := (1 - select1) & (1 - select2) & int(aLarger)
		select4 := (1 - select1) & (1 - select2) & (1 - select3) & int(bLarger)

		constantTimeWordCopy(select1, a.limbs, aHalf.limbs)
		constantTimeWordCopy(select1, u.limbs, uHalf.limbs)
		constantTimeWordCopy(select2, b.limbs, bHalf.limbs)
		constantTimeWordCopy(select2, v.limbs, vHalf.limbs)
		constantTimeWordCopy(select3, a.limbs, aMinusBHalf.limbs)
		constantTimeWordCopy(select3, u.limbs, uMinusVHalf.limbs)
		constantTimeWordCopy(select4, b.limbs, bMinusAHalf.limbs)
		constantTimeWordCopy(select4, v.limbs, vMinusUHalf.limbs)
	}
	z.limbs = u.limbs
	return z
}

// Exp calculates z <- x^y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) Exp(x *Nat, y *Nat, m *Modulus) *Nat {
	limbCount := len(m.nat.limbs)
	var mulScratch, xsquared, zScratch Nat
	xsquared.limbs = make([]Word, limbCount)
	zScratch.limbs = make([]Word, limbCount)
	zScratch.limbs[0] = 1
	// LEAK: limbCount, x's length
	// OK: both should be public information
	copy(xsquared.limbs, x.limbs)
	// LEAK: y's length
	// OK: this should be public
	for i := 0; i < len(y.limbs); i++ {
		yi := y.limbs[i]
		for j := 0; j < _W; j++ {
			mulScratch.ModMul(&zScratch, &xsquared, m)
			selectMultiply := int(yi & 1)
			constantTimeWordCopy(selectMultiply, zScratch.limbs, mulScratch.limbs)
			xsquared.ModMul(&xsquared, &xsquared, m)
			yi >>= 1
		}
	}
	z.limbs = zScratch.limbs
	return z
}

func constantTimeWordEq(x, y Word) int {
	zero := uint64(x ^ y)
	// TODO: Find a better way to do this
	return subtle.ConstantTimeEq(int32(zero), 0) & subtle.ConstantTimeEq(int32(zero>>32), 0)
}

// constantTimeWordCopy copies y into x, if v == 1, otherwise does nothing
//
// Both slices must have the same length.
//
// LEAK: the length of the slices
//
// Otherwise, which branch was taken isn't leaked
func constantTimeWordCopy(v int, x, y []Word) {
	xmask := Word(v - 1)
	ymask := Word(^(v - 1))
	for i := 0; i < len(x); i++ {
		x[i] = (x[i] & xmask) | (y[i] & ymask)
	}
}

// CmpEq compares two natural numbers, returning 1 if they're equal and 0 otherwise
func (z *Nat) CmpEq(x *Nat) int {
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

	var v Word
	// LEAK: size
	// OK: this was calculated using the length of x and z, both public
	for i := 0; i < size; i++ {
		v |= zLimbs[i] ^ xLimbs[i]
	}
	return constantTimeWordEq(v, 0)
}

// FillBytes writes out the big endian bytes of a natural number.
//
// This will always write out the full capacity of the number, without
// any kind trimming.
//
// This will panic if the buffer's length cannot accomodate the capacity of the number
func (z *Nat) FillBytes(buf []byte) []byte {
	z.toInt().FillBytes(buf)
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
	i := length
	// LEAK: Number of limbs
	// OK: The number of limbs is public
	// LEAK: The addresses touched in the out array
	// OK: Every member of out is touched
	for _, x := range z.limbs {
		y := x
		for j := 0; j < _S; j++ {
			i--
			out[i] = byte(y)
			y >>= 8
		}
	}
	return out
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
		for i := 0; i < limbCount; i++ {
			z.limbs[i] = Word(x)
			x >>= _W
		}
	}
	return z
}

// setLeading calculates the number of leading zeros of the top limb of m
//
// This leaks the value, most likely.
func (m *Modulus) setLeading() {
	m.leading = uint(bits.LeadingZeros(uint(m.nat.limbs[len(m.nat.limbs)-1])))
}

// SetUint64 sets the modulus according to an integer
func (z *Modulus) SetUint64(x uint64) *Modulus {
	z.nat.SetUint64(x)
	// edge case for 32 bit limb size
	if _W < 64 && len(z.nat.limbs) > 1 && z.nat.limbs[1] == 0 {
		z.nat.limbs = z.nat.limbs[:1]
	}
	z.setLeading()
	return z
}

// trueSize calculates the actual size necessary for representing these limbs
//
// This is the size with leading zeros removed. This naturally leaks the number
// of such zeros
func trueSize(limbs []Word) int {
	var size int
	for size = len(limbs); size > 0 && limbs[size-1] == 0; size-- {
	}
	return size
}

// SetBytes sets the value of the modulus according to a slice of Big Endian bytes
//
// This will trim the modulus to only use the necessary
func (m *Modulus) SetBytes(bytes []byte) *Modulus {
	// TODO: You could allocate a smaller buffer to begin with, versus using the Nat method
	m.nat.SetBytes(bytes)

	m.nat.limbs = m.nat.limbs[:trueSize(m.nat.limbs)]
	m.setLeading()
	return m
}

// SetNat sets the value of the modulus according to a Nat
//
// This will leak the exact number of bits for the natural number, so this shouldn't be sensitive.
// Using the modulus will continue to leak this.
func (m *Modulus) SetNat(nat Nat) *Modulus {
	// We make a copy here, to avoid any aliasing between buffers
	size := trueSize(nat.limbs)
	m.nat.limbs = m.nat.resizedLimbs(size)
	copy(m.nat.limbs, nat.limbs)
	m.setLeading()
	return m
}

// Bytes returns the big endian bytes making up the modulus
func (m *Modulus) Bytes() []byte {
	return m.nat.Bytes()
}
