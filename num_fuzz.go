// +build gofuzz

package saferith

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"regexp"
)

var zeroHexRegExp = regexp.MustCompile("0(x0+)?")

func FuzzNat(data []byte) int {
	FuzzNatSetBig(data)
	FuzzNatSetBytes(data)
	FuzzNatSetHex(data)
	FuzzNatSetNat(data)
	FuzzNatSetUint64(data)
	FuzzNatUnmarshalBinary(data)

	FuzzNatArithmetic(data)
	FuzzNatBitShifting(data)
	FuzzNatCompare(data)
	FuzzNatResize(data)
	return 0
}

func FuzzModulus(data []byte) int {
	FuzzModulusFromBytes(data)
	FuzzModulusFromHex(data)
	FuzzModulusFromNat(data)
	FuzzModulusFromUint64(data)
	FuzzModulusUnmarshalBinary(data)

	FuzzModulesCompare(data)
	return 0
}

func FuzzNatSetBig(data []byte) int {
	var x big.Int
	x.SetBytes(data)

	l := len(data)
	for i := 0; i < l; i++ {
		var z Nat
		z.SetBig(&x, i)
		runNatFuncs(&z, l)
	}

	return 0
}

func FuzzNatSetBytes(data []byte) int {
	var z Nat
	z.SetBytes(data)
	runNatFuncs(&z, len(data))
	return 0
}

func FuzzNatSetHex(data []byte) int {
	x := string(data)

	var z Nat
	z.SetHex(x)
	runNatFuncs(&z, len(data))
	return 0
}

func FuzzNatSetNat(data []byte) int {
	var x Nat
	x.SetBytes(data)

	var z Nat
	z.SetNat(&x)
	runNatFuncs(&z, len(data))
	return 0
}

func FuzzNatSetUint64(data []byte) int {
	l := len(data)
	if l != 8 {
		return -1
	}

	x := binary.LittleEndian.Uint64(data)

	var z Nat
	z.SetUint64(x)
	runNatFuncs(&z, l)
	return 0
}

func FuzzNatUnmarshalBinary(data []byte) int {
	var z Nat
	if err := z.UnmarshalBinary(data); err != nil {
		panic(err)
	}

	runNatFuncs(&z, len(data))
	return 0
}

func FuzzNatArithmetic(data []byte) int {
	FuzzNatUnaryArithmetic(data)
	FuzzNatAdd(data)
	FuzzNatSub(data)
	FuzzNatMul(data)
	FuzzNatDiv(data)
	FuzzNatExp(data)
	FuzzNatSqrt(data)
	return 0
}

func FuzzNatUnaryArithmetic(data []byte) int {
	x, p, err := getOneNatAndOneMod(data)
	if err != nil {
		return -1
	}

	var z Nat
	z.Mod(x, p)
	z.ModNeg(x, p)
	z.ModInverse(x, p)
	return 0
}

func FuzzNatAdd(data []byte) int {
	if len(data) < 2 {
		return -1
	}

	cap := int(data[0])
	x, y, p, err := getTwoNatsAndOneMod(data[1:])
	if err != nil {
		return -1
	}

	var a Nat
	var b Nat
	a.ModAdd(x, y, p)
	b.ModAdd(y, x, p)
	if a.Eq(&b) != 1 {
		panic("Nat.ModAdd: (x+y)!=(y+x)")
	}

	a.Add(x, y, cap)
	b.Add(y, x, cap)
	if a.Eq(&b) != 1 {
		panic("Nat.Add: (x+y)!=(y+x)")
	}

	return 0
}

func FuzzNatSub(data []byte) int {
	if len(data) < 2 {
		return -1
	}

	cap := int(data[0])
	x, y, p, err := getTwoNatsAndOneMod(data[1:])
	if err != nil {
		return -1
	}

	var z Nat
	z.ModSub(x, y, p)
	z.ModSub(y, x, p)
	z.Sub(x, y, cap)
	z.Sub(y, x, cap)
	return 0
}

func FuzzNatMul(data []byte) int {
	if len(data) < 2 {
		return -1
	}

	cap := int(data[0])
	x, y, p, err := getTwoNatsAndOneMod(data[1:])
	if err != nil {
		return -1
	}

	var a Nat
	var b Nat
	a.ModMul(x, y, p)
	b.ModMul(y, x, p)
	if a.Eq(&b) != 1 {
		panic("Nat.ModMul: (x*y)!=(y*x)")
	}

	a.Mul(x, y, cap)
	b.Mul(y, x, cap)
	if a.Eq(&b) != 1 {
		panic("Nat.Mul: (x*y)!=(y*x)")
	}

	return 0
}

