package kk

import (
	"testing"
)

func TestWhere(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q := Query(input).Where(func(n int) bool { return n%2 == 0 })
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

func TestWhereNoneMatch(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	q := Query(input).Where(func(n int) bool { return n > 10 })
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestWhereAllMatch(t *testing.T) {
	input := []int{2, 4, 6, 8, 10}
	q := Query(input).Where(func(n int) bool { return n%2 == 0 })
	result := Slice(q)

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}
}

func TestTake(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q := Query(input).Take(5)
	result := Slice(q)

	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestTakeMoreThanAvailable(t *testing.T) {
	input := []int{1, 2, 3}
	q := Query(input).Take(10)
	result := Slice(q)

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}
}

func TestTakeZero(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	q := Query(input).Take(0)
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestSkip(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q := Query(input).Skip(5)
	result := Slice(q)

	expected := []int{6, 7, 8, 9, 10}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestSkipMoreThanAvailable(t *testing.T) {
	input := []int{1, 2, 3}
	q := Query(input).Skip(10)
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestSkipZero(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	q := Query(input).Skip(0)
	result := Slice(q)

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}
}

func TestChaining(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q := Query(input).
		Where(func(n int) bool { return n%2 == 0 }).
		Skip(1).
		Take(2)
	result := Slice(q)

	expected := []int{4, 6}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestTakeWhile(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q := Query(input).TakeWhile(func(n int) bool { return n < 5 })
	result := Slice(q)

	expected := []int{1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestTakeWhileNoneMatch(t *testing.T) {
	input := []int{5, 6, 7, 8, 9}
	q := Query(input).TakeWhile(func(n int) bool { return n < 5 })
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestTakeWhileAllMatch(t *testing.T) {
	input := []int{1, 2, 3, 4}
	q := Query(input).TakeWhile(func(n int) bool { return n < 10 })
	result := Slice(q)

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}
}

func TestSkipWhile(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q := Query(input).SkipWhile(func(n int) bool { return n < 5 })
	result := Slice(q)

	expected := []int{5, 6, 7, 8, 9, 10}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestSkipWhileNoneMatch(t *testing.T) {
	input := []int{5, 6, 7, 8, 9}
	q := Query(input).SkipWhile(func(n int) bool { return n < 5 })
	result := Slice(q)

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}
}

func TestSkipWhileAllMatch(t *testing.T) {
	input := []int{1, 2, 3, 4}
	q := Query(input).SkipWhile(func(n int) bool { return n < 10 })
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}
