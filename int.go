package safenum

// Int represents a signed integer of arbitrary size.
//
// Similarly to Nat, each Int comes along with an announced size, representing
// the number of bits need to represent its absolute value. This can be
// larger than its true size, the number of bits actually needed.
type Int struct {
	// This number is represented by (-1)^sign * abs, essentially

	// When 1, this is a negative number, when 0 a positive number.
	//
	// There's a bit of redundancy to note, because -0 and +0 represent the same
	// number. We need to be careful around this edge case.
	sign Choice
	// The absolute value.
	//
	// Not using a point is important, that way the zero value for Int is actually zero.
	abs Nat
}

// SetBytes interprets a number in big-endian form, stores it in z, and returns z.
//
// This number will be positive.
func (z *Int) SetBytes(data []byte) *Int {
	z.sign = 0
	z.abs.SetBytes(data)
	return z
}

// SetUint64 sets the value of z to x.
//
// This number will be positive.
func (z *Int) SetUint64(x uint64) *Int {
	z.sign = 0
	z.abs.SetUint64(x)
	return z
}

// String formats this number as a signed hex string.
//
// This isn't a format that Int knows how to parse. This function exists mainly
// to help debugging, and whatnot.
func (z *Int) String() string {
	sign := ctIfElse(z.sign, Word('-'), Word('+'))
	return string(rune(sign)) + z.abs.String()
}

// Eq checks if this Int has the same value as another Int.
//
// Note that negative zero and positive zero are the same number.
func (z *Int) Eq(x *Int) Choice {
	zero := z.abs.EqZero()
	// If this is zero, then any number as the same sign,
	// otherwise, check that the signs aren't different
	sameSign := zero | (1 ^ z.sign ^ x.sign)
	return sameSign & z.abs.Eq(&x.abs)
}

// Neg calculates z <- -x.
//
// The result has the same announced size.
func (z *Int) Neg(x *Int) *Int {
	z.sign = 1 ^ x.sign
	z.abs.SetNat(&x.abs)
	return z
}

// Mul calculates z <- x * y, returning z.
//
// This will truncate the resulting absolute value, based on the bit capacity passed in.
//
// If cap < 0, then capacity is x.AnnouncedLen() + y.AnnouncedLen().
func (z *Int) Mul(x *Int, y *Int, cap int) *Int {
	// (-1)^sx * ax * (-1)^sy * ay = (-1)^(sx + sy) * ax * ay
	z.sign = x.sign ^ y.sign
	z.abs.Mul(&x.abs, &y.abs, cap)
	return z
}