func FuzzNatDiv(data []byte) int {
	if len(data) < 2 {
		return -1
	}

	cap := int(data[0])
	x, p, err := getOneNatAndOneMod(data[1:])
	if err != nil {
		return -1
	}

	var z Nat
	z.Div(x, p, cap)
	return 0
}

func FuzzNatExp(data []byte) int {
	x, y, p, err := getTwoNatsAndOneMod(data)
	if err != nil {
		return -1
	}

	var z Nat
	z.Exp(x, y, p)
	return 0
}

func FuzzNatSqrt(data []byte) int {
	x, p, err := getOneNatAndOneMod(data)
	if err != nil {
		return -1
	}

	if p.nat.limbs[0]&1 == 0 {
		return -1
	}

	var z Nat
	z.ModSqrt(x, p)
	return 0
}

func FuzzNatBitShifting(data []byte) int {
	if len(data) < 3 {
		return -1
	}

	shift := uint(data[0])
	cap := int(data[1])

	var x Nat
	x.SetBytes(data[2:])

	var z Nat
	z.Rsh(&x, shift, cap)
	z.Lsh(&x, shift, cap)
	return 0
}

func FuzzNatIsUnit(data []byte) int {
	z, p, err := getOneNatAndOneMod(data)
	if err != nil {
		return -1
	}

	z.IsUnit(p)
	return 0
}

func FuzzNatCmp(data []byte) int {
	x, y, _, err := getTwoNatsAndOneMod(data)
	if err != nil {
		return -1
	}

	gt1, eq1, lt1 := x.Cmp(y)
	gt2, eq2, lt2 := y.Cmp(x)
	if eq1 != eq2 {
		panic("Nat.Cmp: (x==y)!=(y==x)")
	}

	if eq1 == 0 {
		if gt1 == gt2 || lt1 == lt2 {
			panic("Nat.Cmp: (x!=y), but !(x>y) or !(x<y)")
		}
	} else {
		if gt1 != gt2 || lt1 != lt2 {
			panic("Nat.Cmp: (x==y), but (x>y) or (x<y)")
		}
	}

	return 0
}

func FuzzNatCmpMod(data []byte) int {
	z, p, err := getOneNatAndOneMod(data)
	if err != nil {
		return -1
	}

	z.CmpMod(p)
	return 0
}

func FuzzNatEq(data []byte) int {
	x, y, _, err := getTwoNatsAndOneMod(data)
	if err != nil {
		return -1
	}

	eq1 := x.Eq(y)
	eq2 := y.Eq(x)
	if eq1 != eq2 {
		panic("Nat.Eq: (x==y)!=(y==x)")
	}

	return 0
}

func FuzzNatCoprime(data []byte) int {
	x, y, _, err := getTwoNatsAndOneMod(data)
	if err != nil {
		return -1
	}

	coprime1 := x.Coprime(y)
	coprime2 := y.Coprime(x)
	if coprime1 != coprime2 {
		panic("Nat.Coprime: (x coprime y)!=(y coprime x)")
	}

	return 0
}

func FuzzNatCompare(data []byte) int {
	FuzzNatIsUnit(data)
	FuzzNatCmp(data)
	FuzzNatCmpMod(data)
	FuzzNatEq(data)
	FuzzNatCoprime(data)
	return 0
}

func FuzzNatResize(data []byte) int {
	if len(data) < 2 {
		return -1
	}

	cap := int(data[0])

	var z Nat
	z.SetBytes(data[1:])
	z.Resize(cap)
	return 0
}

func FuzzModulusFromBytes(data []byte) int {
	if isZero(data) {
		return -1
	}

	p := ModulusFromBytes(data)
	runModulusFuncs(p)
	return 0
}

func FuzzModulusFromHex(data []byte) int {
	if isZero(data) {
		return -1
	}

	x := string(data)
	if p, err := ModulusFromHex(x); err == nil {
		runModulusFuncs(p)
	}

	return 0
}

func FuzzModulusFromNat(data []byte) int {
	if isZero(data) {
		return -1
	}

	var z Nat
	z.SetBytes(data)
	p := ModulusFromNat(&z)
	runModulusFuncs(p)
	return 0
}

