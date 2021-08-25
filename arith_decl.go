// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE_go file.

// +build !math_big_pure_go

package saferith

// implemented in arith_$GOARCH.s
func addVV(z, x, y []Word) (c Word)
func subVV(z, x, y []Word) (c Word)
func addVW(z, x []Word, y Word) (c Word)
func subVW(z, x []Word, y Word) (c Word)
func shlVU(z, x []Word, s uint) (c Word)
func shrVU(z, x []Word, s uint) (c Word)
func addMulVVW(z, x []Word, y Word) (c Word)
