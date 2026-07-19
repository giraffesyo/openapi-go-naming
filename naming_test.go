package naming_test

import (
	"go/ast"
	"go/token"
	"testing"

	naming "github.com/giraffesyo/openapi-go-naming"
)

func TestExported(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"":                   "Unknown",
		"---":                "Unknown",
		"user_id":            "UserID",
		"pet-name":           "PetName",
		"some.dotted.name":   "SomeDottedName",
		"X-Request-ID":       "XRequestID",
		"HTTPServerURL":      "HTTPServerURL",
		"userID":             "UserID",
		"v2_api":             "V2API",
		"get200Response":     "Get200Response",
		"123-response":       "N123Response",
		"type":               "Type",
		"über_status":        "ÜberStatus",
		"用户_id":              "X用户ID",
		string([]byte{0xff}): "Unknown",
		"json_api_http_url":  "JSONAPIHTTPURL",
		"identifier_with_💥":  "IdentifierWith",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if got := naming.Exported(input); got != want {
				t.Fatalf("Exported(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

func TestUnexported(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"":                 "unknown",
		"---":              "unknown",
		"user_id":          "userID",
		"PetStatus":        "petStatus",
		"X-Request-ID":     "xRequestID",
		"HTTPSPort":        "httpsPort",
		"some.dotted.name": "someDottedName",
		"123-response":     "n123Response",
		"type":             "type_",
		"range":            "range_",
		"über_status":      "überStatus",
		"用户_id":            "用户ID",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if got := naming.Unexported(input); got != want {
				t.Fatalf("Unexported(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

func TestConverterAdditionalInitialism(t *testing.T) {
	t.Parallel()

	converter := naming.NewConverter("GPU")
	if got := converter.Exported("gpu_limit"); got != "GPULimit" {
		t.Fatalf("Exported(gpu_limit) = %q, want GPULimit", got)
	}
	if got := converter.Unexported("gpu_limit"); got != "gpuLimit" {
		t.Fatalf("Unexported(gpu_limit) = %q, want gpuLimit", got)
	}
	if got := naming.Exported("gpu_limit"); got != "GpuLimit" {
		t.Fatalf("package converter was mutated: Exported(gpu_limit) = %q", got)
	}

	lowercase := naming.NewConverter("gpu")
	if got := lowercase.Exported("gpu_limit"); got != "GPULimit" {
		t.Fatalf("lowercase initialism produced %q, want GPULimit", got)
	}
}

func TestNilConverterUsesDefaults(t *testing.T) {
	t.Parallel()

	var converter *naming.Converter
	if got := converter.Exported("user_id"); got != "UserID" {
		t.Fatalf("Exported(user_id) = %q, want UserID", got)
	}
}

func TestConversionsAreGoIdentifiers(t *testing.T) {
	t.Parallel()

	inputs := []string{
		"", "---", "123", "123-name", "type", "user_id", "HTTPServer",
		"Δelta-value", "用户", "💥", string([]byte{0xff, 'a'}),
	}
	for _, input := range inputs {
		exported := naming.Exported(input)
		if !token.IsIdentifier(exported) {
			t.Errorf("Exported(%q) produced invalid identifier %q", input, exported)
		}
		if !ast.IsExported(exported) {
			t.Errorf("Exported(%q) produced non-exported identifier %q", input, exported)
		}

		unexported := naming.Unexported(input)
		if !token.IsIdentifier(unexported) {
			t.Errorf("Unexported(%q) produced invalid identifier %q", input, unexported)
		}
		if ast.IsExported(unexported) {
			t.Errorf("Unexported(%q) produced exported identifier %q", input, unexported)
		}
	}
}

func TestConversionAllocations(t *testing.T) {
	if testing.Short() {
		t.Skip("allocation assertion")
	}

	var result string
	allocations := testing.AllocsPerRun(1_000, func() {
		result = naming.Exported("HTTPServer_user_id")
	})
	if result != "HTTPServerUserID" {
		t.Fatalf("unexpected result %q", result)
	}
	if allocations > 1 {
		t.Fatalf("Exported allocated %.2f times, want at most 1", allocations)
	}
}
