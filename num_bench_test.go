package safenum

import (
	"math/big"
	"testing"
)

var resultBig big.Int
var resultNat Nat

const _SIZE = 512

func ones() []byte {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	return bytes
}

// 2^3217 - 1 is a mersenne prime
func largePrime() []byte {
	bytes := make([]byte, 403)
	for i := 1; i < len(bytes); i++ {
		bytes[i] = 0xFF
	}
	bytes[0] = 0x01
	return bytes
}

func BenchmarkAddBig(b *testing.B) {
	var x big.Int
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Add(&x, &x)
		resultBig = z
	}
}

func _benchmarkModAddBig(m *big.Int, b *testing.B) {
	var x big.Int
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Add(&x, &x)
		z.Mod(&x, m)
		resultBig = z
	}
}

func BenchmarkModAddBig(b *testing.B) {
	var m big.Int
	m.SetUint64(13)
	_benchmarkModAddBig(&m, b)
}

func BenchmarkLargeModAddBig(b *testing.B) {
	var m big.Int
	m.SetBytes(largePrime())
	_benchmarkModAddBig(&m, b)
}

func BenchmarkMulBig(b *testing.B) {
	var x big.Int
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Mul(&x, &x)
		resultBig = z
	}
}

func _benchmarkModMulBig(m *big.Int, b *testing.B) {
	var x big.Int
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Mul(&x, &x)
		z.Mod(&x, m)
		resultBig = z
	}
}

func BenchmarkModMulBig(b *testing.B) {
	var m big.Int
	m.SetUint64(13)
	_benchmarkModMulBig(&m, b)
}

func BenchmarkLargeModMulBig(b *testing.B) {
	var m big.Int
	m.SetBytes(largePrime())
	_benchmarkModMulBig(&m, b)
}

func _benchmarkModBig(m *big.Int, b *testing.B) {
	var x big.Int
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Mod(&x, m)
		resultBig = z
	}
}

func BenchmarkModBig(b *testing.B) {
	var m big.Int
	m.SetUint64(13)
	_benchmarkModBig(&m, b)
}

func BenchmarkLargeModBig(b *testing.B) {
	var m big.Int
	m.SetBytes(largePrime())
	_benchmarkModBig(&m, b)
}

func _benchmarkModInverseBig(m *big.Int, b *testing.B) {
	var x big.Int
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.ModInverse(&x, m)
		resultBig = z
	}
}

func BenchmarkModInverseBig(b *testing.B) {
	var m big.Int
	m.SetUint64(13)
	_benchmarkModInverseBig(&m, b)
}

func BenchmarkLargeModInverseBig(b *testing.B) {
	var m big.Int
	m.SetBytes(largePrime())
	_benchmarkModInverseBig(&m, b)
}

func _benchmarkExpBig(m *big.Int, b *testing.B) {
	var x big.Int
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Exp(&x, &x, m)
		resultBig = z
	}
}

func BenchmarkExpBig(b *testing.B) {
	var m big.Int
	m.SetUint64(13)
	_benchmarkExpBig(&m, b)
}

func BenchmarkLargeExpBig(b *testing.B) {
	var m big.Int
	m.SetBytes(largePrime())
	_benchmarkExpBig(&m, b)
}

func BenchmarkSetBytesBig(b *testing.B) {
	bytes := ones()
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.SetBytes(bytes)
		resultBig = z
	}
}

func BenchmarkAddNat(b *testing.B) {
	var x Nat
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Add(&x, &x, _SIZE*8)
		resultNat = z
	}
}

func _benchmarkModAddNat(m *Modulus, b *testing.B) {
	var x Nat
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Add(&x, &x, _SIZE*8)
		z.Mod(&x, m)
		resultNat = z
	}
}

func BenchmarkModAddNat(b *testing.B) {
	m := ModulusFromUint64(13)
	_benchmarkModAddNat(&m, b)
}

func BenchmarkLargeModAddNat(b *testing.B) {
	m := ModulusFromBytes(largePrime())
	_benchmarkModAddNat(&m, b)
}

func BenchmarkMulNat(b *testing.B) {
	var x Nat
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Mul(&x, &x, _SIZE*2*8)
		resultNat = z
	}
}

func _benchmarkModMulNat(m *Modulus, b *testing.B) {
	var x Nat
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z Nat
		z.ModMul(&x, &x, m)
		resultNat = z
	}
}

func BenchmarkModMulNat(b *testing.B) {
	m := ModulusFromUint64(13)
	_benchmarkModMulNat(&m, b)
}

func BenchmarkLargeModMulNat(b *testing.B) {
	m := ModulusFromBytes(largePrime())
	_benchmarkModMulNat(&m, b)
}

func _benchmarkModNat(m *Modulus, b *testing.B) {
	var x Nat
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Mod(&x, m)
		resultNat = z
	}
}

func BenchmarkModNat(b *testing.B) {
	m := ModulusFromUint64(13)
	_benchmarkModNat(&m, b)
}

func BenchmarkLargeModNat(b *testing.B) {
	m := ModulusFromBytes(largePrime())
	_benchmarkModNat(&m, b)
}

func _benchmarkModInverseNat(m *Modulus, b *testing.B) {
	var x Nat
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z Nat
		z.ModInverse(&x, m)
		resultNat = z
	}
}

func BenchmarkModInverseNat(b *testing.B) {
	m := ModulusFromUint64(13)
	_benchmarkModInverseNat(&m, b)
}

func BenchmarkLargeModInverseNat(b *testing.B) {
	m := ModulusFromBytes(largePrime())
	_benchmarkModInverseNat(&m, b)
}

func _benchmarkExpNat(m *Modulus, b *testing.B) {
	var x Nat
	x.SetBytes(ones())
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Exp(&x, &x, m)
		resultNat = z
	}
}

func BenchmarkExpNat(b *testing.B) {
	m := ModulusFromUint64(13)
	_benchmarkExpNat(&m, b)
}

func BenchmarkLargeExpNat(b *testing.B) {
	m := ModulusFromBytes(largePrime())
	_benchmarkExpNat(&m, b)
}

func BenchmarkSetBytesNat(b *testing.B) {
	bytes := ones()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.SetBytes(bytes)
		resultNat = z
	}
}
