package kk

// Mapped transforms each item to a new type.
// This is a function (not a method) because it returns a different type.
func Mapped[T any, R any](q *KKQuery[T], fn func(T) R) *KKQuery[R] {
	return &KKQuery[R]{
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

// Flattened transforms each item to a slice and flattens the results.
// This is a function (not a method) because it returns a different type.
func Flattened[T any, R any](q *KKQuery[T], fn func(T) []R) *KKQuery[R] {
	return &KKQuery[R]{
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
