/*
Package keyboard can be used to read key presses from the keyboard, while in a terminal application. It's crossplatform and keypresses can be combined to check for ctrl+c, alt+4, ctrl-shift, alt+ctrl+right, etc.

Works nicely with https://atomicgo.dev/cursor

```go

	keyboard.StartListener()
	defer keyboard.StopListener()

	for {
		keyInfo, _ := keyboard.GetKey()
		key := keyInfo.Code

		if key == keys.CtrlC {
			break
		}

		fmt.Println("\r", keyInfo.String())
	}

```
*/
package keyboard
