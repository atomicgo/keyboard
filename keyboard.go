package keyboard

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"atomicgo.dev/keyboard/keys"
)

// ErrAlreadyListening is returned by Listen when another Listen call is still running.
// Only one listener can own the terminal and the simulated key presses at a time.
var ErrAlreadyListening = errors.New("keyboard: a listener is already running")

// errNoTerminal is returned by openTerminal when stdin is not an interactive terminal.
var errNoTerminal = errors.New("keyboard: stdin is not a terminal")

// errCanceled is returned by terminal.Read after terminal.Cancel was called.
var errCanceled = errors.New("keyboard: read canceled")

// readerShutdownTimeout bounds how long Listen waits for the terminal reader to
// return after it was canceled. On Unix the reader always wakes up immediately;
// the timeout only guards platforms where canceling a blocked read is best effort.
const readerShutdownTimeout = time.Second

var (
	listening   atomic.Bool
	mockChannel = make(chan keys.Key)

	// openTerminal is a variable so tests can replace the system terminal.
	openTerminal = openSystemTerminal
)

// terminal is an interactive input device that has been switched to raw mode.
type terminal interface {
	// Read blocks until input is available. It returns errCanceled after Cancel was called.
	Read(p []byte) (int, error)
	// Cancel unblocks a pending Read. It is safe to call from another goroutine.
	Cancel()
	// Close restores the original terminal mode and releases all resources.
	Close() error
}

type readResult struct {
	keys []keys.Key
	err  error
}

// Listen calls a callback function when a key is pressed.
//
// The terminal is switched to raw mode while Listen runs, so Ctrl+C is delivered
// to the callback as keys.CtrlC instead of interrupting the program. The previous
// terminal mode is always restored before Listen returns, including when the
// callback returns an error or panics.
//
// Key presses sent with SimulateKeyPress are delivered as well. If stdin is not a
// terminal, only simulated key presses are delivered.
//
// Listen can be called any number of times, one after another. It returns
// ErrAlreadyListening if it is called while another Listen call is still running,
// for example from inside the callback.
//
// Simple example:
//
//	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
//		if key.Code == keys.CtrlC {
//			return true, nil // Stop listener by returning true on Ctrl+C
//		}
//
//		fmt.Println("\r" + key.String()) // Print every key press
//		return false, nil // Return false to continue listening
//	})
func Listen(onKeyPress func(key keys.Key) (stop bool, err error)) (err error) {
	if !listening.CompareAndSwap(false, true) {
		return ErrAlreadyListening
	}
	defer listening.Store(false)

	var terminalKeys <-chan readResult

	term, openErr := openTerminal()

	switch {
	case openErr == nil:
		results := make(chan readResult)
		done := make(chan struct{})
		readerDone := make(chan struct{})

		go func() {
			defer close(readerDone)

			readTerminal(term, results, done)
		}()

		defer func() {
			close(done)
			term.Cancel()

			select {
			case <-readerDone:
			case <-time.After(readerShutdownTimeout):
			}

			if closeErr := term.Close(); closeErr != nil && err == nil {
				err = fmt.Errorf("keyboard: failed to restore terminal: %w", closeErr)
			}
		}()

		terminalKeys = results

	case errors.Is(openErr, errNoTerminal):
		// Without a terminal only simulated key presses can arrive.
	default:
		return openErr
	}

	for {
		var pressed []keys.Key

		select {
		case key := <-mockChannel:
			pressed = []keys.Key{key}
		case result := <-terminalKeys:
			if result.err != nil {
				return result.err
			}

			pressed = result.keys
		}

		for _, key := range pressed {
			stop, callbackErr := onKeyPress(key)
			if callbackErr != nil {
				return callbackErr
			}

			if stop {
				return nil
			}
		}
	}
}

