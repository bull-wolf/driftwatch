package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a single drift check recorded in history.
type Entry struct {
	Timestamp time.Time    `json:"timestamp"`
	Service   string       `json:"service"`
	Drifts    []drift.Drift `json:"drifts"`
	HasDrift  bool         `json:"has_drift"`
}

// Record appends a new history entry for the given service and drifts.
// Entries are stored as newline-delimited JSON in a per-service file.
func Record(dir, service string, drifts []drift.Drift) error {
	if service == "" {
		return fmt.Errorf("service name must not be empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create history dir: %w", err)
	}

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Service:   service,
		Drifts:    drifts,
		HasDrift:  len(drifts) > 0,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}

	path := filepath.Join(dir, service+".history.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open history file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// Read returns all history entries for the given service.
func Read(dir, service string) ([]Entry, error) {
	path := filepath.Join(dir, service+".history.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("read history file: %w", err)
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("unmarshal entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
