package safenum

type word uint32

// Nat represents an arbitrary sized natural number.
type Nat struct {
	limbs []word
}

// Modulus represents a natural number, acting as some kind of modulus.
//
// Moduli are quite important, since they're implicitly used throughout this package
// to bound the size of the calculation, allowing things to be constant time.
type Modulus Nat

// Mod calculates z <- x mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) Mod(x Nat, m Modulus) *Nat {
	panic("unimplemented")
}

// ModAdd calculates z <- x + y mod m
//
// The capacity of the resulting number matches the capacity of the modulus.
func (z *Nat) ModAdd(x Nat, y Nat, m Modulus) *Nat {
	panic("unimplemented")
}

// ModMul calculates z <- x * y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModMul(x Nat, y Nat, m Modulus) *Nat {
	panic("unimplemented")
}

// ModInverse calculates z <- x^-1 mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) ModInverse(x Nat, m Modulus) *Nat {
	panic("unimplemented")
}

// Exp calculates z <- x^y mod m
//
// The capacity of the resulting number matches the capacity of the modulus
func (z *Nat) Exp(x Nat, y Nat, m Modulus) *Nat {
	panic("unimplemented")
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
