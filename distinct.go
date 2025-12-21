package kk

// Distinct removes duplicate items (requires comparable type).
func (q *Query[T]) Distinct() *Query[T] {
	return &Query[T]{
		iterate: func() Iterator[T] {
			iter := q.iterate()
			seen := make(map[any]bool)
			return func() (T, bool) {
				for {
					item, ok := iter()
					if !ok {
						var zero T
						return zero, false
					}
					if !seen[item] {
						seen[item] = true
						return item, true
					}
				}
			}
		},
	}
}

// DistinctBy removes duplicate items based on a key function.
func DistinctBy[T any, K comparable](q *Query[T], keyFn func(T) K) *Query[T] {
	return &Query[T]{
		iterate: func() Iterator[T] {
			iter := q.iterate()
			seen := make(map[K]bool)
			return func() (T, bool) {
				for {
					item, ok := iter()
					if !ok {
						var zero T
						return zero, false
					}
					key := keyFn(item)
					if !seen[key] {
						seen[key] = true
						return item, true
					}
				}
			}
		},
	}
}
