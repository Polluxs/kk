package kk

import (
	"testing"
)

func TestGroupedBy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	groups := Slice(GroupedBy(From(input), func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	}))

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	// First group encountered should be "odd" (1 is first)
	if groups[0].Key != "odd" {
		t.Errorf("expected first group key 'odd', got '%s'", groups[0].Key)
	}
	if len(groups[0].Items) != 3 {
		t.Errorf("expected 3 odd items, got %d", len(groups[0].Items))
	}

	// Second group should be "even"
	if groups[1].Key != "even" {
		t.Errorf("expected second group key 'even', got '%s'", groups[1].Key)
	}
	if len(groups[1].Items) != 3 {
		t.Errorf("expected 3 even items, got %d", len(groups[1].Items))
	}
}

func TestGroupByEmpty(t *testing.T) {
	input := []int{}
	groups := Slice(GroupedBy(From(input), func(n int) int { return n % 2 }))

	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %v", groups)
	}
}

func TestGroupBySingleGroup(t *testing.T) {
	input := []int{2, 4, 6, 8}
	groups := Slice(GroupedBy(From(input), func(n int) string { return "even" }))

	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}

	if groups[0].Key != "even" {
		t.Errorf("expected key 'even', got '%s'", groups[0].Key)
	}

	if len(groups[0].Items) != 4 {
		t.Errorf("expected 4 items, got %d", len(groups[0].Items))
	}
}

func TestGroupByPreservesOrder(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	groups := Slice(GroupedBy(From(input), func(n int) int { return n % 2 }))

	// First group should be "1" (odd), second should be "0" (even)
	if groups[0].Key != 1 {
		t.Errorf("expected first group key 1, got %d", groups[0].Key)
	}
	if groups[1].Key != 0 {
		t.Errorf("expected second group key 0, got %d", groups[1].Key)
	}

	// Items within groups should preserve order
	oddItems := groups[0].Items
	if oddItems[0] != 1 || oddItems[1] != 3 || oddItems[2] != 5 {
		t.Errorf("odd items not in order: %v", oddItems)
	}

	evenItems := groups[1].Items
	if evenItems[0] != 2 || evenItems[1] != 4 {
		t.Errorf("even items not in order: %v", evenItems)
	}
}

type TestPerson struct {
	Name    string
	Age     int
	Country string
}

func TestGroupByStruct(t *testing.T) {
	people := []TestPerson{
		{Name: "Alice", Age: 30, Country: "US"},
		{Name: "Bob", Age: 25, Country: "UK"},
		{Name: "Charlie", Age: 35, Country: "US"},
		{Name: "Diana", Age: 28, Country: "UK"},
		{Name: "Eve", Age: 32, Country: "CA"},
	}

	groups := Slice(GroupedBy(From(people), func(p TestPerson) string {
		return p.Country
	}))

	if len(groups) != 3 {
		t.Errorf("expected 3 country groups, got %d", len(groups))
	}

	// Verify group contents
	countryMap := make(map[string][]TestPerson)
	for _, g := range groups {
		countryMap[g.Key] = g.Items
	}

	if len(countryMap["US"]) != 2 {
		t.Errorf("expected 2 US people, got %d", len(countryMap["US"]))
	}
	if len(countryMap["UK"]) != 2 {
		t.Errorf("expected 2 UK people, got %d", len(countryMap["UK"]))
	}
	if len(countryMap["CA"]) != 1 {
		t.Errorf("expected 1 CA person, got %d", len(countryMap["CA"]))
	}
}

func TestGroupByWithChaining(t *testing.T) {
	people := []TestPerson{
		{Name: "Alice", Age: 30, Country: "US"},
		{Name: "Bob", Age: 25, Country: "UK"},
		{Name: "Charlie", Age: 35, Country: "US"},
		{Name: "Diana", Age: 28, Country: "UK"},
		{Name: "Eve", Age: 32, Country: "CA"},
	}

	// Filter first, then group
	groups := Slice(GroupedBy(
		From(people).Where(func(p TestPerson) bool { return p.Age >= 30 }),
		func(p TestPerson) string { return p.Country },
	))

	// Only Alice (30, US), Charlie (35, US), Eve (32, CA) pass the filter
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	countryMap := make(map[string]int)
	for _, g := range groups {
		countryMap[g.Key] = len(g.Items)
	}

	if countryMap["US"] != 2 {
		t.Errorf("expected 2 US people age >= 30, got %d", countryMap["US"])
	}
	if countryMap["CA"] != 1 {
		t.Errorf("expected 1 CA person age >= 30, got %d", countryMap["CA"])
	}
}

func TestGroupByStrings(t *testing.T) {
	words := []string{"apple", "banana", "apricot", "blueberry", "avocado"}
	groups := Slice(GroupedBy(From(words), func(s string) byte {
		return s[0] // group by first letter
	}))

	if len(groups) != 2 {
		t.Errorf("expected 2 groups (a and b), got %d", len(groups))
	}

	letterMap := make(map[byte][]string)
	for _, g := range groups {
		letterMap[g.Key] = g.Items
	}

	if len(letterMap['a']) != 3 {
		t.Errorf("expected 3 'a' words, got %d", len(letterMap['a']))
	}
	if len(letterMap['b']) != 2 {
		t.Errorf("expected 2 'b' words, got %d", len(letterMap['b']))
	}
}

func TestGroupByWithMap(t *testing.T) {
	people := []TestPerson{
		{Name: "Alice", Age: 30, Country: "US"},
		{Name: "Bob", Age: 25, Country: "UK"},
		{Name: "Charlie", Age: 35, Country: "US"},
	}

	// Group by country, then map to count per country
	groups := GroupedBy(From(people), func(p TestPerson) string {
		return p.Country
	})

	type CountryCount struct {
		Country string
		Count   int
	}

	counts := Slice(Mapped(groups, func(g Group[string, TestPerson]) CountryCount {
		return CountryCount{Country: g.Key, Count: len(g.Items)}
	}))

	if len(counts) != 2 {
		t.Errorf("expected 2 country counts, got %d", len(counts))
	}
}

func TestGroupByAllSameKey(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	groups := Slice(GroupedBy(From(input), func(n int) string { return "all" }))

	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}

	if len(groups[0].Items) != 5 {
		t.Errorf("expected 5 items in group, got %d", len(groups[0].Items))
	}
}

func TestGroupByAllUniqueKeys(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	groups := Slice(GroupedBy(From(input), func(n int) int { return n }))

	if len(groups) != 5 {
		t.Errorf("expected 5 groups (one per item), got %d", len(groups))
	}

	for i, g := range groups {
		if len(g.Items) != 1 {
			t.Errorf("expected 1 item in group %d, got %d", i, len(g.Items))
		}
	}
}
