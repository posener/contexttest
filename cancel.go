package contexttest

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WithCancel can be tested to have the behavior of context.WithCancel.
type WithCancel func(ctx context.Context) (context.Context, context.CancelFunc)

// TestWithCancel tests the behavior of a context.WithCancel function.
func TestWithCancel(wc WithCancel) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Basic", wc.testBasic)
		t.Run("Tree", wc.testTree)
		t.Run("ParentAlreadyCanceled", wc.testParentAlreadyCanceled)
		t.Run("Concurrency", wc.testConcurrency)
	}
}

// Basic functionality test.
func (withCancel WithCancel) testBasic(t *testing.T) {
	t.Parallel()
	ctx, cancel := withCancel(context.Background())
	require.NotNil(t, ctx)
	require.NotNil(t, cancel)
	assertNotCanceled(t, ctx)

	cancel()
	assertCanceled(t, ctx)
}

// Context cancellation propagates only to contexts in the descendants tree,
// and to all of them.
func (withCancel WithCancel) testTree(t *testing.T) {
	t.Parallel()

	ctx0, cancel0 := withCancel(context.Background())
	ctx00, cancel00 := withCancel(ctx0)
	ctx01, cancel01 := withCancel(ctx0)
	ctx000, cancel001 := withCancel(ctx00)

	cancel00()
	assertNotCanceled(t, ctx0)
	assertCanceled(t, ctx00)
	assertNotCanceled(t, ctx01)
	assertCanceled(t, ctx000)

	cancel0()
	assertCanceled(t, ctx0)
	assertCanceled(t, ctx00)
	assertCanceled(t, ctx01)
	assertCanceled(t, ctx000)

	cancel01()
	cancel001()
	assertCanceled(t, ctx0)
	assertCanceled(t, ctx00)
	assertCanceled(t, ctx01)
	assertCanceled(t, ctx000)
}

// Allow a child context to be cancelled after its parent was canceled.
func (withCancel WithCancel) testParentAlreadyCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := withCancel(context.Background())
	cancel()
	assertCanceled(t, ctx)

	ctx, cancel = withCancel(ctx)
	assertCanceled(t, ctx)
	cancel()
}

// Function can be accessed from multiple goroutines without a race
// condition.
func (withCancel WithCancel) testConcurrency(t *testing.T) {
	t.Parallel()

	const n = 1000

	ctx, cancel := withCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(3)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			ctx.Err()
		}
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			ctx.Err()
		}
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			withCancel(ctx)
		}
		wg.Done()
	}(ctx)

	for i := 0; i < n; i++ {
		ctx, cancel = withCancel(ctx)
		cancel()
	}

	wg.Wait()
}

func assertCanceled(t *testing.T, ctx context.Context) {
	t.Helper()
	assert.Equal(t, context.Canceled, ctx.Err())
	select {
	case <-ctx.Done():
	default:
		t.Error("context not done")
	}
}

func assertNotCanceled(t *testing.T, ctx context.Context) {
	t.Helper()
	assert.Nil(t, ctx.Err())
	select {
	case <-ctx.Done():
		t.Error("context done")
	default:
	}
}
