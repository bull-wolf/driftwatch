// Package audit provides append-only audit logging for driftwatch events.
//
// Each service gets its own audit log file stored as newline-delimited JSON
// (JSONL) under a configurable directory. Events capture the service name,
// a timestamp, an event kind, and a freeform details string.
//
// Supported event kinds:
//
//	- EventKindDriftDetected  — one or more fields drifted from the manifest
//	- EventKindPolicyViolated — a policy rule was violated during evaluation
//	- EventKindClean          — no drift or violations were found
//
// Usage:
//
//	err := audit.Record("/var/driftwatch/audit", "auth-service",
//		audit.EventKindDriftDetected, "image: expected nginx:1.25, got nginx:1.24")
//
//	events, err := audit.Read("/var/driftwatch/audit", "auth-service")
package audit
