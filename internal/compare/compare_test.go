package compare_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/compare"
)

func TestFields_NoDrift(t *testing.T) {
	expected := map[string]string{"image": "nginx:1.21", "replicas": "3"}
	actual := map[string]string{"image": "nginx:1.21", "replicas": "3"}

	result := compare.Fields("web", expected, actual)

	if result.IsDrifted() {
		t.Errorf("expected no drift, got %d drifted fields", len(result.Drifted))
	}
}

func TestFields_SingleDrift(t *testing.T) {
	expected := map[string]string{"image": "nginx:1.21"}
	actual := map[string]string{"image": "nginx:1.99"}

	result := compare.Fields("web", expected, actual)

	if !result.IsDrifted() {
		t.Fatal("expected drift to be detected")
	}
	if len(result.Drifted) != 1 {
		t.Fatalf("expected 1 drifted field, got %d", len(result.Drifted))
	}
	if result.Drifted[0].Name != "image" {
		t.Errorf("expected field name 'image', got %q", result.Drifted[0].Name)
	}
	if result.Drifted[0].Expected != "nginx:1.21" {
		t.Errorf("unexpected Expected value: %q", result.Drifted[0].Expected)
	}
	if result.Drifted[0].Actual != "nginx:1.99" {
		t.Errorf("unexpected Actual value: %q", result.Drifted[0].Actual)
	}
}

func TestFields_MissingActualKey(t *testing.T) {
	expected := map[string]string{"replicas": "2"}
	actual := map[string]string{}

	result := compare.Fields("api", expected, actual)

	if !result.IsDrifted() {
		t.Fatal("expected drift when key is missing from actual")
	}
	if result.Drifted[0].Actual != "" {
		t.Errorf("expected empty actual value for missing key, got %q", result.Drifted[0].Actual)
	}
}

func TestResult_Summary_NoDrift(t *testing.T) {
	r := compare.Result{Service: "auth"}
	s := r.Summary()
	if !strings.Contains(s, "no drift") {
		t.Errorf("expected summary to mention 'no drift', got %q", s)
	}
}

func TestResult_Summary_WithDrift(t *testing.T) {
	r := compare.Result{
		Service: "auth",
		Drifted: []compare.Field{
			{Name: "image", Expected: "v1", Actual: "v2"},
			{Name: "replicas", Expected: "3", Actual: "1"},
		},
	}
	s := r.Summary()
	if !strings.Contains(s, "auth") {
		t.Errorf("expected summary to contain service name, got %q", s)
	}
	if !strings.Contains(s, "image") || !strings.Contains(s, "replicas") {
		t.Errorf("expected summary to list drifted fields, got %q", s)
	}
}
