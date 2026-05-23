// Package policy provides rule-based evaluation of detected configuration drift.
//
// A Policy is loaded from a YAML file and contains a list of Rules. Each Rule
// targets a specific manifest field and defines conditions that constitute a
// policy violation — for example, disallowing "latest" image tags.
//
// Usage:
//
//	p, err := policy.Load("policy.yaml")
//	if err != nil { ... }
//
//	violations := p.Evaluate("my-service", map[string]string{
//		"image": "nginx:latest",
//	})
//	for _, v := range violations {
//		fmt.Printf("[%s] %s: %s\n", v.Severity, v.RuleName, v.Message)
//	}
//
// Violations carry a Severity of either "warn" or "error", allowing callers
// to decide whether to surface them as advisory notices or hard failures.
package policy
