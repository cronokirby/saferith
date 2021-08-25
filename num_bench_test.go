package saferith

import (
	"math/big"
	"testing"
)

var resultBig big.Int
var resultNat Nat

const _SIZE = 256

func ones() []byte {
	bytes := make([]byte, _SIZE)
	for i := 0; i < _SIZE; i++ {
		bytes[i] = 1
	}
	return bytes
}

func doubleOnes() []byte {
	bytes := make([]byte, 2*_SIZE)
	for i := 0; i < 2*_SIZE; i++ {
		bytes[i] = 1
	}
	return bytes
}

// a modulus of 2048 bits
func modulus2048() []byte {
	bytes := make([]byte, 256)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = 0xFD
	}
	return bytes
}

// an even modulus of 2048 bits
func modulus2048Even() []byte {
	bytes := make([]byte, 256)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = 0xFE
	}
	return bytes
}

// A 256 bit prime that's 3 mod 4
func prime3Mod4() []byte {
	bytes := make([]byte, 32)
	bytes[0] = 4
	bytes[31] = 0x4F
	return bytes
}

// A 256 bit prime that's 1 mod 4
func prime1Mod4() []byte {
	bytes := make([]byte, 32)
	bytes[0] = 4
	bytes[31] = 0x99
	return bytes
}

func BenchmarkAddBig(b *testing.B) {
	b.StopTimer()

	var x big.Int
	x.SetBytes(ones())

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Add(&x, &x)
		resultBig = z
	}
}

func _benchmarkModAddBig(m *big.Int, b *testing.B) {
	b.StopTimer()

	x := new(big.Int).SetBytes(ones())
	x.Mod(x, m)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Add(x, x)
		z.Mod(x, m)
		resultBig = z
	}
}

func BenchmarkModAddBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetUint64(13)
	_benchmarkModAddBig(&m, b)
}

func BenchmarkLargeModAddBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetBytes(modulus2048())
	_benchmarkModAddBig(&m, b)
}

func BenchmarkMulBig(b *testing.B) {
	b.StopTimer()

	var x big.Int
	x.SetBytes(ones())

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Mul(&x, &x)
		resultBig = z
	}
}

func _benchmarkModMulBig(m *big.Int, b *testing.B) {
	b.StopTimer()

	x := new(big.Int).SetBytes(ones())
	x.Mod(x, m)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Mul(x, x)
		z.Mod(x, m)
		resultBig = z
	}
}

func BenchmarkModMulBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetUint64(13)
	_benchmarkModMulBig(&m, b)
}

func BenchmarkLargeModMulBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetBytes(modulus2048())
	_benchmarkModMulBig(&m, b)
}

func _benchmarkModBig(m *big.Int, b *testing.B) {
	b.StopTimer()

	var x big.Int
	x.SetBytes(doubleOnes())

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Mod(&x, m)
		resultBig = z
	}
}

func BenchmarkModBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetUint64(13)
	_benchmarkModBig(&m, b)
}

func BenchmarkLargeModBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetBytes(modulus2048())
	_benchmarkModBig(&m, b)
}

func _benchmarkModInverseBig(m *big.Int, b *testing.B) {
	b.StopTimer()

	x := new(big.Int).SetBytes(ones())
	x.Mod(x, m)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.ModInverse(x, m)
		resultBig = z
	}
}

func BenchmarkModInverseBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetUint64(13)
	_benchmarkModInverseBig(&m, b)
}

func BenchmarkLargeModInverseBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetBytes(modulus2048())
	_benchmarkModInverseBig(&m, b)
}

func _benchmarkExpBig(m *big.Int, b *testing.B) {
	b.StopTimer()

	x := new(big.Int).SetBytes(ones())
	x.Mod(x, m)
	y := new(big.Int).SetBytes(ones())

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Exp(x, y, m)
		resultBig = z
	}
}

func BenchmarkExpBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetUint64(13)
	_benchmarkExpBig(&m, b)
}

func BenchmarkLargeExpBig(b *testing.B) {
	b.StopTimer()

	var m big.Int
	m.SetBytes(modulus2048())
	_benchmarkExpBig(&m, b)
}

func BenchmarkSetBytesBig(b *testing.B) {
	b.StopTimer()

	bytes := ones()

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.SetBytes(bytes)
		resultBig = z
	}
}

func BenchmarkModSqrt3Mod4Big(b *testing.B) {
	b.StopTimer()

	p := new(big.Int).SetBytes(prime3Mod4())
	// This is a large square modulo p
	x := new(big.Int).Sub(p, new(big.Int).SetUint64(5))

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var z big.Int
		z.ModSqrt(x, p)
		resultBig = z
	}
}

func BenchmarkAddNat(b *testing.B) {
	b.StopTimer()

	var x Nat
	x.SetBytes(ones())

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Add(&x, &x, _SIZE*8)
		resultNat = z
	}
}

func _benchmarkModAddNat(m *Modulus, b *testing.B) {
	b.StopTimer()

	x := new(Nat).SetBytes(ones())
	x.Mod(x, m)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.ModAdd(x, x, m)
		resultNat = z
	}
}

func BenchmarkModAddNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromUint64(13)
	_benchmarkModAddNat(m, b)
}

func BenchmarkLargeModAddNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromBytes(modulus2048())
	_benchmarkModAddNat(m, b)
}

func _benchmarkModNegNat(m *Modulus, b *testing.B) {
	b.StopTimer()

	x := new(Nat).SetBytes(ones())
	x.Mod(x, m)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.ModNeg(x, m)
		resultNat = z
	}
}

func BenchmarkModNegNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromUint64(13)
	_benchmarkModNegNat(m, b)
}

func BenchmarkLargeModNegNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromBytes(modulus2048())
	_benchmarkModNegNat(m, b)
}

func BenchmarkMulNat(b *testing.B) {
	b.StopTimer()

	var x Nat
	x.SetBytes(ones())

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Mul(&x, &x, _SIZE*2*8)
		resultNat = z
	}
}

func _benchmarkModMulNat(m *Modulus, b *testing.B) {
	b.StopTimer()

	x := new(Nat).SetBytes(ones())
	x.Mod(x, m)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.ModMul(x, x, m)
		resultNat = z
	}
}

func BenchmarkModMulNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromUint64(13)
	_benchmarkModMulNat(m, b)
}

func BenchmarkLargeModMulNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromBytes(modulus2048())
	_benchmarkModMulNat(m, b)
}

func BenchmarkLargeModMulNatEven(b *testing.B) {
	b.StopTimer()

	m := ModulusFromBytes(modulus2048Even())
	_benchmarkModMulNat(m, b)
}

func _benchmarkModNat(m *Modulus, b *testing.B) {
	b.StopTimer()

	x := new(Nat).SetBytes(doubleOnes())

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Mod(x, m)
		resultNat = z
	}
}

func BenchmarkModNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromUint64(13)
	_benchmarkModNat(m, b)
}

func BenchmarkLargeModNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromBytes(modulus2048())
	_benchmarkModNat(m, b)
}

func _benchmarkModInverseNat(m *Modulus, b *testing.B) {
	b.StopTimer()

	x := new(Nat).SetBytes(ones())
	x.Mod(x, m)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.ModInverse(x, m)
		resultNat = z
	}
}

func BenchmarkModInverseNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromUint64(13)
	_benchmarkModInverseNat(m, b)
}

func BenchmarkLargeModInverseNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromBytes(modulus2048())
	_benchmarkModInverseNat(m, b)
}

func _benchmarkModInverseEvenNat(m *Modulus, b *testing.B) {
	b.StopTimer()

	var x Nat
	x.SetBytes(ones())

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.ModInverse(&x, m)
		resultNat = z
	}
}

func BenchmarkModInverseEvenNat(b *testing.B) {
	b.StopTimer()
	m := ModulusFromUint64(14)
	_benchmarkModInverseEvenNat(m, b)
}

func BenchmarkLargeModInverseEvenNat(b *testing.B) {
	b.StopTimer()
	var one, m Nat
	m.SetBytes(modulus2048())
	one.SetUint64(1)
	m.Add(&m, &one, 2048)
	_benchmarkModInverseEvenNat(ModulusFromNat(&m), b)
}

func _benchmarkExpNat(m *Modulus, b *testing.B) {
	b.StopTimer()

	x := new(Nat).SetBytes(ones())
	y := new(Nat).SetBytes(ones())
	x.Mod(x, m)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Exp(x, y, m)
		resultNat = z
	}
}

func BenchmarkExpNat(b *testing.B) {
	b.StopTimer()
	m := ModulusFromUint64(13)
	_benchmarkExpNat(m, b)
}

func BenchmarkLargeExpNat(b *testing.B) {
	b.StopTimer()
	m := ModulusFromBytes(modulus2048())
	_benchmarkExpNat(m, b)
}

func BenchmarkLargeExpNatEven(b *testing.B) {
	b.StopTimer()
	m := ModulusFromBytes(modulus2048Even())
	_benchmarkExpNat(m, b)
}

func BenchmarkSetBytesNat(b *testing.B) {
	b.StopTimer()

	bytes := ones()

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.SetBytes(bytes)
		resultNat = z
	}
}

func BenchmarkMontgomeryMul(b *testing.B) {
	b.StopTimer()
	x := new(Nat).SetBytes(ones())
	y := new(Nat).SetBytes(ones())
	scratch := new(Nat).SetBytes(ones())
	out := new(Nat).SetBytes(ones())
	m := ModulusFromBytes(modulus2048())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		montgomeryMul(x.limbs, y.limbs, out.limbs, scratch.limbs, m)
	}
}

func BenchmarkModSqrt3Mod4Nat(b *testing.B) {
	b.StopTimer()

	p := new(Nat).SetBytes(prime3Mod4())
	// This is a large square modulo p
	x := new(Nat).Sub(p, new(Nat).SetUint64(5), 256)
	pMod := ModulusFromNat(p)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var z Nat
		z.ModSqrt(x, pMod)
		resultNat = z
	}
}

func BenchmarkModSqrt1Mod4Nat(b *testing.B) {
	b.StopTimer()

	p := new(Nat).SetBytes(prime1Mod4())
	// This is a large square modulo p
	x := new(Nat).Sub(p, new(Nat).SetUint64(6), 256)
	pMod := ModulusFromNat(p)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var z Nat
		z.ModSqrt(x, pMod)
		resultNat = z
	}
}

func _benchmarkDivNat(m *Modulus, b *testing.B) {
	b.StopTimer()

	x := new(Nat).SetBytes(doubleOnes())

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		var z Nat
		z.Div(x, m, m.BitLen())
		resultNat = z
	}
}

func BenchmarkDivNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromUint64(13)
	_benchmarkDivNat(m, b)
}

func BenchmarkLargeDivNat(b *testing.B) {
	b.StopTimer()

	m := ModulusFromBytes(modulus2048())
	_benchmarkDivNat(m, b)
}
