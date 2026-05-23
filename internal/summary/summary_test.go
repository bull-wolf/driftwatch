package summary_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/summary"
)

func sampleResults() map[string][]drift.Result {
	return map[string][]drift.Result{
		"auth-service": {
			{Service: "auth-service", Field: "image", Expected: "v1", Actual: "v2"},
			{Service: "auth-service", Field: "replicas", Expected: "3", Actual: "2"},
		},
		"payment-service": {
			{Service: "payment-service", Field: "image", Expected: "v1", Actual: "v3"},
		},
		"clean-service": {},
	}
}

func TestBuild_TotalServices(t *testing.T) {
	r := summary.Build(sampleResults())
	if r.TotalServices != 3 {
		t.Errorf("expected 3 total services, got %d", r.TotalServices)
	}
}

func TestBuild_DriftedAndClean(t *testing.T) {
	r := summary.Build(sampleResults())
	if r.DriftedServices != 2 {
		t.Errorf("expected 2 drifted services, got %d", r.DriftedServices)
	}
	if r.CleanServices != 1 {
		t.Errorf("expected 1 clean service, got %d", r.CleanServices)
	}
}

func TestBuild_TotalDrifts(t *testing.T) {
	r := summary.Build(sampleResults())
	if r.TotalDrifts != 3 {
		t.Errorf("expected 3 total drifts, got %d", r.TotalDrifts)
	}
}

func TestBuild_ByField(t *testing.T) {
	r := summary.Build(sampleResults())
	if r.ByField["image"] != 2 {
		t.Errorf("expected image drift count 2, got %d", r.ByField["image"])
	}
	if r.ByField["replicas"] != 1 {
		t.Errorf("expected replicas drift count 1, got %d", r.ByField["replicas"])
	}
}

func TestBuild_ByService(t *testing.T) {
	r := summary.Build(sampleResults())
	if r.ByService["auth-service"] != 2 {
		t.Errorf("expected auth-service drift count 2, got %d", r.ByService["auth-service"])
	}
	if _, ok := r.ByService["clean-service"]; ok {
		t.Error("clean-service should not appear in ByService map")
	}
}

func TestHasDrift_True(t *testing.T) {
	r := summary.Build(sampleResults())
	if !r.HasDrift() {
		t.Error("expected HasDrift to return true")
	}
}

func TestHasDrift_False(t *testing.T) {
	r := summary.Build(map[string][]drift.Result{"svc": {}})
	if r.HasDrift() {
		t.Error("expected HasDrift to return false for clean services")
	}
}

func TestDriftRate_Empty(t *testing.T) {
	r := summary.Build(map[string][]drift.Result{})
	if r.DriftRate() != 0 {
		t.Errorf("expected drift rate 0 for empty results, got %f", r.DriftRate())
	}
}

func TestDriftRate_Partial(t *testing.T) {
	r := summary.Build(sampleResults())
	expected := 2.0 / 3.0
	if r.DriftRate() != expected {
		t.Errorf("expected drift rate %f, got %f", expected, r.DriftRate())
	}
}
