package safenum

// Int represents a signed integer of arbitrary size.
//
// Similarly to Nat, each Int comes along with an announced size, representing
// the number of bits need to represent its absolute value. This can be
// larger than its true size, the number of bits actually needed.
type Int struct {
	// This number is represented by (-1)^negative * abs, essentially

	// When 1, this is a negative number, when 0 a positive number.
	//
	// There's a bit of redundancy to note, because -0 and +0 represent the same
	// number. We need to be careful around this edge case.
	negative Choice
	// The absolute value.
	abs *Nat
}
