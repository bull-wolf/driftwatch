package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceManifest represents the source-of-truth configuration for a service.
type ServiceManifest struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Image       string            `yaml:"image"`
	Replicas    int               `yaml:"replicas"`
	Environment map[string]string `yaml:"environment"`
	Ports       []int             `yaml:"ports"`
	Labels      map[string]string `yaml:"labels"`
}

// Load reads and parses a manifest file from the given path.
func Load(path string) (*ServiceManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest file %q: %w", path, err)
	}

	var m ServiceManifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest file %q: %w", path, err)
	}

	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("invalid manifest %q: %w", path, err)
	}

	return &m, nil
}

// Validate checks that required fields are present and valid.
func (m *ServiceManifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("manifest must have a non-empty name")
	}
	if m.Image == "" {
		return fmt.Errorf("manifest must have a non-empty image")
	}
	if m.Replicas < 0 {
		return fmt.Errorf("replicas must be >= 0, got %d", m.Replicas)
	}
	return nil
}
