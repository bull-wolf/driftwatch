package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/reporter"
)

func drifts(pairs ...string) []drift.DriftResult {
	var results []drift.DriftResult
	for i := 0; i+2 <= len(pairs); i += 3 {
		results = append(results, drift.DriftResult{
			Field:    pairs[i],
			Expected: pairs[i+1],
			Actual:   pairs[i+2],
		})
	}
	return results
}

func TestWriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := reporter.NewWriter(&buf, reporter.FormatText)
	r := reporter.Report{Service: "auth-service", Drifts: nil}

	if err := w.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[OK]") {
		t.Errorf("expected [OK] in output, got: %q", got)
	}
	if !strings.Contains(got, "auth-service") {
		t.Errorf("expected service name in output, got: %q", got)
	}
}

func TestWriteText_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	w := reporter.NewWriter(&buf, reporter.FormatText)
	r := reporter.Report{
		Service: "auth-service",
		Drifts:  drifts("image", "nginx:1.25", "nginx:1.24"),
	}

	if err := w.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[DRIFT]") {
		t.Errorf("expected [DRIFT] in output, got: %q", got)
	}
	if !strings.Contains(got, "nginx:1.24") {
		t.Errorf("expected actual value in output, got: %q", got)
	}
}

func TestWriteJSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	w := reporter.NewWriter(&buf, reporter.FormatJSON)
	r := reporter.Report{
		Service: "payment-service",
		Drifts:  drifts("replicas", "3", "1"),
	}

	if err := w.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"service"`) {
		t.Errorf("expected \"service\" key in JSON, got: %q", got)
	}
	if !strings.Contains(got, `"payment-service"`) {
		t.Errorf("expected service name in JSON, got: %q", got)
	}
	if !strings.Contains(got, `"drift_count"`) {
		t.Errorf("expected drift_count in JSON, got: %q", got)
	}
}

func TestWriteJSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := reporter.NewWriter(&buf, reporter.FormatJSON)
	r := reporter.Report{Service: "worker", Drifts: []drift.DriftResult{}}

	if err := w.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"drift_count": 0`) {
		t.Errorf("expected drift_count 0 in JSON, got: %q", got)
	}
}
