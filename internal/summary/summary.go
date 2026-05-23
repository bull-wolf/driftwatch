// Package summary provides aggregation utilities for drift results
// across multiple services, producing rollup statistics useful for
// dashboards and high-level reporting.
package summary

import "github.com/yourorg/driftwatch/internal/drift"

// Report holds aggregated drift statistics across all inspected services.
type Report struct {
	TotalServices  int            `json:"total_services"`
	DriftedServices int           `json:"drifted_services"`
	CleanServices  int            `json:"clean_services"`
	TotalDrifts    int            `json:"total_drifts"`
	ByField        map[string]int `json:"by_field"`
	ByService      map[string]int `json:"by_service"`
}

// Build computes a Report from a map of service name to its detected drifts.
func Build(results map[string][]drift.Result) Report {
	r := Report{
		ByField:   make(map[string]int),
		ByService: make(map[string]int),
	}

	r.TotalServices = len(results)

	for service, drifts := range results {
		if len(drifts) == 0 {
			r.CleanServices++
			continue
		}

		r.DriftedServices++
		r.ByService[service] = len(drifts)
		r.TotalDrifts += len(drifts)

		for _, d := range drifts {
			r.ByField[d.Field]++
		}
	}

	return r
}

// HasDrift returns true when at least one service has drifted.
func (r Report) HasDrift() bool {
	return r.DriftedServices > 0
}

// DriftRate returns the fraction of services that have drifted (0.0–1.0).
// Returns 0 when there are no services.
func (r Report) DriftRate() float64 {
	if r.TotalServices == 0 {
		return 0
	}
	return float64(r.DriftedServices) / float64(r.TotalServices)
}
