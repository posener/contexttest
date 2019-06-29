package contexttest

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	grace        = 20 * time.Millisecond
	shortTimeout = 100 * time.Millisecond
	longTimeout  = 500 * time.Millisecond
)

// WithDeadline can be tested to have the behavior of context.WithDeadline.
type WithDeadline func(ctx context.Context, deadline time.Time) (context.Context, context.CancelFunc)

// TestWithDeadline tests the behavior of a context.WithDeadline function.
func TestWithDeadline(withDeadline WithDeadline) func(t *testing.T) {

	// For testing WithDeadline cancel functionality.
	withCancel := func(ctx context.Context) (context.Context, context.CancelFunc) {
		return withDeadline(ctx, time.Now().Add(longTimeout))
	}

	return func(t *testing.T) {
		t.Parallel()
		t.Run("Cancel", TestWithCancel(withCancel))
		t.Run("Basic", withDeadline.testBasic)
		t.Run("Tree", withDeadline.testTree)
	}
}

// Basic deadline functionality test.
func (withDeadline WithDeadline) testBasic(t *testing.T) {
	t.Parallel()
	deadline := time.Now().Add(shortTimeout)
	ctx, cancel := withDeadline(context.Background(), deadline)
	require.NotNil(t, ctx)
	require.NotNil(t, cancel)
	assertNotCanceled(t, ctx)

	defer cancel()
	assertDeadlined(t, ctx, deadline)
}

// Context deadline propagates only to contexts in the descendants tree,
// and to all of them.
func (withDeadline WithDeadline) testTree(t *testing.T) {
	t.Parallel()

	t0 := time.Now().Add(shortTimeout)
	t1 := t0.Add(shortTimeout)
	t2 := t1.Add(shortTimeout)

	ctx0, cancel0 := withDeadline(context.Background(), t1)
	ctx00, cancel00 := withDeadline(ctx0, t0)
	ctx01, cancel01 := withDeadline(ctx0, t2)
	ctx000, cancel001 := withDeadline(ctx00, t2)

	defer cancel0()
	defer cancel00()
	defer cancel01()
	defer cancel001()

	var wg sync.WaitGroup
	wg.Add(4)

	// Performs assertDeadline and release a waitgroup.
	assertDeadlinedWg := func(ctx context.Context, deadline time.Time) {
		t.Helper()
		assertDeadlined(t, ctx, deadline)
		wg.Done()
	}

	// ctx0 should be cancelled with its own timeout, t1.
	go assertDeadlinedWg(ctx0, t1)
	// ctx00 should be cancelled with its own timeout, t0.
	go assertDeadlinedWg(ctx00, t0)
	// ctx01 should be deadlined with t1 even though it was set with t2
	// because it is the child of ctx0.
	go assertDeadlinedWg(ctx01, t1)
	// ctx000 should be deadlined with t0 even though it was set with t2
	// because it is the child of ctx00.
	go assertDeadlinedWg(ctx000, t0)

	wg.Wait()
}

func assertDeadlined(t *testing.T, ctx context.Context, deadline time.Time) {
	t.Helper()
	assert.Nil(t, ctx.Err())
	select {
	case <-ctx.Done():
	case <-time.After(time.Until(deadline) + grace):
		t.Error("context not deadlined")
	}
	assert.Error(t, context.DeadlineExceeded, ctx.Err())
}
