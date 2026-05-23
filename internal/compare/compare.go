// Package compare provides utilities for comparing manifest fields
// against live service configurations to produce structured drift results.
package compare

import (
	"fmt"
	"strings"
)

// Field represents a single comparable field between manifest and live state.
type Field struct {
	Name     string
	Expected string
	Actual   string
}

// Result holds the outcome of comparing a single service.
type Result struct {
	Service string
	Drifted []Field
}

// IsDrifted returns true if any fields differ.
func (r Result) IsDrifted() bool {
	return len(r.Drifted) > 0
}

// Summary returns a human-readable one-line summary of the result.
func (r Result) Summary() string {
	if !r.IsDrifted() {
		return fmt.Sprintf("%s: no drift detected", r.Service)
	}
	fields := make([]string, 0, len(r.Drifted))
	for _, f := range r.Drifted {
		fields = append(fields, f.Name)
	}
	return fmt.Sprintf("%s: drift in [%s]", r.Service, strings.Join(fields, ", "))
}

// Fields compares two maps of string key-value pairs and returns a Result
// containing any fields where expected and actual values differ.
func Fields(service string, expected, actual map[string]string) Result {
	result := Result{Service: service}
	for key, expVal := range expected {
		actVal, ok := actual[key]
		if !ok || actVal != expVal {
			result.Drifted = append(result.Drifted, Field{
				Name:     key,
				Expected: expVal,
				Actual:   actVal,
			})
		}
	}
	return result
}
