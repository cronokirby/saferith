// +build gofuzz

package saferith

import (
	"encoding/binary"
	"errors"
	"math/big"
)

func FuzzInt(data []byte) int {
	FuzzIntSetBig(data)
	FuzzIntSetBytes(data)
	FuzzIntSetNat(data)
	FuzzIntSetInt(data)
	FuzzIntSetUint64(data)
	FuzzIntUnmarshalBinary(data)

	FuzzIntArithmetic(data)
	FuzzIntResize(data)
	return 0
}

func FuzzIntSetBig(data []byte) int {
	var x big.Int
	x.SetBytes(data)

	l := len(data)
	for size := 0; size < l; size++ {
		var i Int
		i.SetBig(&x, size)
		runAllInt(&i)
	}

	return 0
}

func FuzzIntSetBytes(data []byte) int {
	var i Int
	i.SetBytes(data)
	runAllInt(&i)
	return 0
}

func FuzzIntSetNat(data []byte) int {
	var n Nat
	n.SetBytes(data)

	var i Int
	i.SetNat(&n)
	runAllInt(&i)
	return 0
}

func FuzzIntSetInt(data []byte) int {
	var x Int
	x.SetBytes(data)

	var i Int
	i.SetInt(&x)
	runAllInt(&i)
	return 0
}

func FuzzIntSetUint64(data []byte) int {
	l := len(data)
	if l != 8 {
		return -1
	}

	x := binary.LittleEndian.Uint64(data)

	var i Int
	i.SetUint64(x)
	runAllInt(&i)
	return 0
}

func FuzzIntUnmarshalBinary(data []byte) int {
	if len(data) < 1 {
		return -1
	}

	var i Int
	if err := i.UnmarshalBinary(data); err != nil {
		panic(err)
	}

	runAllInt(&i)
	return 0
}

func FuzzIntArithmetic(data []byte) int {
	FuzzIntModularArithmetic(data)
	FuzzIntSetModSymmetric(data)
	FuzzIntExpI(data)
	FuzzIntAdd(data)
	FuzzIntMul(data)
	FuzzIntNeg(data)
	return 0
}

func FuzzIntModularArithmetic(data []byte) int {
	if isZero(data) {
		return -1
	}

	p := ModulusFromBytes(data)

	var i Int
	i.CheckInRange(p)
	i.Mod(p)
	return 0
}

func FuzzIntSetModSymmetric(data []byte) int {
	z, p, err := getOneNatAndOneMod(data)
	if err != nil {
		return -1
	}

	var i Int
	i.SetModSymmetric(z, p)
	return 0
}

func FuzzIntExpI(data []byte) int {
	i, x, p, err := getOneIntAndOneNatAndOneMod(data)
	if err != nil {
		return -1
	}

	var z Nat
	z.ExpI(x, i, p)
	return 0
}

func FuzzIntAdd(data []byte) int {
	if len(data) < 2 {
		return -1
	}

	cap := int(data[0])
	x, y, err := getTwoInts(data[1:])
	if err != nil {
		return -1
	}

	var a Int
	var b Int
	a.Add(x, y, cap)
	b.Add(y, x, cap)
	if a.Eq(&b) != 1 {
		panic("Int.Add: (x+y)!=(y+x)")
	}

	return 0
}

func FuzzIntMul(data []byte) int {
	if len(data) < 2 {
		return -1
	}

	cap := int(data[0])
	x, y, err := getTwoInts(data[1:])
	if err != nil {
		return -1
	}

	var a Int
	var b Int
	a.Mul(x, y, cap)
	b.Mul(y, x, cap)
	if a.Eq(&b) != 1 {
		panic("Int.Mul: (x*y)!=(y*x)")
	}

	return 0
}

func FuzzIntNeg(data []byte) int {
	var i Int
	i.SetBytes(data)

	yes := Choice(1)
	i.Neg(yes)

	no := Choice(0)
	i.Neg(no)
	return 0
}

func FuzzIntResize(data []byte) int {
	if len(data) < 2 {
		return -1
	}

	cap := int(data[0])

	var i Int
	i.SetBytes(data[1:])
	i.Resize(cap)
	return 0
}

// Check all methods of an Int that require no Int or Modulus as input
func runAllInt(i *Int) {
	i.Abs()
	i.AnnouncedLen()
	i.Big()
	i.Clone()
	i.IsNegative()
	i.TrueLen()
	i.String()

	if _, err := i.MarshalBinary(); err != nil {
		panic(err)
	}

	if i.Eq(i) != 1 {
		panic("Int.Eq: i!=i")
	}
}

// Convert a byte array into two Ints and one Modulus
func getTwoInts(data []byte) (*Int, *Int, error) {
	l := len(data)
	if l < 3 {
		return nil, nil, errors.New("too few bytes")
	}

	chunk := int(l / 3)
	a := 0 + chunk
	b := a + chunk

	var x Int
	var y Int
	x.SetBytes(data[0 : a-1])
	y.SetBytes(data[a : b-1])
	return &x, &y, nil
}

// Convert a byte array into one Nat, one Int, and one Modulus
func getOneIntAndOneNatAndOneMod(data []byte) (*Int, *Nat, *Modulus, error) {
	l := len(data)
	if l < 3 {
		return nil, nil, nil, errors.New("too few bytes")
	}

	chunk := int(l / 3)
	a := 0 + chunk
	b := a + chunk
	c := b + chunk

	var i Int
	i.SetBytes(data[0 : a-1])

	var z Nat
	z.SetBytes(data[a : b-1])

	pBytes := data[b : c-1]
	if isZero(pBytes) {
		return nil, nil, nil, errors.New("modulus cannot be zero")
	}
	p := ModulusFromBytes(pBytes)

	return &i, &z, p, nil
}