// readTerminal reads key presses from term and sends them to results until done is closed.
func readTerminal(term terminal, results chan<- readResult, done <-chan struct{}) {
	var buf [256]byte

	for {
		n, err := term.Read(buf[:])
		if errors.Is(err, errCanceled) {
			return
		}

		var result readResult

		if err != nil {
			result.err = fmt.Errorf("keyboard: failed to read from terminal: %w", err)
		} else {
			// Every read from a raw terminal holds complete key presses, so the
			// buffer is parsed as if no more input followed.
			result.keys, _ = parseKeys(buf[:n], true, true)
		}

		if len(result.keys) == 0 && result.err == nil {
			continue
		}

		select {
		case results <- result:
		case <-done:
			return
		}

		if result.err != nil {
			return
		}
	}
}

// ListenReader calls a callback function for every key press read from r.
//
// Unlike Listen, it does not touch the terminal and ignores simulated key presses.
// It is meant for tests and for input that is not an interactive terminal, such as
// a pipe. Every rune and every escape sequence in the input is delivered as a
// separate key, and both "\r" and "\n" are delivered as keys.Enter.
//
// ListenReader returns nil when the callback stops it, the callback's error if it
// returns one, and io.EOF if r is exhausted before the callback stops it.
//
// Example:
//
//	keyboard.ListenReader(strings.NewReader("y\r"), func(key keys.Key) (stop bool, err error) {
//		fmt.Println(key) // Prints "y", then "enter"
//		return key.Code == keys.Enter, nil
//	})
func ListenReader(r io.Reader, onKeyPress func(key keys.Key) (stop bool, err error)) error {
	var (
		pending []byte
		chunk   [256]byte
	)

	for {
		n, readErr := r.Read(chunk[:])
		pending = append(pending, chunk[:n]...)
		atEOF := errors.Is(readErr, io.EOF)

		parsed, consumed := parseKeys(pending, atEOF, false)
		pending = append(pending[:0], pending[consumed:]...)

		for _, key := range parsed {
			if key.Code == keys.CtrlJ {
				key.Code = keys.Enter
			}

			stop, callbackErr := onKeyPress(key)
			if callbackErr != nil {
				return callbackErr
			}

			if stop {
				return nil
			}
		}

		if atEOF {
			return io.EOF
		}

		if readErr != nil {
			return fmt.Errorf("keyboard: failed to read input: %w", readErr)
		}
	}
}

// SimulateKeyPress simulate a key press. It can be used to mock user stdin and test your application.
// It accepts runes, strings (one key press per rune), keys.KeyCode and keys.Key values,
// and blocks until a running Listen call has received every key press.
//
// Example:
//
//	go func() {
//		keyboard.SimulateKeyPress("Hello")             // Simulate key press for every letter in string
//		keyboard.SimulateKeyPress(keys.Enter)          // Simulate key press for Enter
//		keyboard.SimulateKeyPress(keys.CtrlShiftRight) // Simulate key press for Ctrl+Shift+Right
//		keyboard.SimulateKeyPress('x')                 // Simulate key press for a single rune
//		keyboard.SimulateKeyPress('x', keys.Down, 'a') // Simulate key presses for multiple inputs
//	}()
func SimulateKeyPress(input ...any) error {
	for _, in := range input {
		switch key := in.(type) {
		case keys.Key:
			mockChannel <- key
		case rune:
			mockChannel <- keys.Key{Code: keys.RuneKey, Runes: []rune{key}}
		case string:
			for _, r := range key {
				mockChannel <- keys.Key{Code: keys.RuneKey, Runes: []rune{r}}
			}

		case keys.KeyCode:
			mockChannel <- keys.Key{Code: key}
		default:
			return fmt.Errorf("keyboard: cannot simulate key press of type %T", in)
		}
	}

	return nil
}

// cancelOnce is embedded by terminal implementations to make Cancel idempotent.
type cancelOnce struct {
	once     sync.Once
	canceled atomic.Bool
}

func (c *cancelOnce) cancel(f func()) {
	c.once.Do(func() {
		c.canceled.Store(true)
		f()
	})
}
