package tag_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/tag"
)

func TestParseTags_RoundTrip_FromEnvString(t *testing.T) {
	raw := "env=prod,team=platform,region=us-east-1"
	tags, err := tag.ParseTags(raw)
	if err != nil {
		t.Fatalf("ParseTags error: %v", err)
	}
	if len(tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tags))
	}
}

func TestApply_TagsAreCopiedPerResult(t *testing.T) {
	results := []drift.Result{
		{Service: "auth-service", Field: "image"},
		{Service: "billing-service", Field: "replicas"},
	}
	original := tag.Tags{"env": "prod"}
	tagged := tag.Apply(results, original)

	// Mutate original; tagged results must not be affected.
	original["env"] = "staging"

	for _, tr := range tagged {
		if tr.Tags["env"] != "prod" {
			t.Errorf("tag isolation failed: got %q", tr.Tags["env"])
		}
	}
}

func TestSampleTagsFile_Exists(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "tags", "sample-tags.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read sample-tags.yaml: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "env: prod") {
		t.Errorf("expected 'env: prod' in sample-tags.yaml")
	}
	if !strings.Contains(content, "team: platform") {
		t.Errorf("expected 'team: platform' in sample-tags.yaml")
	}
}

func TestFilter_MultiTagMatch(t *testing.T) {
	results := []drift.Result{
		{Service: "auth-service", Field: "image"},
	}
	tags, _ := tag.ParseTags("env=prod,team=platform")
	tagged := tag.Apply(results, tags)

	matched := tag.Filter(tagged, tag.Tags{"env": "prod", "team": "platform"})
	if len(matched) != 1 {
		t.Errorf("expected 1 match, got %d", len(matched))
	}

	noMatch := tag.Filter(tagged, tag.Tags{"env": "prod", "team": "security"})
	if len(noMatch) != 0 {
		t.Errorf("expected 0 matches, got %d", len(noMatch))
	}
}
