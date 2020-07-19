package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"text/scanner"
	"time"

	"github.com/micmonay/keybd_event"
)

func main() {

	// Initializing keyboard device
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// Waiting 2 seconds because Linux
	// See here for more info:
	// https://www.kernel.org/doc/html/v4.12/input/uinput.html#keyboard-events
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	// Changing focused window i.e. alt+tabbing
	kb.SetKeys(keybd_event.VK_J)
	kb.HasSuper(true)

	err = kb.Launching()
	if err != nil {
		panic(err)
	}

	// Setting the super key back to false so all other keypresses aren't prefixed with super+key
	kb.HasSuper(false)

	rsPass := retrievePass()

	// Reading the unicode characters of the password string
	var b scanner.Scanner
	b.Init(strings.NewReader(rsPass))

	// Converting the unicode values to a string and passing that to our conversion function
	for i := 0; i < len(rsPass); i++ {
		currentPos := b.Next()
		strOfPos := string(currentPos)
		inputVal, shiftBool := mapStringToInput(strOfPos)

		// If true is returned, the value needs the SHIFT key pressed to be typed correctly
		if shiftBool == true {
			kb.SetKeys(inputVal)
			kb.HasSHIFT(true)

			// Press the selected keys
			err = kb.Launching()
			if err != nil {
				panic(err)
			}
		} else {
			kb.SetKeys(inputVal)
			kb.HasSHIFT(false)

			// Press the selected keys
			err = kb.Launching()
			if err != nil {
				panic(err)
			}
		}

	}

}

func retrievePass() (passOut string) {

	// Telling go what command to run with it's args
	cmd := exec.Command("pass", "passEntryForRuneScapeHere")

	// Redirecting the command's output into a buffer
	cmdOut := &bytes.Buffer{}
	cmd.Stdout = cmdOut

	// Running the command and catching errors
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}

	// Sending the buffered output to get returned as a string
	output := cmdOut.Bytes()

	if len(output) > 0 {
		// Converting the bytes object to a string
		// Splitting on the newline
		// Returning the line with the password
		outString := string(output)
		split := strings.Split(outString, "\n")
		str := split[0]
		return str
	} else {

		// TODO: Return something more helpful than this, if we get here, retrieving the password failed
		return ""
	}

}

// This function will take the character input of a string and return it's int value
// The int value can be directly translated to a keyboard press via kb.SetKeys()
// This script also returns a bool which will set hasSHIFT to true
// The bool will handle capital letters and special chars requiring shift
func mapStringToInput(inputString string) (keypressOut int, shiftPressNeeded bool) {

	const upperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const upperSpecial = "~!@#$%^&*()_+{}|:\"<>?"

	m := make(map[string]int)

	// Setting keyboard values to pass to KeyPress
	m["1"] = 2
	m["2"] = 3
	m["3"] = 4
	m["4"] = 5
	m["5"] = 6
	m["6"] = 7
	m["7"] = 8
	m["8"] = 9
	m["9"] = 10
	m["0"] = 11
	m["q"] = 16
	m["w"] = 17
	m["e"] = 18
	m["r"] = 19
	m["t"] = 20
	m["y"] = 21
	m["u"] = 22
	m["i"] = 23
	m["o"] = 24
	m["p"] = 25
	m["a"] = 30
	m["s"] = 31
	m["d"] = 32
	m["f"] = 33
	m["g"] = 34
	m["h"] = 35
	m["j"] = 36
	m["k"] = 37
	m["l"] = 38
	m["z"] = 44
	m["x"] = 45
	m["c"] = 46
	m["v"] = 47
	m["b"] = 48
	m["n"] = 49
	m["m"] = 50

	// Special characters
	m["-"] = 12  // Upper _
	m["="] = 13  // Upper +
	m["["] = 26  // Upper {
	m["]"] = 27  // upper }
	m[";"] = 39  // Upper :
	m["'"] = 40  // Upper "
	m["`"] = 41  // Upper ~
	m["\\"] = 43 // Upper |
	m[","] = 51  // Upper <
	m["."] = 52  // Upper >
	m["/"] = 53  // Upper ?
	m[" "] = 57

	if strings.Contains(upperCase, inputString) {

		lowerCase := strings.ToLower(inputString)
		return m[lowerCase], true

	}

	if strings.Contains(upperSpecial, inputString) {

		if inputString == "!" {
			return m["1"], true
		}
		if inputString == "@" {
			return m["2"], true
		}
		if inputString == "#" {
			return m["3"], true
		}
		if inputString == "$" {
			return m["4"], true
		}
		if inputString == "%" {
			return m["5"], true
		}
		if inputString == "^" {
			return m["6"], true
		}
		if inputString == "&" {
			return m["7"], true
		}
		if inputString == "*" {
			return m["8"], true
		}
		if inputString == "(" {
			return m["9"], true
		}
		if inputString == ")" {
			return m["0"], true
		}
		if inputString == "_" {
			return m["-"], true
		}
		if inputString == "+" {
			return m["="], true
		}
		if inputString == "{" {
			return m["["], true
		}
		if inputString == "}" {
			return m["]"], true
		}
		if inputString == "|" {
			return m["\\"], true
		}
		if inputString == ":" {
			return m[";"], true
		}
		if inputString == "\"" {
			return m["'"], true
		}
		if inputString == "<" {
			return m[","], true
		}
		if inputString == ">" {
			return m["."], true
		}
		if inputString == "?" {
			return m["/"], true
		}
	}

	return m[inputString], false
}
