package contexttest_test

import (
	"context"
	"testing"

	"github.com/posener/contexttest"
)

func TestStandardLib(t *testing.T) {
	t.Run("WithValue", contexttest.TestWithValue(context.WithValue))
	t.Run("WithCancel", contexttest.TestWithCancel(context.WithCancel))
	t.Run("WithDeadline", contexttest.TestWithDeadline(context.WithDeadline))
	t.Run("WithTimeout", contexttest.TestWithTimeout(context.WithTimeout))
}
