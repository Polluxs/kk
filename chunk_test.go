package kk

import (
	"testing"
)

func TestChunk(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q := Chunk(From(input), 3)
	result := ToSlice(q)

	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, batch := range result {
		if len(batch) != len(expected[i]) {
			t.Errorf("batch %d: expected length %d, got %d", i, len(expected[i]), len(batch))
		}
		for j, v := range batch {
			if v != expected[i][j] {
				t.Errorf("batch %d, index %d: expected %d, got %d", i, j, expected[i][j], v)
			}
		}
	}
}

func TestChunkExact(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	q := Chunk(From(input), 2)
	result := ToSlice(q)

	expected := [][]int{{1, 2}, {3, 4}, {5, 6}}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, batch := range result {
		if len(batch) != len(expected[i]) {
			t.Errorf("batch %d: expected length %d, got %d", i, len(expected[i]), len(batch))
		}
	}
}

func TestChunkEmpty(t *testing.T) {
	input := []int{}
	q := Chunk(From(input), 3)
	result := ToSlice(q)

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestChunkSizeOne(t *testing.T) {
	input := []int{1, 2, 3}
	q := Chunk(From(input), 1)
	result := ToSlice(q)

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}

	for i, batch := range result {
		if len(batch) != 1 {
			t.Errorf("batch %d: expected length 1, got %d", i, len(batch))
		}
		if batch[0] != i+1 {
			t.Errorf("batch %d: expected %d, got %d", i, i+1, batch[0])
		}
	}
}

func TestChunkLargerThanInput(t *testing.T) {
	input := []int{1, 2, 3}
	q := Chunk(From(input), 10)
	result := ToSlice(q)

	if len(result) != 1 {
		t.Errorf("expected length 1, got %d", len(result))
	}

	if len(result[0]) != 3 {
		t.Errorf("expected batch length 3, got %d", len(result[0]))
	}
}

func TestChunkWithChaining(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q := Chunk(From(input).Where(func(n int) bool { return n%2 == 0 }), 2)
	result := ToSlice(q)

	expected := [][]int{{2, 4}, {6, 8}, {10}}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, batch := range result {
		if len(batch) != len(expected[i]) {
			t.Errorf("batch %d: expected length %d, got %d", i, len(expected[i]), len(batch))
		}
	}
}
