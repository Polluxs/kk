package kk

import (
	"testing"
)

func TestDistinct(t *testing.T) {
	input := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	q := From(input).Distinct()
	result := ToSlice(q)

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

func TestDistinctEmpty(t *testing.T) {
	input := []int{}
	q := From(input).Distinct()
	result := ToSlice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestDistinctAllUnique(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	q := From(input).Distinct()
	result := ToSlice(q)

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}
}

func TestDistinctStrings(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b", "c"}
	q := From(input).Distinct()
	result := ToSlice(q)

	expected := []string{"a", "b", "c"}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

type Person struct {
	Name string
	Age  int
}

func TestDistinctBy(t *testing.T) {
	input := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 30},
		{Name: "David", Age: 25},
	}
	q := DistinctBy(From(input), func(p Person) int { return p.Age })
	result := ToSlice(q)

	if len(result) != 2 {
		t.Errorf("expected length 2, got %d", len(result))
	}

	if result[0].Name != "Alice" || result[1].Name != "Bob" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestDistinctByEmpty(t *testing.T) {
	input := []Person{}
	q := DistinctBy(From(input), func(p Person) int { return p.Age })
	result := ToSlice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestDistinctByAllUnique(t *testing.T) {
	input := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	q := DistinctBy(From(input), func(p Person) int { return p.Age })
	result := ToSlice(q)

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}
}
