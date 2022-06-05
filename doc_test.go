package keyboard

import (
	"fmt"
	"os"

	"atomicgo.dev/keyboard/keys"
)

func ExampleListen_simple() {
	Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.CtrlC {
			return true, nil // Stop listener by returning true on Ctrl+C
		}

		fmt.Println("\r" + key.String()) // Print every key press
		return false, nil                // Return false to continue listening
	})
}

func ExampleListen_advanced() {
	// Stop keyboard listener on Escape key press or CTRL+C.
	// Exit application on "q" key press.
	// Print every rune key press.
	// Print every other key press.
	Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC, keys.Escape:
			return true, nil // Return true to stop listener
		case keys.RuneKey: // Check if key is a rune key (a, b, c, 1, 2, 3, ...)
			if key.String() == "q" { // Check if key is "q"
				fmt.Println("\rQuitting application")
				os.Exit(0) // Exit application
			}
			fmt.Printf("\rYou pressed the rune key: %s\n", key)
		default:
			fmt.Printf("\rYou pressed: %s\n", key)
		}

		return false, nil // Return false to continue listening
	})
}

func ExampleSimulateKeyPress() {
	go func() {
		SimulateKeyPress("Hello")             // Simulate key press for every letter in string
		SimulateKeyPress(keys.Enter)          // Simulate key press for Enter
		SimulateKeyPress(keys.CtrlShiftRight) // Simulate key press for Ctrl+Shift+Right
		SimulateKeyPress('x')                 // Simulate key press for a single rune
		SimulateKeyPress('x', keys.Down, 'a') // Simulate key presses for multiple inputs

		SimulateKeyPress(keys.Escape) // Simulate key press for Escape, which quits the program
	}()

	Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.Escape || key.Code == keys.CtrlC {
			os.Exit(0) // Exit program on Escape
		}

		fmt.Println("\r" + key.String()) // Print every key press
		return false, nil                // Return false to continue listening
	})
}
