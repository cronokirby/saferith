package safenum

// limbCount returns the number of limbs needed to accomodate bits.
func limbCount(bits int) int {
	return (bits + _W - 1) / _W
}

// limbMask returns the mask used for the final limb of a Nat with this number of bits.
//
// Note that this function will leak the number of bits. For our library, this isn't
// a problem, since we always call this function with announced sizes.
func limbMask(bits int) Word {
	remaining := bits % _W
	allOnes := ^Word(0)
	if remaining == 0 {
		return allOnes
	}
	return 1 ^ (allOnes << remaining)
}
