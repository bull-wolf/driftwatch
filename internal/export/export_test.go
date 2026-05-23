package export_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/export"
)

var sampleDrifts = []drift.Result{
	{Service: "auth-service", Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"},
	{Service: "auth-service", Field: "replicas", Expected: "3", Actual: "2"},
}

func TestWriteCSV_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	ex := export.New(&buf)
	if err := ex.WriteCSV(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 header line, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "service") {
		t.Errorf("expected header row, got: %s", lines[0])
	}
}

func TestWriteCSV_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	ex := export.New(&buf)
	if err := ex.WriteCSV(sampleDrifts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if !strings.Contains(lines[1], "auth-service") {
		t.Errorf("expected service name in row, got: %s", lines[1])
	}
	if !strings.Contains(lines[1], "nginx:1.24") {
		t.Errorf("expected actual value in row, got: %s", lines[1])
	}
}

func TestWriteMarkdown_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	ex := export.New(&buf)
	if err := ex.WriteMarkdown(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Service |") {
		t.Errorf("expected markdown header, got: %s", out)
	}
}

func TestWriteMarkdown_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	ex := export.New(&buf)
	if err := ex.WriteMarkdown(sampleDrifts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "auth-service") {
		t.Errorf("expected service name in output, got: %s", out)
	}
	if !strings.Contains(out, "replicas") {
		t.Errorf("expected field name in output, got: %s", out)
	}
	if !strings.Contains(out, "nginx:1.25") {
		t.Errorf("expected expected value in output, got: %s", out)
	}
}
