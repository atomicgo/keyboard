//go:build windows

package keyboard

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

// windowsTerminal reads raw input from the console input buffer (CONIN$).
// Cancel interrupts a pending Read with CancelIoEx.
type windowsTerminal struct {
	cancelOnce

	file   *os.File
	handle windows.Handle
	state  *term.State
}

func openSystemTerminal() (terminal, error) {
	if !term.IsTerminal(int(windows.Stdin)) {
		return nil, errNoTerminal
	}

	file, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("keyboard: failed to open console input: %w", err)
	}

	handle := windows.Handle(file.Fd())

	// Raw mode turns off processed input, so Ctrl+C arrives as a key press
	// instead of a signal, and turns on virtual terminal input sequences.
	state, err := term.MakeRaw(int(handle))
	if err != nil {
		_ = file.Close()

		return nil, fmt.Errorf("keyboard: failed to set raw mode: %w", err)
	}

	return &windowsTerminal{file: file, handle: handle, state: state}, nil
}

func (t *windowsTerminal) Read(p []byte) (int, error) {
	if t.canceled.Load() {
		return 0, errCanceled
	}

	n, err := t.file.Read(p)
	if err != nil && t.canceled.Load() {
		return 0, errCanceled
	}

	return n, err
}

func (t *windowsTerminal) Cancel() {
	t.cancel(func() {
		_ = windows.CancelIoEx(t.handle, nil)
	})
}

func (t *windowsTerminal) Close() error {
	return errors.Join(term.Restore(int(t.handle), t.state), t.file.Close())
}
