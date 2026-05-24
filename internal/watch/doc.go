// Package watch polls a directory of manifest files at a configurable interval
// and invokes a callback whenever a file's modification time changes.
//
// Typical usage:
//
//	cfg := watch.Config{
//		Dir:      "./manifests",
//		Interval: 30 * time.Second,
//		OnChange: func(e watch.Event) error {
//			// re-run drift detection for e.Service
//			return nil
//		},
//	}
//	if err := watch.Run(ctx, cfg); err != nil {
//		log.Fatal(err)
//	}
//
// The watcher skips files present on the first poll (initial population)
// and only fires for subsequent modifications, avoiding false positives
// on startup.
package watch
