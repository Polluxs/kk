package kk

// Map transforms each item to a new type.
// This is a function (not a method) because it returns a different type.
func Map[T any, R any](q *Query[T], fn func(T) R) *Query[R] {
	return &Query[R]{
		iterate: func() Iterator[R] {
			iter := q.iterate()
			return func() (R, bool) {
				item, ok := iter()
				if !ok {
					var zero R
					return zero, false
				}
				return fn(item), true
			}
		},
	}
}

// FlatMap transforms each item to a slice and flattens the results.
// This is a function (not a method) because it returns a different type.
func FlatMap[T any, R any](q *Query[T], fn func(T) []R) *Query[R] {
	return &Query[R]{
		iterate: func() Iterator[R] {
			iter := q.iterate()
			var current []R
			index := 0
			return func() (R, bool) {
				for {
					// Return items from current slice if available
					if index < len(current) {
						item := current[index]
						index++
						return item, true
					}

					// Get next source item
					item, ok := iter()
					if !ok {
						var zero R
						return zero, false
					}

					// Transform to slice and reset index
					current = fn(item)
					index = 0
				}
			}
		},
	}
}
