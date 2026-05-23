// Package tag provides utilities for tagging drift results with
// metadata labels such as environment, team, or severity.
package tag

import (
	"fmt"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Tags is a map of key-value label pairs.
type Tags map[string]string

// TaggedResult wraps a drift result with associated metadata tags.
type TaggedResult struct {
	Drift drift.Result
	Tags  Tags
}

// Apply attaches the provided tags to each drift result and returns
// a slice of TaggedResult values.
func Apply(results []drift.Result, tags Tags) []TaggedResult {
	tagged := make([]TaggedResult, 0, len(results))
	for _, r := range results {
		tagged = append(tagged, TaggedResult{
			Drift: r,
			Tags:  copyTags(tags),
		})
	}
	return tagged
}

// Filter returns only those TaggedResults where all provided tags match.
func Filter(results []TaggedResult, match Tags) []TaggedResult {
	var out []TaggedResult
	for _, r := range results {
		if matchesTags(r.Tags, match) {
			out = append(out, r)
		}
	}
	return out
}

// ParseTags parses a comma-separated list of key=value pairs into Tags.
// Example: "env=prod,team=platform"
func ParseTags(raw string) (Tags, error) {
	tags := Tags{}
	if strings.TrimSpace(raw) == "" {
		return tags, nil
	}
	for _, part := range strings.Split(raw, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 || strings.TrimSpace(kv[0]) == "" {
			return nil, fmt.Errorf("invalid tag format %q: expected key=value", part)
		}
		tags[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return tags, nil
}

func copyTags(src Tags) Tags {
	dst := make(Tags, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func matchesTags(have, want Tags) bool {
	for k, v := range want {
		if have[k] != v {
			return false
		}
	}
	return true
}
