package safenum

import (
	"math/big"
	"math/bits"
)

// NOTE: We define a type alias for our limbs, to make integration with
// big's internal routines easier later.

// Word represents the type of limbs of a natural number
type Word = uint

const (
	// Word size in bytes
	_S = bits.UintSize / 8
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
	limbs []Word
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
	numBytes := (cap + 8 - 1) / 8
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
	numBytes := (cap + 8 - 1) / 8
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

// ensureLimbCapacity makes sure that a Nat has capacity for a certain number of limbs
//
// This will modify the slice contained inside the natural, but won't change the size of
// the slice, so it doesn't affect the value of the natural.
//
// LEAK: Probably the current number of limbs, and size
// OK: both of these should be public
func (z *Nat) ensureLimbCapacity(size int) {
	if cap(z.limbs) < size {
		newLimbs := make([]Word, size)
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
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().SetUint64(x))
	return z
}
