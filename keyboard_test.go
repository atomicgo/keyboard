package keyboard_test

import (
	"testing"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/MarvinJWendt/testza"
)

func TestMocking(t *testing.T) {
	simulateErrs := make(chan error, 4)

	go func() {
		defer close(simulateErrs)

		if err := keyboard.SimulateKeyPress('a'); err != nil {
			simulateErrs <- err

			return
		}

		if err := keyboard.SimulateKeyPress("b"); err != nil {
			simulateErrs <- err

			return
		}

		if err := keyboard.SimulateKeyPress("c"); err != nil {
			simulateErrs <- err

			return
		}

		if err := keyboard.SimulateKeyPress(keys.Enter); err != nil {
			simulateErrs <- err
		}
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

	for err := range simulateErrs {
		t.Errorf("failed to simulate key press: %v", err)
	}

	testza.AssertTrue(t, aPressed, "A | %s", keyList)
	testza.AssertTrue(t, bPressed, "B | %s", keyList)
	testza.AssertTrue(t, cPressed, "C | %s", keyList)
	testza.AssertTrue(t, enterPressed, "Enter | %s", keyList)
}
