package group_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/group"
)

// Integration tests exercise Apply with realistic multi-service data
// to confirm grouping invariants hold across all supported dimensions.

func buildLargeResultSet() []drift.DriftResult {
	services := []string{"auth", "payments", "gateway"}
	fields := []string{"image", "replicas", "env.LOG_LEVEL"}
	var out []drift.DriftResult
	for _, svc := range services {
		for i, f := range fields {
			out = append(out, drift.DriftResult{
				Service:  svc,
				Field:    f,
				Drifted:  i%2 == 0,
				Expected: "expected",
				Actual:   "actual",
			})
		}
	}
	return out
}

func TestGroupByService_TotalCountPreserved(t *testing.T) {
	results := buildLargeResultSet()
	res, err := group.Apply(results, group.ByService)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	total := 0
	for _, v := range res.Groups {
		total += len(v)
	}
	if total != len(results) {
		t.Errorf("count mismatch: got %d, want %d", total, len(results))
	}
}

func TestGroupByField_TotalCountPreserved(t *testing.T) {
	results := buildLargeResultSet()
	res, err := group.Apply(results, group.ByField)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	total := 0
	for _, v := range res.Groups {
		total += len(v)
	}
	if total != len(results) {
		t.Errorf("count mismatch: got %d, want %d", total, len(results))
	}
}

func TestGroupBySeverity_ExhaustivePartition(t *testing.T) {
	results := buildLargeResultSet()
	res, err := group.Apply(results, group.BySeverity)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range res.Keys() {
		if key != "drifted" && key != "clean" {
			t.Errorf("unexpected severity key: %q", key)
		}
	}
}
