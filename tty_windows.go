//go:build windows
// +build windows

package keyboard

import (
	"os"

	"github.com/containerd/console"
)

func restoreInput() error {
	if windowsStdin != nil {
		os.Stdin = windowsStdin
	}

	return nil
}

func initInput() error {
	windowsStdin = os.Stdin
	os.Stdin = input

	con = console.Current()

	return nil
}

func openInputTTY() (*os.File, error) {
	f, err := os.OpenFile("CONIN$", os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}
