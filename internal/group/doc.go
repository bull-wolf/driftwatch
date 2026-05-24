// Package group provides grouping utilities for drift results.
//
// Results produced by the drift detector can be organised along
// three dimensions:
//
//   - ByService  — one bucket per service name
//   - ByField    — one bucket per inspected field (image, replicas, …)
//   - BySeverity — two buckets: "drifted" and "clean"
//
// Example:
//
//	res, err := group.Apply(drifts, group.ByField)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, key := range res.Keys() {
//		fmt.Printf("%s: %d results\n", key, len(res.Groups[key]))
//	}
package group
