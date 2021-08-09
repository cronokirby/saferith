package safenum

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func (*Int) Generate(r *rand.Rand, size int) reflect.Value {
	bytes := make([]byte, 16*_S)
	r.Read(bytes)
	i := new(Int).SetBytes(bytes)
	if r.Int()&1 == 1 {
		i.Neg(1)
	}
	return reflect.ValueOf(i)
}

func testIntEqualReflexive(z *Int) bool {
	return z.Eq(z) == 1
}

func TestIntEqualReflexive(t *testing.T) {
	err := quick.Check(testIntEqualReflexive, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testIntMulCommutative(x, y *Int) bool {
	way1 := new(Int).Mul(x, y, -1)
	way2 := new(Int).Mul(y, x, -1)
	return way1.Eq(way2) == 1
}

func TestIntMulCommutative(t *testing.T) {
	err := quick.Check(testIntMulCommutative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testIntMulZeroIsZero(x *Int) bool {
	zero := new(Int)
	timesZero := new(Int).Mul(zero, x, -1)
	return timesZero.Eq(zero) == 1
}

func TestIntMulZeroIsZero(t *testing.T) {
	err := quick.Check(testIntMulZeroIsZero, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testIntMulNegativeOneIsNeg(x *Int) bool {
	minusOne := new(Int).SetUint64(1).Neg(1)

	way1 := new(Int).SetInt(x).Neg(1)
	way2 := new(Int).Mul(x, minusOne, -1)
	return way1.Eq(way2) == 1
}

func TestIntMulNegativeOneIsNeg(t *testing.T) {
	err := quick.Check(testIntMulNegativeOneIsNeg, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testIntModAddNegReturnsZero(x *Int, m Modulus) bool {
	a := new(Int).SetInt(x).Neg(1).Mod(&m)
	b := x.Mod(&m)
	return b.ModAdd(a, b, &m).EqZero() == 1
}

func TestIntModAddNegReturnsZero(t *testing.T) {
	err := quick.Check(testIntModAddNegReturnsZero, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testIntModRoundtrip(x Nat, m Modulus) bool {
	xModM := new(Nat).Mod(&x, &m)
	i := new(Int).SetModSymmetric(xModM, &m)
	if i.CheckInRange(&m) != 1 {
		return false
	}
	roundTrip := i.Mod(&m)
	return xModM.Eq(roundTrip) == 1
}

func TestIntModRoundtrip(t *testing.T) {
	err := quick.Check(testIntModRoundtrip, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testIntAddNegZero(i *Int) bool {
	zero := new(Int)
	neg := new(Int).SetInt(i).Neg(1)
	shouldBeZero := new(Int).Add(i, neg, -1)
	return shouldBeZero.Eq(zero) == 1
}

func TestIntAddNegZero(t *testing.T) {
	err := quick.Check(testIntAddNegZero, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testIntAddCommutative(x *Int, y *Int) bool {
	way1 := new(Int).Add(x, y, -1)
	way2 := new(Int).Add(x, y, -1)
	return way1.Eq(way2) == 1
}

func TestIntAddCommutative(t *testing.T) {
	err := quick.Check(testIntAddCommutative, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testIntAddZeroIdentity(x *Int) bool {
	zero := new(Int)
	shouldBeX := new(Int).Add(x, zero, -1)
	return shouldBeX.Eq(x) == 1
}

func TestIntAddZeroIdentity(t *testing.T) {
	err := quick.Check(testIntAddZeroIdentity, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestCheckInRangeExamples(t *testing.T) {
	x := new(Int).SetUint64(0)
	m := ModulusFromUint64(13)
	if x.CheckInRange(m) != 1 {
		t.Errorf("expected zero to be in range of modulus")
	}
}

func TestIntAddExamples(t *testing.T) {
	x := new(Int).SetUint64(3).Resize(8)
	y := new(Int).SetUint64(4).Neg(1).Resize(8)
	expected := new(Int).SetUint64(1).Neg(1)
	actual := new(Int).Add(x, y, -1)
	if expected.Eq(actual) != 1 {
		t.Errorf("%+v != %+v", expected, actual)
	}
}

func testIntMarshalBinaryRoundTrip(x *Int) bool {
	out, err := x.MarshalBinary()
	if err != nil {
		return false
	}
	y := new(Int)
	err = y.UnmarshalBinary(out)
	if err != nil {
		return false
	}
	return x.Eq(y) == 1
}

func TestIntMarshalBinaryRoundTrip(t *testing.T) {
	err := quick.Check(testIntMarshalBinaryRoundTrip, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func testInvalidInt(expected []byte) bool {
	x := new(Int)
	err := x.UnmarshalBinary(expected)
	// empty slice is invalid, so we expect an error
	if len(expected) == 0 {
		return err != nil
	}
	expectedBytes := expected[1:]
	expectedSign := Choice(expected[0]) & 1
	actualBytes := x.Abs().Bytes()
	actualSign := x.sign
	return (expectedSign == actualSign) && bytes.Equal(expectedBytes, actualBytes)
}

func TestInvalidInt(t *testing.T) {
	err := quick.Check(testInvalidInt, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}
