package safenum

import "testing"

func TestUint64Creation(t *testing.T) {
	var x, y Nat
	x.SetUint64(0)
	y.SetUint64(0)
	if x.Cmp(y) != 0 {
		t.Errorf("%+v != %+v", x, y)
	}
	x.SetUint64(1)
	if x.Cmp(y) == 0 {
		t.Errorf("%+v == %+v", x, y)
	}
	x.SetUint64(0x1111)
	y.SetUint64(0x1111)
	if x.Cmp(y) != 0 {
		t.Errorf("%+v != %+v", x, y)
	}
}

func TestAddExamples(t *testing.T) {
	var x, y, z Nat
	x.SetUint64(100)
	y.SetUint64(100)
	z.SetUint64(200)
	x = *x.Add(x, y, 8)
	if x.Cmp(z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
	z.SetUint64(300 - 256)
	x = *x.Add(x, y, 8)
	if x.Cmp(z) != 0 {
		t.Errorf("%+v != %+v", x, z)
	}
}
