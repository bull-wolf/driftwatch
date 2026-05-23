package export_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/export"
)

func TestWriteMarkdown_MatchesTestdata(t *testing.T) {
	drifts := []drift.Result{
		{Service: "auth-service", Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"},
		{Service: "auth-service", Field: "replicas", Expected: "3", Actual: "2"},
	}

	var buf bytes.Buffer
	ex := export.New(&buf)
	if err := ex.WriteMarkdown(drifts); err != nil {
		t.Fatalf("WriteMarkdown error: %v", err)
	}

	tdPath := filepath.Join("..", "..", "testdata", "export", "expected-output.md")
	expected, err := os.ReadFile(tdPath)
	if err != nil {
		t.Fatalf("could not read testdata file: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	want := strings.TrimSpace(string(expected))
	if got != want {
		t.Errorf("markdown output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteCSV_HeaderAlwaysPresent(t *testing.T) {
	var buf bytes.Buffer
	ex := export.New(&buf)
	if err := ex.WriteCSV([]drift.Result{}); err != nil {
		t.Fatalf("WriteCSV error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "service,field,expected,actual,timestamp") {
		t.Errorf("expected CSV header as first line, got: %q", out)
	}
}
