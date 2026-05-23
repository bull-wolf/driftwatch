package diff_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/diff"
	"github.com/driftwatch/internal/drift"
)

func sampleDrifts() []drift.DriftResult {
	return []drift.DriftResult{
		{Service: "auth-service", Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.19"},
		{Service: "auth-service", Field: "replicas", Expected: 3, Actual: 1},
	}
}

func TestSummarize_ReturnsSameCount(t *testing.T) {
	drifts := sampleDrifts()
	summaries := diff.Summarize(drifts)
	if len(summaries) != len(drifts) {
		t.Fatalf("expected %d summaries, got %d", len(drifts), len(summaries))
	}
}

func TestSummarize_FieldsPopulated(t *testing.T) {
	summaries := diff.Summarize(sampleDrifts())
	s := summaries[0]
	if s.Service != "auth-service" {
		t.Errorf("expected service auth-service, got %s", s.Service)
	}
	if s.Field != "image" {
		t.Errorf("expected field image, got %s", s.Field)
	}
	if s.Expected != "nginx:1.25" {
		t.Errorf("unexpected Expected value: %s", s.Expected)
	}
	if s.Actual != "nginx:1.19" {
		t.Errorf("unexpected Actual value: %s", s.Actual)
	}
}

func TestSummarize_LineContainsAllParts(t *testing.T) {
	summaries := diff.Summarize(sampleDrifts())
	line := summaries[0].Line
	for _, part := range []string{"auth-service", "image", "nginx:1.25", "nginx:1.19"} {
		if !strings.Contains(line, part) {
			t.Errorf("line %q missing part %q", line, part)
		}
	}
}

func TestFormat_NoDrift(t *testing.T) {
	out := diff.Format(nil)
	if out != "no drift detected" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithDrift(t *testing.T) {
	summaries := diff.Summarize(sampleDrifts())
	out := diff.Format(summaries)
	if !strings.Contains(out, "auth-service") {
		t.Errorf("expected service name in output, got: %s", out)
	}
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestSummarize_EmptyInput(t *testing.T) {
	summaries := diff.Summarize([]drift.DriftResult{})
	if len(summaries) != 0 {
		t.Errorf("expected empty summaries, got %d", len(summaries))
	}
}
