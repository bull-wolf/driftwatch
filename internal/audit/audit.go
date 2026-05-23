// Package audit provides functionality for recording and retrieving
// audit events related to drift detection runs and policy evaluations.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// EventKind classifies the type of audit event.
type EventKind string

const (
	EventKindDriftDetected  EventKind = "drift_detected"
	EventKindPolicyViolated EventKind = "policy_violated"
	EventKindClean         EventKind = "clean"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp   time.Time `json:"timestamp"`
	ServiceName string    `json:"service_name"`
	Kind        EventKind `json:"kind"`
	Details     string    `json:"details"`
}

// Record appends an audit event to the service's audit log file.
// Events are stored as newline-delimited JSON under dir/<service>.audit.jsonl.
func Record(dir, serviceName string, kind EventKind, details string) error {
	if serviceName == "" {
		return fmt.Errorf("audit: service name must not be empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("audit: create directory: %w", err)
	}
	event := Event{
		Timestamp:   time.Now().UTC(),
		ServiceName: serviceName,
		Kind:        kind,
		Details:     details,
	}
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	path := filepath.Join(dir, serviceName+".audit.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// Read returns all audit events recorded for the given service.
// Returns an empty slice if no audit file exists yet.
func Read(dir, serviceName string) ([]Event, error) {
	path := filepath.Join(dir, serviceName+".audit.jsonl")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Event{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read file: %w", err)
	}
	var events []Event
	for _, line := range splitLines(string(data)) {
		var e Event
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("audit: parse line: %w", err)
		}
		events = append(events, e)
	}
	return events, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if line := s[start:i]; line != "" {
				lines = append(lines, line)
			}
			start = i + 1
		}
	}
	if tail := s[start:]; tail != "" {
		lines = append(lines, tail)
	}
	return lines
}
