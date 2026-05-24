package group_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/group"
)

var sampleResults = []drift.DriftResult{
	{Service: "auth", Field: "image", Drifted: true, Expected: "v1", Actual: "v2"},
	{Service: "auth", Field: "replicas", Drifted: false, Expected: "3", Actual: "3"},
	{Service: "payments", Field: "image", Drifted: true, Expected: "v1", Actual: "v3"},
	{Service: "payments", Field: "replicas", Drifted: true, Expected: "2", Actual: "1"},
}

func TestApply_ByService(t *testing.T) {
	res, err := group.Apply(sampleResults, group.ByService)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Dimension != group.ByService {
		t.Errorf("expected dimension %q, got %q", group.ByService, res.Dimension)
	}
	if len(res.Groups["auth"]) != 2 {
		t.Errorf("expected 2 results for auth, got %d", len(res.Groups["auth"]))
	}
	if len(res.Groups["payments"]) != 2 {
		t.Errorf("expected 2 results for payments, got %d", len(res.Groups["payments"]))
	}
}

func TestApply_ByField(t *testing.T) {
	res, err := group.Apply(sampleResults, group.ByField)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups["image"]) != 2 {
		t.Errorf("expected 2 image results, got %d", len(res.Groups["image"]))
	}
	if len(res.Groups["replicas"]) != 2 {
		t.Errorf("expected 2 replicas results, got %d", len(res.Groups["replicas"]))
	}
}

func TestApply_BySeverity(t *testing.T) {
	res, err := group.Apply(sampleResults, group.BySeverity)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups["drifted"]) != 3 {
		t.Errorf("expected 3 drifted, got %d", len(res.Groups["drifted"]))
	}
	if len(res.Groups["clean"]) != 1 {
		t.Errorf("expected 1 clean, got %d", len(res.Groups["clean"]))
	}
}

func TestApply_UnknownDimension(t *testing.T) {
	_, err := group.Apply(sampleResults, group.By("unknown"))
	if err == nil {
		t.Fatal("expected error for unknown dimension, got nil")
	}
}

func TestApply_EmptyResults(t *testing.T) {
	res, err := group.Apply(nil, group.ByService)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(res.Groups))
	}
}

func TestResult_Keys_ReturnsAllGroups(t *testing.T) {
	res, _ := group.Apply(sampleResults, group.ByService)
	keys := res.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}
