package saferith

import (
	"math/rand"
	"testing"
)

var result Word

func BenchmarkLimbMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := rand.Int()
		result = limbMask(x)
	}
}
