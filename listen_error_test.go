package keyboard

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"atomicgo.dev/keyboard/keys"
)

// listenTimeout bounds every Listen call in these tests: a regression that
// swallows read errors makes Listen spin forever instead of returning.
const listenTimeout = 5 * time.Second

// useFakeInput makes Listen read from f instead of the real terminal and skips
// raw-mode setup, so the tests never touch the developer's console.
//
// stdin gets its own non-console file: initInput calls Fd() on it, which
// switches the file to blocking mode, and a blocking Read on f could not be
// interrupted by the Close in stop().
func useFakeInput(t *testing.T, f *os.File) {
	t.Helper()

	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}

	prevStdin, prevInputTTY, prevMocking := stdin, inputTTY, mocking
	stdin, inputTTY, mocking = devNull, f, true

	t.Cleanup(func() {
		stdin, inputTTY, mocking = prevStdin, prevInputTTY, prevMocking

		_ = f.Close()
		_ = devNull.Close()
	})
}

// listenAsync runs Listen with a timeout and returns its error.
func listenAsync(t *testing.T, onKeyPress func(key keys.Key) (bool, error)) error {
	t.Helper()

	done := make(chan error, 1)

	go func() {
		done <- Listen(onKeyPress)
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(listenTimeout):
		t.Fatalf("Listen did not return within %s: read errors are being swallowed (busy loop)", listenTimeout)

		return nil
	}
}

// exhaustedPipe returns a pipe whose write end is already closed:
// every Read returns io.EOF immediately, the way a hung-up tty keeps failing.
func exhaustedPipe(t *testing.T) *os.File {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_ = w.Close()

	return r
}

func TestGetKeyPressReturnsReadError(t *testing.T) {
	useFakeInput(t, exhaustedPipe(t))

	_, err := getKeyPress()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("getKeyPress() error = %v, want io.EOF", err)
	}
}

// TestListenStopWhileReadBlocked guards the normal shutdown path: a Read
// interrupted by our own stop() must not be reported as an error.
// It runs before the busy-loop regressions below on purpose: a Listen leaked
// by a failed test would compete for the simulated key.
func TestListenStopWhileReadBlocked(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = w.Close() })

	// Write end stays open: Read blocks like an idle terminal until stop() closes r.
	useFakeInput(t, r)

	go func() {
		_ = SimulateKeyPress(keys.Enter)
	}()

	err = listenAsync(t, func(key keys.Key) (bool, error) {
		return key.Code == keys.Enter, nil
	})
	if err != nil {
		t.Fatalf("Listen() error = %v, want nil: a read interrupted by our own stop is a normal shutdown", err)
	}
}

func TestListenReturnsErrorWhenInputExhausted(t *testing.T) {
	useFakeInput(t, exhaustedPipe(t))

	calls := 0

	err := listenAsync(t, func(keys.Key) (bool, error) {
		calls++

		return false, nil
	})

	if !errors.Is(err, io.EOF) {
		t.Fatalf("Listen() error = %v, want io.EOF", err)
	}

	if calls != 0 {
		t.Fatalf("callback called %d times without input: a failed read must not be reported as a key press", calls)
	}
}

func TestListenReturnsErrorOnClosedInput(t *testing.T) {
	r := exhaustedPipe(t)
	_ = r.Close() // stale descriptor left behind by a previous Listen

	useFakeInput(t, r)

	err := listenAsync(t, func(keys.Key) (bool, error) {
		return false, nil
	})

	if !errors.Is(err, os.ErrClosed) {
		t.Fatalf("Listen() error = %v, want os.ErrClosed", err)
	}
}
