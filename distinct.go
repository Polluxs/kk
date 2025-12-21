package kk

// Distinct removes duplicate items (requires comparable type).
func (q *KKQuery[T]) Distinct() *KKQuery[T] {
	return &KKQuery[T]{
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
func DistinctBy[T any, K comparable](q *KKQuery[T], keyFn func(T) K) *KKQuery[T] {
	return &KKQuery[T]{
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
