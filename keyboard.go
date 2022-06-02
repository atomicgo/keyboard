package keyboard

import (
	"fmt"
	"os"

	"github.com/containerd/console"

	"atomicgo.dev/keyboard/keys"
)

var windowsStdin *os.File
var con console.Console
var input = os.Stdin
var inputTTY *os.File
var mockChannel = make(chan keys.Key)

// StartListener starts the keyboard listener
//
// Example:
//  keyboard.StartListener()
//
//  for {
//    keyInfo, _ := keyboard.GetKey()
//    key := keyInfo.Code
//
//    if key == keys.CtrlC {
//      break
//    }
//
//    fmt.Println("\r", keyInfo.String())
//  }
//
//  keyboard.StopListener()
func startListener() error {
	err := initInput()
	if err != nil {
		return err
	}

	if con != nil {
		err := con.SetRaw()
		if err != nil {
			return fmt.Errorf("failed to set raw mode: %w", err)
		}
	}

	inputTTY, err = openInputTTY()
	if err != nil {
		return err
	}

	return nil
}

// StopListener stops the keyboard listener
//
// Example:
//  keyboard.StartListener()
//
//  for {
//    keyInfo, _ := keyboard.GetKey()
//    key := keyInfo.Code
//
//    if key == keys.CtrlC {
//      break
//    }
//
//    fmt.Println("\r", keyInfo.String())
//  }
//
//  keyboard.StopListener()
func stopListener() error {
	if con != nil {
		err := con.Reset()
		if err != nil {

			return fmt.Errorf("failed to reset console: %w", err)
		}
	}

	return restoreInput()
}

// GetKey blocks until a key is pressed and returns the key info.
//
// Example:
//  keyboard.StartListener()
//
//  for {
//    keyInfo, _ := keyboard.GetKey()
//    key := keyInfo.Code
//
//    if key == keys.CtrlC {
//      break
//    }
//
//    fmt.Println("\r", keyInfo.String())
//  }
//
//  keyboard.StopListener()
// func GetKey() (keys.Key, error) {
// 	return getKeyPress(inputTTY)
// }

func GetListener() (listener chan keys.Key, cancel chan bool) {
	listener = make(chan keys.Key)
	cancel = make(chan bool)

	go func() {
		Listen(func(keyInfo keys.Key) (stop bool, err error) {
			select {
			case <-cancel:
				return true, nil
			default:
				listener <- keyInfo
				return false, nil
			}
		})
	}()

	return
}

func Listen(onKeyPress func(keyInfo keys.Key) (stop bool, err error)) error {
	cancel := make(chan bool)

	err := startListener()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case c := <-cancel:
				if c {
					return
				}
			case keyInfo := <-mockChannel:
				onKeyPress(keyInfo)
			}
		}
	}()

	for {
		key, err := getKeyPress(inputTTY)
		if err != nil {
			return err
		}

		stop, err := onKeyPress(key)
		if err != nil {
			return err
		}

		if stop {
			break
		}
	}

	err = stopListener()
	if err != nil {
		return err
	}

	cancel <- true

	return nil
}

func MockKey(key keys.Key) error {
	mockChannel <- key

	return nil
}
