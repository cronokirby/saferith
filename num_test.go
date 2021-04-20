package safenum

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func (Nat) Generate(r *rand.Rand, size int) reflect.Value {
	bytes := make([]byte, 64)
	r.Read(bytes)
	var n Nat
	n.SetBytes(bytes)
	return reflect.ValueOf(n)
}

func (Modulus) Generate(r *rand.Rand, size int) reflect.Value {
	bytes := make([]byte, 32)
	r.Read(bytes)
	bytes[len(bytes)-1] |= 1
	n := ModulusFromBytes(bytes)
	return reflect.ValueOf(*n)
}

func testAddZeroIdentity(n Nat) bool {
	var x, zero Nat
	zero.SetUint64(0)
	x.Add(&n, &zero, uint(len(n.limbs)*_W))
	if n.Cmp(&x) != 0 {
		return false
	}
	x.Add(&zero, &n, uint(len(n.limbs)*_W))
	return n.Cmp(&x) == 0
}

func TestAddZeroIdentity(t *testing.T) {
	err := quick.Check(testAddZeroIdentity, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testAddCommutative(a Nat, b Nat) bool {
	var aPlusB, bPlusA Nat
	for _, x := range []uint{256, 128, 64, 32, 8} {
		aPlusB.Add(&a, &b, x)
		bPlusA.Add(&b, &a, x)
		if aPlusB.Cmp(&bPlusA) != 0 {
			return false
		}
	}
	return true
}

func TestAddCommutative(t *testing.T) {
	err := quick.Check(testAddCommutative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testAddAssociative(a Nat, b Nat, c Nat) bool {
	var order1, order2 Nat
	for _, x := range []uint{256, 128, 64, 32, 8} {
		order1 = *order1.Add(&a, &b, x)
		order1.Add(&order1, &c, x)
		order2 = *order2.Add(&b, &c, x)
		order2.Add(&a, &order2, x)
		if order1.Cmp(&order2) != 0 {
			return false
		}
	}
	return true
}

func TestAddAssociative(t *testing.T) {
	err := quick.Check(testAddAssociative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testMulCommutative(a Nat, b Nat) bool {
	var aTimesB, bTimesA Nat
	for _, x := range []uint{256, 128, 64, 32, 8} {
		aTimesB.Mul(&a, &b, x)
		bTimesA.Mul(&b, &a, x)
		if aTimesB.Cmp(&bTimesA) != 0 {
			return false
		}
	}
	return true
}

func TestMulCommutative(t *testing.T) {
	err := quick.Check(testMulCommutative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testMulAssociative(a Nat, b Nat, c Nat) bool {
	var order1, order2 Nat
	for _, x := range []uint{256, 128, 64, 32, 8} {
		order1 = *order1.Mul(&a, &b, x)
		order1.Mul(&order1, &c, x)
		order2 = *order2.Mul(&b, &c, x)
		order2.Mul(&a, &order2, x)
		if order1.Cmp(&order2) != 0 {
			return false
		}
	}
	return true
}

func TestMulAssociative(t *testing.T) {
	err := quick.Check(testMulAssociative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testMulOneIdentity(n Nat) bool {
	var x, one Nat
	one.SetUint64(1)
	x.Mul(&n, &one, uint(len(n.limbs)*_W))
	if n.Cmp(&x) != 0 {
		return false
	}
	x.Mul(&one, &n, uint(len(n.limbs)*_W))
	return n.Cmp(&x) == 0
}

func TestMulOneIdentity(t *testing.T) {
	err := quick.Check(testMulOneIdentity, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModIdempotent(a Nat, m Modulus) bool {
	var way1, way2 Nat
	way1.Mod(&a, &m)
	way2.Mod(&way1, &m)
	return way1.Cmp(&way2) == 0
}

func TestModIdempotent(t *testing.T) {
	err := quick.Check(testModIdempotent, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModAddCommutative(a Nat, b Nat, m Modulus) bool {
	var aPlusB, bPlusA Nat
	aPlusB.ModAdd(&a, &b, &m)
	bPlusA.ModAdd(&b, &a, &m)
	return aPlusB.Cmp(&bPlusA) == 0
}

func TestModAddCommutative(t *testing.T) {
	err := quick.Check(testModAddCommutative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModAddAssociative(a Nat, b Nat, c Nat, m Modulus) bool {
	var order1, order2 Nat
	order1 = *order1.ModAdd(&a, &b, &m)
	order1.ModAdd(&order1, &c, &m)
	order2 = *order2.ModAdd(&b, &c, &m)
	order2.ModAdd(&a, &order2, &m)
	return order1.Cmp(&order2) == 0
}

func TestModAddAssociative(t *testing.T) {
	err := quick.Check(testModAddAssociative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModAddModSubInverse(a Nat, b Nat, m Modulus) bool {
	var c Nat
	c.ModAdd(&a, &b, &m)
	c.ModSub(&c, &b, &m)
	expected := new(Nat)
	expected.Mod(&a, &m)
	return c.Cmp(expected) == 0
}

func TestModAddModSubInverse(t *testing.T) {
	err := quick.Check(testModAddModSubInverse, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModMulCommutative(a Nat, b Nat, m Modulus) bool {
	var aPlusB, bPlusA Nat
	aPlusB.ModMul(&a, &b, &m)
	bPlusA.ModMul(&b, &a, &m)
	return aPlusB.Cmp(&bPlusA) == 0
}

func TestModMulCommutative(t *testing.T) {
	err := quick.Check(testModMulCommutative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModMulAssociative(a Nat, b Nat, c Nat, m Modulus) bool {
	var order1, order2 Nat
	order1 = *order1.ModMul(&a, &b, &m)
	order1.ModMul(&order1, &c, &m)
	order2 = *order2.ModMul(&b, &c, &m)
	order2.ModMul(&a, &order2, &m)
	return order1.Cmp(&order2) == 0
}

func TestModMulAssociative(t *testing.T) {
	err := quick.Check(testModMulAssociative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModInverseMultiplication(a Nat) bool {
	var scratch, one, zero Nat
	zero.SetUint64(0)
	one.SetUint64(1)
	for _, x := range []uint64{3, 5, 7, 13, 19, 47, 97} {
		m := ModulusFromUint64(x)
		scratch.Mod(&a, m)
		if scratch.Cmp(&zero) == 0 {
			continue
		}
		scratch.ModInverse(&a, m)
		scratch.ModMul(&scratch, &a, m)
		if scratch.Cmp(&one) != 0 {
			return false
		}
	}
	return true
}

func TestModInverseMultiplication(t *testing.T) {
	err := quick.Check(testModInverseMultiplication, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModInverseMinusOne(a Nat) bool {
	// Clear out the lowest bit
	a.limbs[0] &= ^Word(1)
	var zero Nat
	zero.SetUint64(0)
	if a.Cmp(&zero) == 0 {
		return true
	}
	var one Nat
	one.SetUint64(1)
	var z Nat
	z.Add(&a, &one, uint(len(a.limbs)*_W+1))
	m := ModulusFromNat(z)
	z.ModInverse(&a, m)
	return z.Cmp(&a) == 0
}

func TestModInverseMinusOne(t *testing.T) {
	err := quick.Check(testModInverseMinusOne, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModInverseEvenMinusOne(a Nat) bool {
	// Set the lowest bit
	a.limbs[0] |= 1
	var zero Nat
	zero.SetUint64(0)
	if a.Cmp(&zero) == 0 {
		return true
	}
	var one Nat
	one.SetUint64(1)
	var z Nat
	z.Add(&a, &one, uint(len(a.limbs)*_W+1))
	z.ModInverseEven(&a, &z)
	return z.Cmp(&a) == 0
}

func TestModInverseEvenMinusOne(t *testing.T) {
	err := quick.Check(testModInverseEvenMinusOne, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModInverseEvenOne(a Nat) bool {
	// Clear the lowest bit
	a.limbs[0] &= ^Word(1)
	var zero Nat
	zero.SetUint64(0)
	if a.Cmp(&zero) == 0 {
		return true
	}
	var one Nat
	one.SetUint64(1)
	var z Nat
	z.ModInverseEven(&one, &a)
	return z.Cmp(&one) == 0
}

func TestModInverseEvenOne(t *testing.T) {
	err := quick.Check(testModInverseEvenOne, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testExpAddition(x Nat, a Nat, b Nat, m Modulus) bool {
	var expA, expB, aPlusB, way1, way2 Nat
	expA.Exp(&x, &a, &m)
	expB.Exp(&x, &b, &m)
	// Enough bits to hold the full amount
	aPlusB.Add(&a, &b, uint(len(a.limbs)*_W)+1)
	way1.ModMul(&expA, &expB, &m)
	way2.Exp(&x, &aPlusB, &m)
	return way1.Cmp(&way2) == 0
}

func TestExpAddition(t *testing.T) {
	err := quick.Check(testExpAddition, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestUint64Creation(t *testing.T) {
	var x, y Nat
	x.SetUint64(0)
	y.SetUint64(0)
	if x.Cmp(&y) != 0 {
		t.Errorf("%+v != %+v", x, y)
	}
	x.SetUint64(1)
	if x.Cmp(&y) == 0 {
		t.Errorf("%+v == %+v", x, y)
	}
	x.SetUint64(0x1111)
	y.SetUint64(0x1111)
	if x.Cmp(&y) != 0 {
		t.Errorf("%+v != %+v", x, y)
	}
}

func TestAddExamples(t *testing.T) {
	var x, y, z Nat
	x.SetUint64(100)
	y.SetUint64(100)
	z.SetUint64(200)
	x = *x.Add(&x, &y, 8)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	z.SetUint64(300 - 256)
	x = *x.Add(&x, &y, 8)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(0xf3e5487232169930)
	y.SetUint64(0)
	z.SetUint64(0xf3e5487232169930)
	var x2 Nat
	x2.Add(&x, &y, 128)
	if x2.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestSubExamples(t *testing.T) {
	x := new(Nat).SetUint64(100)
	y := new(Nat).SetUint64(200)
	y.Sub(y, x, 8)
	if y.Cmp(x) != 0 {
		t.Errorf("%+v != %+v", y, x)
	}
}

func TestMulExamples(t *testing.T) {
	var x, y, z Nat
	x.SetUint64(10)
	y.SetUint64(10)
	z.SetUint64(100)
	x = *x.Mul(&x, &y, 8)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	z.SetUint64(232)
	x = *x.Mul(&x, &y, 8)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModAddExamples(t *testing.T) {
	m := ModulusFromUint64(13)
	var x, y, z Nat
	x.SetUint64(40)
	y.SetUint64(40)
	x = *x.ModAdd(&x, &y, m)
	z.SetUint64(2)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModMulExamples(t *testing.T) {
	var x, y, z Nat
	m := ModulusFromUint64(13)
	x.SetUint64(40)
	y.SetUint64(40)
	x = *x.ModMul(&x, &y, m)
	z.SetUint64(1)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	m = ModulusFromBytes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 1})
	x.SetUint64(1)
	x = *x.ModMul(&x, &x, m)
	z.SetUint64(1)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	m = ModulusFromBytes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 1})
	x.SetUint64(16390320477281102916)
	y.SetUint64(13641051446569424315)
	x = *x.ModMul(&x, &y, m)
	z.SetUint64(12559215458690093993)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModExamples(t *testing.T) {
	var x, test Nat
	x.SetUint64(40)
	m := ModulusFromUint64(13)
	x.Mod(&x, m)
	test.SetUint64(1)
	if x.Cmp(&test) != 0 {
		t.Errorf("%+v != %+v", x, test)
	}
	m = ModulusFromBytes([]byte{13, 0, 0, 0, 0, 0, 0, 0, 1})
	x.SetBytes([]byte{41, 0, 0, 0, 0, 0, 0, 0, 0})
	x.Mod(&x, m)
	test.SetBytes([]byte{1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFD})
	if x.Cmp(&test) != 0 {
		t.Errorf("%+v != %+v", x, test)
	}
}

func TestModInverseExamples(t *testing.T) {
	var x, z Nat
	x.SetUint64(2)
	m := ModulusFromUint64(13)
	x = *x.ModInverse(&x, m)
	z.SetUint64(7)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(16359684999990746055)
	m = ModulusFromUint64(7)
	x = *x.ModInverse(&x, m)
	z.SetUint64(3)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(461423694560)
	m = ModulusFromUint64(461423694561)
	z.ModInverse(&x, m)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestExpExamples(t *testing.T) {
	var x, y, z Nat
	x.SetUint64(3)
	y.SetUint64(345)
	m := ModulusFromUint64(13)
	x = *x.Exp(&x, &y, m)
	z.SetUint64(1)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	m = ModulusFromBytes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 1})
	x.SetUint64(1)
	y.SetUint64(2)
	x = *x.Exp(&x, &y, m)
	z.SetUint64(1)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestSetBytesExamples(t *testing.T) {
	var x, z Nat
	x.SetBytes([]byte{0x12, 0x34, 0x56})
	z.SetUint64(0x123456)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetBytes([]byte{0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF})
	z.SetUint64(0xAABBCCDDEEFF)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestFillBytesExamples(t *testing.T) {
	var x Nat
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD}
	x.SetBytes(expected)
	buf := make([]byte, 8)
	x.FillBytes(buf)
	if !bytes.Equal(expected, buf) {
		t.Errorf("%+v != %+v", expected, buf)
	}
}

func TestBytesExamples(t *testing.T) {
	var x Nat
	expected := []byte{0x11, 0x22, 0x33, 0x44, 0xAA, 0xBB, 0xCC, 0xDD}
	x.SetBytes(expected)
	out := x.Bytes()
	if !bytes.Equal(expected, out) {
		t.Errorf("%+v != %+v", expected, out)
	}
}

func TestModInverseEvenExamples(t *testing.T) {
	var z, x, m Nat
	x.SetUint64(9)
	m.SetUint64(10)
	x.ModInverseEven(&x, &m)
	z.SetUint64(9)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(1)
	m.SetUint64(10)
	x.ModInverseEven(&x, &m)
	z.SetUint64(1)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(19)
	x.ModInverseEven(&x, &m)
	z.SetUint64(9)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(99)
	x.ModInverseEven(&x, &m)
	z.SetUint64(9)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(999)
	m.SetUint64(1000)
	x.ModInverseEven(&x, &m)
	z.SetUint64(999)
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	// There's an edge case when the modulus is much larger than the input,
	// in which case when we do m^-1 mod x, we need to first calculate the remainder
	// of m.
	x.SetUint64(3)
	m.SetBytes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0})
	x.ModInverseEven(&x, &m)
	z.SetBytes([]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAB})
	if x.Cmp(&z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModSubExamples(t *testing.T) {
	m := ModulusFromUint64(13)
	x := new(Nat).SetUint64(0)
	y := new(Nat).SetUint64(1)
	x.ModSub(x, y, m)
	z := new(Nat).SetUint64(12)
	if x.Cmp(z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}
