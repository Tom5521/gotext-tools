package util_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
)

func TestROverSeq(t *testing.T) {
	seq := util.Seq[int](func(yield func(int) bool) {
		for i := 1; i <= 5; i++ {
			if !yield(i) {
				break
			}
		}
	})

	expected := []int{1, 2, 3, 4, 5}
	idx := 0

	for v := range util.ROverSeq(seq) {
		if v != expected[idx] {
			t.Errorf("Expected %d, but got %d", expected[idx], v)
		}
		idx++
	}
}

func TestROverSeq2(t *testing.T) {
	seq2 := util.Seq2[string, int](func(yield func(string, int) bool) {
		values := map[string]int{"a": 1, "b": 2, "c": 3}
		for k, v := range values {
			if !yield(k, v) {
				break
			}
		}
	})

	expected := map[int]bool{1: true, 2: true, 3: true}

	for f := range util.ROverSeq2(seq2) {
		if !expected[f.V] {
			t.Errorf("Unexpected value: %d", f.V)
		}
		delete(expected, f.V)
	}

	if len(expected) > 0 {
		t.Errorf("Missing values in output: %v", expected)
	}
}
