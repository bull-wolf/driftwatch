package drift

import (
	"fmt"

	"github.com/driftwatch/internal/manifest"
)

// DriftResult holds the outcome of comparing a deployed service against its manifest.
type DriftResult struct {
	ServiceName string
	HasDrift    bool
	Diffs       []string
}

// DeployedState represents the observed live state of a service.
type DeployedState struct {
	Name        string
	Image       string
	Replicas    int
	Environment string
}

// Detect compares a deployed service state against the source-of-truth manifest.
// It returns a DriftResult describing any discrepancies found.
func Detect(deployed DeployedState, m manifest.Manifest) DriftResult {
	result := DriftResult{
		ServiceName: m.Name,
		HasDrift:    false,
		Diffs:       []string{},
	}

	if deployed.Image != m.Image {
		result.Diffs = append(result.Diffs, fmt.Sprintf(
			"image mismatch: deployed=%q manifest=%q",
			deployed.Image, m.Image,
		))
	}

	if deployed.Replicas != m.Replicas {
		result.Diffs = append(result.Diffs, fmt.Sprintf(
			"replicas mismatch: deployed=%d manifest=%d",
			deployed.Replicas, m.Replicas,
		))
	}

	if deployed.Environment != m.Environment {
		result.Diffs = append(result.Diffs, fmt.Sprintf(
			"environment mismatch: deployed=%q manifest=%q",
			deployed.Environment, m.Environment,
		))
	}

	if len(result.Diffs) > 0 {
		result.HasDrift = true
	}

	return result
}
