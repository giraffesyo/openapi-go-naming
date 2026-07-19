package naming_test

import (
	"go/ast"
	"go/token"
	"testing"

	naming "github.com/giraffesyo/openapi-go-naming"
)

func FuzzExported(f *testing.F) {
	for _, seed := range []string{"", "user_id", "123-name", "HTTPServer", "用户", "💥"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		first := naming.Exported(input)
		second := naming.Exported(input)
		if first != second {
			t.Fatalf("non-deterministic conversion: %q then %q", first, second)
		}
		if !token.IsIdentifier(first) || !ast.IsExported(first) {
			t.Fatalf("Exported(%q) produced %q", input, first)
		}
	})
}

func FuzzUnexported(f *testing.F) {
	for _, seed := range []string{"", "user_id", "123-name", "type", "HTTPServer", "用户", "💥"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		first := naming.Unexported(input)
		second := naming.Unexported(input)
		if first != second {
			t.Fatalf("non-deterministic conversion: %q then %q", first, second)
		}
		if !token.IsIdentifier(first) || ast.IsExported(first) {
			t.Fatalf("Unexported(%q) produced %q", input, first)
		}
	})
}
