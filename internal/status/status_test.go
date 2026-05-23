package status_test

import (
	"strings"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/policy"
	"github.com/example/driftwatch/internal/status"
)

var noDrifts []drift.Result
var noViolations []policy.Violation

func sampleDrifts() []drift.Result {
	return []drift.Result{
		{ServiceName: "auth-service", Field: "image", Expected: "v1", Actual: "v2"},
		{ServiceName: "auth-service", Field: "replicas", Expected: "3", Actual: "1"},
	}
}

func sampleViolations() []policy.Violation {
	return []policy.Violation{
		{Field: "image", Message: "image tag must not be latest"},
	}
}

func TestBuild_CleanState(t *testing.T) {
	r := status.Build("auth-service", noDrifts, noViolations)
	if r.State != status.StateClean {
		t.Errorf("expected clean, got %s", r.State)
	}
	if r.DriftCount != 0 {
		t.Errorf("expected 0 drifts, got %d", r.DriftCount)
	}
}

func TestBuild_DriftedState(t *testing.T) {
	r := status.Build("auth-service", sampleDrifts(), noViolations)
	if r.State != status.StateDrifted {
		t.Errorf("expected drifted, got %s", r.State)
	}
	if r.DriftCount != 2 {
		t.Errorf("expected 2 drifts, got %d", r.DriftCount)
	}
}

func TestBuild_ViolatedState_WinsOverDrift(t *testing.T) {
	r := status.Build("auth-service", sampleDrifts(), sampleViolations())
	if r.State != status.StateViolated {
		t.Errorf("expected violated, got %s", r.State)
	}
	if len(r.Violations) != 1 {
		t.Errorf("expected 1 violation message, got %d", len(r.Violations))
	}
}

func TestBuild_ServiceNamePreserved(t *testing.T) {
	r := status.Build("auth-service", noDrifts, noViolations)
	if r.ServiceName != "auth-service" {
		t.Errorf("expected auth-service, got %s", r.ServiceName)
	}
}

func TestBuild_SummaryContainsServiceName(t *testing.T) {
	r := status.Build("auth-service", noDrifts, noViolations)
	if !strings.Contains(r.Summary, "auth-service") {
		t.Errorf("summary missing service name: %s", r.Summary)
	}
}

func TestBuild_SummaryMentionsDriftCount(t *testing.T) {
	r := status.Build("auth-service", sampleDrifts(), noViolations)
	if !strings.Contains(r.Summary, "2") {
		t.Errorf("summary should mention drift count: %s", r.Summary)
	}
}

func TestBuild_CheckedAtIsRecent(t *testing.T) {
	before := time.Now().UTC()
	r := status.Build("auth-service", noDrifts, noViolations)
	after := time.Now().UTC()
	if r.CheckedAt.Before(before) || r.CheckedAt.After(after) {
		t.Errorf("CheckedAt %v not within expected range", r.CheckedAt)
	}
}
