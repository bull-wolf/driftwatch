// Package schedule provides periodic drift-check scheduling for driftwatch.
// It allows running drift detection at a fixed interval and invoking a
// callback with the results of each check.
package schedule

import (
	"context"
	"fmt"
	"time"
)

// CheckFunc is the function signature for a drift-check callback.
// It receives the service name being checked and returns an error if the
// check itself fails (not if drift is found).
type CheckFunc func(ctx context.Context, serviceName string) error

// Config holds the configuration for a scheduled drift check.
type Config struct {
	// ServiceName is the name of the service to check.
	ServiceName string
	// Interval is how often the check should run.
	Interval time.Duration
	// Check is the function called on each tick.
	Check CheckFunc
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("schedule: service name must not be empty")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("schedule: interval must be positive, got %s", c.Interval)
	}
	if c.Check == nil {
		return fmt.Errorf("schedule: check function must not be nil")
	}
	return nil
}

// Run starts a blocking loop that calls cfg.Check on every cfg.Interval tick.
// It returns when ctx is cancelled, forwarding the context error.
// Any error returned by cfg.Check is passed to the optional onErr handler;
// if onErr is nil, check errors are silently ignored.
func Run(ctx context.Context, cfg Config, onErr func(error)) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := cfg.Check(ctx, cfg.ServiceName); err != nil && onErr != nil {
				onErr(err)
			}
		}
	}
}
