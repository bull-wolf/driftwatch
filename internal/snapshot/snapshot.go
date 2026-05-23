package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of a deployed service state.
type Snapshot struct {
	ServiceName string            `json:"service_name"`
	Image       string            `json:"image"`
	Replicas    int               `json:"replicas"`
	Env         map[string]string `json:"env,omitempty"`
	CapturedAt  time.Time         `json:"captured_at"`
}

// Save writes a snapshot to disk as JSON under the given directory.
// The file is named <service_name>.snapshot.json.
func Save(dir string, s Snapshot) error {
	if s.ServiceName == "" {
		return fmt.Errorf("snapshot: service_name must not be empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("snapshot: create directory: %w", err)
	}
	s.CapturedAt = time.Now().UTC()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := filepath.Join(dir, s.ServiceName+".snapshot.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write file: %w", err)
	}
	return nil
}

// Load reads a snapshot from disk for the given service name.
func Load(dir, serviceName string) (Snapshot, error) {
	path := filepath.Join(dir, serviceName+".snapshot.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read file %q: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return s, nil
}
