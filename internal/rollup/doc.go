// Package rollup provides cross-service aggregation of drift detection results.
//
// It accepts a map of service names to their detected drifts and produces a
// consolidated Result containing totals, per-field counts, severity breakdowns,
// and the top drifting fields across all services.
//
// Typical usage:
//
//	results := map[string][]drift.Drift{
//		"auth-service":    detectedDrifts,
//		"payment-service": otherDrifts,
//	}
//	r := rollup.Aggregate(results)
//	fmt.Println(rollup.Summary(r))
package rollup
