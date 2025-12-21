package kk

// KKQuery represents a lazy sequence of items that can be filtered, transformed, and executed.
type KKQuery[T any] struct {
	iterate func() Iterator[T]
}

// Iterator represents a function that returns the next item and whether there are more items.
type Iterator[T any] func() (T, bool)

// Query creates a KKQuery from a slice.
func Query[T any](slice []T) *KKQuery[T] {
	return &KKQuery[T]{
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

// QueryChan creates a KKQuery from a channel.
func QueryChan[T any](ch <-chan T) *KKQuery[T] {
	return &KKQuery[T]{
		iterate: func() Iterator[T] {
			return func() (T, bool) {
				item, ok := <-ch
				return item, ok
			}
		},
	}
}

// Slice materializes the query to a slice.
func Slice[T any](q *KKQuery[T]) []T {
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
