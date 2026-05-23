package tag_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/tag"
)

func sampleResults() []drift.Result {
	return []drift.Result{
		{Service: "auth-service", Field: "image", Expected: "v1", Actual: "v2"},
		{Service: "payment-service", Field: "replicas", Expected: "3", Actual: "1"},
	}
}

func TestApply_AttachesTagsToAllResults(t *testing.T) {
	results := sampleResults()
	tags := tag.Tags{"env": "prod", "team": "platform"}

	tagged := tag.Apply(results, tags)

	if len(tagged) != len(results) {
		t.Fatalf("expected %d tagged results, got %d", len(results), len(tagged))
	}
	for _, tr := range tagged {
		if tr.Tags["env"] != "prod" {
			t.Errorf("expected env=prod, got %q", tr.Tags["env"])
		}
		if tr.Tags["team"] != "platform" {
			t.Errorf("expected team=platform, got %q", tr.Tags["team"])
		}
	}
}

func TestApply_EmptyResults(t *testing.T) {
	tagged := tag.Apply(nil, tag.Tags{"env": "staging"})
	if len(tagged) != 0 {
		t.Errorf("expected empty slice, got %d", len(tagged))
	}
}

func TestFilter_MatchingTag(t *testing.T) {
	results := sampleResults()
	tagged := tag.Apply(results, tag.Tags{"env": "prod"})

	matched := tag.Filter(tagged, tag.Tags{"env": "prod"})
	if len(matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(matched))
	}
}

func TestFilter_NoMatchingTag(t *testing.T) {
	results := sampleResults()
	tagged := tag.Apply(results, tag.Tags{"env": "prod"})

	matched := tag.Filter(tagged, tag.Tags{"env": "staging"})
	if len(matched) != 0 {
		t.Errorf("expected 0 matched, got %d", len(matched))
	}
}

func TestParseTags_Valid(t *testing.T) {
	tags, err := tag.ParseTags("env=prod,team=platform")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tags["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", tags["env"])
	}
	if tags["team"] != "platform" {
		t.Errorf("expected team=platform, got %q", tags["team"])
	}
}

func TestParseTags_Empty(t *testing.T) {
	tags, err := tag.ParseTags("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected empty tags, got %v", tags)
	}
}

func TestParseTags_Invalid(t *testing.T) {
	_, err := tag.ParseTags("env-prod")
	if err == nil {
		t.Error("expected error for malformed tag, got nil")
	}
}
