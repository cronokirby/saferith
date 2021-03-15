package safenum

import (
	"math/big"
	"testing"
)

var resultInt big.Int

func BenchmarkAddBig(b *testing.B) {
	bytes := make([]byte, 100)
	for i := 0; i < 100; i++ {
		bytes[i] = 1
	}
	var x big.Int
	x.SetBytes(bytes)
	for n := 0; n < b.N; n++ {
		var z big.Int
		z.Add(&x, &x)
		z.Add(&x, &x)
		z.Add(&x, &x)
		resultInt = z
	}
}
