//go:build ignore
// +build ignore

package util_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
)

func BenchmarkPJWHash(b *testing.B) {
	text := "abcdefg"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		util.PJWHash(text)
	}
}
