package safenum

import "math/big"

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
	limbs []big.Word
}

func fromInt(i *big.Int) Nat {
	return Nat{limbs: i.Bits()}
}

func (z Nat) toInt() *big.Int {
	var ret big.Int
	ret.SetBits(z.limbs)
	return &ret
}

// Mod calculates z <- x mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) Mod(x Nat, m Nat) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().Mod(x.toInt(), m.toInt()))
	return z
}

// ModAdd calculates z <- x + y mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) ModAdd(x Nat, y Nat, m Nat) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().Add(x.toInt(), y.toInt()))
	*z = fromInt(z.toInt().Mod(z.toInt(), m.toInt()))
	return z
}

// Add calculates z <- x + y, modulo 2^cap
//
// The capacity is given in bits, and also controls the size of the result.
func (z *Nat) Add(x Nat, y Nat, cap uint) *Nat {
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
func (z *Nat) ModMul(x Nat, y Nat, m Nat) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().Mul(x.toInt(), y.toInt()))
	*z = fromInt(z.toInt().Mod(z.toInt(), m.toInt()))
	return z
}

// Mul calculates z <- x * y, modulo 2^cap
//
// The capacity is given in bits, and also controls the size of the result.
func (z *Nat) Mul(x Nat, y Nat, cap uint) *Nat {
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
func (z *Nat) ModInverse(x Nat, m Nat) *Nat {
	// TODO: Use actual implementation
	*z = fromInt(z.toInt().ModInverse(x.toInt(), m.toInt()))
	return z
}

// Exp calculates z <- x^y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) Exp(x Nat, y Nat, m Nat) *Nat {
	*z = fromInt(z.toInt().Exp(x.toInt(), y.toInt(), m.toInt()))
	return z
}

// TODO: What should the semantics here be for equivalent values with differing capacities?
// Is it possible to do normalized comparison in consant time?
// NOTE: It would be unsound to do normalized comparison, since this would leak the number
// of leading zeroes

// Cmp compares two natural numbers, returning 1 if they're equal and 0 otherwise
func (z Nat) Cmp(x Nat) int {
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

// SetBytes interprets a number in big-endian format, stores it in z, and returns z.
//
// The exact length of the buffer must be public information! This length also dictates
// the capacity of the number returned, and thus the resulting timings for operations
// involving that number.
func (z *Nat) SetBytes(buf []byte) *Nat {
	*z = fromInt(z.toInt().SetBytes(buf))
	return z
}

// SetUint64 sets z to x, and returns z
//
// This will have the exact same capacity as a 64 bit number
func (z *Nat) SetUint64(x uint64) *Nat {
	// TODO: Use an actual implementation
	*z = fromInt(z.toInt().SetUint64(x))
	return z
}
