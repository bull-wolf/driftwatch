// Package group provides utilities for grouping drift results by
// a specified dimension such as service, field, or severity.
package group

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// By defines the dimension along which results are grouped.
type By string

const (
	ByService  By = "service"
	ByField    By = "field"
	BySeverity By = "severity"
)

// Result holds a grouped set of drift results keyed by the group value.
type Result struct {
	Dimension By
	Groups    map[string][]drift.DriftResult
}

// Apply groups the provided drift results by the given dimension.
// Returns an error if the dimension is not recognised.
func Apply(results []drift.DriftResult, by By) (*Result, error) {
	switch by {
	case ByService:
		return &Result{Dimension: by, Groups: groupBy(results, func(r drift.DriftResult) string {
			return r.Service
		})}, nil
	case ByField:
		return &Result{Dimension: by, Groups: groupBy(results, func(r drift.DriftResult) string {
			return r.Field
		})}, nil
	case BySeverity:
		return &Result{Dimension: by, Groups: groupBy(results, func(r drift.DriftResult) string {
			if r.Drifted {
				return "drifted"
			}
			return "clean"
		})}, nil
	default:
		return nil, fmt.Errorf("group: unknown dimension %q", by)
	}
}

// Keys returns the sorted list of group keys in the result.
func (r *Result) Keys() []string {
	keys := make([]string, 0, len(r.Groups))
	for k := range r.Groups {
		keys = append(keys, k)
	}
	return keys
}

func groupBy(results []drift.DriftResult, key func(drift.DriftResult) string) map[string][]drift.DriftResult {
	m := make(map[string][]drift.DriftResult)
	for _, r := range results {
		k := key(r)
		m[k] = append(m[k], r)
	}
	return m
}
