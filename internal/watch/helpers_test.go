package watch_test

import (
	"context"
	"testing"
	"time"
)

// newCtxWithTimeout returns a context and cancel func, automatically
// failing the test if the deadline is exceeded.
func newCtxWithTimeout(t *testing.T, d time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), d)
	t.Cleanup(func() {
		cancel()
	})
	return ctx, cancel
}
