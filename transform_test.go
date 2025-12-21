package kk

import (
	"strconv"
	"strings"
	"testing"
)

func TestMapped(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	q := Mapped(From(input), func(n int) int { return n * 2 })
	result := Slice(q)

	expected := []int{2, 4, 6, 8, 10}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestMapTypeChange(t *testing.T) {
	input := []int{1, 2, 3}
	q := Mapped(From(input), func(n int) string { return strconv.Itoa(n) })
	result := Slice(q)

	expected := []string{"1", "2", "3"}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

func TestMapEmpty(t *testing.T) {
	input := []int{}
	q := Mapped(From(input), func(n int) int { return n * 2 })
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestMapChaining(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	q := Mapped(From(input).Where(func(n int) bool { return n%2 == 0 }), func(n int) int { return n * 10 })
	result := Slice(q)

	expected := []int{20, 40}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestMapToStruct(t *testing.T) {
	type Result struct {
		Value   int
		Doubled int
	}

	input := []int{1, 2, 3}
	q := Mapped(From(input), func(n int) Result {
		return Result{Value: n, Doubled: n * 2}
	})
	result := Slice(q)

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}

	if result[0].Value != 1 || result[0].Doubled != 2 {
		t.Errorf("unexpected first result: %v", result[0])
	}
}

func TestFlattened(t *testing.T) {
	input := []int{1, 2, 3}
	q := Flattened(From(input), func(n int) []int {
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = n
		}
		return result
	})
	result := Slice(q)

	// 1 -> [1], 2 -> [2, 2], 3 -> [3, 3, 3]
	expected := []int{1, 2, 2, 3, 3, 3}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestFlatMapStrings(t *testing.T) {
	input := []string{"hello world", "foo bar baz"}
	q := Flattened(From(input), func(s string) []string {
		return strings.Split(s, " ")
	})
	result := Slice(q)

	expected := []string{"hello", "world", "foo", "bar", "baz"}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

func TestFlatMapEmpty(t *testing.T) {
	input := []int{}
	q := Flattened(From(input), func(n int) []int { return []int{n, n} })
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestFlatMapWithEmptyResults(t *testing.T) {
	input := []int{1, 2, 3}
	q := Flattened(From(input), func(n int) []int {
		if n == 2 {
			return []int{} // Empty slice for 2
		}
		return []int{n}
	})
	result := Slice(q)

	expected := []int{1, 3}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestFlatMapTypeChange(t *testing.T) {
	input := []int{1, 2}
	q := Flattened(From(input), func(n int) []string {
		result := make([]string, n)
		for i := 0; i < n; i++ {
			result[i] = strconv.Itoa(n)
		}
		return result
	})
	result := Slice(q)

	expected := []string{"1", "2", "2"}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

func TestNestedMapped(t *testing.T) {
	input := []int{1, 2, 3}
	q := Mapped(Mapped(From(input), func(n int) int { return n * 2 }), func(n int) string { return strconv.Itoa(n) })
	result := Slice(q)

	expected := []string{"2", "4", "6"}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}
