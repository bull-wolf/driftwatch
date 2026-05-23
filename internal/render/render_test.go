package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/render"
)

func sampleDrifts() []drift.Result {
	return []drift.Result{
		{Service: "auth-service", Field: "image", Expected: "v1.0", Actual: "v1.1"},
		{Service: "auth-service", Field: "replicas", Expected: "3", Actual: "2"},
	}
}

const simpleTmpl = `Service: {{.Service}}
Drifted: {{.HasDrift}}
{{range .Drifts}}- {{.Field}}: {{.Expected}} -> {{.Actual}}
{{end}}`

func TestRender_NoDrift(t *testing.T) {
	r, err := render.New(simpleTmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	if err := r.Render(&buf, "auth-service", nil); err != nil {
		t.Fatalf("render error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Service: auth-service") {
		t.Errorf("expected service name in output, got: %s", out)
	}
	if !strings.Contains(out, "Drifted: false") {
		t.Errorf("expected HasDrift=false in output, got: %s", out)
	}
}

func TestRender_WithDrift(t *testing.T) {
	r, err := render.New(simpleTmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	if err := r.Render(&buf, "auth-service", sampleDrifts()); err != nil {
		t.Fatalf("render error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Drifted: true") {
		t.Errorf("expected HasDrift=true, got: %s", out)
	}
	if !strings.Contains(out, "image: v1.0 -> v1.1") {
		t.Errorf("expected image drift line, got: %s", out)
	}
}

func TestNew_InvalidTemplate(t *testing.T) {
	_, err := render.New("{{.Unclosed")
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
}

func TestRender_TimestampPresent(t *testing.T) {
	tmplWithTime := `ts={{.Timestamp.IsZero}}`
	r, err := render.New(tmplWithTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	if err := r.Render(&buf, "svc", nil); err != nil {
		t.Fatalf("render error: %v", err)
	}
	if !strings.Contains(buf.String(), "ts=false") {
		t.Errorf("expected non-zero timestamp, got: %s", buf.String())
	}
}
