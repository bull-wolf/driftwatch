package policy_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/policy"
)

func TestLoad_DefaultPolicyFile(t *testing.T) {
	p, err := policy.Load("../../testdata/policies/default-policy.yaml")
	if err != nil {
		t.Fatalf("failed to load default policy: %v", err)
	}
	if len(p.Rules) == 0 {
		t.Fatal("expected at least one rule in default policy")
	}
}

func TestEvaluate_DefaultPolicy_ImageViolation(t *testing.T) {
	p, err := policy.Load("../../testdata/policies/default-policy.yaml")
	if err != nil {
		t.Fatalf("failed to load default policy: %v", err)
	}

	drifted := map[string]string{
		"image": "nginx:latest",
	}
	violations := p.Evaluate("auth-service", drifted)
	if len(violations) == 0 {
		t.Fatal("expected at least one violation for nginx:latest image")
	}
	found := false
	for _, v := range violations {
		if v.RuleName == "no-latest-image" && v.Severity == "error" {
			found = true
		}
	}
	if !found {
		t.Error("expected violation with rule 'no-latest-image' and severity 'error'")
	}
}

func TestEvaluate_DefaultPolicy_NoViolations(t *testing.T) {
	p, err := policy.Load("../../testdata/policies/default-policy.yaml")
	if err != nil {
		t.Fatalf("failed to load default policy: %v", err)
	}

	drifted := map[string]string{
		"image": "nginx:1.25.3",
	}
	violations := p.Evaluate("auth-service", drifted)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}
