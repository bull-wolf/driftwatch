package audit_test

import (
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/audit"
)

const testdataDir = "../../testdata/audit"

func TestRead_TestdataFile(t *testing.T) {
	events, err := audit.Read(testdataDir, "auth-service")
	if err != nil {
		t.Fatalf("Read: unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events from testdata, got %d", len(events))
	}
}

func TestRead_TestdataFirstEvent_IsDriftDetected(t *testing.T) {
	events, err := audit.Read(testdataDir, "auth-service")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if events[0].Kind != audit.EventKindDriftDetected {
		t.Errorf("Kind: got %q, want %q", events[0].Kind, audit.EventKindDriftDetected)
	}
	if events[0].ServiceName != "auth-service" {
		t.Errorf("ServiceName: got %q, want %q", events[0].ServiceName, "auth-service")
	}
}

func TestRead_TestdataSecondEvent_IsClean(t *testing.T) {
	events, err := audit.Read(testdataDir, "auth-service")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if events[1].Kind != audit.EventKindClean {
		t.Errorf("Kind: got %q, want %q", events[1].Kind, audit.EventKindClean)
	}
}

func TestRead_TestdataFile_PathResolution(t *testing.T) {
	abs, err := filepath.Abs(testdataDir)
	if err != nil {
		t.Fatalf("Abs: %v", err)
	}
	events, err := audit.Read(abs, "auth-service")
	if err != nil {
		t.Fatalf("Read with absolute path: %v", err)
	}
	if len(events) == 0 {
		t.Error("expected events from absolute path, got none")
	}
}
