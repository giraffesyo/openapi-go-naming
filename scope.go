package naming

import (
	"strconv"
	"strings"
	"sync"
)

// Scope allocates unique identifiers within a namespace. Its zero value is
// ready to use, and its methods are safe for concurrent use. A Scope must not
// be copied after first use.
type Scope struct {
	mu      sync.Mutex
	entries map[string]uint64
}

// NewScope returns a scope with reserved identifiers already marked as used.
func NewScope(reserved ...string) *Scope {
	scope := &Scope{entries: make(map[string]uint64, len(reserved)+1)}
	for _, identifier := range reserved {
		scope.entries[identifier] = 2
	}
	return scope
}

// Reserve marks identifiers as unavailable. Reserving an identifier more than
// once has no effect.
func (s *Scope) Reserve(identifiers ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.initialize(len(identifiers) + 1)
	for _, identifier := range identifiers {
		if s.entries[identifier] < 2 {
			s.entries[identifier] = 2
		}
	}
}

// Unique returns identifier when it is available. Otherwise, it appends the
// lowest available decimal suffix beginning with 2.
func (s *Scope) Unique(identifier string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.initialize(1)
	suffix, exists := s.entries[identifier]
	if !exists {
		s.entries[identifier] = 2
		return identifier
	}

	if suffix < 2 {
		suffix = 2
	}
	for {
		candidate := suffixed(identifier, suffix)
		if _, candidateExists := s.entries[candidate]; !candidateExists {
			s.entries[candidate] = 2
			s.entries[identifier] = suffix + 1
			return candidate
		}
		suffix++
	}
}

func (s *Scope) initialize(capacity int) {
	if s.entries == nil {
		s.entries = make(map[string]uint64, capacity)
	}
}

func suffixed(identifier string, suffix uint64) string {
	var digits [20]byte
	encoded := strconv.AppendUint(digits[:0], suffix, 10)

	var result strings.Builder
	result.Grow(len(identifier) + len(encoded))
	result.WriteString(identifier)
	result.Write(encoded)
	return result.String()
}
