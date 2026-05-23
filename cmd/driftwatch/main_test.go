package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the driftwatch binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "driftwatch")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return binPath
}

func TestMain_NoManifestFlag(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code when --manifest is missing")
	}
	if !strings.Contains(string(out), "--manifest flag is required") {
		t.Errorf("expected error message about --manifest, got: %s", out)
	}
}

func TestMain_FileNotFound(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "--manifest", "/nonexistent/path.yaml")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code for missing manifest file")
	}
	if !strings.Contains(string(out), "error loading manifest") {
		t.Errorf("expected loading error message, got: %s", out)
	}
}

func TestMain_ValidManifestTextOutput(t *testing.T) {
	bin := buildBinary(t)
	manifestPath := filepath.Join("..", "..", "testdata", "manifests", "auth-service.yaml")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Skip("testdata manifest not found, skipping integration test")
	}
	cmd := exec.Command(bin, "--manifest", manifestPath, "--format", "text")
	out, err := cmd.CombinedOutput()
	// Exit code 0 (no drift) or 2 (drift found) are both acceptable outcomes.
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok || exitErr.ExitCode() != 2 {
			t.Fatalf("unexpected error running driftwatch: %v\n%s", err, out)
		}
	}
}

func TestMain_InvalidFormat(t *testing.T) {
	bin := buildBinary(t)
	manifestPath := filepath.Join("..", "..", "testdata", "manifests", "auth-service.yaml")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Skip("testdata manifest not found, skipping integration test")
	}
	cmd := exec.Command(bin, "--manifest", manifestPath, "--format", "xml")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code for invalid format")
	}
	if !strings.Contains(string(out), "unknown format") {
		t.Errorf("expected unknown format error, got: %s", out)
	}
}
