# openapi-go-naming

[![Go Reference](https://pkg.go.dev/badge/github.com/giraffesyo/openapi-go-naming.svg)](https://pkg.go.dev/github.com/giraffesyo/openapi-go-naming)
[![CI](https://github.com/giraffesyo/openapi-go-naming/actions/workflows/ci.yml/badge.svg?branch=canary)](https://github.com/giraffesyo/openapi-go-naming/actions/workflows/ci.yml)

Fast, deterministic conversion of OpenAPI identifiers into idiomatic,
collision-safe Go names.

**Zero dependencies.** Built entirely on Go's standard library.
Requires Go 1.25 or newer.

    import naming "github.com/giraffesyo/openapi-go-naming"

    typeName := naming.Exported("user_profile")   // UserProfile
    fieldName := naming.Exported("account_id")    // AccountID
    paramName := naming.Unexported("account_id")  // accountID

Every conversion returns a valid Go identifier. Exported names are guaranteed
to be exported, unexported names avoid Go keywords, Unicode is supported, and
identifiers beginning with digits receive a stable prefix.

| Input | `Exported` | `Unexported` |
|---|---|---|
| `user_id` | `UserID` | `userID` |
| `HTTPServerURL` | `HTTPServerURL` | `httpServerURL` |
| `some.dotted-name` | `SomeDottedName` | `someDottedName` |
| `123-response` | `N123Response` | `n123Response` |
| `type` | `Type` | `type_` |
| `用户_id` | `X用户ID` | `用户ID` |

## Initialisms

The default converter recognizes conventional Go initialisms such as `API`,
`HTTP`, `ID`, `JSON`, `URL`, and `UUID`. Add project-specific initialisms with
an immutable converter:

    converter := naming.NewConverter("GPU")

    converter.Exported("gpu_limit")   // GPULimit
    converter.Unexported("gpu_limit") // gpuLimit

Converters are safe for concurrent use.

## Collision handling

`Scope` allocates unique names without changing the first occurrence:

    scope := naming.NewScope("User2")

    scope.Unique("User") // User
    scope.Unique("User") // User3
    scope.Unique("User") // User4

The zero value is ready to use, and all operations are safe for concurrent
use.

## Performance

Conversion uses a single-pass tokenizer and writes directly into one output
buffer. Common ASCII conversions allocate only the returned string. Run the
included benchmarks with:

    go test -bench . -benchmem -run '^$' ./...

## License

MIT.
