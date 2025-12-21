package kk

import (
	"context"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

// Integration tests that combine multiple operations

type User struct {
	ID      int
	Name    string
	Age     int
	Active  bool
	Country string
	Score   int
}

func getTestUsers() []User {
	return []User{
		{ID: 1, Name: "Alice", Age: 30, Active: true, Country: "US", Score: 85},
		{ID: 2, Name: "Bob", Age: 25, Active: false, Country: "UK", Score: 90},
		{ID: 3, Name: "Charlie", Age: 35, Active: true, Country: "US", Score: 75},
		{ID: 4, Name: "Diana", Age: 28, Active: true, Country: "UK", Score: 95},
		{ID: 5, Name: "Eve", Age: 32, Active: false, Country: "US", Score: 80},
		{ID: 6, Name: "Frank", Age: 45, Active: true, Country: "CA", Score: 70},
		{ID: 7, Name: "Grace", Age: 27, Active: true, Country: "UK", Score: 88},
		{ID: 8, Name: "Henry", Age: 38, Active: false, Country: "CA", Score: 92},
		{ID: 9, Name: "Ivy", Age: 29, Active: true, Country: "US", Score: 78},
		{ID: 10, Name: "Jack", Age: 33, Active: true, Country: "UK", Score: 82},
	}
}

func TestIntegration_FilterSortTake(t *testing.T) {
	// Get active users from US, sorted by score descending, take top 2
	users := getTestUsers()
	q := SortedByDesc(
		Query(users).
			Where(func(u User) bool { return u.Active }).
			Where(func(u User) bool { return u.Country == "US" }),
		func(u User) int { return u.Score },
	)
	result := Slice(q.KKQuery.Take(2))

	if len(result) != 2 {
		t.Errorf("expected 2 users, got %d", len(result))
	}

	// Should be Alice (85) and Ivy (78)
	if result[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", result[0].Name)
	}
	if result[1].Name != "Ivy" {
		t.Errorf("expected Ivy, got %s", result[1].Name)
	}
}

func TestIntegration_MapAndAggregate(t *testing.T) {
	// Get total score of active users
	users := getTestUsers()
	totalScore := Sum(
		Mapped(
			Query(users).Where(func(u User) bool { return u.Active }),
			func(u User) int { return u.Score },
		),
		func(n int) int { return n },
	)

	// Active users: Alice(85), Charlie(75), Diana(95), Frank(70), Grace(88), Ivy(78), Jack(82) = 573
	if totalScore != 573 {
		t.Errorf("expected total score 573, got %d", totalScore)
	}
}

func TestIntegration_GroupByCountry(t *testing.T) {
	// Count active users by country using DistinctBy
	users := getTestUsers()

	// Get distinct countries of active users
	activeByCountry := DistinctBy(
		Query(users).Where(func(u User) bool { return u.Active }),
		func(u User) string { return u.Country },
	)
	countries := Slice(activeByCountry)

	// Should have US, UK, CA
	if len(countries) != 3 {
		t.Errorf("expected 3 countries, got %d", len(countries))
	}
}

func TestIntegration_ChunkAndParallel(t *testing.T) {
	// Process users in batches
	users := getTestUsers()
	var processed atomic.Int32

	err := ParallelByBatch(
		context.Background(),
		Query(users).Where(func(u User) bool { return u.Active }),
		3, // batch size
		2, // concurrent batches
		func(ctx context.Context, batch []User) error {
			for range batch {
				processed.Add(1)
			}
			time.Sleep(10 * time.Millisecond) // Simulate work
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 7 active users
	if processed.Load() != 7 {
		t.Errorf("expected 7 processed, got %d", processed.Load())
	}
}

func TestIntegration_TransformCollectAndProcess(t *testing.T) {
	// Transform users to DTOs, collect results in parallel
	users := getTestUsers()

	type UserDTO struct {
		ID    int
		Label string
	}

	results, err := ParallelResult(
		context.Background(),
		Query(users).Where(func(u User) bool { return u.Active }).Take(3),
		2,
		func(ctx context.Context, u User) (UserDTO, error) {
			return UserDTO{
				ID:    u.ID,
				Label: u.Name + " (" + strconv.Itoa(u.Age) + ")",
			}, nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	// Results should maintain order
	if results[0].Label != "Alice (30)" {
		t.Errorf("expected 'Alice (30)', got '%s'", results[0].Label)
	}
}

func TestIntegration_SetOperations(t *testing.T) {
	// Test union of different filtered sets
	users := getTestUsers()

	usUsers := Query(users).Where(func(u User) bool { return u.Country == "US" })
	highScorers := Query(users).Where(func(u User) bool { return u.Score >= 85 })

	// Union: users from US OR with score >= 85
	union := Slice(usUsers.Union(highScorers))

	// US users: Alice, Charlie, Eve, Ivy (4)
	// High scorers: Alice(85), Bob(90), Diana(95), Grace(88), Henry(92) (5)
	// Union (distinct): Alice, Charlie, Eve, Ivy, Bob, Diana, Grace, Henry = 8
	if len(union) != 8 {
		t.Errorf("expected 8 users in union, got %d", len(union))
	}

	// Intersection: users from US AND with score >= 85
	intersection := Slice(usUsers.Intersect(highScorers))

	// Only Alice is from US and has score >= 85
	if len(intersection) != 1 {
		t.Errorf("expected 1 user in intersection, got %d", len(intersection))
	}
	if intersection[0].Name != "Alice" {
		t.Errorf("expected Alice in intersection, got %s", intersection[0].Name)
	}
}

func TestIntegration_ComplexChaining(t *testing.T) {
	// Complex query: active users, skip first, take next 5, sort by age, then by name
	users := getTestUsers()

	q := ThenBy(
		SortedBy(
			Query(users).
				Where(func(u User) bool { return u.Active }).
				Skip(1).
				Take(5),
			func(u User) int { return u.Age },
		),
		func(u User) string { return u.Name },
	)

	result := Slice(q.KKQuery)

	if len(result) != 5 {
		t.Errorf("expected 5 users, got %d", len(result))
	}

	// Should be sorted by age, then by name
	prevAge := 0
	for i, u := range result {
		if u.Age < prevAge {
			t.Errorf("users not sorted by age at index %d", i)
		}
		prevAge = u.Age
	}
}

func TestIntegration_FlatMapWithParallel(t *testing.T) {
	// Each user has multiple tasks, flatten and process in parallel
	users := getTestUsers()[:3] // Just use first 3 users

	type Task struct {
		UserID int
		TaskID int
	}

	// Create tasks for each user
	tasks := Flattened(
		Query(users), func(u User) []Task {
			return []Task{
				{UserID: u.ID, TaskID: 1},
				{UserID: u.ID, TaskID: 2},
			}
		},
	)

	var processed atomic.Int32
	err := Parallel(
		context.Background(), tasks, 3, func(ctx context.Context, t Task) error {
			processed.Add(1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 3 users * 2 tasks = 6 tasks
	if processed.Load() != 6 {
		t.Errorf("expected 6 tasks processed, got %d", processed.Load())
	}
}

func TestIntegration_PerKeyRateLimiting(t *testing.T) {
	// Process users with per-country rate limiting
	users := getTestUsers()
	var processed atomic.Int32

	err := ParallelByKey(
		context.Background(),
		Query(users),
		10, // max total
		1,  // max per country (simulate rate limit)
		func(u User) string { return u.Country },
		func(ctx context.Context, u User) error {
			processed.Add(1)
			time.Sleep(5 * time.Millisecond)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if processed.Load() != 10 {
		t.Errorf("expected 10 processed, got %d", processed.Load())
	}
}

func TestIntegration_AnyAllFirst(t *testing.T) {
	users := getTestUsers()

	// Any active user over 40?
	hasOldActive := Any(
		Query(users), func(u User) bool {
			return u.Active && u.Age > 40
		},
	)
	if !hasOldActive {
		t.Error("expected to find active user over 40")
	}

	// All active users have score > 60?
	allActiveHighScore := All(
		Query(users).Where(func(u User) bool { return u.Active }),
		func(u User) bool { return u.Score > 60 },
	)
	if !allActiveHighScore {
		t.Error("expected all active users to have score > 60")
	}

	// First user from Canada
	canadian, found := First(Query(users).Where(func(u User) bool { return u.Country == "CA" }))
	if !found {
		t.Error("expected to find Canadian user")
	}
	if canadian.Name != "Frank" {
		t.Errorf("expected Frank, got %s", canadian.Name)
	}
}

func TestIntegration_ConcatAndDistinct(t *testing.T) {
	// Combine users from two sources and get distinct countries
	users1 := []User{
		{ID: 1, Country: "US"},
		{ID: 2, Country: "UK"},
	}
	users2 := []User{
		{ID: 3, Country: "UK"},
		{ID: 4, Country: "CA"},
	}

	combined := Query(users1).Concat(Query(users2))
	countries := Mapped(combined, func(u User) string { return u.Country })
	distinctCountries := Slice(countries.Distinct())

	if len(distinctCountries) != 3 {
		t.Errorf("expected 3 distinct countries, got %d", len(distinctCountries))
	}
}

func TestIntegration_PipelineExample(t *testing.T) {
	// Simulate the README example: filter users, take 100, process in parallel
	users := getTestUsers()

	var emailsSent atomic.Int32

	// Build a query with method chaining
	q := Query(users).
		Where(func(u User) bool { return u.Active }).
		Take(100)

	// Execute in parallel
	err := Parallel(
		context.Background(), q, 10, func(ctx context.Context, u User) error {
			emailsSent.Add(1)
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 7 active users
	if emailsSent.Load() != 7 {
		t.Errorf("expected 7 emails sent, got %d", emailsSent.Load())
	}
}
