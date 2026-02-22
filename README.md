<div align="center">
  <h1>kk</h1>
  <p><strong>A fluent, QUERY-style library for parallel processing in Go.</strong></p>
</div>


> **⚠️ Warning:** This package is in active development and may introduce breaking changes between versions.

## Install

```bash
go get github.com/polluxs/kk
```

## Quick Start

```go
// Build a query with method chaining
q := kk.From(users).
    Where(func(u User) bool { return u.Active }).
    Take(100)

// Execute in parallel
err := kk.Parallel(q, ctx, 10, process)
```

---

## API

### Methods (chainable)

Called on the query, returns a new query of the same type.

| Method | Description |
|:-------|:------------|
| `.Where(predicate)` | Filter items |
| `.Take(n)` | First n items |
| `.Skip(n)` | Skip first n items |
| `.TakeWhile(predicate)` | Take while condition is true |
| `.SkipWhile(predicate)` | Skip while condition is true |
| `.Distinct()` | Remove duplicates |
| `.DistinctBy(keyFn)` | Remove duplicates by key |
| `.SortedBy(keyFn)` | Sort ascending |
| `.SortedByDesc(keyFn)` | Sort descending |
| `.ThenBy(keyFn)` | Secondary sort |
| `.Chunk(size)` | Split into batches |
| `.Concat(other)` | Combine queries |
| `.Except(other)` | Items not in other |
| `.Intersect(other)` | Items in both |
| `.Union(other)` | Items in either (distinct) |

### Functions (terminal)

Package-level functions that transform, execute, or aggregate.

| Function | Description |
|:---------|:------------|
| `kk.From(slice)` | Create query from slice |
| `kk.FromChan(ch)` | Create query from channel |
| `kk.QueryMapKeys(m)` | Create query from map keys |
| `kk.Mapped(q, fn)` | Transform each item to new type |
| `kk.Flattened(q, fn)` | Transform and flatten |
| `kk.GroupedBy(q, keyFn)` | Group items by key |
| `kk.Parallel(q, ctx, n, fn)` | Process items in parallel |
| `kk.ParallelResult(q, ctx, n, fn)` | Process and collect results |
| `kk.ParallelByKey(q, ctx, n, perKey, keyFn, fn)` | Parallel with per-key limit |
| `kk.ParallelByBatch(q, ctx, size, n, fn)` | Process in batches |
| `kk.ParallelByBatchChan(ctx, ch, size, n, fn)` | Stream batches from channel |
| `kk.Count(q)` | Count items |
| `kk.Sum(q, fn)` | Sum values |
| `kk.First(q)` | First item |
| `kk.Any(q, predicate)` | Any match? |
| `kk.All(q, predicate)` | All match? |
| `kk.Slice(q)` | Materialize to slice |
| `kk.Print(q)` | Print items (debug) |

---

## Examples

### Filter and process in parallel

```go
q := kk.From(users).
    Where(func(u User) bool { return u.Active }).
    Take(1000)

err := kk.Parallel(q, ctx, 20, func(ctx context.Context, u User) error {
    return sendEmail(ctx, u.Email)
})
```

### Transform and collect results

```go
q := kk.From(urls).Where(isValid)
responses, err := kk.ParallelResult(
    kk.Mapped(q, toRequest),
    ctx, 10, fetch,
)
```

### Group and aggregate

```go
// Group users by country, count per group
groups := kk.GroupedBy(kk.From(users), func(u User) string {
    return u.Country
})

for _, g := range kk.Slice(groups) {
    fmt.Printf("%s: %d users\n", g.Key, len(g.Items))
}
```

### Per-key rate limiting

```go
// Max 50 total, max 2 per customer
q := kk.From(orders).Where(isPending)
err := kk.ParallelByKey(q, ctx, 50, 2,
    func(o Order) string { return o.CustomerID },
    processOrder,
)
```

### Batch processing

```go
q := kk.From(records).Where(isValid).Take(10000)
err := kk.ParallelByBatch(q, ctx, 100, 4, func(ctx context.Context, batch []Record) error {
    return db.BulkInsert(ctx, batch)
})
```

### Streaming batch processing from a channel

```go
ch := make(chan Record)
go produceRecords(ch) // closes ch when done

err := kk.ParallelByBatchChan(ctx, ch, 100, 4, func(ctx context.Context, batch []Record) error {
    return db.BulkInsert(ctx, batch)
})
```

### Debug a query

```go
kk.Print(kk.From(users).Where(active).Take(5))
```

---

## ⚠️ WARNING: Go Generics Suck

You might wonder why `Mapped` and `ParallelResult` are awkward package functions instead of nice chainable methods like everything else.

**Because Go generics are half-baked.**

Go's generics have a fundamental limitation: **methods cannot introduce new type parameters**. This means:

```go
// This works - T stays T
func (q *Query[T]) Where(fn func(T) bool) *Query[T]  ✓

// This doesn't compile - R is new
func (q *Query[T]) Mapped[R any](fn func(T) R) *Query[R]  ✗
```

So we're forced to write:
```go
kk.Mapped(query, func(u User) UserDTO { ... })  // ugly
```

Instead of:
```go
query.Mapped(func(u User) UserDTO { ... })  // what we want
```

This affects every operation that transforms to a different type:
- `Mapped` → stuck as function
- `Flattened` → stuck as function
- `ParallelResult` → stuck as function
- `GroupedBy` → stuck as function
- Any future `Select`, `Zip` → all stuck as functions

**9 out of 10 times, you're mapping to your own custom types**. And if ```[ANY]``` was OK for me, I'd write Typescript. 

So we use: `kk.Mapped(kk.Mapped(query, fn1), fn2)` nesting instead of `query.Mapped(fn1).Mapped(fn2)`.

I did what I could with it 🤷

---

For consistency, all terminal operations (`Parallel*`, `Sum`, `Count`, etc.) are also functions. This gives a clear mental model:

- **Methods** = build the query (chainable)
- **Functions** = transform type or execute

---

## License

MIT
