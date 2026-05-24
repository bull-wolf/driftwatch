// Package watch provides file-system watching for manifest changes,
// triggering drift detection when source-of-truth files are updated.
package watch

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Event represents a detected file change.
type Event struct {
	Path    string
	Service string
	At      time.Time
}

// Config holds configuration for the watcher.
type Config struct {
	Dir      string
	Interval time.Duration
	OnChange func(Event) error
}

// Validate checks that the Config is well-formed.
func (c Config) Validate() error {
	if c.Dir == "" {
		return fmt.Errorf("watch: dir must not be empty")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("watch: interval must be positive")
	}
	if c.OnChange == nil {
		return fmt.Errorf("watch: OnChange callback must not be nil")
	}
	return nil
}

// Run starts polling the directory for manifest file changes.
// It blocks until ctx is cancelled.
func Run(ctx context.Context, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	mtimes := map[string]time.Time{}
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := poll(cfg, mtimes); err != nil {
				return err
			}
		}
	}
}

func poll(cfg Config, mtimes map[string]time.Time) error {
	entries, err := os.ReadDir(cfg.Dir)
	if err != nil {
		return fmt.Errorf("watch: read dir %q: %w", cfg.Dir, err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(cfg.Dir, e.Name())
		info, err := e.Info()
		if err != nil {
			continue
		}
		mod := info.ModTime()
		if prev, seen := mtimes[path]; !seen || mod.After(prev) {
			mtimes[path] = mod
			if !seen {
				continue // skip initial population
			}
			svc := serviceNameFromPath(path)
			if err := cfg.OnChange(Event{Path: path, Service: svc, At: mod}); err != nil {
				return err
			}
		}
	}
	return nil
}

func serviceNameFromPath(path string) string {
	base := filepath.Base(path)
	for _, ext := range []string{".yaml", ".yml", ".json"} {
		if len(base) > len(ext) && base[len(base)-len(ext):] == ext {
			return base[:len(base)-len(ext)]
		}
	}
	return base
}
