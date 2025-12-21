package kk

// Concat combines two queries into one.
func (q *KKQuery[T]) Concat(other *KKQuery[T]) *KKQuery[T] {
	return &KKQuery[T]{
		iterate: func() Iterator[T] {
			iter1 := q.iterate()
			iter2 := other.iterate()
			first := true
			return func() (T, bool) {
				if first {
					item, ok := iter1()
					if ok {
						return item, true
					}
					first = false
				}
				return iter2()
			}
		},
	}
}

// Except returns items in the first query that are not in the second.
func (q *KKQuery[T]) Except(other *KKQuery[T]) *KKQuery[T] {
	return &KKQuery[T]{
		iterate: func() Iterator[T] {
			// Materialize the second query to check membership
			otherSet := make(map[any]bool)
			for _, item := range Slice(other) {
				otherSet[item] = true
			}

			iter := q.iterate()
			seen := make(map[any]bool)
			return func() (T, bool) {
				for {
					item, ok := iter()
					if !ok {
						var zero T
						return zero, false
					}
					// Skip if in other set or already seen
					if !otherSet[item] && !seen[item] {
						seen[item] = true
						return item, true
					}
				}
			}
		},
	}
}

// Intersect returns items that are in both queries (distinct).
func (q *KKQuery[T]) Intersect(other *KKQuery[T]) *KKQuery[T] {
	return &KKQuery[T]{
		iterate: func() Iterator[T] {
			// Materialize the second query to check membership
			otherSet := make(map[any]bool)
			for _, item := range Slice(other) {
				otherSet[item] = true
			}

			iter := q.iterate()
			seen := make(map[any]bool)
			return func() (T, bool) {
				for {
					item, ok := iter()
					if !ok {
						var zero T
						return zero, false
					}
					// Return if in other set and not already seen
					if otherSet[item] && !seen[item] {
						seen[item] = true
						return item, true
					}
				}
			}
		},
	}
}

// Union returns items that are in either query (distinct).
func (q *KKQuery[T]) Union(other *KKQuery[T]) *KKQuery[T] {
	return &KKQuery[T]{
		iterate: func() Iterator[T] {
			iter1 := q.iterate()
			iter2 := other.iterate()
			seen := make(map[any]bool)
			first := true
			return func() (T, bool) {
				for {
					if first {
						item, ok := iter1()
						if ok {
							if !seen[item] {
								seen[item] = true
								return item, true
							}
							continue
						}
						first = false
					}
					item, ok := iter2()
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
