package naming_test

import (
	"testing"

	naming "github.com/giraffesyo/openapi-go-naming"
)

var benchmarkResult string

func BenchmarkExported(b *testing.B) {
	for _, input := range []string{
		"user_id",
		"HTTPServer_user_id",
		"com.example.api.resource_identifier",
		"über_status_用户_id",
	} {
		b.Run(input, func(b *testing.B) {
			for b.Loop() {
				benchmarkResult = naming.Exported(input)
			}
		})
	}
}

func BenchmarkUnexported(b *testing.B) {
	for b.Loop() {
		benchmarkResult = naming.Unexported("HTTPServer_user_id")
	}
}

func BenchmarkScopeUnique(b *testing.B) {
	for b.Loop() {
		scope := naming.NewScope()
		benchmarkResult = scope.Unique("Identifier")
	}
}
