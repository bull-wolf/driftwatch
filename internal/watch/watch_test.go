package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/driftwatch/internal/watch"
)

func TestValidate_MissingDir(t *testing.T) {
	cfg := watch.Config{Interval: time.Second, OnChange: func(watch.Event) error { return nil }}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing dir")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := watch.Config{Dir: "/tmp", OnChange: func(watch.Event) error { return nil }}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestValidate_NilOnChange(t *testing.T) {
	cfg := watch.Config{Dir: "/tmp", Interval: time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for nil OnChange")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := watch.Config{
		Dir:      "/tmp",
		Interval: time.Second,
		OnChange: func(watch.Event) error { return nil },
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_DetectsChangedFile(t *testing.T) {
	dir := t.TempDir()
	manifest := filepath.Join(dir, "auth-service.yaml")

	// create initial file so watcher populates mtimes
	if err := os.WriteFile(manifest, []byte("initial"), 0644); err != nil {
		t.Fatal(err)
	}

	var callCount atomic.Int32
	var gotEvent watch.Event

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg := watch.Config{
		Dir:      dir,
		Interval: 50 * time.Millisecond,
		OnChange: func(e watch.Event) error {
			gotEvent = e
			callCount.Add(1)
			cancel()
			return nil
		},
	}

	// allow first poll to populate mtimes, then modify file
	go func() {
		time.Sleep(120 * time.Millisecond)
		_ = os.WriteFile(manifest, []byte("updated"), 0644)
	}()

	if err := watch.Run(ctx, cfg); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if callCount.Load() == 0 {
		t.Fatal("expected OnChange to be called at least once")
	}
	if gotEvent.Service != "auth-service" {
		t.Errorf("expected service %q, got %q", "auth-service", gotEvent.Service)
	}
}

func TestRun_CancelledImmediately(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := watch.Config{
		Dir:      dir,
		Interval: 10 * time.Millisecond,
		OnChange: func(watch.Event) error { return nil },
	}
	if err := watch.Run(ctx, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