func FuzzModulusFromUint64(data []byte) int {
	if len(data) != 8 || isZero(data) {
		return -1
	}

	x := binary.LittleEndian.Uint64(data)
	p := ModulusFromUint64(x)
	runModulusFuncs(p)
	return 0
}

func FuzzModulusUnmarshalBinary(data []byte) int {
	if isZero(data) {
		return -1
	}

	var p Modulus
	if err := p.UnmarshalBinary(data); err != nil {
		panic(err)
	}

	runModulusFuncs(&p)
	return 0
}

func FuzzModulesCompare(data []byte) int {
	l := len(data)
	if l < 2 {
		return -1
	}

	chunkSize := int(l / 2)
	p := ModulusFromBytes(data[0 : chunkSize-1])
	q := ModulusFromBytes(data[chunkSize:])

	gt1, eq1, lt1 := p.Cmp(q)
	gt2, eq2, lt2 := q.Cmp(p)
	if eq1 != eq2 {
		panic("Modulus.Cmp: (p==q)!=(q==p)")
	}

	if eq1 == 0 {
		if gt1 == gt2 || lt1 == lt2 {
			panic("Modules.Cmp: (p!=q), but !(p>q) or !(p<q)")
		}
	} else {
		if gt1 != gt2 || lt1 != lt2 {
			panic("Modules.Cmp: (p==q), but (p>q) or (p<q)")
		}
	}

	return 0
}

// Run methods of a Nat that require no Nat or Modulus as input
func runNatFuncs(z *Nat, l int) {
	z.AnnouncedLen()
	z.Big()
	z.Bytes()
	z.Clone()
	z.EqZero()
	z.Hex()
	z.String()
	z.TrueLen()
	z.Uint64()

	if _, err := z.MarshalBinary(); err != nil {
		panic(err)
	}

	if z.Eq(z) != 1 {
		panic("Nat.Eq: z!=z")
	}

	gt, eq, lt := z.Cmp(z)
	if gt != 0 || eq != 1 || lt != 0 {
		panic(fmt.Sprintf("Nat.Cmp: z!=z, gt=%b,eq=%b,lt=%b", gt, eq, lt))
	}

	for i := 0; i < l; i++ {
		z.Byte(i)
	}

	buf := make([]byte, l)
	z.FillBytes(buf)
}

// Run methods of a Modulus that require no Nat or Modulus as input
func runModulusFuncs(p *Modulus) {
	p.Big()
	p.BitLen()
	p.Bytes()
	p.Hex()
	p.Nat()
	p.String()

	if _, err := p.MarshalBinary(); err != nil {
		panic(err)
	}

	gt, eq, lt := p.Cmp(p)
	if gt != 0 || eq != 1 || lt != 0 {
		panic(fmt.Sprintf("Modulus.Cmp: p!=p, gt=%b,eq=%b,lt=%b", gt, eq, lt))
	}
}

// Convert a byte array into two Nats and one Modulus
func getTwoNatsAndOneMod(data []byte) (*Nat, *Nat, *Modulus, error) {
	l := len(data)
	if l < 3 {
		return nil, nil, nil, errors.New("too few bytes")
	}

	chunkSize := int(l / 3)
	a := 0 + chunkSize
	b := a + chunkSize
	c := b + chunkSize

	var x Nat
	var y Nat
	x.SetBytes(data[0 : a-1])
	y.SetBytes(data[a : b-1])

	pBytes := data[b : c-1]
	if isZero(pBytes) {
		return nil, nil, nil, errors.New("modulus cannot be zero")
	}
	p := ModulusFromBytes(pBytes)

	return &x, &y, p, nil
}

// Convert a byte array into one Nat and one Modulus
func getOneNatAndOneMod(data []byte) (*Nat, *Modulus, error) {
	l := len(data)
	if l < 2 {
		return nil, nil, errors.New("too few bytes")
	}

	chunkSize := int(l / 2)
	a := 0 + chunkSize
	b := a + chunkSize

	var z Nat
	z.SetBytes(data[0 : a-1])

	pBytes := data[a : b-1]
	if isZero(pBytes) {
		return nil, nil, errors.New("modulus cannot be zero")
	}
	p := ModulusFromBytes(pBytes)

	return &z, p, nil
}

// Check if a byte array is all zeros
func isZero(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	s := string(data)
	if zeroHexRegExp.MatchString(s) {
		return true
	}

	// if data != 0x00... return false
	for _, b := range data {
		if b != 0x0 {
			return false
		}
	}

	return true
}
