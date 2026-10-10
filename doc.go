/*
Package keyboard can be used to read key presses from the keyboard, while in a terminal application. It's crossplatform and keypresses can be combined to check for ctrl+c, alt+4, ctrl-shift, alt+ctrl+right, etc.
It can also be used to simulate (mock) keypresses for CI testing, or to read keypresses from any io.Reader.

The terminal is switched to raw mode while Listen runs and is always restored afterwards,
even if the callback returns an error or panics. Ctrl+C is delivered to the callback as keys.CtrlC.

Works nicely with https://atomicgo.dev/cursor

## Simple Usage

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.CtrlC {
			return true, nil // Stop listener by returning true on Ctrl+C
		}

		fmt.Println("\r" + key.String()) // Print every key press
		return false, nil // Return false to continue listening
	})

## Advanced Usage

	// Stop keyboard listener on Escape key press or CTRL+C.
	// Exit application on "q" key press.
	// Print every rune key press.
	// Print every other key press.
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
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

## Simulate Key Presses (for mocking in tests)

	go func() {
		keyboard.SimulateKeyPress("Hello")             // Simulate key press for every letter in string
		keyboard.SimulateKeyPress(keys.Enter)          // Simulate key press for Enter
		keyboard.SimulateKeyPress(keys.CtrlShiftRight) // Simulate key press for Ctrl+Shift+Right
		keyboard.SimulateKeyPress('x')                 // Simulate key press for a single rune
		keyboard.SimulateKeyPress('x', keys.Down, 'a') // Simulate key presses for multiple inputs

		keyboard.SimulateKeyPress(keys.Escape) // Simulate key press for Escape, which quits the program
	}()

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.Escape || key.Code == keys.CtrlC {
			os.Exit(0) // Exit program on Escape
		}

		fmt.Println("\r" + key.String()) // Print every key press
		return false, nil                // Return false to continue listening
	})

## Read Key Presses From an io.Reader

	// Every rune and escape sequence is delivered as its own key. "\r" and "\n" are delivered as keys.Enter.
	err := keyboard.ListenReader(strings.NewReader("y\r"), func(key keys.Key) (stop bool, err error) {
		fmt.Println(key) // Prints "y", then "enter"
		return key.Code == keys.Enter, nil
	})
*/
package keyboard
