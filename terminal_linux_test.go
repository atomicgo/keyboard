package keyboard

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

// openPTY opens a pseudo terminal pair and returns the controller and the terminal side.
func openPTY(t *testing.T) (int, int) {
	t.Helper()

	controller, err := unix.Open("/dev/ptmx", unix.O_RDWR|unix.O_NOCTTY|unix.O_CLOEXEC, 0)
	if err != nil {
		t.Skipf("pseudo terminals are not available: %v", err)
	}

	t.Cleanup(func() { _ = unix.Close(controller) })

	err = unix.IoctlSetPointerInt(controller, unix.TIOCSPTLCK, 0)
	if err != nil {
		t.Fatalf("failed to unlock pty: %v", err)
	}

	number, err := unix.IoctlGetInt(controller, unix.TIOCGPTN)
	if err != nil {
		t.Fatalf("failed to get pty number: %v", err)
	}

	name := "/dev/pts/" + strconv.Itoa(number)

	tty, err := unix.Open(name, unix.O_RDWR|unix.O_NOCTTY|unix.O_CLOEXEC, 0)
	if err != nil {
		t.Fatalf("failed to open %s: %v", name, err)
	}

	t.Cleanup(func() { _ = unix.Close(tty) })

	return controller, tty
}

func TestUnixTerminal(t *testing.T) {
	controller, tty := openPTY(t)

	before, err := unix.IoctlGetTermios(tty, unix.TCGETS)
	if err != nil {
		t.Fatal(err)
	}

	ut, err := newUnixTerminal(tty, false)
	if err != nil {
		t.Fatal(err)
	}

	raw, err := unix.IoctlGetTermios(tty, unix.TCGETS)
	if err != nil {
		t.Fatal(err)
	}

	if raw.Lflag&(unix.ICANON|unix.ECHO|unix.ISIG) != 0 {
		t.Error("terminal is not in raw mode")
	}

	// Ctrl+C has to arrive as a byte instead of raising SIGINT.
	_, err = unix.Write(controller, []byte{0x03})
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 16)

	n, err := ut.Read(buf)
	if err != nil || n != 1 || buf[0] != 0x03 {
		t.Errorf("Read() = %q, %v, want ctrl+c", buf[:n], err)
	}

	// Cancel has to unblock a pending read.
	readErr := make(chan error, 1)

	go func() {
		_, cancelErr := ut.Read(buf)
		readErr <- cancelErr
	}()

	time.Sleep(50 * time.Millisecond)
	ut.Cancel()
	ut.Cancel()

	select {
	case cancelErr := <-readErr:
		if !errors.Is(cancelErr, errCanceled) {
			t.Errorf("canceled Read returned %v, want errCanceled", cancelErr)
		}

	case <-time.After(5 * time.Second):
		t.Fatal("Cancel did not unblock Read")
	}

	err = ut.Close()
	if err != nil {
		t.Fatalf("Close returned an error: %v", err)
	}

	after, err := unix.IoctlGetTermios(tty, unix.TCGETS)
	if err != nil {
		t.Fatal(err)
	}

	if after.Lflag != before.Lflag || after.Iflag != before.Iflag || after.Oflag != before.Oflag {
		t.Error("terminal mode was not restored")
	}

	if !term.IsTerminal(tty) {
		t.Error("terminal file descriptor was closed although it is not owned")
	}
}
