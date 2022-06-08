package keyboard_test

import (
	"testing"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/MarvinJWendt/testza"
)

func TestMocking(t *testing.T) {
	go func() {
		keyboard.SimulateKeyPress('a')
		keyboard.SimulateKeyPress("b")
		keyboard.SimulateKeyPress("c")
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	var aPressed, bPressed, cPressed, enterPressed bool
	var keyList []keys.Key

	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		keyList = append(keyList, key)
		switch key.Code {
		case keys.RuneKey:
			switch key.String() {
			case "a":
				println("a pressed")
				aPressed = true
				return false, nil
			case "b":
				println("b pressed")
				bPressed = true
				return false, nil
			case "c":
				println("c pressed")
				cPressed = true
				return false, nil
			}
		case keys.Enter:
			println("enter pressed")
			enterPressed = true
			return true, nil
		}

		return false, nil
	})

	testza.AssertNoError(t, err)

	testza.AssertTrue(t, aPressed, "A | %s", keyList)
	testza.AssertTrue(t, bPressed, "B | %s", keyList)
	testza.AssertTrue(t, cPressed, "C | %s", keyList)
	testza.AssertTrue(t, enterPressed, "Enter | %s", keyList)
}
