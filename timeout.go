package contexttest

import (
	"context"
	"testing"
	"time"
)

// WithTimeout can be tested to have the behavior of context.WithTimeout.
type WithTimeout func(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc)

// TestWithTimeout tests the behavior of a context.WithTimeout function.
func TestWithTimeout(withTimeout WithTimeout) func(t *testing.T) {

	// Convert the WithTimeout function to a WithDeadline function and test it.
	withDeadline := func(ctx context.Context, timeout time.Time) (context.Context, context.CancelFunc) {
		return withTimeout(ctx, timeout.Sub(time.Now()))
	}

	return TestWithDeadline(withDeadline)
}
