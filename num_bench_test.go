package safenum

import (
	"math/big"
	"testing"
)

var resultBig big.Int
var resultNat Nat

const _SIZE = 512

func BenchmarkAddBig(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	var x big.Int
	x.SetBytes(bytes)
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Add(&x, &x)
		resultBig = z
	}
}

func BenchmarkModAddBig(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
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
		resultBig = z
	}
}

func BenchmarkMulBig(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	var x big.Int
	x.SetBytes(bytes)
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Mul(&x, &x)
		resultBig = z
	}
}

func BenchmarkModMulBig(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	var x big.Int
	var m big.Int
	x.SetBytes(bytes)
	m.SetUint64(13)
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Mul(&x, &x)
		z.Mod(&x, &m)
		resultBig = z
	}
}

func BenchmarkModBig(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	var x big.Int
	var m big.Int
	m.SetUint64(13)
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Mod(&x, &m)
		resultBig = z
	}
}

func BenchmarkAddNat(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	var x Nat
	x.SetBytes(bytes)
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Add(&x, &x, _SIZE*8)
		resultNat = z
	}
}

func BenchmarkModAddNat(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
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

func BenchmarkMulNat(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	var x Nat
	x.SetBytes(bytes)
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Mul(&x, &x, _SIZE*2*8)
		resultNat = z
	}
}

func BenchmarkModMulNat(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	var x Nat
	var m Nat
	x.SetBytes(bytes)
	m.SetUint64(13)
	for n := 0; n < b.N; n++ {
		var z Nat
		z.ModMul(&x, &x, &m)
		resultNat = z
	}
}

func BenchmarkModNat(b *testing.B) {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	var x Nat
	var m Nat
	m.SetUint64(13)
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Mod(&x, &m)
		resultNat = z
	}
}
