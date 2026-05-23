package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/manifest"
)

func writeTempManifest(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp manifest: %v", err)
	}
	return path
}

func TestLoad_ValidManifest(t *testing.T) {
	content := `
name: auth-service
version: "1.2.3"
image: registry.example.com/auth-service:1.2.3
replicas: 3
environment:
  LOG_LEVEL: info
  DB_HOST: postgres.internal
ports:
  - 8080
labels:
  team: platform
`
	path := writeTempManifest(t, content)

	m, err := manifest.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.Name != "auth-service" {
		t.Errorf("expected name %q, got %q", "auth-service", m.Name)
	}
	if m.Replicas != 3 {
		t.Errorf("expected replicas 3, got %d", m.Replicas)
	}
	if m.Environment["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", m.Environment["LOG_LEVEL"])
	}
}

func TestLoad_MissingName(t *testing.T) {
	content := `image: registry.example.com/svc:latest
replicas: 1
`
	path := writeTempManifest(t, content)

	_, err := manifest.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing name, got nil")
	}
}

func TestLoad_NegativeReplicas(t *testing.T) {
	content := `name: bad-service
image: registry.example.com/svc:latest
replicas: -1
`
	path := writeTempManifest(t, content)

	_, err := manifest.Load(path)
	if err == nil {
		t.Fatal("expected validation error for negative replicas, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := manifest.Load("/nonexistent/path/manifest.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
