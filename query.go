package kk

// Query represents a lazy sequence of items that can be filtered, transformed, and executed.
type Query[T any] struct {
	iterate func() Iterator[T]
}

// Iterator represents a function that returns the next item and whether there are more items.
type Iterator[T any] func() (T, bool)

// From creates a Query from a slice.
func From[T any](slice []T) *Query[T] {
	return &Query[T]{
		iterate: func() Iterator[T] {
			index := 0
			return func() (T, bool) {
				if index >= len(slice) {
					var zero T
					return zero, false
				}
				item := slice[index]
				index++
				return item, true
			}
		},
	}
}

// FromChan creates a Query from a channel.
func FromChan[T any](ch <-chan T) *Query[T] {
	return &Query[T]{
		iterate: func() Iterator[T] {
			return func() (T, bool) {
				item, ok := <-ch
				return item, ok
			}
		},
	}
}

// ToSlice materializes the query to a slice.
func ToSlice[T any](q *Query[T]) []T {
	var result []T
	iter := q.iterate()
	for {
		item, ok := iter()
		if !ok {
			break
		}
		result = append(result, item)
	}
	return result
}
