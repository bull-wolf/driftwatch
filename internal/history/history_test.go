package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/history"
)

func TestRecord_AndRead_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	drifts := []drift.Drift{
		{Field: "image", Deployed: "nginx:1.19", Expected: "nginx:1.21"},
	}

	if err := history.Record(dir, "auth-service", drifts); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := history.Read(dir, "auth-service")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.Service != "auth-service" {
		t.Errorf("service = %q, want auth-service", e.Service)
	}
	if !e.HasDrift {
		t.Error("expected HasDrift = true")
	}
	if len(e.Drifts) != 1 || e.Drifts[0].Field != "image" {
		t.Errorf("unexpected drifts: %+v", e.Drifts)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	dir := t.TempDir()

	for i := 0; i < 3; i++ {
		if err := history.Record(dir, "svc", nil); err != nil {
			t.Fatalf("Record iteration %d: %v", i, err)
		}
	}

	entries, err := history.Read(dir, "svc")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestRecord_NoDrift(t *testing.T) {
	dir := t.TempDir()

	if err := history.Record(dir, "svc", []drift.Drift{}); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := history.Read(dir, "svc")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if entries[0].HasDrift {
		t.Error("expected HasDrift = false for empty drifts")
	}
}

func TestRecord_EmptyServiceName(t *testing.T) {
	dir := t.TempDir()
	err := history.Record(dir, "", nil)
	if err == nil {
		t.Error("expected error for empty service name")
	}
}

func TestRead_NoFile(t *testing.T) {
	dir := t.TempDir()
	entries, err := history.Read(dir, "nonexistent")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecord_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/history"

	if err := history.Record(dir, "svc", nil); err != nil {
		t.Fatalf("Record: %v", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}

	_ = time.Now() // suppress import warning
}
