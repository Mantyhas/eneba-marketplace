package search

import "strings"

// Pattern escapes LIKE wildcards so user input is always literal.
func Pattern(q string) string {
	q = strings.ToLower(strings.TrimSpace(q))
	q = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(q)
	return "%" + q + "%"
}
