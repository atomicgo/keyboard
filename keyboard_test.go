package keyboard

import (
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"atomicgo.dev/keyboard/keys"
)

func TestMain(m *testing.M) {
	// Never touch the terminal the tests run in. Tests that need terminal input
	// install a fakeTerminal instead.
	openTerminal = noTerminal

	os.Exit(m.Run())
}

func noTerminal() (terminal, error) {
	return nil, errNoTerminal
}

// fakeTerminal is an in-memory terminal that records whether it was restored.
type fakeTerminal struct {
	cancelOnce

	input   chan []byte
	readErr error
	done    chan struct{}

	mu     sync.Mutex
	closed int
}

func newFakeTerminal(input ...string) *fakeTerminal {
	ft := &fakeTerminal{
		input: make(chan []byte, len(input)),
		done:  make(chan struct{}),
	}

	for _, in := range input {
		ft.input <- []byte(in)
	}

	return ft
}

func (ft *fakeTerminal) Read(p []byte) (int, error) {
	select {
	case in := <-ft.input:
		return copy(p, in), nil
	case <-ft.done:
		return 0, errCanceled
	default:
	}

	if ft.readErr != nil {
		return 0, ft.readErr
	}

	select {
	case in := <-ft.input:
		return copy(p, in), nil
	case <-ft.done:
		return 0, errCanceled
	}
}

func (ft *fakeTerminal) Cancel() {
	ft.cancel(func() { close(ft.done) })
}

func (ft *fakeTerminal) Close() error {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	ft.closed++

	return nil
}

func (ft *fakeTerminal) closeCount() int {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	return ft.closed
}

func useTerminal(t *testing.T, ft *fakeTerminal) {
	t.Helper()

	openTerminal = func() (terminal, error) { return ft, nil }

	t.Cleanup(func() { openTerminal = noTerminal })
}

// collect returns a callback that records every key and stops at stopCode.
func collect(got *[]keys.Key, stopCode keys.KeyCode) func(keys.Key) (bool, error) {
	return func(key keys.Key) (bool, error) {
		*got = append(*got, key)

		return key.Code == stopCode, nil
	}
}

func simulate(t *testing.T, input ...any) {
	t.Helper()

	errs := make(chan error, 1)

	go func() {
		errs <- SimulateKeyPress(input...)
	}()

	t.Cleanup(func() {
		select {
		case err := <-errs:
			if err != nil {
				t.Errorf("SimulateKeyPress: %v", err)
			}

		case <-time.After(5 * time.Second):
			t.Error("simulated key presses were not consumed")
		}
	})
}

func rk(s string) keys.Key {
	return keys.Key{Code: keys.RuneKey, Runes: []rune(s)}
}

func TestSimulateKeyPress(t *testing.T) {
	simulate(t, 'a', "bc", keys.Key{Code: keys.Up, AltPressed: true}, keys.Enter)

	var got []keys.Key

	if err := Listen(collect(&got, keys.Enter)); err != nil {
		t.Fatalf("Listen returned an error: %v", err)
	}

	want := []keys.Key{rk("a"), rk("b"), rk("c"), {Code: keys.Up, AltPressed: true}, {Code: keys.Enter}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got keys %v, want %v", got, want)
	}
}

func TestSimulateKeyPressUnsupportedType(t *testing.T) {
	if err := SimulateKeyPress(3.14); err == nil {
		t.Error("expected an error for an unsupported input type")
	}
}

func TestListenCanBeCalledRepeatedly(t *testing.T) {
	ft := newFakeTerminal()
	useTerminal(t, ft)

	for i := range 5 {
		simulate(t, "x", keys.Enter)

		var got []keys.Key

		if err := Listen(collect(&got, keys.Enter)); err != nil {
			t.Fatalf("Listen #%d returned an error: %v", i, err)
		}

		if len(got) != 2 {
			t.Fatalf("Listen #%d got %v, want two keys", i, got)
		}
	}
}

func TestListenRestoresTerminalOnCallbackError(t *testing.T) {
	ft := newFakeTerminal("q")
	useTerminal(t, ft)

	wantErr := errors.New("callback failed")

	err := Listen(func(keys.Key) (bool, error) {
		return false, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Errorf("got error %v, want %v", err, wantErr)
	}

	if ft.closeCount() != 1 {
		t.Errorf("terminal was restored %d times, want 1", ft.closeCount())
	}
}

func TestListenRestoresTerminalOnPanic(t *testing.T) {
	ft := newFakeTerminal("q")
	useTerminal(t, ft)

	func() {
		defer func() {
			if recover() == nil {
				t.Error("expected the callback panic to propagate")
			}
		}()

		_ = Listen(func(keys.Key) (bool, error) {
			panic("boom")
		})
	}()

	if ft.closeCount() != 1 {
		t.Errorf("terminal was restored %d times, want 1", ft.closeCount())
	}

	// The panic must not leave the listener marked as running.
	simulate(t, keys.Enter)

	if err := Listen(collect(new([]keys.Key), keys.Enter)); err != nil {
		t.Errorf("Listen after a panic returned an error: %v", err)
	}
}

func TestListenDeliversTerminalInput(t *testing.T) {
	ft := newFakeTerminal("\x1b[A", "a", "\x1b[1;5C", "\x03")
	useTerminal(t, ft)

	var got []keys.Key

	if err := Listen(collect(&got, keys.CtrlC)); err != nil {
		t.Fatalf("Listen returned an error: %v", err)
	}

	want := []keys.Key{{Code: keys.Up}, rk("a"), {Code: keys.CtrlRight}, {Code: keys.CtrlC}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got keys %v, want %v", got, want)
	}

	if ft.closeCount() != 1 {
		t.Errorf("terminal was restored %d times, want 1", ft.closeCount())
	}
}

func TestListenReturnsTerminalReadError(t *testing.T) {
	ft := newFakeTerminal()
	ft.readErr = io.ErrUnexpectedEOF
	useTerminal(t, ft)

	err := Listen(collect(new([]keys.Key), keys.Enter))
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("got error %v, want %v", err, io.ErrUnexpectedEOF)
	}

	if ft.closeCount() != 1 {
		t.Errorf("terminal was restored %d times, want 1", ft.closeCount())
	}
}

