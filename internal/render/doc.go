// Package render provides Go-template-based rendering for drift detection reports.
//
// Usage:
//
//	const tmpl = `Service: {{.Service}}\nDrifted: {{.HasDrift}}`
//
//	r, err := render.New(tmpl)
//	if err != nil { ... }
//
//	var buf bytes.Buffer
//	if err := r.Render(&buf, "auth-service", drifts); err != nil { ... }
//
// ReportData fields available in templates:
//
//	.Service   — name of the service being reported
//	.Timestamp — UTC time the report was generated
//	.Drifts    — slice of drift.Result values
//	.HasDrift  — bool, true when len(Drifts) > 0
package render
