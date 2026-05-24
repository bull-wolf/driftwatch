package rollup_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/rollup"
)

func sampleResults() map[string][]drift.Drift {
	return map[string][]drift.Drift{
		"auth-service": {
			{Field: "image", Expected: "v1", Actual: "v2"},
			{Field: "replicas", Expected: "2", Actual: "3"},
		},
		"payment-service": {
			{Field: "env", Expected: "prod", Actual: "staging"},
		},
		"clean-service": {},
	}
}

func TestAggregate_TotalServices(t *testing.T) {
	r := rollup.Aggregate(sampleResults())
	if r.TotalServices != 3 {
		t.Errorf("expected 3 total services, got %d", r.TotalServices)
	}
}

func TestAggregate_DriftedAndClean(t *testing.T) {
	r := rollup.Aggregate(sampleResults())
	if r.DriftedCount != 2 {
		t.Errorf("expected 2 drifted, got %d", r.DriftedCount)
	}
	if r.CleanCount != 1 {
		t.Errorf("expected 1 clean, got %d", r.CleanCount)
	}
}

func TestAggregate_ByField(t *testing.T) {
	r := rollup.Aggregate(sampleResults())
	if r.ByField["image"] != 1 {
		t.Errorf("expected image count 1, got %d", r.ByField["image"])
	}
	if r.ByField["env"] != 1 {
		t.Errorf("expected env count 1, got %d", r.ByField["env"])
	}
}

func TestAggregate_BySeverity(t *testing.T) {
	r := rollup.Aggregate(sampleResults())
	if r.BySeverity["high"] != 2 {
		t.Errorf("expected 2 high severity, got %d", r.BySeverity["high"])
	}
	if r.BySeverity["medium"] != 1 {
		t.Errorf("expected 1 medium severity, got %d", r.BySeverity["medium"])
	}
}

func TestAggregate_TopDriftFields(t *testing.T) {
	r := rollup.Aggregate(sampleResults())
	if len(r.TopDriftFields) == 0 {
		t.Fatal("expected at least one top drift field")
	}
	// image and replicas both appear once; env once — all tied, sorted alpha
	found := map[string]bool{}
	for _, f := range r.TopDriftFields {
		found[f] = true
	}
	if !found["env"] && !found["image"] && !found["replicas"] {
		t.Errorf("unexpected top fields: %v", r.TopDriftFields)
	}
}

func TestAggregate_EmptyResults(t *testing.T) {
	r := rollup.Aggregate(map[string][]drift.Drift{})
	if r.TotalServices != 0 || r.DriftedCount != 0 || r.CleanCount != 0 {
		t.Errorf("expected all zeros for empty input, got %+v", r)
	}
}

func TestSummary_Format(t *testing.T) {
	r := rollup.Aggregate(sampleResults())
	s := rollup.Summary(r)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
	expected := "2/3 services drifted, 3 total drift(s)"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
