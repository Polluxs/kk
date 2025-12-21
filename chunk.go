package kk

// Chunk splits items into batches of the specified size.
// This is a function (not a method) because it returns a different type.
func Chunk[T any](q *KKQuery[T], size int) *KKQuery[[]T] {
	return &KKQuery[[]T]{
		iterate: func() Iterator[[]T] {
			iter := q.iterate()
			done := false
			return func() ([]T, bool) {
				if done {
					return nil, false
				}
				batch := make([]T, 0, size)
				for i := 0; i < size; i++ {
					item, ok := iter()
					if !ok {
						done = true
						break
					}
					batch = append(batch, item)
				}
				if len(batch) == 0 {
					return nil, false
				}
				return batch, true
			}
		},
	}
}
