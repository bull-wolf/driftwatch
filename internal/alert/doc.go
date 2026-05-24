// Package alert implements threshold-based alerting for the driftwatch tool.
//
// It evaluates detected drift results and policy violations against
// user-defined thresholds and produces structured Alert events indicating
// whether the current state is within acceptable bounds.
//
// Thresholds:
//
//	MaxDrifts      — warn when total drift count exceeds this value
//	MaxViolations  — warn when total violation count exceeds this value
//	CriticalDrifts — escalate alert level to critical above this value
//
// Usage:
//
//	t := alert.Threshold{MaxDrifts: 3, CriticalDrifts: 10, MaxViolations: 1}
//	a := alert.Evaluate(drifts, violations, t)
//	if a != nil {
//	    fmt.Printf("[%s] %s\n", a.Level, a.Message)
//	}
package alert
