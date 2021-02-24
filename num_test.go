package safenum

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func (Nat) Generate(r *rand.Rand, size int) reflect.Value {
	var n Nat
	n.SetUint64(r.Uint64())
	return reflect.ValueOf(n)
}

func testAddZeroIdentity(n Nat) bool {
	var x, zero Nat
	zero.SetUint64(0)
	x.Add(n, zero, 128)
	if n.Cmp(x) != 1 {
		return false
	}
	x.Add(zero, n, 128)
	if n.Cmp(x) != 1 {
		return false
	}
	return true
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
		aPlusB.Add(a, b, x)
		bPlusA.Add(b, a, x)
		if aPlusB.Cmp(bPlusA) != 1 {
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
		order1 = *order1.Add(a, b, x)
		order1.Add(order1, c, x)
		order2 = *order2.Add(b, c, x)
		order2.Add(a, order2, x)
		if order1.Cmp(order2) != 1 {
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
		aTimesB.Mul(a, b, x)
		bTimesA.Mul(b, a, x)
		if aTimesB.Cmp(bTimesA) != 1 {
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
		order1 = *order1.Mul(a, b, x)
		order1.Mul(order1, c, x)
		order2 = *order2.Mul(b, c, x)
		order2.Mul(a, order2, x)
		if order1.Cmp(order2) != 1 {
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
	x.Mul(n, one, 128)
	if n.Cmp(x) != 1 {
		return false
	}
	x.Mul(one, n, 128)
	if n.Cmp(x) != 1 {
		return false
	}
	return true
}

func TestMulOneIdentity(t *testing.T) {
	err := quick.Check(testMulOneIdentity, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testModAndTruncationMatch(a Nat) bool {
	var way1, way2, m, zero Nat
	zero.SetUint64(0)
	for _, x := range []uint{48, 32, 16, 8} {
		way1.Add(a, zero, x)
		m.SetUint64(1 << x)
		way2.Mod(a, m)
		if way1.Cmp(way2) != 1 {
			return false
		}
	}
	return true
}

func TestModAndTruncationMatch(t *testing.T) {
	err := quick.Check(testModAndTruncationMatch, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestUint64Creation(t *testing.T) {
	var x, y Nat
	x.SetUint64(0)
	y.SetUint64(0)
	if x.Cmp(y) != 1 {
		t.Errorf("%+v != %+v", x, y)
	}
	x.SetUint64(1)
	if x.Cmp(y) == 1 {
		t.Errorf("%+v == %+v", x, y)
	}
	x.SetUint64(0x1111)
	y.SetUint64(0x1111)
	if x.Cmp(y) != 1 {
		t.Errorf("%+v != %+v", x, y)
	}
}

func TestAddExamples(t *testing.T) {
	var x, y, z Nat
	x.SetUint64(100)
	y.SetUint64(100)
	z.SetUint64(200)
	x = *x.Add(x, y, 8)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	z.SetUint64(300 - 256)
	x = *x.Add(x, y, 8)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetUint64(0xf3e5487232169930)
	y.SetUint64(0)
	z.SetUint64(0xf3e5487232169930)
	var x2 Nat
	x2.Add(x, y, 128)
	if x2.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestMulExamples(t *testing.T) {
	var x, y, z Nat
	x.SetUint64(10)
	y.SetUint64(10)
	z.SetUint64(100)
	x = *x.Mul(x, y, 8)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	z.SetUint64(232)
	x = *x.Mul(x, y, 8)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModAddExamples(t *testing.T) {
	var x, y, m, z Nat
	m.SetUint64(13)
	x.SetUint64(40)
	y.SetUint64(40)
	x = *x.ModAdd(x, y, m)
	z.SetUint64(2)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModMulExamples(t *testing.T) {
	var x, y, m, z Nat
	m.SetUint64(13)
	x.SetUint64(40)
	y.SetUint64(40)
	x = *x.ModMul(x, y, m)
	z.SetUint64(1)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModExamples(t *testing.T) {
	var x, z, m Nat
	x.SetUint64(40)
	m.SetUint64(13)
	x = *x.Mod(x, m)
	z.SetUint64(1)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestModInverseExamples(t *testing.T) {
	var x, z, m Nat
	x.SetUint64(2)
	m.SetUint64(13)
	x = *x.ModInverse(x, m)
	z.SetUint64(7)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestExpExamples(t *testing.T) {
	var x, y, z, m Nat
	x.SetUint64(3)
	y.SetUint64(345)
	m.SetUint64(13)
	x = *x.Exp(x, y, m)
	z.SetUint64(1)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestSetBytesExamples(t *testing.T) {
	var x, z Nat
	x.SetBytes([]byte{0x12, 0x34, 0x56})
	z.SetUint64(0x123456)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
	x.SetBytes([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF})
	z.SetUint64(0xAABBCCDDEEFF)
	if x.Cmp(z) != 1 {
		t.Errorf("%+v != %+v", x, z)
	}
}

func TestFillBytesExamples(t *testing.T) {
	var x Nat
	expected := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	x.SetBytes(expected)
	buf := make([]byte, 4)
	x.FillBytes(buf)
	if !bytes.Equal(expected, buf) {
		t.Errorf("%+v != %+v", expected, buf)
	}
}
