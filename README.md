<h1 align="center">AtomicGo | keyboard</h1>

<p align="center">

<a href="https://github.com/atomicgo/keyboard/releases">
<img src="https://img.shields.io/github/v/release/atomicgo/keyboard?style=flat-square" alt="Latest Release">
</a>

<a href="https://codecov.io/gh/atomicgo/keyboard" target="_blank">
<img src="https://img.shields.io/github/workflow/status/atomicgo/keyboard/Go?label=tests&style=flat-square" alt="Tests">
</a>

<a href="https://codecov.io/gh/atomicgo/keyboard" target="_blank">
<img src="https://img.shields.io/codecov/c/gh/atomicgo/keyboard?color=magenta&logo=codecov&style=flat-square" alt="Coverage">
</a>

<a href="https://codecov.io/gh/atomicgo/keyboard">
<!-- unittestcount:start --><img src="https://img.shields.io/badge/Unit_Tests-0-magenta?style=flat-square" alt="Unit test count"><!-- unittestcount:end -->
</a>

<a href="https://github.com/atomicgo/keyboard/issues">
<img src="https://img.shields.io/github/issues/atomicgo/keyboard.svg?style=flat-square" alt="Issues">
</a>

<a href="https://opensource.org/licenses/MIT" target="_blank">
<img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
</a>

</p>

---

<p align="center">
<strong><a href="#install">Get The Module</a></strong>
|
<strong><a href="https://pkg.go.dev/atomicgo.dev/keyboard#section-documentation" target="_blank">Documentation</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CONTRIBUTING.md" target="_blank">Contributing</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CODE_OF_CONDUCT.md" target="_blank">Code of Conduct</a></strong>
</p>

---

<p align="center">
  <img src="https://raw.githubusercontent.com/atomicgo/atomicgo/main/assets/header.png" alt="AtomicGo">
</p>

## Description

Package keyboard can be used to read key presses from the keyboard, while in a
terminal application. It's crossplatform and keypresses can be combined.

Example:

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

## Install

<p align="center">
<table>
<tbody>
<td align="center">
<img width="2000" height="0"><br>
  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
<img width="2000" height="0">
</td>
</tbody>
</table>
</p>
<h3  align="center"><pre>go get atomicgo.dev/keyboard</pre></h3>
<p align="center">
<table>
<tbody>
<td align="center">
<img width="2000" height="0"><br>
   ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
<img width="2000" height="0">
</td>
</tbody>
</table>
</p>

```go
// Add this to your imports
import "atomicgo.dev/keyboard"
```

## Usage

#### func  GetKey

```go
func GetKey() (keys.Key, error)
```
GetKey blocks until a key is pressed and returns the key info.

Example:

    keyboard.StartListener()

    for {
      keyInfo, _ := keyboard.GetKey()
      key := keyInfo.Code

      if key == keys.CtrlC {
        break
      }

      fmt.Println("\r", keyInfo.String())
    }

    keyboard.StopListener()

#### func  StartListener

```go
func StartListener() error
```
StartListener starts the keyboard listener

Example:

    keyboard.StartListener()

    for {
      keyInfo, _ := keyboard.GetKey()
      key := keyInfo.Code

      if key == keys.CtrlC {
        break
      }

      fmt.Println("\r", keyInfo.String())
    }

    keyboard.StopListener()

#### func  StopListener

```go
func StopListener() error
```
StopListener stops the keyboard listener

Example:

    keyboard.StartListener()

    for {
      keyInfo, _ := keyboard.GetKey()
      key := keyInfo.Code

      if key == keys.CtrlC {
        break
      }

      fmt.Println("\r", keyInfo.String())
    }

    keyboard.StopListener()

---

> [AtomicGo.dev](https://atomicgo.dev) &nbsp;&middot;&nbsp;
> with ❤️ by [@MarvinJWendt](https://github.com/MarvinJWendt) |
> [MarvinJWendt.com](https://marvinjwendt.com)
