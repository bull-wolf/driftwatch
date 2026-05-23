package schedule_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/driftwatch/internal/schedule"
)

func TestValidate_MissingServiceName(t *testing.T) {
	cfg := schedule.Config{
		Interval: time.Second,
		Check:    func(_ context.Context, _ string) error { return nil },
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty service name, got nil")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := schedule.Config{
		ServiceName: "auth-service",
		Check:       func(_ context.Context, _ string) error { return nil },
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero interval, got nil")
	}
}

func TestValidate_NilCheck(t *testing.T) {
	cfg := schedule.Config{
		ServiceName: "auth-service",
		Interval:    time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for nil check, got nil")
	}
}

func TestRun_CallsCheckOnTick(t *testing.T) {
	var callCount atomic.Int32

	cfg := schedule.Config{
		ServiceName: "auth-service",
		Interval:    20 * time.Millisecond,
		Check: func(_ context.Context, svc string) error {
			if svc != "auth-service" {
				t.Errorf("unexpected service name: %s", svc)
			}
			callCount.Add(1)
			return nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	err := schedule.Run(ctx, cfg, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if callCount.Load() < 2 {
		t.Errorf("expected at least 2 calls, got %d", callCount.Load())
	}
}

func TestRun_OnErrCalledOnCheckFailure(t *testing.T) {
	checkErr := errors.New("drift check failed")
	var errCount atomic.Int32

	cfg := schedule.Config{
		ServiceName: "auth-service",
		Interval:    20 * time.Millisecond,
		Check: func(_ context.Context, _ string) error {
			return checkErr
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	schedule.Run(ctx, cfg, func(err error) { //nolint:errcheck
		if !errors.Is(err, checkErr) {
			t.Errorf("unexpected error: %v", err)
		}
		errCount.Add(1)
	})

	if errCount.Load() < 1 {
		t.Errorf("expected onErr to be called at least once, got %d", errCount.Load())
	}
}

func TestRun_InvalidConfig_ReturnsError(t *testing.T) {
	cfg := schedule.Config{} // all zero
	err := schedule.Run(context.Background(), cfg, nil)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
