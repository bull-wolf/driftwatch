package filter

import (
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Options holds filtering criteria for drift results.
type Options struct {
	// Services restricts results to only the named services.
	// An empty slice means no restriction.
	Services []string

	// Fields restricts results to only drifts on the given field names.
	// An empty slice means no restriction.
	Fields []string

	// OnlyDrifted, when true, excludes results that have zero drifts.
	OnlyDrifted bool
}

// Apply filters a slice of drift.Result according to the given Options.
// It returns a new slice containing only the results that match all criteria.
func Apply(results []drift.Result, opts Options) []drift.Result {
	serviceSet := toSet(opts.Services)
	fieldSet := toSet(opts.Fields)

	var filtered []drift.Result
	for _, r := range results {
		if len(serviceSet) > 0 && !serviceSet[strings.ToLower(r.Service)] {
			continue
		}

		matched := filterDrifts(r.Drifts, fieldSet)

		if opts.OnlyDrifted && len(matched) == 0 {
			continue
		}

		filtered = append(filtered, drift.Result{
			Service: r.Service,
			Drifts:  matched,
		})
	}
	return filtered
}

func filterDrifts(drifts []drift.Drift, fieldSet map[string]bool) []drift.Drift {
	if len(fieldSet) == 0 {
		return drifts
	}
	var out []drift.Drift
	for _, d := range drifts {
		if fieldSet[strings.ToLower(d.Field)] {
			out = append(out, d)
		}
	}
	return out
}

func toSet(items []string) map[string]bool {
	if len(items) == 0 {
		return nil
	}
	s := make(map[string]bool, len(items))
	for _, v := range items {
		s[strings.ToLower(v)] = true
	}
	return s
}
