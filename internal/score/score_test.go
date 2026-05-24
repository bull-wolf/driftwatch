package score_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/policy"
	"github.com/driftwatch/internal/score"
)

func sampleDrifts(n int) []drift.Drift {
	out := make([]drift.Drift, n)
	for i := range out {
		out[i] = drift.Drift{Field: "image", Expected: "v1", Actual: "v2"}
	}
	return out
}

func sampleViolations(n int) []policy.Violation {
	out := make([]policy.Violation, n)
	for i := range out {
		out[i] = policy.Violation{Field: "replicas", Message: "too low"}
	}
	return out
}

func TestCompute_NoDriftNoViolations(t *testing.T) {
	r := score.Compute("auth-service", nil, nil)
	if r.Score != 100 {
		t.Errorf("expected score 100, got %d", r.Score)
	}
}

func TestCompute_SingleDrift(t *testing.T) {
	r := score.Compute("auth-service", sampleDrifts(1), nil)
	want := 100 - score.DriftPenalty
	if r.Score != want {
		t.Errorf("expected score %d, got %d", want, r.Score)
	}
}

func TestCompute_SingleViolation(t *testing.T) {
	r := score.Compute("auth-service", nil, sampleViolations(1))
	want := 100 - score.ViolationPenalty
	if r.Score != want {
		t.Errorf("expected score %d, got %d", want, r.Score)
	}
}

func TestCompute_ClampedToZero(t *testing.T) {
	r := score.Compute("auth-service", sampleDrifts(20), sampleViolations(10))
	if r.Score != 0 {
		t.Errorf("expected clamped score 0, got %d", r.Score)
	}
}

func TestCompute_SummaryContainsService(t *testing.T) {
	r := score.Compute("auth-service", sampleDrifts(1), nil)
	if !strings.Contains(r.Summary, "auth-service") {
		t.Errorf("summary missing service name: %s", r.Summary)
	}
}

func TestCompute_CountsAreCorrect(t *testing.T) {
	r := score.Compute("svc", sampleDrifts(3), sampleViolations(2))
	if r.Drifts != 3 {
		t.Errorf("expected 3 drifts, got %d", r.Drifts)
	}
	if r.Violations != 2 {
		t.Errorf("expected 2 violations, got %d", r.Violations)
	}
}

func TestGrade_Boundaries(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{100, "A"}, {90, "A"}, {89, "B"}, {75, "B"},
		{74, "C"}, {50, "C"}, {49, "D"}, {25, "D"},
		{24, "F"}, {0, "F"},
	}
	for _, tc := range cases {
		got := score.Grade(tc.score)
		if got != tc.want {
			t.Errorf("Grade(%d) = %q, want %q", tc.score, got, tc.want)
		}
	}
}
