// Package export provides functionality for exporting drift results
// to various output formats such as CSV and Markdown.
package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Format represents a supported export format.
type Format string

const (
	FormatCSV      Format = "csv"
	FormatMarkdown Format = "markdown"
)

// Exporter writes drift results to an output stream.
type Exporter struct {
	w io.Writer
}

// New creates a new Exporter writing to w.
func New(w io.Writer) *Exporter {
	return &Exporter{w: w}
}

// WriteCSV writes drift results as CSV rows to the underlying writer.
func (e *Exporter) WriteCSV(drifts []drift.Result) error {
	cw := csv.NewWriter(e.w)
	if err := cw.Write([]string{"service", "field", "expected", "actual", "timestamp"}); err != nil {
		return fmt.Errorf("export: write csv header: %w", err)
	}
	ts := time.Now().UTC().Format(time.RFC3339)
	for _, d := range drifts {
		row := []string{d.Service, d.Field, d.Expected, d.Actual, ts}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("export: write csv row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}

// WriteMarkdown writes drift results as a Markdown table to the underlying writer.
func (e *Exporter) WriteMarkdown(drifts []drift.Result) error {
	var sb strings.Builder
	sb.WriteString("| Service | Field | Expected | Actual |\n")
	sb.WriteString("|---------|-------|----------|--------|\n")
	for _, d := range drifts {
		row := fmt.Sprintf("| %s | %s | %s | %s |\n", d.Service, d.Field, d.Expected, d.Actual)
		sb.WriteString(row)
	}
	_, err := fmt.Fprint(e.w, sb.String())
	if err != nil {
		return fmt.Errorf("export: write markdown: %w", err)
	}
	return nil
}
