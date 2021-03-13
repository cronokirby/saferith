// Parts of this file come from Go's standard library
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found at https://github.com/golang/go/blob/master/LICENSE

// NOTE: When integrating this library into the standard library, you can get
// rid of the functions that already have a counterpart, since those
// are already safe, at least with our assumptions about the shape of slices

package safenum

import "math/bits"

// Many of the loops in this file are of the form
//   for i := 0; i < len(z) && i < len(x) && i < len(y); i++
// i < len(z) is the real condition.
// However, checking i < len(x) && i < len(y) as well is faster than
// having the compiler do a bounds check in the body of the loop;
// remarkably it is even faster than hoisting the bounds check
// out of the loop, by doing something like
//   _, _ = x[len(z)-1], y[len(z)-1]
// There are other ways to hoist the bounds check out of the loop,
// but the compiler's BCE isn't powerful enough for them (yet?).
// See the discussion in CL 164966.

// Add two slices of Word, returning the carry you end up with
//
// LEAK: The lengths of x, y, and z
func addVV(z, x, y []Word) (c Word) {
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x) && i < len(y); i++ {
		zi, cc := bits.Add(uint(x[i]), uint(y[i]), uint(c))
		z[i] = Word(zi)
		c = Word(cc)
	}
	return
}

// Subtract one slice of Word from another, returning the carry you end up with
//
// The carry is 1 if the result underflows, so to speak
//
// LEAK: The lengths of x, y, and z
func subVV(z, x, y []Word) (c Word) {
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x) && i < len(y); i++ {
		zi, cc := bits.Sub(uint(x[i]), uint(y[i]), uint(c))
		z[i] = Word(zi)
		c = Word(cc)
	}
	return
}

// Shift x by s, outputting the result in z
//
// The carry consists of all the bits that were shifted out
//
// The length of z and x must match
//
// LEAK: the length of z and x, whether or not s is 0
func shlVU(z, x []Word, s uint) (c Word) {
	if s == 0 {
		copy(z, x)
		return
	}
	if len(z) == 0 {
		return
	}
	s &= _W - 1 // hint to the compiler that shifts by s don't need guard code
	ŝ := _W - s
	ŝ &= _W - 1 // ditto
	c = x[len(z)-1] >> ŝ
	for i := len(z) - 1; i > 0; i-- {
		z[i] = x[i]<<s | x[i-1]>>ŝ
	}
	z[0] = x[0] << s
	return
}

// Perform a logical right shift by s bits, from x into z
//
// The carry consists of the bits that were shifted out, but adjusted
// by _W - s bits to the left. Shifting right by _W - s can get the bits
// as they were in their original position.
//
// The length of z and x must match
//
// LEAK: the length of z and x, whether or not s is 0
func shrVU(z, x []Word, s uint) (c Word) {
	if s == 0 {
		copy(z, x)
		return
	}
	if len(z) == 0 {
		return
	}
	s &= _W - 1 // hint to the compiler that shifts by s don't need guard code
	ŝ := _W - s
	ŝ &= _W - 1 // ditto
	c = x[0] << ŝ
	for i := 0; i < len(z)-1; i++ {
		z[i] = x[i]>>s | x[i+1]<<ŝ
	}
	z[len(z)-1] = x[len(z)-1] >> s
	return
}

func mulAddWWW(x, y, c Word) (z1, z0 Word) {
	hi, lo := bits.Mul(uint(x), uint(y))
	var cc uint
	lo, cc = bits.Add(lo, uint(c), 0)
	return Word(hi + cc), Word(lo)
}

// z += x * y
//
// LEAK: The length of z and x
func addMulVVW(z, x []Word, y Word) (c Word) {
	// The comment near the top of this file discusses this for loop condition.
	for i := 0; i < len(z) && i < len(x); i++ {
		z1, z0 := mulAddWWW(x[i], y, z[i])
		lo, cc := bits.Add(uint(z0), uint(c), 0)
		c, z[i] = Word(cc), Word(lo)
		c += z1
	}
	return
}
