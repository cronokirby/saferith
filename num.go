package safenum

import (
	"math/big"
)

// Nat represents an arbitrary sized natural number.
//
// Different methods on Nats will talk about a "capacity". The capacity represents
// the announced size of some number. Operations may vary in time *only* relative
// to this capacity, and not to the actual value of the number.
//
// The capacity of a number is usually inherited through whatever method was used to
// create the number in the first place.
type Nat struct {
	// TODO: Once we don't rely on math/big at all, use our own word type
	limbs []uint32
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

// Mod calculates z <- x mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) Mod(x *Nat, m *Nat) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().Mod(x.toInt(), m.toInt()))
	return z
}

// ModAdd calculates z <- x + y mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) ModAdd(x *Nat, y *Nat, m *Nat) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().Add(x.toInt(), y.toInt()))
	*z = fromInt(z.toInt().Mod(z.toInt(), m.toInt()))
	return z
}

// Add calculates z <- x + y, modulo 2^cap
//
// The capacity is given in bits, and also controls the size of the result.
func (z *Nat) Add(x *Nat, y *Nat, cap uint) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().Add(x.toInt(), y.toInt()))
	bytes := z.toInt().Bytes()
	numBytes := (cap + 8 - 1) >> 3
	if int(numBytes) < len(bytes) {
		*z = fromInt(z.toInt().SetBytes(bytes[len(bytes)-int(numBytes):]))
	}
	return z
}

// ModMul calculates z <- x * y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModMul(x *Nat, y *Nat, m *Nat) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().Mul(x.toInt(), y.toInt()))
	*z = fromInt(z.toInt().Mod(z.toInt(), m.toInt()))
	return z
}

// Mul calculates z <- x * y, modulo 2^cap
//
// The capacity is given in bits, and also controls the size of the result.
func (z *Nat) Mul(x *Nat, y *Nat, cap uint) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().Mul(x.toInt(), y.toInt()))
	bytes := z.toInt().Bytes()
	numBytes := (cap + 8 - 1) >> 3
	if int(numBytes) < len(bytes) {
		*z = fromInt(z.toInt().SetBytes(bytes[len(bytes)-int(numBytes):]))
	}
	return z
}

// ModInverse calculates z <- x^-1 mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModInverse(x *Nat, m *Nat) *Nat {
	// TODO: Use actual implementation
	*z = fromInt(z.toInt().ModInverse(x.toInt(), m.toInt()))
	return z
}

// Exp calculates z <- x^y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) Exp(x *Nat, y *Nat, m *Nat) *Nat {
	*z = fromInt(z.toInt().Exp(x.toInt(), y.toInt(), m.toInt()))
	return z
}

// CmpEq compares two natural numbers, returning 1 if they're equal and 0 otherwise
func (z *Nat) CmpEq(x *Nat) int {
	// TODO: Use actual implementation
	ret := z.toInt().Cmp(x.toInt())
	if ret == 0 {
		return 1
	}
	return 0
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

func extendFront(buf []byte, to int) []byte {
	// TODO: Scrutinize this
	if len(buf) >= to {
		return buf
	}

	shift := to - len(buf)
	if cap(buf) < to {
		newBuf := make([]byte, to)
		copy(newBuf[shift:], buf)
		return newBuf
	}

	newBuf := buf[:to]
	copy(newBuf[shift:], buf)
	for i := 0; i < shift; i++ {
		newBuf[i] = 0
	}
	return newBuf
}

func (z *Nat) ensureLimbCapacity(size int) {
	if cap(z.limbs) < size {
		newLimbs := make([]uint32, size)
		copy(newLimbs, z.limbs)
		z.limbs = newLimbs
	}
}

func (z *Nat) resizedLimbs(size int) []uint32 {
	z.ensureLimbCapacity(size)
	return z.limbs[:size]
}

// SetBytes interprets a number in big-endian format, stores it in z, and returns z.
//
// The exact length of the buffer must be public information! This length also dictates
// the capacity of the number returned, and thus the resulting timings for operations
// involving that number.
func (z *Nat) SetBytes(buf []byte) *Nat {
	// TODO: Scrutinize the implementation a bit more here
	// We pad the front so that we have a multiple of 4
	// Padding the front is adding extra zeros to the BE representation
	necessary := (len(buf) + 3) &^ 0b11
	buf = extendFront(buf, necessary)
	limbCount := necessary >> 2
	z.limbs = z.resizedLimbs(limbCount)
	j := necessary
	for i := 0; i < limbCount; i++ {
		z.limbs[i] = 0
		j--
		z.limbs[i] |= uint32(buf[j])
		j--
		z.limbs[i] |= uint32(buf[j]) << 8
		j--
		z.limbs[i] |= uint32(buf[j]) << 16
		j--
		z.limbs[i] |= uint32(buf[j]) << 24
	}
	return z
}

// Bytes creates a slice containing the contents of this Nat, in big endian
//
// This will always fill the output byte slice based on the announced length of this Nat.
func (z *Nat) Bytes() []byte {
	length := len(z.limbs) << 2
	out := make([]byte, length)
	i := length
	// LEAK: Number of limbs
	// OK: The number of limbs is public
	// LEAK: The addresses touched in the out array
	// OK: Every member of out is touched
	for _, x := range z.limbs {
		i--
		out[i] = byte(x)
		i--
		out[i] = byte(x >> 8)
		i--
		out[i] = byte(x >> 16)
		i--
		out[i] = byte(x >> 24)
	}
	return out
}

// SetUint64 sets z to x, and returns z
//
// This will have the exact same capacity as a 64 bit number
func (z *Nat) SetUint64(x uint64) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().SetUint64(x))
	return z
}
