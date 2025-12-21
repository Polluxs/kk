package kk

// Where filters items based on a predicate.
func (q *Query_[T]) Where(predicate func(T) bool) *Query_[T] {
	return &Query_[T]{
		iterate: func() Iterator[T] {
			iter := q.iterate()
			return func() (T, bool) {
				for {
					item, ok := iter()
					if !ok {
						var zero T
						return zero, false
					}
					if predicate(item) {
						return item, true
					}
				}
			}
		},
	}
}

// Take returns the first n items.
func (q *Query_[T]) Take(n int) *Query_[T] {
	return &Query_[T]{
		iterate: func() Iterator[T] {
			iter := q.iterate()
			count := 0
			return func() (T, bool) {
				if count >= n {
					var zero T
					return zero, false
				}
				item, ok := iter()
				if !ok {
					var zero T
					return zero, false
				}
				count++
				return item, true
			}
		},
	}
}

// Skip skips the first n items.
func (q *Query_[T]) Skip(n int) *Query_[T] {
	return &Query_[T]{
		iterate: func() Iterator[T] {
			iter := q.iterate()
			skipped := false
			return func() (T, bool) {
				if !skipped {
					for i := 0; i < n; i++ {
						_, ok := iter()
						if !ok {
							var zero T
							return zero, false
						}
					}
					skipped = true
				}
				return iter()
			}
		},
	}
}

// TakeWhile returns items while the predicate is true.
func (q *Query_[T]) TakeWhile(predicate func(T) bool) *Query_[T] {
	return &Query_[T]{
		iterate: func() Iterator[T] {
			iter := q.iterate()
			done := false
			return func() (T, bool) {
				if done {
					var zero T
					return zero, false
				}
				item, ok := iter()
				if !ok {
					var zero T
					return zero, false
				}
				if !predicate(item) {
					done = true
					var zero T
					return zero, false
				}
				return item, true
			}
		},
	}
}

// SkipWhile skips items while the predicate is true.
func (q *Query_[T]) SkipWhile(predicate func(T) bool) *Query_[T] {
	return &Query_[T]{
		iterate: func() Iterator[T] {
			iter := q.iterate()
			skipping := true
			return func() (T, bool) {
				for {
					item, ok := iter()
					if !ok {
						var zero T
						return zero, false
					}
					if skipping {
						if !predicate(item) {
							skipping = false
							return item, true
						}
						continue
					}
					return item, true
				}
			}
		},
	}
}
