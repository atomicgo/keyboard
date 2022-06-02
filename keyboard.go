package keyboard

import (
	"fmt"
	"os"

	"atomicgo.dev/keyboard/keys"
	"github.com/containerd/console"
)

var windowsStdin *os.File
var con console.Console
var input = os.Stdin
var inputTTY *os.File

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
func StartListener() error {
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
func StopListener() error {
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
func GetKey() (keys.Key, error) {
	return getKeyPress(inputTTY)
}
