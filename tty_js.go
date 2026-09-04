//go:build js

package keyboard

import (
	"errors"
	"os"
)

// initInput is a no-op on js/wasm: there is no terminal to configure.
func initInput() error { return nil }

// restoreInput is a no-op on js/wasm.
func restoreInput() error { return nil }

// closeInput is a no-op on js/wasm.
func closeInput() {}

// openInputTTY reports that no controlling terminal exists on js/wasm;
// callers supply input through explicit readers instead.
func openInputTTY() (*os.File, error) {
	return nil, errors.New("keyboard: no /dev/tty on js/wasm")
}
