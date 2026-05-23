package baseline_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/drift"
)

func sampleDrifts() []drift.Drift {
	return []drift.Drift{
		{ServiceName: "auth-service", Field: "image", Expected: "auth:v1", Actual: "auth:v2"},
		{ServiceName: "auth-service", Field: "replicas", Expected: "3", Actual: "1"},
	}
}

func TestSave_AndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	drifts := sampleDrifts()

	if err := baseline.Save(dir, "auth-service", drifts); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	entry, err := baseline.Load(dir, "auth-service")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if entry.ServiceName != "auth-service" {
		t.Errorf("ServiceName = %q, want %q", entry.ServiceName, "auth-service")
	}
	if len(entry.Drifts) != len(drifts) {
		t.Fatalf("len(Drifts) = %d, want %d", len(entry.Drifts), len(drifts))
	}
	if entry.Drifts[0].Field != "image" {
		t.Errorf("Drifts[0].Field = %q, want %q", entry.Drifts[0].Field, "image")
	}
	if entry.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
	if entry.CapturedAt.After(time.Now().Add(time.Second)) {
		t.Error("CapturedAt is in the future")
	}
}

func TestSave_EmptyServiceName(t *testing.T) {
	dir := t.TempDir()
	if err := baseline.Save(dir, "", nil); err == nil {
		t.Error("expected error for empty service name, got nil")
	}
}

func TestLoad_NoBaseline(t *testing.T) {
	dir := t.TempDir()
	_, err := baseline.Load(dir, "missing-service")
	if err == nil {
		t.Error("expected error when no baseline exists, got nil")
	}
}

func TestLoad_EmptyServiceName(t *testing.T) {
	dir := t.TempDir()
	_, err := baseline.Load(dir, "")
	if err == nil {
		t.Error("expected error for empty service name, got nil")
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	nestedDir := dir + "/nested/baselines"

	if err := baseline.Save(nestedDir, "auth-service", sampleDrifts()); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Error("expected nested directory to be created")
	}
}

func TestSave_NoDrifts(t *testing.T) {
	dir := t.TempDir()

	if err := baseline.Save(dir, "auth-service", []drift.Drift{}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	entry, err := baseline.Load(dir, "auth-service")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(entry.Drifts) != 0 {
		t.Errorf("expected 0 drifts, got %d", len(entry.Drifts))
	}
}
