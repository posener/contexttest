package contexttest

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WithValue can be tested to have the behavior of context.WithValue.
type WithValue func(ctx context.Context, key, val interface{}) context.Context

// TestWithValue tests the behavior of a context.WithValue function.
func TestWithValue(wv WithValue) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		t.Run("Basic", wv.testBasic)
		t.Run("Tree", wv.testTree)
		t.Run("Override", wv.testOverride)
		t.Run("Concurrency", wv.testConcurrency)
	}
}

// Basic functionality test.
func (withValue WithValue) testBasic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = withValue(ctx, 1, 2)
	require.NotNil(t, ctx)
	assert.Equal(t, 2, ctx.Value(1))
}

// Values appear only on their descendants tree.
func (withValue WithValue) testTree(t *testing.T) {
	t.Parallel()

	ctx0 := context.Background()

	ctx01 := withValue(ctx0, 1, 1)
	assert.Equal(t, 1, ctx01.Value(1))
	assert.Equal(t, nil, ctx0.Value(1))

	ctx02 := withValue(ctx0, 2, 2)
	assert.Equal(t, nil, ctx02.Value(1))
	assert.Equal(t, 2, ctx02.Value(2))
	assert.Equal(t, nil, ctx0.Value(2))
	assert.Equal(t, 1, ctx01.Value(1))
	assert.Equal(t, nil, ctx01.Value(2))

	ctx021 := withValue(ctx02, 3, 3)
	assert.Equal(t, nil, ctx021.Value(1))
	assert.Equal(t, 2, ctx021.Value(2))
	assert.Equal(t, 3, ctx021.Value(3))
}

// When two values have the same key, the parent context will hold the
// old value, and the child will hold the new value.
func (withValue WithValue) testOverride(t *testing.T) {
	t.Parallel()

	ctx0 := withValue(context.Background(), 0, 0)
	ctx1 := withValue(ctx0, 0, 1)

	assert.Equal(t, 0, ctx0.Value(0))
	assert.Equal(t, 1, ctx1.Value(0))
}

// Function can be accessed from multiple goroutines without a race
// condition.
func (withValue WithValue) testConcurrency(t *testing.T) {
	t.Parallel()

	const n = 1000

	ctx := withValue(context.Background(), 0, 0)

	var wg sync.WaitGroup
	wg.Add(3)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			ctx.Value(0)
		}
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			ctx.Value(0)
		}
		wg.Done()
	}(ctx)

	go func(ctx context.Context) {
		for i := 0; i < n; i++ {
			withValue(ctx, i, i)
		}
		wg.Done()
	}(ctx)

	for i := 0; i < n; i++ {
		ctx = withValue(ctx, i, i)
	}

	wg.Wait()
}
