// Package resolve provides field resolution utilities for live service state.
//
// It allows callers to extract individual field values from a ServiceState
// by name, supporting top-level fields (name, image, replicas) as well as
// namespaced fields for environment variables (env.<KEY>) and labels (label.<KEY>).
//
// Example usage:
//
//	state := resolve.ServiceState{
//		Name:     "auth-service",
//		Image:    "auth:v2.0.0",
//		Replicas: 2,
//	}
//
//	val, err := resolve.Field(state, "image")
//	// val == "auth:v2.0.0"
//
//	all := resolve.AllFields(state)
//	// all["replicas"] == "2"
package resolve
