package safenum

// limbCountFromBits returns the number of limbs needed to accomodate bits.
func limbCountFromBits(bits int) int {
	return (bits + _W - 1) / _W
}
