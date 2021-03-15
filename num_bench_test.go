package safenum

import (
	"math/big"
	"testing"
)

var resultInt big.Int
var resultNat Nat

func BenchmarkAddBig(b *testing.B) {
	bytes := make([]byte, 100)
	for i := 0; i < 100; i++ {
		bytes[i] = 1
	}
	var x big.Int
	x.SetBytes(bytes)
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Add(&x, &x)
		resultInt = z
	}
}

func BenchmarkAddBigMod(b *testing.B) {
	bytes := make([]byte, 100)
	for i := 0; i < 100; i++ {
		bytes[i] = 1
	}
	var x big.Int
	var m big.Int
	x.SetBytes(bytes)
	m.SetUint64(13)
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Add(&x, &x)
		z.Mod(&x, &m)
		resultInt = z
	}
}

func BenchmarkAddNat(b *testing.B) {
	bytes := make([]byte, 100)
	for i := 0; i < 100; i++ {
		bytes[i] = 1
	}
	var x Nat
	x.SetBytes(bytes)
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Add(&x, &x, 100*64)
		resultNat = z
	}
}

func BenchmarkAddNatMod(b *testing.B) {
	bytes := make([]byte, 100)
	for i := 0; i < 100; i++ {
		bytes[i] = 1
	}
	var x Nat
	var m Nat
	x.SetBytes(bytes)
	m.SetUint64(13)
	for n := 0; n < b.N; n++ {
		var z Nat
		z.ModAdd(&x, &x, &m)
		resultNat = z
	}
}
