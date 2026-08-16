// Package testutil provides helpers shared across the package tests.
package testutil

import (
	"os"
	"testing"
)

// WriteFile creates a file at path with a single byte of content, failing the
// test on error.
func WriteFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}

// EqualStrings reports whether a and b contain the same strings in the same
// order.
func EqualStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
