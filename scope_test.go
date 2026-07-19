package naming_test

import (
	"fmt"
	"sync"
	"testing"

	naming "github.com/giraffesyo/openapi-go-naming"
)

func TestScopeUnique(t *testing.T) {
	t.Parallel()

	scope := naming.NewScope()
	for index, want := range []string{"User", "User2", "User3"} {
		if got := scope.Unique("User"); got != want {
			t.Fatalf("Unique(User) call %d = %q, want %q", index+1, got, want)
		}
	}
	if got := scope.Unique("Pet"); got != "Pet" {
		t.Fatalf("Unique(Pet) = %q, want Pet", got)
	}
}

func TestScopeSkipsReservedSuffixes(t *testing.T) {
	t.Parallel()

	scope := naming.NewScope("User2", "User4")
	if got := scope.Unique("User"); got != "User" {
		t.Fatalf("first Unique(User) = %q, want User", got)
	}
	if got := scope.Unique("User"); got != "User3" {
		t.Fatalf("second Unique(User) = %q, want User3", got)
	}
	if got := scope.Unique("User"); got != "User5" {
		t.Fatalf("third Unique(User) = %q, want User5", got)
	}
}

func TestScopeZeroValue(t *testing.T) {
	t.Parallel()

	var scope naming.Scope
	if got := scope.Unique("Value"); got != "Value" {
		t.Fatalf("Unique(Value) = %q, want Value", got)
	}
	scope.Reserve("Other")
	if got := scope.Unique("Other"); got != "Other2" {
		t.Fatalf("Unique(Other) = %q, want Other2", got)
	}
}

func TestScopeConcurrentUse(t *testing.T) {
	t.Parallel()

	const workers = 128
	scope := naming.NewScope()
	results := make(chan string, workers)

	var group sync.WaitGroup
	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			results <- scope.Unique("Name")
		}()
	}
	group.Wait()
	close(results)

	seen := make(map[string]struct{}, workers)
	for result := range results {
		if _, exists := seen[result]; exists {
			t.Fatalf("duplicate result %q", result)
		}
		seen[result] = struct{}{}
	}
	if len(seen) != workers {
		t.Fatalf("got %d unique results, want %d", len(seen), workers)
	}
}

func TestScopeFirstUseAllocations(t *testing.T) {
	if testing.Short() {
		t.Skip("allocation assertion")
	}

	allocations := testing.AllocsPerRun(1_000, func() {
		scope := naming.NewScope()
		if got := scope.Unique("Name"); got != "Name" {
			panic(fmt.Sprintf("unexpected result %q", got))
		}
	})
	if allocations > 3 {
		t.Fatalf("first scope use allocated %.2f times, want at most 3", allocations)
	}
}
