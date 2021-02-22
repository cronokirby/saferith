package safenum

import (
	"math/big"
)

type word uint32

// Nat represents an arbitrary sized natural number.
//
// Different methods on Nats will talk about a "capacity". The capacity represents
// the announced size of some number. Operations may vary in time *only* relative
// to this capacity, and not to the actual value of the number.
//
// The capacity of a number is usually inherited through whatever method was used to
// create the number in the first place.
type Nat struct {
	// TODO: Don't rely math/big
	i big.Int
}

// Mod calculates z <- x mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) Mod(x Nat, m Nat) *Nat {
	panic("unimplemented")
}

// ModAdd calculates z <- x + y mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) ModAdd(x Nat, y Nat, m Nat) *Nat {
	z.i = *z.i.Add(&x.i, &y.i)
	z.i = *z.i.Mod(&z.i, &m.i)
	return z
}

// Add calculates z <- x + y, modulo 2^cap
//
// The capacity is given in bits, and also controls the size of the result.
func (z *Nat) Add(x Nat, y Nat, cap uint) *Nat {
	// TODO: Use an actual implementation
	z.i = *z.i.Add(&x.i, &y.i)
	bytes := z.i.Bytes()
	numBytes := (cap + 8 - 1) >> 3
	z.i.SetBytes(bytes[len(bytes)-int(numBytes):])
	return z
}

// ModMul calculates z <- x * y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModMul(x Nat, y Nat, m Nat) *Nat {
	panic("unimplemented")
}

// Mul calculates z <- x * y, modulo 2^cap
//
// The capacity is given in bits, and also controls the size of the result.
func (z *Nat) Mul(x Nat, y Nat, cap uint) *Nat {
	// TODO: Use an actual implementation
	z.i = *z.i.Mul(&x.i, &y.i)
	bytes := z.i.Bytes()
	numBytes := (cap + 8 - 1) >> 3
	z.i.SetBytes(bytes[len(bytes)-int(numBytes):])
	return z
}

// ModInverse calculates z <- x^-1 mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModInverse(x Nat, m Nat) *Nat {
	panic("unimplemented")
}

// Exp calculates z <- x^y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) Exp(x Nat, y Nat, m Nat) *Nat {
	panic("unimplemented")
}

// TODO: What should the semantics here be for equivalent values with differing capacities?
// Is it possible to do normalized comparison in consant time?

// Cmp compares two natural numbers, returning 1 if they're equal and 0 otherwise
func (z Nat) Cmp(x Nat) int {
	// TODO: Use actual implementation
	ret := z.i.Cmp(&x.i)
	if ret < 0 {
		return -ret
	}
	return ret
}

// FillBytes writes out the big endian bytes of a natural number.
//
// This will always write out the full capacity of the number, without
// any kind trimming.
//
// This will panic if the buffer's length cannot accomodate the capacity of the number
func (z *Nat) FillBytes(buf []byte) []byte {
	panic("unimplemented")
}

// SetBytes interprets a number in big-endian format, stores it in z, and returns z.
//
// The exact length of the buffer must be public information! This length also dictates
// the capacity of the number returned, and thus the resulting timings for operations
// involving that number.
func (z *Nat) SetBytes(buf []byte) *Nat {
	panic("unimplemented")
}

// SetUint64 sets z to x, and returns z
//
// This will have the exact same capacity as a 64 bit number
func (z *Nat) SetUint64(x uint64) *Nat {
	z.i.SetUint64(x)
	return z
}
