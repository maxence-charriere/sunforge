//go:generate go run gen.go
//go:generate go fmt

// Package vi provides reusable UI components and styles for go-app applications.
// It includes components for page layouts, navigation, media, and loading states.
package vi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	styles string
)

// An app.Library implementation.
type Library struct{}

func (l Library) Styles() (string, string) {
	return "/vi.css", styles
}

// PathAtIndex returns the path part at the given index.
func PathAtIndex(path string, idx int) string {
	if path == "" {
		return ""
	}

	s := strings.Split(path[1:], "/")
	if idx >= len(s) {
		return ""
	}

	return s[idx]
}

// PrettyPrint prints the given value in a JSON format.
func PrettyPrint(v any) {
	b, _ := json.MarshalIndent(v, "", "    ")
	fmt.Println((string(b)))
}

// Slugify converts the given string into a simplified, URL-friendly slug.
// It transforms all characters to lowercase, replaces spaces with hyphens,
// and converts newlines into underscores. Useful for creating identifiers
// from titles or names that can be used safely in URLs.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return strings.ReplaceAll(s, "\n", "_")
}

// Unslugify reverses the transformations applied by Slugify and restores
// a human-readable title. It replaces certain patterns and symbols back
// to spaces or newlines, capitalizes words according to U.S. English title
// conventions, and adjusts possessive forms like "'s".
func Unslugify(s string) string {
	s = strings.ReplaceAll(s, "-s-", "'s ")
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", "\n")
	s = Title(s)
	s = strings.ReplaceAll(s, "'S", "'s")
	return s
}

// Title returns the input string with each word capitalized according to
// American English title casing rules.
func Title(s string) string {
	return cases.Title(language.AmericanEnglish).String(s)
}

// URLWithFragment returns a copy of the provided URL with its Fragment field
// set to the given value. The original URL is not modified.
//
// The returned URL preserves all components of the input URL (scheme, host,
// path, query, and user info) but replaces the fragment.
func URLWithFragment(u *url.URL, fragment string) *url.URL {
	if u == nil {
		return nil
	}
	modified := *u
	modified.Fragment = fragment
	return &modified
}

// Chunk slices s into subslices of at most size n and returns them as [][]T.
//
// The subslices share the same underlying array as s. The last chunk may have
// fewer than n elements if len(s) is not a multiple of n.
func Chunk[T any](s []T, n int) [][]T {
	if n <= 0 || len(s) == 0 {
		return nil
	}

	chunks := make([][]T, 0, (len(s)+n-1)/n)
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
