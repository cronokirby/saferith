package safenum

import (
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
		i.Neg(i)
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
	minusOne := new(Int).SetUint64(1)
	minusOne.Neg(minusOne)

	way1 := new(Int).Neg(x)
	way2 := new(Int).Mul(x, minusOne, -1)
	return way1.Eq(way2) == 1
}

func TestIntMulNegativeOneIsNeg(t *testing.T) {
	err := quick.Check(testIntMulNegativeOneIsNeg, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}
