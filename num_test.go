package safenum

import (
	"bytes"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func (Nat) Generate(r *rand.Rand, size int) reflect.Value {
	bytes := make([]byte, 16*_S)
	r.Read(bytes)
	var n Nat
	n.SetBytes(bytes)
	return reflect.ValueOf(n)
}

func (Modulus) Generate(r *rand.Rand, size int) reflect.Value {
	bytes := make([]byte, 8*_S)
	r.Read(bytes)
	// Ensure that our number isn't 0, but being even is ok
	bytes[len(bytes)-1] |= 0b10
	n := ModulusFromBytes(bytes)
	return reflect.ValueOf(*n)
}

func testBigConversion(x Nat) bool {
	if !x.checkInvariants() {
		return false
	}
	xBig := x.Big()
	xNatAgain := new(Nat).SetBig(xBig, x.AnnouncedLen())
	if !xNatAgain.checkInvariants() {
		return false
	}
	return x.Eq(xNatAgain) == 1
}

func TestBigConversion(t *testing.T) {
	err := quick.Check(testBigConversion, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testByteVsBytes(x Nat) bool {
	if !x.checkInvariants() {
		return false
	}
	bytes := x.Bytes()
	for i := 0; i < len(bytes); i++ {
		if x.Byte(i) != bytes[len(bytes)-i-1] {
			return false
		}
	}
	return true
}

func TestByteVsBytes(t *testing.T) {
	err := quick.Check(testByteVsBytes, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testSetBytesRoundTrip(expected []byte) bool {
	x := new(Nat).SetBytes(expected)
	actual := x.Bytes()
	return bytes.Equal(expected, actual)
}

func TestSetBytesRoundTrip(t *testing.T) {
	err := quick.Check(testSetBytesRoundTrip, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testAddZeroIdentity(n Nat) bool {
	if !n.checkInvariants() {
		return false
	}
	var x, zero Nat
	zero.SetUint64(0)
	x.Add(&n, &zero, len(n.limbs)*_W)
	if !x.checkInvariants() {
		return false
	}
	if n.Eq(&x) != 1 {
		return false
	}
	x.Add(&zero, &n, len(n.limbs)*_W)
	if !x.checkInvariants() {
		return false
	}
	return n.Eq(&x) == 1
}

func TestAddZeroIdentity(t *testing.T) {
	err := quick.Check(testAddZeroIdentity, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testAddCommutative(a Nat, b Nat) bool {
	if !(a.checkInvariants() && b.checkInvariants()) {
		return false
	}
	var aPlusB, bPlusA Nat
	for _, x := range []int{256, 128, 64, 32, 8} {
		aPlusB.Add(&a, &b, x)
		bPlusA.Add(&b, &a, x)
		if !(aPlusB.checkInvariants() && bPlusA.checkInvariants()) {
			return false
		}
		if aPlusB.Eq(&bPlusA) != 1 {
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

func testCondAssign(a Nat, b Nat) bool {
	if !(a.checkInvariants() && b.checkInvariants()) {
		return false
	}
	shouldBeA := new(Nat).SetNat(&a)
	shouldBeB := new(Nat).SetNat(&a)
	shouldBeA.CondAssign(0, &b)
	shouldBeB.CondAssign(1, &b)
	if !(shouldBeA.checkInvariants() && shouldBeB.checkInvariants()) {
		return false
	}
	return shouldBeA.Eq(&a) == 1 && shouldBeB.Eq(&b) == 1
}

func TestCondAssign(t *testing.T) {
	err := quick.Check(testCondAssign, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testAddAssociative(a Nat, b Nat, c Nat) bool {
	if !(a.checkInvariants() && b.checkInvariants() && c.checkInvariants()) {
		return false
	}
	var order1, order2 Nat
	for _, x := range []int{256, 128, 64, 32, 8} {
		order1 = *order1.Add(&a, &b, x)
		order1.Add(&order1, &c, x)
		order2 = *order2.Add(&b, &c, x)
		order2.Add(&a, &order2, x)
		if !(order1.checkInvariants() && order2.checkInvariants()) {
			return false
		}
		if order1.Eq(&order2) != 1 {
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

func testModAddNegIsSub(a Nat, b Nat, m Modulus) bool {
	if !(a.checkInvariants() && b.checkInvariants()) {
		return false
	}
	subbed := new(Nat).ModSub(&a, &b, &m)
	negated := new(Nat).ModNeg(&b, &m)
	addWithNegated := new(Nat).ModAdd(&a, negated, &m)
	if !(subbed.checkInvariants() && negated.checkInvariants() && addWithNegated.checkInvariants()) {
		return false
	}
	return subbed.Eq(addWithNegated) == 1
}

func TestModAddNegIsSub(t *testing.T) {
	err := quick.Check(testModAddNegIsSub, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testMulCommutative(a Nat, b Nat) bool {
	if !(a.checkInvariants() && b.checkInvariants()) {
		return false
	}
	var aTimesB, bTimesA Nat
	for _, x := range []int{256, 128, 64, 32, 8} {
		aTimesB.Mul(&a, &b, x)
		bTimesA.Mul(&b, &a, x)
		if !(aTimesB.checkInvariants() && bTimesA.checkInvariants()) {
			return false
		}
		if aTimesB.Eq(&bTimesA) != 1 {
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
	if !(a.checkInvariants() && b.checkInvariants() && c.checkInvariants()) {
		return false
	}
	var order1, order2 Nat
	for _, x := range []int{256, 128, 64, 32, 8} {
		order1 = *order1.Mul(&a, &b, x)
		order1.Mul(&order1, &c, x)
		order2 = *order2.Mul(&b, &c, x)
		order2.Mul(&a, &order2, x)
		if !(order1.checkInvariants() && order2.checkInvariants()) {
			return false
		}
		if order1.Eq(&order2) != 1 {
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
	if !n.checkInvariants() {
		return false
	}
	var x, one Nat
	one.SetUint64(1)
	x.Mul(&n, &one, len(n.limbs)*_W)
	if !x.checkInvariants() {
		return false
	}
	if n.Eq(&x) != 1 {
		return false
	}
	x.Mul(&one, &n, len(n.limbs)*_W)
	if !x.checkInvariants() {
		return false
	}
	return n.Eq(&x) == 1
}

func TestMulOneIdentity(t *testing.T) {
	err := quick.Check(testMulOneIdentity, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModIdempotent(a Nat, m Modulus) bool {
	if !a.checkInvariants() {
		return false
	}
	var way1, way2 Nat
	way1.Mod(&a, &m)
	way2.Mod(&way1, &m)
	if !(way1.checkInvariants() && way2.checkInvariants()) {
		return false
	}
	return way1.Eq(&way2) == 1
}

func TestModIdempotent(t *testing.T) {
	err := quick.Check(testModIdempotent, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModAddCommutative(a Nat, b Nat, m Modulus) bool {
	if !(a.checkInvariants() && b.checkInvariants()) {
		return false
	}
	var aPlusB, bPlusA Nat
	aPlusB.ModAdd(&a, &b, &m)
	bPlusA.ModAdd(&b, &a, &m)
	if !(aPlusB.checkInvariants() && bPlusA.checkInvariants()) {
		return false
	}
	return aPlusB.Eq(&bPlusA) == 1
}

func TestModAddCommutative(t *testing.T) {
	err := quick.Check(testModAddCommutative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModAddAssociative(a Nat, b Nat, c Nat, m Modulus) bool {
	if !(a.checkInvariants() && b.checkInvariants() && c.checkInvariants()) {
		return false
	}
	var order1, order2 Nat
	order1 = *order1.ModAdd(&a, &b, &m)
	order1.ModAdd(&order1, &c, &m)
	order2 = *order2.ModAdd(&b, &c, &m)
	order2.ModAdd(&a, &order2, &m)
	if !(order1.checkInvariants() && order2.checkInvariants()) {
		return false
	}
	return order1.Eq(&order2) == 1
}

func TestModAddAssociative(t *testing.T) {
	err := quick.Check(testModAddAssociative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModAddModSubInverse(a Nat, b Nat, m Modulus) bool {
	if !(a.checkInvariants() && b.checkInvariants()) {
		return false
	}
	var c Nat
	c.ModAdd(&a, &b, &m)
	c.ModSub(&c, &b, &m)
	expected := new(Nat)
	expected.Mod(&a, &m)
	if !(c.checkInvariants() && expected.checkInvariants()) {
		return false
	}
	return c.Eq(expected) == 1
}

func TestModAddModSubInverse(t *testing.T) {
	err := quick.Check(testModAddModSubInverse, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModMulCommutative(a Nat, b Nat, m Modulus) bool {
	if !(a.checkInvariants() && b.checkInvariants()) {
		return false
	}
	var aPlusB, bPlusA Nat
	aPlusB.ModMul(&a, &b, &m)
	bPlusA.ModMul(&b, &a, &m)
	if !(aPlusB.checkInvariants() && bPlusA.checkInvariants()) {
		return false
	}
	return aPlusB.Eq(&bPlusA) == 1
}

func TestModMulCommutative(t *testing.T) {
	err := quick.Check(testModMulCommutative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModMulAssociative(a Nat, b Nat, c Nat, m Modulus) bool {
	if !(a.checkInvariants() && b.checkInvariants() && c.checkInvariants()) {
		return false
	}
	var order1, order2 Nat
	order1 = *order1.ModMul(&a, &b, &m)
	order1.ModMul(&order1, &c, &m)
	order2 = *order2.ModMul(&b, &c, &m)
	order2.ModMul(&a, &order2, &m)
	if !(order1.checkInvariants() && order2.checkInvariants()) {
		return false
	}
	return order1.Eq(&order2) == 1
}

func TestModMulAssociative(t *testing.T) {
	err := quick.Check(testModMulAssociative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModInverseMultiplication(a Nat) bool {
	if !a.checkInvariants() {
		return false
	}
	var scratch, one, zero Nat
	zero.SetUint64(0)
	one.SetUint64(1)
	for _, x := range []uint64{3, 5, 7, 13, 19, 47, 97} {
		m := ModulusFromUint64(x)
		scratch.Mod(&a, m)
		if scratch.Eq(&zero) == 1 {
			continue
		}
		scratch.ModInverse(&a, m)
		scratch.ModMul(&scratch, &a, m)
		if !scratch.checkInvariants() {
			return false
		}
		if scratch.Eq(&one) != 1 {
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
	if !a.checkInvariants() {
		return false
	}
	// Clear out the lowest bit
	a.limbs[0] &= ^Word(1)
	var zero Nat
	zero.SetUint64(0)
	if a.Eq(&zero) == 1 {
		return true
	}
	var one Nat
	one.SetUint64(1)
	z := new(Nat).Add(&a, &one, a.AnnouncedLen()+1)
	m := ModulusFromNat(z)
	z.ModInverse(&a, m)
	if !z.checkInvariants() {
		return false
	}
	return z.Eq(&a) == 1
}

func TestModInverseMinusOne(t *testing.T) {
	err := quick.Check(testModInverseMinusOne, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModInverseEvenMinusOne(a Nat) bool {
	if !a.checkInvariants() {
		return false
	}
	// Set the lowest bit
	a.limbs[0] |= 1
	var zero Nat
	zero.SetUint64(0)
	if a.Eq(&zero) == 1 {
		return true
	}
	var one Nat
	one.SetUint64(1)
	var z Nat
	z.Add(&a, &one, a.AnnouncedLen()+1)
	if !z.checkInvariants() {
		return false
	}
	z2 := new(Nat).ModInverse(&a, ModulusFromNat(&z))
	if !z2.checkInvariants() {
		return false
	}
	return z2.Eq(&a) == 1
}

func TestModInverseEvenMinusOne(t *testing.T) {
	err := quick.Check(testModInverseEvenMinusOne, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModInverseEvenOne(a Nat) bool {
	if !a.checkInvariants() {
		return false
	}
	// Clear the lowest bit
	a.limbs[0] &= ^Word(1)
	var zero Nat
	zero.SetUint64(0)
	if a.Eq(&zero) == 1 {
		return true
	}
	var one Nat
	one.SetUint64(1)
	var z Nat
	m := ModulusFromNat(&a)
	z.ModInverse(&one, m)
	if !z.checkInvariants() {
		return false
	}
	return z.Eq(&one) == 1
}

func TestModInverseEvenOne(t *testing.T) {
	err := quick.Check(testModInverseEvenOne, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testExpAddition(x Nat, a Nat, b Nat, m Modulus) bool {
	if !(x.checkInvariants() && a.checkInvariants() && b.checkInvariants()) {
		return false
	}
	var expA, expB, aPlusB, way1, way2 Nat
	expA.Exp(&x, &a, &m)
	expB.Exp(&x, &b, &m)
	// Enough bits to hold the full amount
	aPlusB.Add(&a, &b, len(a.limbs)*_W+1)
	way1.ModMul(&expA, &expB, &m)
	way2.Exp(&x, &aPlusB, &m)
	if !(way1.checkInvariants() && way2.checkInvariants() && aPlusB.checkInvariants()) {
		return false
	}
	return way1.Eq(&way2) == 1
}

func TestExpAddition(t *testing.T) {
	err := quick.Check(testExpAddition, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testSqrtRoundTrip(x *Nat, p *Modulus) bool {
	xSquared := x.ModMul(x, x, p)
	xRoot := new(Nat).ModSqrt(xSquared, p)
	if !(xRoot.checkInvariants() && xSquared.checkInvariants()) {
		return false
	}
	xRoot.ModMul(xRoot, xRoot, p)
	if !xRoot.checkInvariants() {
		return false
	}
	return xRoot.Eq(xSquared) == 1
}

func testModSqrt(x Nat) bool {
	if !x.checkInvariants() {
		return false
	}
	p := ModulusFromBytes([]byte{
		13,
	})
	if !testSqrtRoundTrip(&x, p) {
		return false
	}
	p = ModulusFromUint64((1 << 61) - 1)
	if !testSqrtRoundTrip(&x, p) {
		return false
	}
	p = ModulusFromBytes([]byte{
		0x1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	})
	if !testSqrtRoundTrip(&x, p) {
		return false
	}
	p = ModulusFromBytes([]byte{
		0x3,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfb,
	})
	if !testSqrtRoundTrip(&x, p) {
		return false
	}
	// 2^224 - 2^96 + 1
	p = ModulusFromBytes([]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 1,
	})
	return testSqrtRoundTrip(&x, p)
}

func TestModSqrt(t *testing.T) {
	err := quick.Check(testModSqrt, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testMultiplyThenDivide(x Nat, m Modulus) bool {
	if !x.checkInvariants() {
		return false
	}
	mNat := &m.nat

	xm := new(Nat).Mul(&x, mNat, x.AnnouncedLen()+mNat.AnnouncedLen())
	divided := new(Nat).Div(xm, &m, x.AnnouncedLen())
	if divided.Eq(&x) != 1 {
		return false
	}
	// Adding m - 1 shouldn't change the result either
	xm.Add(xm, new(Nat).Sub(mNat, new(Nat).SetUint64(1), xm.AnnouncedLen()), xm.AnnouncedLen())
	divided = new(Nat).Div(xm, &m, x.AnnouncedLen())
	if !(divided.checkInvariants() && xm.checkInvariants()) {
		return false
	}
	return divided.Eq(&x) == 1
}

func TestMultiplyThenDivide(t *testing.T) {
	err := quick.Check(testMultiplyThenDivide, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestUint64Creation(t *testing.T) {
	var x, y Nat
	x.SetUint64(0)
	y.SetUint64(0)
	if x.Eq(&y) != 1 {
		t.Errorf("%+v != %+v", x, y)
	}
	x.SetUint64(1)
	if x.Eq(&y) == 1 {
		t.Errorf("%+v == %+v", x, y)
	}
	x.SetUint64(0x1111)
	y.SetUint64(0x1111)
	if x.Eq(&y) != 1 {
		t.Errorf("%+v != %+v", x, y)
	}
}

func TestAddExamples(t *testing.T) {
	var x, y, z Nat
	x.SetUint64(100)
	y.SetUint64(100)
	z.SetUint64(200)
	x = *x.Add(&x, &y, 8)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	z.SetUint64(300 - 256)
	x = *x.Add(&x, &y, 8)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(0xf3e5487232169930)
	y.SetUint64(0)
	z.SetUint64(0xf3e5487232169930)
	var x2 Nat
	x2.Add(&x, &y, 128)
	if x2.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestSubExamples(t *testing.T) {
	x := new(Nat).SetUint64(100)
	y := new(Nat).SetUint64(200)
	y.Sub(y, x, 8)
	if y.Eq(x) != 1 {
		t.Errorf("%+v != %+v", y, x)
	}
}

func TestMulExamples(t *testing.T) {
	var x, y, z Nat
	x.SetUint64(10)
	y.SetUint64(10)
	z.SetUint64(100)
	x = *x.Mul(&x, &y, 8)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	z.SetUint64(232)
	x = *x.Mul(&x, &y, 8)
	if x.Eq(&z) != 1 {
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
	if x.Eq(&z) != 1 {
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
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	m = ModulusFromBytes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 1})
	x.SetUint64(1)
	x = *x.ModMul(&x, &x, m)
	z.SetUint64(1)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	m = ModulusFromBytes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 1})
	x.SetUint64(16390320477281102916)
	y.SetUint64(13641051446569424315)
	x = *x.ModMul(&x, &y, m)
	z.SetUint64(12559215458690093993)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModExamples(t *testing.T) {
	var x, test Nat
	x.SetUint64(40)
	m := ModulusFromUint64(13)
	x.Mod(&x, m)
	test.SetUint64(1)
	if x.Eq(&test) != 1 {
		t.Errorf("%+v != %+v", x, test)
	}
	m = ModulusFromBytes([]byte{13, 0, 0, 0, 0, 0, 0, 0, 1})
	x.SetBytes([]byte{41, 0, 0, 0, 0, 0, 0, 0, 0})
	x.Mod(&x, m)
	test.SetBytes([]byte{1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFD})
	if x.Eq(&test) != 1 {
		t.Errorf("%+v != %+v", x, test)
	}
}

func TestModInverseExamples(t *testing.T) {
	var x, z Nat
	x.SetUint64(2)
	m := ModulusFromUint64(13)
	x = *x.ModInverse(&x, m)
	z.SetUint64(7)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(16359684999990746055)
	m = ModulusFromUint64(7)
	x = *x.ModInverse(&x, m)
	z.SetUint64(3)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(461423694560)
	m = ModulusFromUint64(461423694561)
	z.ModInverse(&x, m)
	if x.Eq(&z) != 1 {
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
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	m = ModulusFromBytes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 1})
	x.SetUint64(1)
	y.SetUint64(2)
	x = *x.Exp(&x, &y, m)
	z.SetUint64(1)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestSetBytesExamples(t *testing.T) {
	var x, z Nat
	x.SetBytes([]byte{0x12, 0x34, 0x56})
	z.SetUint64(0x123456)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetBytes([]byte{0x00, 0x00, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF})
	z.SetUint64(0xAABBCCDDEEFF)
	if x.Eq(&z) != 1 {
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

func TestByteExample(t *testing.T) {
	x := new(Nat).SetBytes([]byte{8, 7, 6, 5, 4, 3, 2, 1, 0})
	for i := 0; i <= 8; i++ {
		expected := byte(i)
		actual := x.Byte(i)
		if expected != actual {
			t.Errorf("%+v != %+v", expected, actual)
		}
	}
}

func TestModInverseEvenExamples(t *testing.T) {
	var z, x Nat
	x.SetUint64(9)
	m := ModulusFromUint64(10)
	x.ModInverse(&x, m)
	z.SetUint64(9)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(1)
	m = ModulusFromUint64(10)
	x.ModInverse(&x, m)
	z.SetUint64(1)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(19)
	x.ModInverse(&x, m)
	z.SetUint64(9)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(99)
	x.ModInverse(&x, m)
	z.SetUint64(9)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(999)
	m = ModulusFromUint64(1000)
	x.ModInverse(&x, m)
	z.SetUint64(999)
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	// There's an edge case when the modulus is much larger than the input,
	// in which case when we do m^-1 mod x, we need to first calculate the remainder
	// of m.
	x.SetUint64(3)
	m = ModulusFromBytes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0})
	x.ModInverse(&x, m)
	z.SetBytes([]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAB})
	if x.Eq(&z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModSubExamples(t *testing.T) {
	m := ModulusFromUint64(13)
	x := new(Nat).SetUint64(0)
	y := new(Nat).SetUint64(1)
	x.ModSub(x, y, m)
	z := new(Nat).SetUint64(12)
	if x.Eq(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModNegExamples(t *testing.T) {
	m := ModulusFromUint64(13)
	x := new(Nat).SetUint64(0)
	x.ModNeg(x, m)
	z := new(Nat).SetUint64(0)
	if x.Eq(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(1)
	x.ModNeg(x, m)
	z.SetUint64(12)
	if x.Eq(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModSqrtExamples(t *testing.T) {
	m := ModulusFromUint64(13)
	x := new(Nat).SetUint64(4)
	x.ModSqrt(x, m)
	z := new(Nat).SetUint64(11)
	if x.Eq(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestBigExamples(t *testing.T) {
	theBytes := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88}
	x := new(Nat).SetBytes(theBytes)
	expected := new(big.Int).SetBytes(theBytes)
	actual := x.Big()
	if expected.Cmp(actual) != 0 {
		t.Errorf("%+v != %+v", expected, actual)
	}
	expectedNat := x
	actualNat := new(Nat).SetBig(expected, len(theBytes)*8)
	if expectedNat.Eq(actualNat) != 1 {
		t.Errorf("%+v != %+v", expectedNat, actualNat)
	}
}

func TestDivExamples(t *testing.T) {
	x := &Nat{announced: 3 * _W, limbs: []Word{0, 64, 64}}
	n := &Nat{announced: 2 * _W, limbs: []Word{1, 1}}
	nMod := ModulusFromNat(n)

	expectedNat := &Nat{announced: 2 * _W, limbs: []Word{0, 64}}
	actualNat := new(Nat).Div(x, nMod, 2*_W)
	if expectedNat.Eq(actualNat) != 1 {
		t.Errorf("%+v != %+v", expectedNat, actualNat)
	}

	nMod = ModulusFromUint64(1)
	actualNat.Div(x, nMod, x.AnnouncedLen())
	if x.Eq(actualNat) != 1 {
		t.Errorf("%+v != %+v", x, actualNat)
	}
}

func TestCoprimeExamples(t *testing.T) {
	x := new(Nat).SetUint64(5 * 7 * 13)
	y := new(Nat).SetUint64(3 * 7 * 11)
	expected := Choice(0)
	actual := x.Coprime(y)
	if expected != actual {
		t.Errorf("%+v != %+v", expected, actual)
	}
	x.SetUint64(2)
	y.SetUint64(13)
	expected = Choice(1)
	actual = x.Coprime(y)
	if expected != actual {
		t.Errorf("%+v != %+v", expected, actual)
	}
	x.SetUint64(13)
	y.SetUint64(2)
	expected = Choice(1)
	actual = x.Coprime(y)
	if expected != actual {
		t.Errorf("%+v != %+v", expected, actual)
	}
	x.SetUint64(2 * 13 * 11)
	y.SetUint64(2 * 5 * 7)
	expected = Choice(0)
	actual = x.Coprime(y)
	if expected != actual {
		t.Errorf("%+v != %+v", expected, actual)
	}
}

func TestTrueLenExamples(t *testing.T) {
	x := new(Nat).SetUint64(0x0000_0000_0000_0001)
	expected := 1
	actual := x.TrueLen()
	if expected != actual {
		t.Errorf("%+v != %+v", expected, actual)
	}
	x.SetUint64(0x0000_0000_0100_0001)
	expected = 25
	actual = x.TrueLen()
	if expected != actual {
		t.Errorf("%+v != %+v", expected, actual)
	}
}
