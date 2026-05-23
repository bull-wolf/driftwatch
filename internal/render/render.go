// Package render provides template-based rendering for drift reports.
package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// ReportData holds the data passed to a report template.
type ReportData struct {
	Service   string
	Timestamp time.Time
	Drifts    []drift.Result
	HasDrift  bool
}

// Renderer renders drift reports using Go templates.
type Renderer struct {
	tmpl *template.Template
}

// New creates a Renderer from the given template string.
func New(tmplSrc string) (*Renderer, error) {
	t, err := template.New("report").Parse(tmplSrc)
	if err != nil {
		return nil, fmt.Errorf("render: parse template: %w", err)
	}
	return &Renderer{tmpl: t}, nil
}

// Render writes the rendered report for the given service and drifts to w.
func (r *Renderer) Render(w io.Writer, service string, drifts []drift.Result) error {
	data := ReportData{
		Service:   service,
		Timestamp: time.Now().UTC(),
		Drifts:    drifts,
		HasDrift:  len(drifts) > 0,
	}
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render: execute template: %w", err)
	}
	_, err := w.Write(buf.Bytes())
	return err
}
