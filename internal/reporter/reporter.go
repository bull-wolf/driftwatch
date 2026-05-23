package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Format represents the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the results of a drift detection run.
type Report struct {
	Service string
	Drifts  []drift.DriftResult
}

// Writer writes drift reports to an output stream.
type Writer struct {
	format Format
	out    io.Writer
}

// NewWriter creates a new Reporter that writes to w in the given format.
func NewWriter(out io.Writer, format Format) *Writer {
	return &Writer{out: out, format: format}
}

// Write outputs the report in the configured format.
func (w *Writer) Write(r Report) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(r)
	default:
		return w.writeText(r)
	}
}

func (w *Writer) writeText(r Report) error {
	if len(r.Drifts) == 0 {
		_, err := fmt.Fprintf(w.out, "[OK] %s: no drift detected\n", r.Service)
		return err
	}

	_, err := fmt.Fprintf(w.out, "[DRIFT] %s: %d issue(s) found\n", r.Service, len(r.Drifts))
	if err != nil {
		return err
	}

	for _, d := range r.Drifts {
		line := fmt.Sprintf("  - field=%s expected=%q actual=%q\n",
			d.Field, d.Expected, d.Actual)
		if _, err := fmt.Fprint(w.out, line); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeJSON(r Report) error {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("  \"service\": %q,\n", r.Service))
	sb.WriteString(fmt.Sprintf("  \"drift_count\": %d,\n", len(r.Drifts)))
	sb.WriteString("  \"drifts\": [\n")
	for i, d := range r.Drifts {
		comma := ","
		if i == len(r.Drifts)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    {\"field\": %q, \"expected\": %q, \"actual\": %q}%s\n",
			d.Field, d.Expected, d.Actual, comma))
	}
	sb.WriteString("  ]\n}\n")
	_, err := fmt.Fprint(w.out, sb.String())
	return err
}
