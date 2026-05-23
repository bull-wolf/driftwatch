package audit_test

import (
	"os"
	"testing"

	"github.com/driftwatch/internal/audit"
)

func TestRecord_AndRead_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	err := audit.Record(dir, "auth-service", audit.EventKindDriftDetected, "image changed")
	if err != nil {
		t.Fatalf("Record: unexpected error: %v", err)
	}
	events, err := audit.Read(dir, "auth-service")
	if err != nil {
		t.Fatalf("Read: unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].ServiceName != "auth-service" {
		t.Errorf("ServiceName: got %q, want %q", events[0].ServiceName, "auth-service")
	}
	if events[0].Kind != audit.EventKindDriftDetected {
		t.Errorf("Kind: got %q, want %q", events[0].Kind, audit.EventKindDriftDetected)
	}
	if events[0].Details != "image changed" {
		t.Errorf("Details: got %q, want %q", events[0].Details, "image changed")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	kinds := []audit.EventKind{audit.EventKindClean, audit.EventKindPolicyViolated, audit.EventKindDriftDetected}
	for _, k := range kinds {
		if err := audit.Record(dir, "svc", k, "detail"); err != nil {
			t.Fatalf("Record(%s): %v", k, err)
		}
	}
	events, err := audit.Read(dir, "svc")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	for i, k := range kinds {
		if events[i].Kind != k {
			t.Errorf("event[%d].Kind: got %q, want %q", i, events[i].Kind, k)
		}
	}
}

func TestRecord_EmptyServiceName(t *testing.T) {
	dir := t.TempDir()
	err := audit.Record(dir, "", audit.EventKindClean, "")
	if err == nil {
		t.Fatal("expected error for empty service name, got nil")
	}
}

func TestRead_NoFile(t *testing.T) {
	dir := t.TempDir()
	events, err := audit.Read(dir, "nonexistent-service")
	if err != nil {
		t.Fatalf("Read: unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestRecord_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/audit"
	if err := audit.Record(dir, "svc", audit.EventKindClean, ""); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("expected directory %q to be created", dir)
	}
}

func TestRecord_TimestampIsSet(t *testing.T) {
	dir := t.TempDir()
	if err := audit.Record(dir, "svc", audit.EventKindClean, ""); err != nil {
		t.Fatalf("Record: %v", err)
	}
	events, _ := audit.Read(dir, "svc")
	if events[0].Timestamp.IsZero() {
		t.Error("expected Timestamp to be set, got zero value")
	}
}
