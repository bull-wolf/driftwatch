// Package status provides a lightweight health summary for a monitored service,
// combining drift results, policy violations, and baseline state into a single
// StatusReport that can be used for dashboards or alerting.
package status

import (
	"fmt"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/policy"
)

// State represents the overall health of a service.
type State string

const (
	StateClean    State = "clean"
	StateDrifted  State = "drifted"
	StateViolated State = "violated"
)

// Report holds a point-in-time health summary for a single service.
type Report struct {
	ServiceName string    `json:"service_name"`
	State       State     `json:"state"`
	DriftCount  int       `json:"drift_count"`
	Violations  []string  `json:"violations,omitempty"`
	CheckedAt   time.Time `json:"checked_at"`
	Summary     string    `json:"summary"`
}

// Build constructs a Report from drift results and policy violations.
// The worst state wins: violated > drifted > clean.
func Build(serviceName string, drifts []drift.Result, violations []policy.Violation) Report {
	state := StateClean

	if len(drifts) > 0 {
		state = StateDrifted
	}

	violationMessages := make([]string, 0, len(violations))
	for _, v := range violations {
		violationMessages = append(violationMessages, v.Message)
		state = StateViolated
	}

	r := Report{
		ServiceName: serviceName,
		State:       state,
		DriftCount:  len(drifts),
		Violations:  violationMessages,
		CheckedAt:   time.Now().UTC(),
	}
	r.Summary = buildSummary(r)
	return r
}

func buildSummary(r Report) string {
	switch r.State {
	case StateClean:
		return fmt.Sprintf("%s: no drift detected", r.ServiceName)
	case StateDrifted:
		return fmt.Sprintf("%s: %d drift(s) detected", r.ServiceName, r.DriftCount)
	case StateViolated:
		return fmt.Sprintf("%s: %d drift(s), %d policy violation(s)", r.ServiceName, r.DriftCount, len(r.Violations))
	default:
		return fmt.Sprintf("%s: unknown state", r.ServiceName)
	}
}
