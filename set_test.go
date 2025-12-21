package kk

import (
	"testing"
)

func TestConcat(t *testing.T) {
	q1 := From([]int{1, 2, 3})
	q2 := From([]int{4, 5, 6})
	result := Slice(q1.Concat(q2))

	expected := []int{1, 2, 3, 4, 5, 6}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestConcatEmpty(t *testing.T) {
	q1 := From([]int{1, 2, 3})
	q2 := From([]int{})
	result := Slice(q1.Concat(q2))

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

func TestConcatBothEmpty(t *testing.T) {
	q1 := From([]int{})
	q2 := From([]int{})
	result := Slice(q1.Concat(q2))

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestExcept(t *testing.T) {
	q1 := From([]int{1, 2, 3, 4, 5})
	q2 := From([]int{3, 4, 5, 6, 7})
	result := Slice(q1.Except(q2))

	expected := []int{1, 2}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestExceptWithDuplicates(t *testing.T) {
	q1 := From([]int{1, 1, 2, 2, 3, 3})
	q2 := From([]int{2})
	result := Slice(q1.Except(q2))

	expected := []int{1, 3}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

func TestExceptEmpty(t *testing.T) {
	q1 := From([]int{1, 2, 3})
	q2 := From([]int{})
	result := Slice(q1.Except(q2))

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}
}

func TestIntersect(t *testing.T) {
	q1 := From([]int{1, 2, 3, 4, 5})
	q2 := From([]int{3, 4, 5, 6, 7})
	result := Slice(q1.Intersect(q2))

	expected := []int{3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestIntersectWithDuplicates(t *testing.T) {
	q1 := From([]int{1, 2, 2, 3, 3, 3})
	q2 := From([]int{2, 3, 4})
	result := Slice(q1.Intersect(q2))

	expected := []int{2, 3}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

func TestIntersectEmpty(t *testing.T) {
	q1 := From([]int{1, 2, 3})
	q2 := From([]int{})
	result := Slice(q1.Intersect(q2))

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestIntersectNoOverlap(t *testing.T) {
	q1 := From([]int{1, 2, 3})
	q2 := From([]int{4, 5, 6})
	result := Slice(q1.Intersect(q2))

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestUnion(t *testing.T) {
	q1 := From([]int{1, 2, 3})
	q2 := From([]int{3, 4, 5})
	result := Slice(q1.Union(q2))

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

func TestUnionWithDuplicates(t *testing.T) {
	q1 := From([]int{1, 1, 2, 2})
	q2 := From([]int{2, 2, 3, 3})
	result := Slice(q1.Union(q2))

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

func TestUnionEmpty(t *testing.T) {
	q1 := From([]int{1, 2, 3})
	q2 := From([]int{})
	result := Slice(q1.Union(q2))

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

func TestUnionBothEmpty(t *testing.T) {
	q1 := From([]int{})
	q2 := From([]int{})
	result := Slice(q1.Union(q2))

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}
