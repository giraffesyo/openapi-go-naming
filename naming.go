// Package naming converts external identifiers into idiomatic Go identifiers.
package naming

import (
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	exportedFallback   = "Unknown"
	unexportedFallback = "unknown"
)

var defaultInitialisms = map[string]string{
	"acl":   "ACL",
	"api":   "API",
	"ascii": "ASCII",
	"cpu":   "CPU",
	"css":   "CSS",
	"dns":   "DNS",
	"eof":   "EOF",
	"guid":  "GUID",
	"html":  "HTML",
	"http":  "HTTP",
	"https": "HTTPS",
	"id":    "ID",
	"ip":    "IP",
	"io":    "IO",
	"json":  "JSON",
	"qps":   "QPS",
	"ram":   "RAM",
	"rpc":   "RPC",
	"sla":   "SLA",
	"smtp":  "SMTP",
	"sql":   "SQL",
	"ssh":   "SSH",
	"ssl":   "SSL",
	"tcp":   "TCP",
	"tls":   "TLS",
	"ttl":   "TTL",
	"udp":   "UDP",
	"ui":    "UI",
	"uid":   "UID",
	"uri":   "URI",
	"url":   "URL",
	"utf8":  "UTF8",
	"uuid":  "UUID",
	"vm":    "VM",
	"xml":   "XML",
	"xmpp":  "XMPP",
	"xsrf":  "XSRF",
	"xss":   "XSS",
}

var standard = &Converter{initialisms: defaultInitialisms}

// Converter converts identifiers using an immutable set of initialisms. It is
// safe for concurrent use.
type Converter struct {
	initialisms map[string]string
}

// NewConverter returns a converter containing the standard Go initialisms plus
// each additional canonical spelling. Matching is case-insensitive.
//
// For example, NewConverter("GPU") converts "gpu_limit" to "GPULimit".
func NewConverter(initialisms ...string) *Converter {
	configured := make(map[string]string, len(defaultInitialisms)+len(initialisms))
	for key, value := range defaultInitialisms {
		configured[key] = value
	}
	for _, initialism := range initialisms {
		if isWord(initialism) {
			configured[strings.ToLower(initialism)] = strings.ToUpper(initialism)
		}
	}
	return &Converter{initialisms: configured}
}

// Exported converts input to an exported Go identifier in PascalCase.
func Exported(input string) string {
	return standard.Exported(input)
}

// Unexported converts input to an unexported Go identifier in camelCase.
func Unexported(input string) string {
	return standard.Unexported(input)
}

// Exported converts input to an exported Go identifier in PascalCase.
func (c *Converter) Exported(input string) string {
	return c.convert(input, true)
}

// Unexported converts input to an unexported Go identifier in camelCase.
func (c *Converter) Unexported(input string) string {
	return c.convert(input, false)
}

func (c *Converter) convert(input string, exported bool) string {
	if c == nil {
		c = standard
	}

	var output strings.Builder
	output.Grow(len(input) + 1)

	words := c.writeWords(&output, input, exported)
	if words == 0 {
		if exported {
			return exportedFallback
		}
		return unexportedFallback
	}

	result := output.String()
	if !exported && token.Lookup(result).IsKeyword() {
		output.WriteByte('_')
		result = output.String()
	}
	return result
}

func (c *Converter) writeWords(output *strings.Builder, input string, exported bool) int {
	wordStart := -1
	previousStart := -1
	var previous rune
	words := 0

	for offset := 0; offset < len(input); {
		current, size := utf8.DecodeRuneInString(input[offset:])
		if !unicode.IsLetter(current) && !unicode.IsDigit(current) {
			if wordStart >= 0 {
				c.writeWord(output, input[wordStart:offset], exported, words == 0)
				words++
				wordStart = -1
			}
			offset += size
			continue
		}

		if wordStart < 0 {
			wordStart = offset
		} else {
			switch {
			case unicode.IsUpper(current) && (unicode.IsLower(previous) || unicode.IsDigit(previous)):
				c.writeWord(output, input[wordStart:offset], exported, words == 0)
				words++
				wordStart = offset
			case unicode.IsLower(current) && unicode.IsUpper(previous) && previousStart > wordStart:
				c.writeWord(output, input[wordStart:previousStart], exported, words == 0)
				words++
				wordStart = previousStart
			}
		}

		previous = current
		previousStart = offset
		offset += size
	}

	if wordStart >= 0 {
		c.writeWord(output, input[wordStart:], exported, words == 0)
		words++
	}
	return words
}

func (c *Converter) writeWord(output *strings.Builder, word string, exported, first bool) {
	firstRune, _ := utf8.DecodeRuneInString(word)
	if first {
		switch {
		case exported && unicode.IsDigit(firstRune):
			output.WriteByte('N')
		case exported && !unicode.IsUpper(unicode.ToUpper(firstRune)):
			output.WriteByte('X')
		case !exported && unicode.IsDigit(firstRune):
			output.WriteByte('n')
		case !exported && unicode.IsUpper(unicode.ToLower(firstRune)):
			output.WriteByte('x')
		}
	}

	if initialism, ok := c.lookupInitialism(word); ok {
		if first && !exported {
			writeLower(output, initialism)
		} else {
			output.WriteString(initialism)
		}
		return
	}

	if first && !exported {
		writeLower(output, word)
		return
	}
	writeTitle(output, word)
}

func (c *Converter) lookupInitialism(word string) (string, bool) {
	initialisms := c.initialisms
	if initialisms == nil {
		initialisms = defaultInitialisms
	}

	if len(word) <= 64 {
		var folded [64]byte
		for i := range len(word) {
			char := word[i]
			if char >= utf8.RuneSelf {
				return lookupFolded(initialisms, word)
			}
			if char >= 'A' && char <= 'Z' {
				char += 'a' - 'A'
			}
			folded[i] = char
		}
		initialism, ok := initialisms[string(folded[:len(word)])]
		return initialism, ok
	}
	return lookupFolded(initialisms, word)
}

func lookupFolded(initialisms map[string]string, word string) (string, bool) {
	initialism, ok := initialisms[strings.ToLower(word)]
	return initialism, ok
}

func writeLower(output *strings.Builder, word string) {
	for _, current := range word {
		output.WriteRune(unicode.ToLower(current))
	}
}

func writeTitle(output *strings.Builder, word string) {
	first := true
	for _, current := range word {
		if first {
			output.WriteRune(unicode.ToUpper(current))
			first = false
			continue
		}
		output.WriteRune(unicode.ToLower(current))
	}
}

func isWord(value string) bool {
	if value == "" {
		return false
	}
	for _, current := range value {
		if !unicode.IsLetter(current) && !unicode.IsDigit(current) {
			return false
		}
	}
	return true
}
