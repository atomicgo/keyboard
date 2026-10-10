//go:build !(aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || windows)

package keyboard

// openSystemTerminal reports that raw terminal input is not supported on this
// platform, so Listen only delivers simulated key presses.
func openSystemTerminal() (terminal, error) {
	return nil, errNoTerminal
}
