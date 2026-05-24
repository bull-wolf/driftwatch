package watch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/watch"
)

// TestServiceNameFromPath_KnownExtensions verifies that service names are
// derived correctly from real testdata manifest filenames.
func TestServiceNameFromPath_KnownExtensions(t *testing.T) {
	cases := []struct {
		filename string
		want     string
	}{
		{"auth-service.yaml", "auth-service"},
		{"payment-service.yml", "payment-service"},
		{"gateway.json", "gateway"},
		{"no-extension", "no-extension"},
	}
	for _, tc := range cases {
		dir := t.TempDir()
		path := filepath.Join(dir, tc.filename)
		if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		// We exercise the public surface indirectly via an OnChange callback.
		// Modify the file after watcher initialises.
		var got string
		done := make(chan struct{})
		cfg := watch.Config{
			Dir:      dir,
			Interval: 20 * time.Millisecond,
			OnChange: func(e watch.Event) error {
				got = e.Service
				close(done)
				return nil
			},
		}
		// run in background
		ctx, cancel := newCtxWithTimeout(t, 2*time.Second)
		defer cancel()
		go func() { _ = watch.Run(ctx, cfg) }()

		// wait for first poll to populate mtimes, then touch the file
		time.Sleep(60 * time.Millisecond)
		if err := os.WriteFile(path, []byte("changed"), 0644); err != nil {
			t.Fatal(err)
		}

		select {
		case <-done:
		case <-ctx.Done():
			t.Fatalf("timeout waiting for change event for %q", tc.filename)
		}
		cancel()
		if got != tc.want {
			t.Errorf("filename %q: expected service %q, got %q", tc.filename, tc.want, got)
		}
	}
}

func TestTestdataManifestDir_Exists(t *testing.T) {
	const dir = "../../testdata/manifests"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("testdata/manifests directory not found: %v", err)
	}
}
