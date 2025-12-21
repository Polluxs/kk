package kk

import (
	"fmt"
)

// Count returns the number of items in the query.
func Count[T any](q *Query[T]) int {
	count := 0
	iter := q.iterate()
	for {
		_, ok := iter()
		if !ok {
			break
		}
		count++
	}
	return count
}

// Sum returns the sum of values produced by the selector function.
func Sum[T any, N Number](q *Query[T], selector func(T) N) N {
	var sum N
	iter := q.iterate()
	for {
		item, ok := iter()
		if !ok {
			break
		}
		sum += selector(item)
	}
	return sum
}

// Number is a constraint that permits numeric types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// First returns the first item, or the zero value if the query is empty.
// The second return value indicates whether an item was found.
func First[T any](q *Query[T]) (T, bool) {
	iter := q.iterate()
	return iter()
}

// Any returns true if any item matches the predicate.
func Any[T any](q *Query[T], predicate func(T) bool) bool {
	iter := q.iterate()
	for {
		item, ok := iter()
		if !ok {
			return false
		}
		if predicate(item) {
			return true
		}
	}
}

// All returns true if all items match the predicate.
// Returns true for empty queries.
func All[T any](q *Query[T], predicate func(T) bool) bool {
	iter := q.iterate()
	for {
		item, ok := iter()
		if !ok {
			return true
		}
		if !predicate(item) {
			return false
		}
	}
}

// Print prints all items in the query (for debugging).
func Print[T any](q *Query[T]) {
	iter := q.iterate()
	for {
		item, ok := iter()
		if !ok {
			break
		}
		fmt.Printf("%v\n", item)
	}
}