func TestListenWhileListening(t *testing.T) {
	simulate(t, keys.Enter)

	var nestedErr error

	err := Listen(func(keys.Key) (bool, error) {
		nestedErr = Listen(collect(new([]keys.Key), keys.Enter))

		return true, nil
	})
	if err != nil {
		t.Fatalf("Listen returned an error: %v", err)
	}

	if !errors.Is(nestedErr, ErrAlreadyListening) {
		t.Errorf("nested Listen returned %v, want %v", nestedErr, ErrAlreadyListening)
	}
}

func TestListenReader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []keys.Key
	}{
		{"runes", "ab\r", []keys.Key{rk("a"), rk("b"), {Code: keys.Enter}}},
		{"newline is enter", "a\n", []keys.Key{rk("a"), {Code: keys.Enter}}},
		{"space", " \r", []keys.Key{{Code: keys.Space, Runes: []rune{' '}}, {Code: keys.Enter}}},
		{"unicode", "äß😀\r", []keys.Key{rk("ä"), rk("ß"), rk("😀"), {Code: keys.Enter}}},
		{"arrows", "\x1b[A\x1b[B\x1b[1;6D\r", []keys.Key{{Code: keys.Up}, {Code: keys.Down}, {Code: keys.CtrlShiftLeft}, {Code: keys.Enter}}},
		{"alt", "\x1bx\x1b\x1b[C\r", []keys.Key{{Code: keys.RuneKey, Runes: []rune{'x'}, AltPressed: true}, {Code: keys.Right, AltPressed: true}, {Code: keys.Enter}}},
		{"backspace and tab", "\x7f\x08\t\r", []keys.Key{{Code: keys.Backspace}, {Code: keys.Backspace}, {Code: keys.Tab}, {Code: keys.Enter}}},
		{"function keys", "\x1bOP\x1b[24~\r", []keys.Key{{Code: keys.F1}, {Code: keys.F12}, {Code: keys.Enter}}},
		{"escape", "\x1b", []keys.Key{{Code: keys.Escape}}},
		{"ctrl+c", "\x03", []keys.Key{{Code: keys.CtrlC}}},
		{"invalid utf-8 is skipped", "a\xffb\r", []keys.Key{rk("a"), rk("b"), {Code: keys.Enter}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readers := map[string]io.Reader{
				"whole":    strings.NewReader(tt.input),
				"one byte": iotest.OneByteReader(strings.NewReader(tt.input)),
			}

			for readerName, r := range readers {
				var got []keys.Key

				err := ListenReader(r, func(key keys.Key) (bool, error) {
					got = append(got, key)

					return false, nil
				})
				if !errors.Is(err, io.EOF) {
					t.Errorf("%s: got error %v, want io.EOF", readerName, err)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s: got keys %v, want %v", readerName, got, tt.want)
				}
			}
		})
	}
}

func TestListenReaderStops(t *testing.T) {
	var got []keys.Key

	err := ListenReader(strings.NewReader("y\rignored"), collect(&got, keys.Enter))
	if err != nil {
		t.Fatalf("ListenReader returned an error: %v", err)
	}

	if want := []keys.Key{rk("y"), {Code: keys.Enter}}; !reflect.DeepEqual(got, want) {
		t.Errorf("got keys %v, want %v", got, want)
	}
}

func TestListenReaderErrors(t *testing.T) {
	wantErr := errors.New("callback failed")

	err := ListenReader(strings.NewReader("a"), func(keys.Key) (bool, error) { return false, wantErr })
	if !errors.Is(err, wantErr) {
		t.Errorf("got error %v, want %v", err, wantErr)
	}

	readErr := errors.New("read failed")

	err = ListenReader(iotest.ErrReader(readErr), collect(new([]keys.Key), keys.Enter))
	if !errors.Is(err, readErr) {
		t.Errorf("got error %v, want %v", err, readErr)
	}
}

func TestParseKeysGroupsPastedRunes(t *testing.T) {
	got, n := parseKeys([]byte("hello world\x1b[A "), true, true)

	want := []keys.Key{rk("hello world"), {Code: keys.Up}, {Code: keys.Space, Runes: []rune{' '}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got keys %v, want %v", got, want)
	}

	if n != 15 {
		t.Errorf("consumed %d bytes, want 15", n)
	}
}

func TestParseKeysWaitsForIncompleteInput(t *testing.T) {
	for _, input := range []string{"\x1b", "\x1b[", "\x1b[1;5", "\xc3"} {
		got, n := parseKeys([]byte("a"+input), false, false)
		if n != 1 || len(got) != 1 {
			t.Errorf("%q: parsed %v and consumed %d bytes, want only the leading rune", input, got, n)
		}
	}
}

func TestKeyString(t *testing.T) {
	tests := map[string]keys.Key{
		"a":         rk("a"),
		"alt+a":     {Code: keys.RuneKey, Runes: []rune{'a'}, AltPressed: true},
		"enter":     {Code: keys.Enter},
		"ctrl+c":    {Code: keys.CtrlC},
		"alt+up":    {Code: keys.Up, AltPressed: true},
		"backspace": {Code: keys.Backspace},
	}

	for want, key := range tests {
		if got := key.String(); got != want {
			t.Errorf("%#v.String() = %q, want %q", key, got, want)
		}
	}
}
