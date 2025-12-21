package kk

import (
	"testing"
)

func TestCount(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result := Count(From(input))

	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func TestCountEmpty(t *testing.T) {
	input := []int{}
	result := Count(From(input))

	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

func TestCountWithFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := Count(From(input).Where(func(n int) bool { return n%2 == 0 }))

	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func TestSumInt(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result := Sum(From(input), func(n int) int { return n })

	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

func TestSumFloat(t *testing.T) {
	input := []float64{1.5, 2.5, 3.0}
	result := Sum(From(input), func(n float64) float64 { return n })

	if result != 7.0 {
		t.Errorf("expected 7.0, got %f", result)
	}
}

func TestSumEmpty(t *testing.T) {
	input := []int{}
	result := Sum(From(input), func(n int) int { return n })

	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

func TestSumWithSelector(t *testing.T) {
	type Item struct {
		Value int
	}
	input := []Item{{Value: 10}, {Value: 20}, {Value: 30}}
	result := Sum(From(input), func(item Item) int { return item.Value })

	if result != 60 {
		t.Errorf("expected 60, got %d", result)
	}
}

func TestFirst(t *testing.T) {
	input := []int{5, 10, 15}
	result, ok := First(From(input))

	if !ok {
		t.Error("expected ok to be true")
	}
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func TestFirstEmpty(t *testing.T) {
	input := []int{}
	result, ok := First(From(input))

	if ok {
		t.Error("expected ok to be false")
	}
	if result != 0 {
		t.Errorf("expected zero value, got %d", result)
	}
}

func TestFirstWithFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result, ok := First(From(input).Where(func(n int) bool { return n > 3 }))

	if !ok {
		t.Error("expected ok to be true")
	}
	if result != 4 {
		t.Errorf("expected 4, got %d", result)
	}
}

func TestAny(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result := Any(From(input), func(n int) bool { return n > 3 })

	if !result {
		t.Error("expected true")
	}
}

func TestAnyNoneMatch(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result := Any(From(input), func(n int) bool { return n > 10 })

	if result {
		t.Error("expected false")
	}
}

func TestAnyEmpty(t *testing.T) {
	input := []int{}
	result := Any(From(input), func(n int) bool { return true })

	if result {
		t.Error("expected false for empty query")
	}
}

func TestAll(t *testing.T) {
	input := []int{2, 4, 6, 8, 10}
	result := All(From(input), func(n int) bool { return n%2 == 0 })

	if !result {
		t.Error("expected true")
	}
}

func TestAllOneFails(t *testing.T) {
	input := []int{2, 4, 5, 8, 10}
	result := All(From(input), func(n int) bool { return n%2 == 0 })

	if result {
		t.Error("expected false")
	}
}

func TestAllEmpty(t *testing.T) {
	input := []int{}
	result := All(From(input), func(n int) bool { return false })

	if !result {
		t.Error("expected true for empty query (vacuous truth)")
	}
}

func TestAllNoneMatch(t *testing.T) {
	input := []int{1, 3, 5, 7, 9}
	result := All(From(input), func(n int) bool { return n%2 == 0 })

	if result {
		t.Error("expected false")
	}
}
