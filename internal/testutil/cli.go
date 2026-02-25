// Package testutil provides shared helpers for tests.
package testutil

import "sync"

var (
	cliMu    sync.Mutex
	globalMu sync.Mutex
)

// LockCLI serializes CLI tests that mutate global state.
func LockCLI() {
	cliMu.Lock()
}

// UnlockCLI releases the CLI test lock.
func UnlockCLI() {
	cliMu.Unlock()
}

// LockGlobal serializes tests that mutate global state.
func LockGlobal() {
	globalMu.Lock()
}

// UnlockGlobal releases the global test lock.
func UnlockGlobal() {
	globalMu.Unlock()
}
