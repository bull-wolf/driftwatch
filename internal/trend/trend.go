// Package trend analyzes drift history over time to identify
// services with recurring or worsening drift patterns.
package trend

import (
	"fmt"
	"time"

	"github.com/driftwatch/internal/history"
)

// Direction indicates whether drift is improving, worsening, or stable.
type Direction string

const (
	Improving Direction = "improving"
	Worsening Direction = "worsening"
	Stable    Direction = "stable"
)

// Summary holds the trend analysis for a single service.
type Summary struct {
	ServiceName    string
	TotalEntries   int
	DriftedEntries int
	Direction      Direction
	DriftRate      float64 // fraction of entries that had drift
}

// Analyze reads the history for the given service and computes a trend summary.
// It compares the drift rate of the first half of entries to the second half
// to determine whether drift is improving or worsening.
func Analyze(serviceName string, dir string) (Summary, error) {
	if serviceName == "" {
		return Summary{}, fmt.Errorf("trend: service name must not be empty")
	}

	entries, err := history.Read(serviceName, dir)
	if err != nil {
		return Summary{}, fmt.Errorf("trend: failed to read history for %q: %w", serviceName, err)
	}

	total := len(entries)
	if total == 0 {
		return Summary{
			ServiceName:  serviceName,
			Direction:    Stable,
		}, nil
	}

	driftedCount := 0
	for _, e := range entries {
		if e.DriftCount > 0 {
			driftedCount++
		}
	}

	driftRate := float64(driftedCount) / float64(total)
	direction := directionFromEntries(entries)

	return Summary{
		ServiceName:    serviceName,
		TotalEntries:   total,
		DriftedEntries: driftedCount,
		Direction:      direction,
		DriftRate:      driftRate,
	}, nil
}

// directionFromEntries splits entries into two halves and compares drift rates.
func directionFromEntries(entries []history.Entry) Direction {
	if len(entries) < 2 {
		return Stable
	}

	mid := len(entries) / 2
	firstHalf := entries[:mid]
	secondHalf := entries[mid:]

	firstRate := driftRateOf(firstHalf)
	secondRate := driftRateOf(secondHalf)

	switch {
	case secondRate > firstRate:
		return Worsening
	case secondRate < firstRate:
		return Improving
	default:
		return Stable
	}
}

func driftRateOf(entries []history.Entry) float64 {
	if len(entries) == 0 {
		return 0
	}
	count := 0
	for _, e := range entries {
		if e.DriftCount > 0 {
			count++
		}
	}
	return float64(count) / float64(len(entries))
}

// RecentDriftWindow returns only entries within the given duration from now.
func RecentDriftWindow(entries []history.Entry, window time.Duration) []history.Entry {
	cutoff := time.Now().UTC().Add(-window)
	var result []history.Entry
	for _, e := range entries {
		if e.RecordedAt.After(cutoff) {
			result = append(result, e)
		}
	}
	return result
}
