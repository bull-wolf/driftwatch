package drift_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/manifest"
)

func baseManifest() manifest.Manifest {
	return manifest.Manifest{
		Name:        "auth-service",
		Image:       "auth-service:v1.2.0",
		Replicas:    3,
		Environment: "production",
	}
}

func TestDetect_NoDrift(t *testing.T) {
	m := baseManifest()
	deployed := drift.DeployedState{
		Name:        "auth-service",
		Image:       "auth-service:v1.2.0",
		Replicas:    3,
		Environment: "production",
	}

	result := drift.Detect(deployed, m)

	if result.HasDrift {
		t.Errorf("expected no drift, got diffs: %v", result.Diffs)
	}
	if len(result.Diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(result.Diffs))
	}
}

func TestDetect_ImageDrift(t *testing.T) {
	m := baseManifest()
	deployed := drift.DeployedState{
		Name:        "auth-service",
		Image:       "auth-service:v1.1.0",
		Replicas:    3,
		Environment: "production",
	}

	result := drift.Detect(deployed, m)

	if !result.HasDrift {
		t.Error("expected drift to be detected")
	}
	if len(result.Diffs) != 1 {
		t.Errorf("expected 1 diff, got %d: %v", len(result.Diffs), result.Diffs)
	}
}

func TestDetect_MultipleDrifts(t *testing.T) {
	m := baseManifest()
	deployed := drift.DeployedState{
		Name:        "auth-service",
		Image:       "auth-service:v0.9.0",
		Replicas:    1,
		Environment: "staging",
	}

	result := drift.Detect(deployed, m)

	if !result.HasDrift {
		t.Error("expected drift to be detected")
	}
	if len(result.Diffs) != 3 {
		t.Errorf("expected 3 diffs, got %d: %v", len(result.Diffs), result.Diffs)
	}
}

func TestDetect_ServiceNameInResult(t *testing.T) {
	m := baseManifest()
	deployed := drift.DeployedState{
		Name:        "auth-service",
		Image:       m.Image,
		Replicas:    m.Replicas,
		Environment: m.Environment,
	}

	result := drift.Detect(deployed, m)

	if result.ServiceName != "auth-service" {
		t.Errorf("expected service name %q, got %q", "auth-service", result.ServiceName)
	}
}
