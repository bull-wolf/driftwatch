// Package diff provides utilities for computing human-readable diffs
// between expected and actual service configuration values.
package diff

import (
	"fmt"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Summary holds a formatted diff summary for a single drift entry.
type Summary struct {
	Service  string
	Field    string
	Expected string
	Actual   string
	Line     string
}

// Summarize converts a slice of DriftResult into human-readable Summary entries.
func Summarize(drifts []drift.DriftResult) []Summary {
	summaries := make([]Summary, 0, len(drifts))
	for _, d := range drifts {
		s := Summary{
			Service:  d.Service,
			Field:    d.Field,
			Expected: fmt.Sprintf("%v", d.Expected),
			Actual:   fmt.Sprintf("%v", d.Actual),
		}
		s.Line = formatLine(s)
		summaries = append(summaries, s)
	}
	return summaries
}

// Format returns a multi-line string representation of all summaries.
func Format(summaries []Summary) string {
	if len(summaries) == 0 {
		return "no drift detected"
	}
	var sb strings.Builder
	for _, s := range summaries {
		sb.WriteString(s.Line)
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatLine(s Summary) string {
	return fmt.Sprintf("[%s] %s: expected=%q actual=%q", s.Service, s.Field, s.Expected, s.Actual)
}
