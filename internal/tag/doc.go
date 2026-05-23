// Package tag provides labelling support for drift results.
//
// Tags are arbitrary key=value metadata pairs (e.g. env=prod, team=platform)
// that can be attached to drift results to enable downstream filtering,
// routing, or reporting based on organisational context.
//
// Usage:
//
//	tags, err := tag.ParseTags("env=prod,team=platform")
//	tagged := tag.Apply(results, tags)
//	filtered := tag.Filter(tagged, tag.Tags{"env": "prod"})
package tag
