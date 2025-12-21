package kk

import (
	"cmp"
	"sort"
)

// OrderedQuery wraps a Query with ordering capabilities for ThenBy chaining.
type OrderedQuery[T any] struct {
	*Query[T]
	comparators []func(T, T) int
}

// OrderBy sorts items in ascending order by a key.
func OrderBy[T any, K cmp.Ordered](q *Query[T], keyFn func(T) K) *OrderedQuery[T] {
	return &OrderedQuery[T]{
		Query: &Query[T]{
			iterate: func() Iterator[T] {
				items := ToSlice(q)
				sort.Slice(items, func(i, j int) bool {
					return keyFn(items[i]) < keyFn(items[j])
				})
				index := 0
				return func() (T, bool) {
					if index >= len(items) {
						var zero T
						return zero, false
					}
					item := items[index]
					index++
					return item, true
				}
			},
		},
		comparators: []func(T, T) int{
			func(a, b T) int {
				ka, kb := keyFn(a), keyFn(b)
				if ka < kb {
					return -1
				}
				if ka > kb {
					return 1
				}
				return 0
			},
		},
	}
}

// OrderByDescending sorts items in descending order by a key.
func OrderByDescending[T any, K cmp.Ordered](q *Query[T], keyFn func(T) K) *OrderedQuery[T] {
	return &OrderedQuery[T]{
		Query: &Query[T]{
			iterate: func() Iterator[T] {
				items := ToSlice(q)
				sort.Slice(items, func(i, j int) bool {
					return keyFn(items[i]) > keyFn(items[j])
				})
				index := 0
				return func() (T, bool) {
					if index >= len(items) {
						var zero T
						return zero, false
					}
					item := items[index]
					index++
					return item, true
				}
			},
		},
		comparators: []func(T, T) int{
			func(a, b T) int {
				ka, kb := keyFn(a), keyFn(b)
				if ka > kb {
					return -1
				}
				if ka < kb {
					return 1
				}
				return 0
			},
		},
	}
}

// ThenBy adds a secondary ascending sort.
func ThenBy[T any, K cmp.Ordered](oq *OrderedQuery[T], keyFn func(T) K) *OrderedQuery[T] {
	newComparators := append(oq.comparators, func(a, b T) int {
		ka, kb := keyFn(a), keyFn(b)
		if ka < kb {
			return -1
		}
		if ka > kb {
			return 1
		}
		return 0
	})

	return &OrderedQuery[T]{
		Query: &Query[T]{
			iterate: func() Iterator[T] {
				items := ToSlice(oq.Query)
				sort.Slice(items, func(i, j int) bool {
					for _, cmp := range newComparators {
						result := cmp(items[i], items[j])
						if result != 0 {
							return result < 0
						}
					}
					return false
				})
				index := 0
				return func() (T, bool) {
					if index >= len(items) {
						var zero T
						return zero, false
					}
					item := items[index]
					index++
					return item, true
				}
			},
		},
		comparators: newComparators,
	}
}

// ThenByDescending adds a secondary descending sort.
func ThenByDescending[T any, K cmp.Ordered](oq *OrderedQuery[T], keyFn func(T) K) *OrderedQuery[T] {
	newComparators := append(oq.comparators, func(a, b T) int {
		ka, kb := keyFn(a), keyFn(b)
		if ka > kb {
			return -1
		}
		if ka < kb {
			return 1
		}
		return 0
	})

	return &OrderedQuery[T]{
		Query: &Query[T]{
			iterate: func() Iterator[T] {
				items := ToSlice(oq.Query)
				sort.Slice(items, func(i, j int) bool {
					for _, cmp := range newComparators {
						result := cmp(items[i], items[j])
						if result != 0 {
							return result < 0
						}
					}
					return false
				})
				index := 0
				return func() (T, bool) {
					if index >= len(items) {
						var zero T
						return zero, false
					}
					item := items[index]
					index++
					return item, true
				}
			},
		},
		comparators: newComparators,
	}
}
