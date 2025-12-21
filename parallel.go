package kk

import (
	"context"
	"sync"
)

// Parallel processes items in parallel with a maximum of n concurrent operations.
// Returns the first error encountered, or nil if all operations succeed.
func Parallel[T any](ctx context.Context, q *Query[T], n int, fn func(context.Context, T) error) error {
	items := ToSlice(q)
	if len(items) == 0 {
		return nil
	}

	// Create a context that can be cancelled on first error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Semaphore for limiting concurrency
	sem := make(chan struct{}, n)

	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once

loop:
	for _, item := range items {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			break loop
		default:
		}

		// Acquire semaphore
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			break loop
		}

		wg.Add(1)
		go func(item T) {
			defer wg.Done()
			defer func() { <-sem }()

			// Check if we should still process
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := fn(ctx, item); err != nil {
				errOnce.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}(item)
	}

	wg.Wait()

	// Check if context was cancelled before we started
	if firstErr == nil && ctx.Err() != nil {
		return ctx.Err()
	}

	return firstErr
}

// ParallelResult processes items in parallel and collects results.
// Returns results and the first error encountered.
// This is a function (not a method) because it returns a different type.
func ParallelResult[T any, R any](ctx context.Context, q *Query[T], n int, fn func(context.Context, T) (R, error)) ([]R, error) {
	items := ToSlice(q)
	if len(items) == 0 {
		return nil, nil
	}

	// Create a context that can be cancelled on first error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Semaphore for limiting concurrency
	sem := make(chan struct{}, n)

	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once
	var mu sync.Mutex

	// Pre-allocate results slice to maintain order
	results := make([]R, len(items))

loop:
	for i, item := range items {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			break loop
		default:
		}

		// Acquire semaphore
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			break loop
		}

		wg.Add(1)
		go func(idx int, item T) {
			defer wg.Done()
			defer func() { <-sem }()

			// Check if we should still process
			select {
			case <-ctx.Done():
				return
			default:
			}

			result, err := fn(ctx, item)
			if err != nil {
				errOnce.Do(func() {
					firstErr = err
					cancel()
				})
				return
			}

			mu.Lock()
			results[idx] = result
			mu.Unlock()
		}(i, item)
	}

	wg.Wait()

	// Check if context was cancelled before we started
	if firstErr == nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if firstErr != nil {
		return nil, firstErr
	}

	return results, nil
}

// ParallelByKey processes items in parallel with both a global limit and per-key limit.
// n is the maximum total concurrent operations.
// perKey is the maximum concurrent operations per key.
// keyFn extracts the key from each item.
func ParallelByKey[T any, K comparable](ctx context.Context, q *Query[T], n int, perKey int, keyFn func(T) K, fn func(context.Context, T) error) error {
	items := ToSlice(q)
	if len(items) == 0 {
		return nil
	}

	// Create a context that can be cancelled on first error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Global semaphore for limiting total concurrency
	globalSem := make(chan struct{}, n)

	// Per-key semaphores
	var keySemMu sync.Mutex
	keySems := make(map[K]chan struct{})

	getKeySem := func(key K) chan struct{} {
		keySemMu.Lock()
		defer keySemMu.Unlock()
		if sem, ok := keySems[key]; ok {
			return sem
		}
		sem := make(chan struct{}, perKey)
		keySems[key] = sem
		return sem
	}

	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once

loop:
	for _, item := range items {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			break loop
		default:
		}

		key := keyFn(item)
		keySem := getKeySem(key)

		// Acquire global semaphore
		select {
		case globalSem <- struct{}{}:
		case <-ctx.Done():
			break loop
		}

		// Acquire per-key semaphore
		select {
		case keySem <- struct{}{}:
		case <-ctx.Done():
			<-globalSem // Release global semaphore
			break loop
		}

		wg.Add(1)
		go func(item T, keySem chan struct{}) {
			defer wg.Done()
			defer func() {
				<-keySem
				<-globalSem
			}()

			// Check if we should still process
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := fn(ctx, item); err != nil {
				errOnce.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}(item, keySem)
	}

	wg.Wait()

	// Check if context was cancelled before we started
	if firstErr == nil && ctx.Err() != nil {
		return ctx.Err()
	}

	return firstErr
}

// ParallelByBatch processes items in batches with parallel batch execution.
// batchSize is the number of items per batch.
// n is the maximum number of concurrent batches.
func ParallelByBatch[T any](ctx context.Context, q *Query[T], batchSize int, n int, fn func(context.Context, []T) error) error {
	// Create batches using Chunk
	batches := ToSlice(Chunk(q, batchSize))
	if len(batches) == 0 {
		return nil
	}

	// Create a context that can be cancelled on first error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Semaphore for limiting batch concurrency
	sem := make(chan struct{}, n)

	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once

loop:
	for _, batch := range batches {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			break loop
		default:
		}

		// Acquire semaphore
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			break loop
		}

		wg.Add(1)
		go func(batch []T) {
			defer wg.Done()
			defer func() { <-sem }()

			// Check if we should still process
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := fn(ctx, batch); err != nil {
				errOnce.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}(batch)
	}

	wg.Wait()

	// Check if context was cancelled before we started
	if firstErr == nil && ctx.Err() != nil {
		return ctx.Err()
	}

	return firstErr
}
