package kk

import (
	"testing"
)

func TestSortedBy(t *testing.T) {
	input := []int{5, 2, 8, 1, 9, 3}
	q := SortedBy(Query(input), func(n int) int { return n })
	result := Slice(q.KKQuery)

	expected := []int{1, 2, 3, 5, 8, 9}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestOrderByEmpty(t *testing.T) {
	input := []int{}
	q := SortedBy(Query(input), func(n int) int { return n })
	result := Slice(q.KKQuery)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestSortedByDesc(t *testing.T) {
	input := []int{5, 2, 8, 1, 9, 3}
	q := SortedByDesc(Query(input), func(n int) int { return n })
	result := Slice(q.KKQuery)

	expected := []int{9, 8, 5, 3, 2, 1}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestOrderByString(t *testing.T) {
	input := []string{"banana", "apple", "cherry", "date"}
	q := SortedBy(Query(input), func(s string) string { return s })
	result := Slice(q.KKQuery)

	expected := []string{"apple", "banana", "cherry", "date"}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

type Employee struct {
	Name       string
	Department string
	Salary     int
}

func TestThenBy(t *testing.T) {
	input := []Employee{
		{Name: "Alice", Department: "IT", Salary: 50000},
		{Name: "Bob", Department: "HR", Salary: 45000},
		{Name: "Charlie", Department: "IT", Salary: 55000},
		{Name: "David", Department: "HR", Salary: 40000},
	}

	q := ThenBy(
		SortedBy(Query(input), func(e Employee) string { return e.Department }),
		func(e Employee) int { return e.Salary },
	)
	result := Slice(q.KKQuery)

	// Should be: HR(David 40k, Bob 45k), IT(Alice 50k, Charlie 55k)
	expectedNames := []string{"David", "Bob", "Alice", "Charlie"}
	for i, e := range result {
		if e.Name != expectedNames[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expectedNames[i], e.Name)
		}
	}
}

func TestThenByDescending(t *testing.T) {
	input := []Employee{
		{Name: "Alice", Department: "IT", Salary: 50000},
		{Name: "Bob", Department: "HR", Salary: 45000},
		{Name: "Charlie", Department: "IT", Salary: 55000},
		{Name: "David", Department: "HR", Salary: 40000},
	}

	q := ThenByDescending(
		SortedBy(Query(input), func(e Employee) string { return e.Department }),
		func(e Employee) int { return e.Salary },
	)
	result := Slice(q.KKQuery)

	// Should be: HR(Bob 45k, David 40k), IT(Charlie 55k, Alice 50k)
	expectedNames := []string{"Bob", "David", "Charlie", "Alice"}
	for i, e := range result {
		if e.Name != expectedNames[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expectedNames[i], e.Name)
		}
	}
}

func TestOrderByWithChaining(t *testing.T) {
	input := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	q := SortedBy(
		Query(input).Where(func(n int) bool { return n > 3 }), func(n int) int { return n },
	)
	result := Slice(q.KKQuery.Take(3))

	expected := []int{4, 5, 6}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}
