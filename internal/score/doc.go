// Package score provides a numeric health scoring system for driftwatch services.
//
// A score of 100 indicates a perfectly healthy service with no detected drift
// and no policy violations. Each drifted field reduces the score by DriftPenalty
// points, and each policy violation reduces it by ViolationPenalty points.
// The score is clamped to a minimum of 0.
//
// Usage:
//
//	result := score.Compute("auth-service", drifts, violations)
//	fmt.Printf("Score: %d (%s)\n", result.Score, score.Grade(result.Score))
package score
