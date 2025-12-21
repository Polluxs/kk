// Package kk provides a fluent, LINQ-style library for parallel processing in Go.
//
// It offers chainable query methods for filtering and transforming data,
// combined with powerful parallel execution functions.
//
// # KKQuery Building
//
// Build queries using method chaining on the KKQuery type:
//
//	q := kk.From(users).
//	    Where(func(u User) bool { return u.Active }).
//	    Take(100)
//
// # Methods (chainable)
//
// Methods that return a new KKQuery of the same type:
//   - Where(predicate) - Filter items
//   - Take(n) - First n items
//   - Skip(n) - Skip first n items
//   - TakeWhile(predicate) - Take while condition is true
//   - SkipWhile(predicate) - Skip while condition is true
//   - Distinct() - Remove duplicates
//   - Concat(other) - Combine queries
//   - Except(other) - Items not in other
//   - Intersect(other) - Items in both
//   - Union(other) - Items in either (distinct)
//
// # Functions (terminal)
//
// Package-level functions that transform, execute, or aggregate:
//   - Query(slice) - Create query from slice
//   - FromChan(ch) - Create query from channel
//   - Map(q, fn) - Transform each item to new type
//   - FlatMap(q, fn) - Transform and flatten
//   - Chunk(q, size) - Split into batches
//   - DistinctBy(q, keyFn) - Remove duplicates by key
//   - OrderBy(q, keyFn) - Sort ascending
//   - OrderByDescending(q, keyFn) - Sort descending
//   - ThenBy(oq, keyFn) - Secondary sort
//   - ThenByDescending(oq, keyFn) - Secondary sort descending
//   - Parallel(q, ctx, n, fn) - Process items in parallel
//   - ParallelResult(q, ctx, n, fn) - Process and collect results
//   - ParallelByKey(q, ctx, n, perKey, keyFn, fn) - Parallel with per-key limit
//   - ParallelByBatch(q, ctx, size, n, fn) - Process in batches
//   - Count(q) - Count items
//   - Sum(q, fn) - Sum values
//   - First(q) - First item
//   - Any(q, predicate) - Any match?
//   - All(q, predicate) - All match?
//   - ToSlice(q) - Materialize to slice
//   - Print(q) - Print items (debug)
//
// # Parallel Execution
//
// Process items in parallel with concurrency control:
//
//	err := kk.Parallel(q, ctx, 20, func(ctx context.Context, u User) error {
//	    return sendEmail(ctx, u.Email)
//	})
//
// # Batch Processing
//
//	err := kk.ParallelByBatch(q, ctx, 100, 4, func(ctx context.Context, batch []Record) error {
//	    return db.BulkInsert(ctx, batch)
//	})
//
// # Per-Key Rate Limiting
//
//	err := kk.ParallelByKey(q, ctx, 50, 2,
//	    func(o Order) string { return o.CustomerID },
//	    processOrder,
//	)
package kk
