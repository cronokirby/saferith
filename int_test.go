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
		i.Neg()
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
