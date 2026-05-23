// Package compare provides field-level comparison between expected manifest
// state and observed live service configuration.
//
// # Overview
//
// Use [Fields] to compare two string maps representing expected vs actual
// service configuration. The returned [Result] contains all fields that
// differ, along with their expected and actual values.
//
// # Example
//
//	result := compare.Fields("auth-service",
//		map[string]string{"image": "nginx:1.21", "replicas": "3"},
//		map[string]string{"image": "nginx:1.99", "replicas": "3"},
//	)
//	if result.IsDrifted() {
//		fmt.Println(result.Summary())
//	}
package compare
