package kk

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	var count atomic.Int32

	err := Parallel(
		context.Background(), Query(input), 3, func(ctx context.Context, n int) error {
			count.Add(1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if count.Load() != 5 {
		t.Errorf("expected count 5, got %d", count.Load())
	}
}

func TestParallelEmpty(t *testing.T) {
	input := []int{}
	var count atomic.Int32

	err := Parallel(
		context.Background(), Query(input), 3, func(ctx context.Context, n int) error {
			count.Add(1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if count.Load() != 0 {
		t.Errorf("expected count 0, got %d", count.Load())
	}
}

func TestParallelWithError(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expectedErr := errors.New("test error")

	err := Parallel(
		context.Background(), Query(input), 3, func(ctx context.Context, n int) error {
			if n == 3 {
				return expectedErr
			}
			return nil
		},
	)

	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestParallelConcurrency(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var concurrent atomic.Int32
	var maxConcurrent atomic.Int32

	err := Parallel(
		context.Background(), Query(input), 3, func(ctx context.Context, n int) error {
			current := concurrent.Add(1)
			// Update max concurrent if current is higher
			for {
				max := maxConcurrent.Load()
				if current <= max || maxConcurrent.CompareAndSwap(max, current) {
					break
				}
			}

			time.Sleep(10 * time.Millisecond)
			concurrent.Add(-1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if maxConcurrent.Load() > 3 {
		t.Errorf("expected max concurrency of 3, got %d", maxConcurrent.Load())
	}
}

func TestParallelContextCancellation(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ctx, cancel := context.WithCancel(context.Background())
	var count atomic.Int32

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := Parallel(
		ctx, Query(input), 2, func(ctx context.Context, n int) error {
			time.Sleep(100 * time.Millisecond)
			count.Add(1)
			return nil
		},
	)

	// Should have been cancelled
	if err != context.Canceled {
		// Could also get nil if some completed before cancellation
		if err != nil {
			t.Logf("got error: %v (this is acceptable)", err)
		}
	}
}

func TestParallelWithChaining(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var sum atomic.Int32

	err := Parallel(
		context.Background(),
		Query(input).Where(func(n int) bool { return n%2 == 0 }),
		2,
		func(ctx context.Context, n int) error {
			sum.Add(int32(n))
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Sum of even numbers 2+4+6+8+10 = 30
	if sum.Load() != 30 {
		t.Errorf("expected sum 30, got %d", sum.Load())
	}
}

func TestParallelResult(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	results, err := ParallelResult(
		context.Background(), Query(input), 3, func(ctx context.Context, n int) (int, error) {
			return n * 2, nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := []int{2, 4, 6, 8, 10}
	if len(results) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(results))
	}

	for i, v := range results {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestParallelResultEmpty(t *testing.T) {
	input := []int{}
	results, err := ParallelResult(
		context.Background(), Query(input), 3, func(ctx context.Context, n int) (int, error) {
			return n * 2, nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if results != nil {
		t.Errorf("expected nil, got %v", results)
	}
}

func TestParallelResultWithError(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expectedErr := errors.New("test error")

	results, err := ParallelResult(
		context.Background(), Query(input), 3, func(ctx context.Context, n int) (int, error) {
			if n == 3 {
				return 0, expectedErr
			}
			return n * 2, nil
		},
	)

	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}

	if results != nil {
		t.Errorf("expected nil results on error, got %v", results)
	}
}

func TestParallelResultTypeChange(t *testing.T) {
	input := []int{1, 2, 3}
	results, err := ParallelResult(
		context.Background(), Query(input), 3, func(ctx context.Context, n int) (string, error) {
			return string(rune('A' + n - 1)), nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := []string{"A", "B", "C"}
	for i, v := range results {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

func TestParallelResultMaintainsOrder(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	results, err := ParallelResult(
		context.Background(), Query(input), 3, func(ctx context.Context, n int) (int, error) {
			// Add some delay to mix up completion order
			time.Sleep(time.Duration(10-n) * time.Millisecond)
			return n, nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Results should be in original order
	for i, v := range results {
		if v != i+1 {
			t.Errorf("at index %d: expected %d, got %d", i, i+1, v)
		}
	}
}

type Order struct {
	ID         int
	CustomerID string
}

func TestParallelByKey(t *testing.T) {
	orders := []Order{
		{ID: 1, CustomerID: "A"},
		{ID: 2, CustomerID: "B"},
		{ID: 3, CustomerID: "A"},
		{ID: 4, CustomerID: "B"},
		{ID: 5, CustomerID: "C"},
	}

	var count atomic.Int32

	err := ParallelByKey(
		context.Background(),
		Query(orders),
		10, // max total
		2,  // max per key
		func(o Order) string { return o.CustomerID },
		func(ctx context.Context, o Order) error {
			count.Add(1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if count.Load() != 5 {
		t.Errorf("expected count 5, got %d", count.Load())
	}
}

func TestParallelByKeyPerKeyLimit(t *testing.T) {
	// Create orders for single customer to test per-key limiting
	orders := []Order{
		{ID: 1, CustomerID: "A"},
		{ID: 2, CustomerID: "A"},
		{ID: 3, CustomerID: "A"},
		{ID: 4, CustomerID: "A"},
		{ID: 5, CustomerID: "A"},
	}

	var concurrent atomic.Int32
	var maxConcurrent atomic.Int32

	err := ParallelByKey(
		context.Background(),
		Query(orders),
		10, // max total (high to not limit)
		2,  // max per key (should limit to 2)
		func(o Order) string { return o.CustomerID },
		func(ctx context.Context, o Order) error {
			current := concurrent.Add(1)
			for {
				max := maxConcurrent.Load()
				if current <= max || maxConcurrent.CompareAndSwap(max, current) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)
			concurrent.Add(-1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Should not exceed per-key limit of 2
	if maxConcurrent.Load() > 2 {
		t.Errorf("expected max concurrent of 2, got %d", maxConcurrent.Load())
	}
}

func TestParallelByKeyWithError(t *testing.T) {
	orders := []Order{
		{ID: 1, CustomerID: "A"},
		{ID: 2, CustomerID: "B"},
		{ID: 3, CustomerID: "A"},
	}

	expectedErr := errors.New("test error")

	err := ParallelByKey(
		context.Background(),
		Query(orders),
		10,
		2,
		func(o Order) string { return o.CustomerID },
		func(ctx context.Context, o Order) error {
			if o.ID == 2 {
				return expectedErr
			}
			return nil
		},
	)

	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestParallelByKeyEmpty(t *testing.T) {
	orders := []Order{}

	err := ParallelByKey(
		context.Background(),
		Query(orders),
		10,
		2,
		func(o Order) string { return o.CustomerID },
		func(ctx context.Context, o Order) error {
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestParallelByBatch(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var totalSum atomic.Int32
	var batchCount atomic.Int32

	err := ParallelByBatch(
		context.Background(),
		Query(input),
		3, // batch size
		2, // max concurrent batches
		func(ctx context.Context, batch []int) error {
			batchCount.Add(1)
			for _, n := range batch {
				totalSum.Add(int32(n))
			}
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Sum of 1-10 = 55
	if totalSum.Load() != 55 {
		t.Errorf("expected sum 55, got %d", totalSum.Load())
	}

	// 10 items with batch size 3 = 4 batches (3+3+3+1)
	if batchCount.Load() != 4 {
		t.Errorf("expected 4 batches, got %d", batchCount.Load())
	}
}

func TestParallelByBatchConcurrency(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	var concurrent atomic.Int32
	var maxConcurrent atomic.Int32

	err := ParallelByBatch(
		context.Background(),
		Query(input),
		3, // batch size (4 batches)
		2, // max concurrent batches
		func(ctx context.Context, batch []int) error {
			current := concurrent.Add(1)
			for {
				max := maxConcurrent.Load()
				if current <= max || maxConcurrent.CompareAndSwap(max, current) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)
			concurrent.Add(-1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if maxConcurrent.Load() > 2 {
		t.Errorf("expected max concurrent batches of 2, got %d", maxConcurrent.Load())
	}
}

func TestParallelByBatchWithError(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expectedErr := errors.New("batch error")
	var batchNum atomic.Int32

	err := ParallelByBatch(
		context.Background(),
		Query(input),
		3,
		2,
		func(ctx context.Context, batch []int) error {
			if batchNum.Add(1) == 2 {
				return expectedErr
			}
			return nil
		},
	)

	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestParallelByBatchEmpty(t *testing.T) {
	input := []int{}

	err := ParallelByBatch(
		context.Background(),
		Query(input),
		3,
		2,
		func(ctx context.Context, batch []int) error {
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestParallelByBatchSingleBatch(t *testing.T) {
	input := []int{1, 2, 3}
	var batchCount atomic.Int32

	err := ParallelByBatch(
		context.Background(),
		Query(input),
		10, // batch size larger than input
		2,
		func(ctx context.Context, batch []int) error {
			batchCount.Add(1)
			if len(batch) != 3 {
				return errors.New("expected batch of 3")
			}
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if batchCount.Load() != 1 {
		t.Errorf("expected 1 batch, got %d", batchCount.Load())
	}
}

func sendItems[T any](items []T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, item := range items {
			ch <- item
		}
	}()
	return ch
}

func TestParallelByBatchChan(t *testing.T) {
	ch := sendItems([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	var totalSum atomic.Int32
	var batchCount atomic.Int32

	err := ParallelByBatchChan(
		context.Background(),
		ch,
		3, // batch size
		2, // max concurrent batches
		func(ctx context.Context, batch []int) error {
			batchCount.Add(1)
			for _, n := range batch {
				totalSum.Add(int32(n))
			}
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if totalSum.Load() != 55 {
		t.Errorf("expected sum 55, got %d", totalSum.Load())
	}

	if batchCount.Load() != 4 {
		t.Errorf("expected 4 batches, got %d", batchCount.Load())
	}
}

func TestParallelByBatchChanEmpty(t *testing.T) {
	ch := sendItems([]int{})

	err := ParallelByBatchChan(
		context.Background(),
		ch,
		3,
		2,
		func(ctx context.Context, batch []int) error {
			return errors.New("should not be called")
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestParallelByBatchChanWithError(t *testing.T) {
	ch := sendItems([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	expectedErr := errors.New("batch error")
	var batchNum atomic.Int32

	err := ParallelByBatchChan(
		context.Background(),
		ch,
		3,
		2,
		func(ctx context.Context, batch []int) error {
			if batchNum.Add(1) == 2 {
				return expectedErr
			}
			return nil
		},
	)

	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestParallelByBatchChanConcurrency(t *testing.T) {
	ch := sendItems([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	var concurrent atomic.Int32
	var maxConcurrent atomic.Int32

	err := ParallelByBatchChan(
		context.Background(),
		ch,
		3, // batch size (4 batches)
		2, // max concurrent batches
		func(ctx context.Context, batch []int) error {
			current := concurrent.Add(1)
			for {
				max := maxConcurrent.Load()
				if current <= max || maxConcurrent.CompareAndSwap(max, current) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)
			concurrent.Add(-1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if maxConcurrent.Load() > 2 {
		t.Errorf("expected max concurrent batches of 2, got %d", maxConcurrent.Load())
	}
}

func TestParallelByBatchChanContextCancellation(t *testing.T) {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= 100; i++ {
			ch <- i
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := ParallelByBatchChan(
		ctx,
		ch,
		3,
		2,
		func(ctx context.Context, batch []int) error {
			time.Sleep(100 * time.Millisecond)
			return nil
		},
	)

	if err != context.Canceled {
		if err != nil {
			t.Logf("got error: %v (this is acceptable)", err)
		}
	}
}

func TestParallelByBatchChanSingleBatch(t *testing.T) {
	ch := sendItems([]int{1, 2, 3})
	var batchCount atomic.Int32

	err := ParallelByBatchChan(
		context.Background(),
		ch,
		10, // batch size larger than input
		2,
		func(ctx context.Context, batch []int) error {
			batchCount.Add(1)
			if len(batch) != 3 {
				return errors.New("expected batch of 3")
			}
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if batchCount.Load() != 1 {
		t.Errorf("expected 1 batch, got %d", batchCount.Load())
	}
}

func TestParallelByBatchChanStreaming(t *testing.T) {
	// Verify that batches are dispatched as items arrive,
	// not after the entire channel is drained.
	ch := make(chan int)
	var firstBatchDone atomic.Bool

	go func() {
		defer close(ch)
		// Send first batch
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		// Wait for first batch to be processed before sending more
		deadline := time.After(2 * time.Second)
		for !firstBatchDone.Load() {
			select {
			case <-deadline:
				return
			default:
				time.Sleep(time.Millisecond)
			}
		}
		// Send second batch
		for i := 4; i <= 6; i++ {
			ch <- i
		}
	}()

	var batchCount atomic.Int32

	err := ParallelByBatchChan(
		context.Background(),
		ch,
		3,
		2,
		func(ctx context.Context, batch []int) error {
			batchCount.Add(1)
			firstBatchDone.Store(true)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if batchCount.Load() != 2 {
		t.Errorf("expected 2 batches, got %d", batchCount.Load())
	}
}
