//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package keyboard

import (
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

// unixTerminal reads raw input from a terminal file descriptor.
//
// Reads wait in select(2) on the terminal and on a wake-up pipe, so Cancel can
// interrupt a pending Read without closing the terminal. select is used instead
// of poll because poll does not support terminal devices on macOS.
type unixTerminal struct {
	cancelOnce

	fd       int
	ownsFD   bool
	state    *term.State
	wakeRead *os.File
	wakeFD   int
	wake     *os.File
}

func openSystemTerminal() (terminal, error) {
	if !term.IsTerminal(unix.Stdin) {
		return nil, errNoTerminal
	}

	// Prefer the controlling terminal, so stdin stays untouched for the program.
	fd, err := unix.Open("/dev/tty", unix.O_RDWR|unix.O_NOCTTY|unix.O_CLOEXEC, 0)
	if err != nil {
		return newUnixTerminal(unix.Stdin, false)
	}

	t, err := newUnixTerminal(fd, true)
	if err != nil {
		_ = unix.Close(fd)

		return nil, err
	}

	return t, nil
}

func newUnixTerminal(fd int, ownsFD bool) (*unixTerminal, error) {
	wakeRead, wake, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("keyboard: failed to create wake-up pipe: %w", err)
	}

	state, err := term.MakeRaw(fd)
	if err != nil {
		_ = wakeRead.Close()
		_ = wake.Close()

		return nil, fmt.Errorf("keyboard: failed to set raw mode: %w", err)
	}

	return &unixTerminal{
		fd:       fd,
		ownsFD:   ownsFD,
		state:    state,
		wakeRead: wakeRead,
		wakeFD:   int(wakeRead.Fd()),
		wake:     wake,
	}, nil
}

func (t *unixTerminal) Read(p []byte) (int, error) {
	for {
		if t.canceled.Load() {
			return 0, errCanceled
		}

		var readSet unix.FdSet

		readSet.Zero()
		readSet.Set(t.fd)
		readSet.Set(t.wakeFD)

		_, err := unix.Select(max(t.fd, t.wakeFD)+1, &readSet, nil, nil, nil)
		if errors.Is(err, unix.EINTR) {
			continue
		}

		if err != nil {
			return 0, err
		}

		if readSet.IsSet(t.wakeFD) {
			return 0, errCanceled
		}

		if !readSet.IsSet(t.fd) {
			continue
		}

		n, err := unix.Read(t.fd, p)
		if errors.Is(err, unix.EINTR) || errors.Is(err, unix.EAGAIN) {
			continue
		}

		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, io.EOF
		}

		return n, nil
	}
}

func (t *unixTerminal) Cancel() {
	t.cancel(func() {
		_, _ = t.wake.Write([]byte{0})
	})
}

func (t *unixTerminal) Close() error {
	err := term.Restore(t.fd, t.state)

	if t.ownsFD {
		err = errors.Join(err, unix.Close(t.fd))
	}

	return errors.Join(err, t.wake.Close(), t.wakeRead.Close())
}
