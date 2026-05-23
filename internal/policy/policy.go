package policy

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Rule defines a single policy rule that flags specific drift conditions.
type Rule struct {
	Name     string   `yaml:"name"`
	Field    string   `yaml:"field"`
	Severity string   `yaml:"severity"` // "warn" or "error"
	DenyList []string `yaml:"deny_list,omitempty"`
	Required bool     `yaml:"required,omitempty"`
}

// Policy holds a collection of rules loaded from a policy file.
type Policy struct {
	Rules []Rule `yaml:"rules"`
}

// Violation represents a policy rule that was triggered by a drift result.
type Violation struct {
	Service  string
	Field    string
	RuleName string
	Severity string
	Message  string
}

// Load reads and parses a YAML policy file from the given path.
func Load(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("policy: read file: %w", err)
	}
	var p Policy
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("policy: parse yaml: %w", err)
	}
	return &p, nil
}

// Evaluate checks a set of drift fields against the policy rules and returns violations.
func (p *Policy) Evaluate(service string, driftedFields map[string]string) []Violation {
	var violations []Violation
	for _, rule := range p.Rules {
		actual, drifted := driftedFields[rule.Field]
		if rule.Required && !drifted {
			continue
		}
		if !drifted {
			continue
		}
		for _, denied := range rule.DenyList {
			if strings.EqualFold(actual, denied) {
				violations = append(violations, Violation{
					Service:  service,
					Field:    rule.Field,
					RuleName: rule.Name,
					Severity: rule.Severity,
					Message:  fmt.Sprintf("field %q has disallowed value %q", rule.Field, actual),
				})
			}
		}
	}
	return violations
}
