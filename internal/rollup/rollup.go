// Package rollup aggregates drift results across multiple services into
// a single consolidated report, grouping by severity and field.
package rollup

import (
	"fmt"
	"sort"

	"github.com/driftwatch/internal/drift"
)

// Result holds aggregated drift data across all scanned services.
type Result struct {
	TotalServices  int
	DriftedCount   int
	CleanCount     int
	ByField        map[string]int
	BySeverity     map[string]int
	TopDriftFields []string
}

// Aggregate combines drift results from multiple services into a single Result.
func Aggregate(results map[string][]drift.Drift) Result {
	r := Result{
		ByField:    make(map[string]int),
		BySeverity: make(map[string]int),
	}

	r.TotalServices = len(results)

	for _, drifts := range results {
		if len(drifts) > 0 {
			r.DriftedCount++
		} else {
			r.CleanCount++
		}
		for _, d := range drifts {
			r.ByField[d.Field]++
			severity := severityFor(d.Field)
			r.BySeverity[severity]++
		}
	}

	r.TopDriftFields = topN(r.ByField, 3)
	return r
}

// Summary returns a human-readable one-line summary of the rollup result.
func Summary(r Result) string {
	return fmt.Sprintf(
		"%d/%d services drifted, %d total drift(s)",
		r.DriftedCount, r.TotalServices, totalDrifts(r.ByField),
	)
}

func totalDrifts(byField map[string]int) int {
	n := 0
	for _, v := range byField {
		n += v
	}
	return n
}

func severityFor(field string) string {
	switch field {
	case "image", "replicas":
		return "high"
	case "env":
		return "medium"
	default:
		return "low"
	}
}

func topN(counts map[string]int, n int) []string {
	type kv struct {
		key string
		val int
	}
	var sorted []kv
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].val != sorted[j].val {
			return sorted[i].val > sorted[j].val
		}
		return sorted[i].key < sorted[j].key
	})
	result := make([]string, 0, n)
	for i := 0; i < n && i < len(sorted); i++ {
		result = append(result, sorted[i].key)
	}
	return result
}
