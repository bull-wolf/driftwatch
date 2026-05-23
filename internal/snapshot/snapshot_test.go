package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/snapshot"
)

func TestSave_AndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()

	orig := snapshot.Snapshot{
		ServiceName: "auth-service",
		Image:       "auth:v1.2.3",
		Replicas:    3,
		Env:         map[string]string{"LOG_LEVEL": "info"},
	}

	if err := snapshot.Save(dir, orig); err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	got, err := snapshot.Load(dir, "auth-service")
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if got.ServiceName != orig.ServiceName {
		t.Errorf("ServiceName: got %q, want %q", got.ServiceName, orig.ServiceName)
	}
	if got.Image != orig.Image {
		t.Errorf("Image: got %q, want %q", got.Image, orig.Image)
	}
	if got.Replicas != orig.Replicas {
		t.Errorf("Replicas: got %d, want %d", got.Replicas, orig.Replicas)
	}
	if got.Env["LOG_LEVEL"] != "info" {
		t.Errorf("Env[LOG_LEVEL]: got %q, want %q", got.Env["LOG_LEVEL"], "info")
	}
	if got.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero after Save")
	}
	if got.CapturedAt.After(time.Now().UTC()) {
		t.Error("CapturedAt should not be in the future")
	}
}

func TestSave_EmptyServiceName(t *testing.T) {
	dir := t.TempDir()
	err := snapshot.Save(dir, snapshot.Snapshot{Image: "img:latest"})
	if err == nil {
		t.Fatal("expected error for empty service_name, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := snapshot.Load(dir, "nonexistent-service")
	if err == nil {
		t.Fatal("expected error for missing snapshot file, got nil")
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/snapshots"

	s := snapshot.Snapshot{ServiceName: "billing", Image: "billing:v2", Replicas: 1}
	if err := snapshot.Save(dir, s); err != nil {
		t.Fatalf("Save() should create missing directories: %v", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
