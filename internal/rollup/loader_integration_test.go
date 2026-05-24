package rollup_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/rollup"
)

func buildLargeResultSet() map[string][]drift.Drift {
	services := []string{
		"svc-a", "svc-b", "svc-c", "svc-d", "svc-e",
	}
	results := make(map[string][]drift.Drift, len(services))
	for i, svc := range services {
		if i%2 == 0 {
			results[svc] = []drift.Drift{
				{Field: "image", Expected: "v1", Actual: "v2"},
			}
		} else {
			results[svc] = []drift.Drift{}
		}
	}
	return results
}

func TestAggregate_LargeSet_TotalPreserved(t *testing.T) {
	results := buildLargeResultSet()
	r := rollup.Aggregate(results)
	if r.TotalServices != 5 {
		t.Errorf("expected 5 services, got %d", r.TotalServices)
	}
	if r.DriftedCount+r.CleanCount != r.TotalServices {
		t.Errorf("drifted + clean should equal total: %d + %d != %d",
			r.DriftedCount, r.CleanCount, r.TotalServices)
	}
}

func TestAggregate_LargeSet_SeverityHigh(t *testing.T) {
	results := buildLargeResultSet()
	r := rollup.Aggregate(results)
	// 3 services with image drift → 3 high severity
	if r.BySeverity["high"] != 3 {
		t.Errorf("expected 3 high severity entries, got %d", r.BySeverity["high"])
	}
}

func TestAggregate_LargeSet_TopFieldIsImage(t *testing.T) {
	results := buildLargeResultSet()
	r := rollup.Aggregate(results)
	if len(r.TopDriftFields) == 0 {
		t.Fatal("expected top drift fields to be populated")
	}
	if r.TopDriftFields[0] != "image" {
		t.Errorf("expected top field to be 'image', got %q", r.TopDriftFields[0])
	}
}

func TestSummary_LargeSet_NonEmpty(t *testing.T) {
	results := buildLargeResultSet()
	r := rollup.Aggregate(results)
	s := rollup.Summary(r)
	if s == "" {
		t.Fatal("expected non-empty summary for large result set")
	}
}
