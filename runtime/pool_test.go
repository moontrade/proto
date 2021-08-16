package runtime

import (
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {
	b := make([]byte, 6999)
	b = GetBytes(71)
	PutBytes(b)
	b = GetBytes(72)
	fmt.Println(cap(b))
	fmt.Println(cap(make([]byte, 191)))
}

func BenchmarkPool(b *testing.B) {
	b.Run("64", benchGet(64))
	b.Run("256", benchGet(256))
	b.Run("384", benchGet(384))
}

func benchGet(n int) func(b *testing.B) {
	return func(b *testing.B) {
		PutBytes(GetBytes(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			PutBytes(GetBytes(n))
		}
	}
}
