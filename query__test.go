package kk

import (
	"testing"
)

func TestFrom(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	q := Query(input)
	result := Slice(q)

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}

	for i, v := range result {
		if v != input[i] {
			t.Errorf("at index %d: expected %d, got %d", i, input[i], v)
		}
	}
}

func TestFromEmpty(t *testing.T) {
	input := []int{}
	q := Query(input)
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestFromChan(t *testing.T) {
	ch := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	close(ch)

	q := QueryChan(ch)
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

func TestFromChanEmpty(t *testing.T) {
	ch := make(chan int)
	close(ch)

	q := QueryChan(ch)
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestQueryMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	q := QueryMapKeys(m)
	result := Slice(q)

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}

	seen := map[string]bool{}
	for _, k := range result {
		seen[k] = true
	}
	for k := range m {
		if !seen[k] {
			t.Errorf("missing key %q", k)
		}
	}
}

func TestQueryMapKeysEmpty(t *testing.T) {
	m := map[string]int{}
	q := QueryMapKeys(m)
	result := Slice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}
