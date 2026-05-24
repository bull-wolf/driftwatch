// Package alert provides threshold-based alerting for drift results.
// It evaluates drift counts and policy violations against configurable
// thresholds and emits structured alert events when limits are exceeded.
package alert

import (
	"fmt"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/policy"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn     Level = "warn"
	LevelCritical Level = "critical"
)

// Threshold defines the conditions that trigger an alert.
type Threshold struct {
	MaxDrifts      int   // alert if total drifts exceed this value
	MaxViolations  int   // alert if total violations exceed this value
	CriticalDrifts int   // escalate to critical above this value
}

// Alert is an emitted alert event.
type Alert struct {
	Level       Level
	Message     string
	DriftCount  int
	Violations  int
	Timestamp   time.Time
}

// Evaluate checks drifts and violations against the given threshold.
// It returns a non-nil Alert if any threshold is breached, or nil if clean.
func Evaluate(drifts []drift.Drift, violations []policy.Violation, t Threshold) *Alert {
	driftCount := len(drifts)
	violationCount := len(violations)

	if driftCount == 0 && violationCount == 0 {
		return nil
	}

	var level Level
	var msgs []string

	if driftCount > t.CriticalDrifts && t.CriticalDrifts > 0 {
		level = LevelCritical
		msgs = append(msgs, fmt.Sprintf("%d drifts exceed critical threshold of %d", driftCount, t.CriticalDrifts))
	} else if driftCount > t.MaxDrifts && t.MaxDrifts > 0 {
		level = LevelWarn
		msgs = append(msgs, fmt.Sprintf("%d drifts exceed warn threshold of %d", driftCount, t.MaxDrifts))
	}

	if violationCount > t.MaxViolations && t.MaxViolations > 0 {
		if level != LevelCritical {
			level = LevelWarn
		}
		msgs = append(msgs, fmt.Sprintf("%d policy violations exceed threshold of %d", violationCount, t.MaxViolations))
	}

	if level == "" {
		return nil
	}

	msg := ""
	for i, m := range msgs {
		if i > 0 {
			msg += "; "
		}
		msg += m
	}

	return &Alert{
		Level:      level,
		Message:    msg,
		DriftCount: driftCount,
		Violations: violationCount,
		Timestamp:  time.Now().UTC(),
	}
}
