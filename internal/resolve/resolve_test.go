package resolve_test

import (
	"testing"

	"github.com/driftwatch/internal/resolve"
)

func sampleState() resolve.ServiceState {
	return resolve.ServiceState{
		Name:     "auth-service",
		Image:    "auth:v1.2.3",
		Replicas: 3,
		Environment: map[string]string{
			"LOG_LEVEL": "info",
			"PORT":      "8080",
		},
		Labels: map[string]string{
			"team": "platform",
		},
	}
}

func TestField_Name(t *testing.T) {
	val, err := resolve.Field(sampleState(), "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "auth-service" {
		t.Errorf("expected auth-service, got %s", val)
	}
}

func TestField_Image(t *testing.T) {
	val, err := resolve.Field(sampleState(), "image")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "auth:v1.2.3" {
		t.Errorf("expected auth:v1.2.3, got %s", val)
	}
}

func TestField_Replicas(t *testing.T) {
	val, err := resolve.Field(sampleState(), "replicas")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "3" {
		t.Errorf("expected 3, got %s", val)
	}
}

func TestField_EnvKey(t *testing.T) {
	val, err := resolve.Field(sampleState(), "env.LOG_LEVEL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "info" {
		t.Errorf("expected info, got %s", val)
	}
}

func TestField_LabelKey(t *testing.T) {
	val, err := resolve.Field(sampleState(), "label.team")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "platform" {
		t.Errorf("expected platform, got %s", val)
	}
}

func TestField_MissingEnvKey(t *testing.T) {
	_, err := resolve.Field(sampleState(), "env.MISSING")
	if err == nil {
		t.Fatal("expected error for missing env key")
	}
}

func TestField_UnsupportedField(t *testing.T) {
	_, err := resolve.Field(sampleState(), "unknown")
	if err == nil {
		t.Fatal("expected error for unsupported field")
	}
}

func TestAllFields_ContainsTopLevel(t *testing.T) {
	fields := resolve.AllFields(sampleState())
	for _, key := range []string{"name", "image", "replicas", "env.LOG_LEVEL", "label.team"} {
		if _, ok := fields[key]; !ok {
			t.Errorf("expected field %q to be present", key)
		}
	}
}
