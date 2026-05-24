// Package score computes a numeric drift health score for a service
// based on the number and severity of detected drifts and policy violations.
package score

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/policy"
)

const (
	// MaxScore is a perfect health score with no drift or violations.
	MaxScore = 100

	// DriftPenalty is subtracted per drifted field.
	DriftPenalty = 10

	// ViolationPenalty is subtracted per policy violation.
	ViolationPenalty = 20
)

// Result holds the computed score and a human-readable summary.
type Result struct {
	Service    string
	Score      int
	Drifts     int
	Violations int
	Summary    string
}

// Compute calculates the health score for a single service given its
// detected drifts and policy violations. The score is clamped to [0, 100].
func Compute(service string, drifts []drift.Drift, violations []policy.Violation) Result {
	penalty := len(drifts)*DriftPenalty + len(violations)*ViolationPenalty
	score := MaxScore - penalty
	if score < 0 {
		score = 0
	}

	summary := fmt.Sprintf(
		"service=%s score=%d drifts=%d violations=%d",
		service, score, len(drifts), len(violations),
	)

	return Result{
		Service:    service,
		Score:      score,
		Drifts:     len(drifts),
		Violations: len(violations),
		Summary:    summary,
	}
}

// Grade returns a letter grade for a given score.
func Grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 50:
		return "C"
	case score >= 25:
		return "D"
	default:
		return "F"
	}
}
