package testutil_test

import (
	"testing"

	"diagram-gen/internal/testutil"
)

func TestLockUnlockCLI(t *testing.T) {
	t.Parallel()

	testutil.LockCLI()
	testutil.UnlockCLI()

	testutil.LockGlobal()
	testutil.UnlockGlobal()
}
