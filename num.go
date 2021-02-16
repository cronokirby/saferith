package safenum

type word uint32

// Nat represents an arbitrary sized natural number.
type Nat struct {
	limbs []word
}
