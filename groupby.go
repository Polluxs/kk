package kk

// Group represents a collection of items that share a common key.
type Group[K comparable, T any] struct {
	Key   K
	Items []T
}

// GroupedBy groups items by a key function and returns a query of groups.
// This is a function (not a method) because it returns a different type.
func GroupedBy[T any, K comparable](q *Query[T], keyFn func(T) K) *Query[Group[K, T]] {
	return &Query[Group[K, T]]{
		iterate: func() Iterator[Group[K, T]] {
			// Materialize all items and group them
			groups := make(map[K][]T)
			var keys []K // maintain insertion order

			iter := q.iterate()
			for {
				item, ok := iter()
				if !ok {
					break
				}
				key := keyFn(item)
				if _, exists := groups[key]; !exists {
					keys = append(keys, key)
				}
				groups[key] = append(groups[key], item)
			}

			index := 0
			return func() (Group[K, T], bool) {
				if index >= len(keys) {
					var zero Group[K, T]
					return zero, false
				}
				key := keys[index]
				index++
				return Group[K, T]{Key: key, Items: groups[key]}, true
			}
		},
	}
}
