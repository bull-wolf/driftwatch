package policy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/policy"
)

func writeTempPolicy(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp policy: %v", err)
	}
	return path
}

func TestLoad_ValidPolicy(t *testing.T) {
	path := writeTempPolicy(t, `
rules:
  - name: no-latest-image
    field: image
    severity: error
    deny_list:
      - "nginx:latest"
`)
	p, err := policy.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(p.Rules))
	}
	if p.Rules[0].Name != "no-latest-image" {
		t.Errorf("expected rule name 'no-latest-image', got %q", p.Rules[0].Name)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := policy.Load("/nonexistent/policy.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestEvaluate_NoViolations(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{Name: "no-latest", Field: "image", Severity: "error", DenyList: []string{"nginx:latest"}},
		},
	}
	violations := p.Evaluate("auth-service", map[string]string{"image": "nginx:1.25"})
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestEvaluate_WithViolation(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{Name: "no-latest", Field: "image", Severity: "error", DenyList: []string{"nginx:latest"}},
		},
	}
	violations := p.Evaluate("auth-service", map[string]string{"image": "nginx:latest"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Severity != "error" {
		t.Errorf("expected severity 'error', got %q", violations[0].Severity)
	}
	if violations[0].Service != "auth-service" {
		t.Errorf("expected service 'auth-service', got %q", violations[0].Service)
	}
}

func TestEvaluate_UndriftedFieldSkipped(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{Name: "no-latest", Field: "image", Severity: "warn", DenyList: []string{"nginx:latest"}},
		},
	}
	// Field not in drifted map — should produce no violations
	violations := p.Evaluate("svc", map[string]string{"replicas": "3"})
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}
