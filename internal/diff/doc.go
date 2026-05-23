// Package diff provides utilities for generating human-readable summaries
// of configuration drift between expected manifest values and actual
// deployed service state.
//
// # Overview
//
// Given a slice of drift.DriftResult values produced by the detector,
// the diff package formats each entry into a structured Summary that
// captures the service name, field, expected value, and actual value.
//
// # Usage
//
//	summaries := diff.Summarize(drifts)
//	fmt.Println(diff.Format(summaries))
//
// Output example:
//
//	[auth-service] image: expected="nginx:1.25" actual="nginx:1.19"
//	[auth-service] replicas: expected="3" actual="1"
package diff
