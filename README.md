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

<p align="center">
<table>
<tbody>
<td align="center">
<img width="2000" height="0"><br>
  ------------------------------------------------------------------------------------------------------------------------------
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
   ------------------------------------------------------------------------------------------------------------------------------
<img width="2000" height="0">
</td>
</tbody>
</table>
</p>

## Description

Package keyboard can be used to read key presses from the keyboard, while in a
terminal application. It's crossplatform and keypresses can be combined to check
for ctrl+c, alt+4, ctrl-shift, alt+ctrl+right, etc.

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


## Usage

#### func  GetListener

```go
func GetListener() (listener chan keys.Key, cancel chan bool)
```

#### func  Listen

```go
func Listen(onKeyPress func(keyInfo keys.Key) (stop bool, err error)) error
```

#### func  MockKey

```go
func MockKey(key keys.Key) error
```

---

> [AtomicGo.dev](https://atomicgo.dev) &nbsp;&middot;&nbsp;
> with ❤️ by [@MarvinJWendt](https://github.com/MarvinJWendt) |
> [MarvinJWendt.com](https://marvinjwendt.com)
