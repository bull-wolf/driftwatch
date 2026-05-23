package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a saved baseline for a service at a point in time.
type Entry struct {
	ServiceName string        `json:"service_name"`
	CapturedAt  time.Time     `json:"captured_at"`
	Drifts      []drift.Drift `json:"drifts"`
}

// Save writes the current drift results as a baseline for the given service.
// Baselines are stored under dir/<service>.baseline.json.
func Save(dir, serviceName string, drifts []drift.Drift) error {
	if serviceName == "" {
		return fmt.Errorf("service name must not be empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create baseline dir: %w", err)
	}

	entry := Entry{
		ServiceName: serviceName,
		CapturedAt:  time.Now().UTC(),
		Drifts:      drifts,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal baseline: %w", err)
	}

	path := filepath.Join(dir, serviceName+".baseline.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write baseline file: %w", err)
	}
	return nil
}

// Load reads the baseline entry for the given service from dir.
// Returns an error if no baseline exists.
func Load(dir, serviceName string) (*Entry, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("service name must not be empty")
	}

	path := filepath.Join(dir, serviceName+".baseline.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no baseline found for service %q", serviceName)
		}
		return nil, fmt.Errorf("read baseline file: %w", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("unmarshal baseline: %w", err)
	}
	return &entry, nil
}
