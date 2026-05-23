package filter_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/filter"
)

var sampleResults = []drift.Result{
	{
		Service: "auth-service",
		Drifts: []drift.Drift{
			{Field: "image", Expected: "auth:v1", Actual: "auth:v2"},
			{Field: "replicas", Expected: "2", Actual: "3"},
		},
	},
	{
		Service: "payment-service",
		Drifts:  []drift.Drift{},
	},
	{
		Service: "gateway",
		Drifts: []drift.Drift{
			{Field: "image", Expected: "gw:stable", Actual: "gw:edge"},
		},
	},
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{})
	if len(got) != len(sampleResults) {
		t.Fatalf("expected %d results, got %d", len(sampleResults), len(got))
	}
}

func TestApply_FilterByService(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{
		Services: []string{"auth-service"},
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Service != "auth-service" {
		t.Errorf("unexpected service %q", got[0].Service)
	}
}

func TestApply_FilterByField(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{
		Fields: []string{"image"},
	})
	// all three services pass service filter; only image drifts kept
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
	for _, r := range got {
		for _, d := range r.Drifts {
			if d.Field != "image" {
				t.Errorf("unexpected field %q in service %q", d.Field, r.Service)
			}
		}
	}
}

func TestApply_OnlyDrifted(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{OnlyDrifted: true})
	if len(got) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(got))
	}
	for _, r := range got {
		if len(r.Drifts) == 0 {
			t.Errorf("service %q should have been excluded", r.Service)
		}
	}
}

func TestApply_ServiceAndField_Combined(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{
		Services:    []string{"auth-service"},
		Fields:      []string{"replicas"},
		OnlyDrifted: true,
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if len(got[0].Drifts) != 1 || got[0].Drifts[0].Field != "replicas" {
		t.Errorf("unexpected drifts: %+v", got[0].Drifts)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	got := filter.Apply(nil, filter.Options{OnlyDrifted: true})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d", len(got))
	}
}
