package keyboard_test

import (
	"sync"
	"testing"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

var wg sync.WaitGroup

func TestMocking(t *testing.T) {
	wg.Add(1)

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code != keys.Down {
			t.Errorf("down key should be pressed, but %s is pressed", key)
		}

		wg.Done()

		return true, nil
	})

	go func() {
		keyboard.SimulateKeyPress(keys.Down)
	}()

	wg.Wait()
}
