package drift

import (
	"fmt"

	"github.com/driftwatch/internal/manifest"
)

// DriftResult describes a single field that has drifted from the manifest.
type DriftResult struct {
	Service  string
	Field    string
	Expected string
	Actual   string
}

// LiveState represents the observed runtime state of a deployed service.
type LiveState struct {
	Name     string
	Image    string
	Replicas int
	Env      map[string]string
}

// Detect compares a manifest against the live state and returns any drifts found.
func Detect(m manifest.Manifest, live LiveState) []DriftResult {
	var results []DriftResult

	if m.Image != live.Image {
		results = append(results, DriftResult{
			Service:  m.Name,
			Field:    "image",
			Expected: m.Image,
			Actual:   live.Image,
		})
	}

	if m.Replicas != live.Replicas {
		results = append(results, DriftResult{
			Service:  m.Name,
			Field:    "replicas",
			Expected: fmt.Sprintf("%d", m.Replicas),
			Actual:   fmt.Sprintf("%d", live.Replicas),
		})
	}

	for key, expectedVal := range m.Env {
		actualVal, ok := live.Env[key]
		if !ok {
			results = append(results, DriftResult{
				Service:  m.Name,
				Field:    fmt.Sprintf("env.%s", key),
				Expected: expectedVal,
				Actual:   "<missing>",
			})
			continue
		}
		if expectedVal != actualVal {
			results = append(results, DriftResult{
				Service:  m.Name,
				Field:    fmt.Sprintf("env.%s", key),
				Expected: expectedVal,
				Actual:   actualVal,
			})
		}
	}

	return results
}
