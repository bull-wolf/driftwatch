package alert_test

import (
	"testing"

	"github.com/driftwatch/internal/alert"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/policy"
)

func sampleDrifts(n int) []drift.Drift {
	result := make([]drift.Drift, n)
	for i := range result {
		result[i] = drift.Drift{Field: "image", Expected: "v1", Actual: "v2"}
	}
	return result
}

func sampleViolations(n int) []policy.Violation {
	result := make([]policy.Violation, n)
	for i := range result {
		result[i] = policy.Violation{Field: "replicas", Message: "below minimum"}
	}
	return result
}

func TestEvaluate_NoDriftsNoViolations_ReturnsNil(t *testing.T) {
	a := alert.Evaluate(nil, nil, alert.Threshold{MaxDrifts: 1})
	if a != nil {
		t.Errorf("expected nil alert, got %+v", a)
	}
}

func TestEvaluate_DriftsBelowThreshold_ReturnsNil(t *testing.T) {
	a := alert.Evaluate(sampleDrifts(2), nil, alert.Threshold{MaxDrifts: 5})
	if a != nil {
		t.Errorf("expected nil alert for drifts below threshold, got %+v", a)
	}
}

func TestEvaluate_DriftsExceedWarnThreshold(t *testing.T) {
	a := alert.Evaluate(sampleDrifts(3), nil, alert.Threshold{MaxDrifts: 2})
	if a == nil {
		t.Fatal("expected non-nil alert")
	}
	if a.Level != alert.LevelWarn {
		t.Errorf("expected warn level, got %s", a.Level)
	}
	if a.DriftCount != 3 {
		t.Errorf("expected DriftCount 3, got %d", a.DriftCount)
	}
}

func TestEvaluate_DriftsExceedCriticalThreshold(t *testing.T) {
	a := alert.Evaluate(sampleDrifts(10), nil, alert.Threshold{MaxDrifts: 2, CriticalDrifts: 5})
	if a == nil {
		t.Fatal("expected non-nil alert")
	}
	if a.Level != alert.LevelCritical {
		t.Errorf("expected critical level, got %s", a.Level)
	}
}

func TestEvaluate_ViolationsExceedThreshold(t *testing.T) {
	a := alert.Evaluate(nil, sampleViolations(3), alert.Threshold{MaxViolations: 1})
	if a == nil {
		t.Fatal("expected non-nil alert")
	}
	if a.Level != alert.LevelWarn {
		t.Errorf("expected warn, got %s", a.Level)
	}
	if a.Violations != 3 {
		t.Errorf("expected 3 violations, got %d", a.Violations)
	}
}

func TestEvaluate_MessageIsNonEmpty(t *testing.T) {
	a := alert.Evaluate(sampleDrifts(5), nil, alert.Threshold{MaxDrifts: 2})
	if a == nil {
		t.Fatal("expected alert")
	}
	if a.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestEvaluate_TimestampSet(t *testing.T) {
	a := alert.Evaluate(sampleDrifts(2), nil, alert.Threshold{MaxDrifts: 1})
	if a == nil {
		t.Fatal("expected alert")
	}
	if a.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
