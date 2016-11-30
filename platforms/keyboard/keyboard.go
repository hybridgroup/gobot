package keyboard

import (
	"os"
	"os/exec"
)

type bytes [3]byte

// KeyEvent contains data about a keyboard event
type KeyEvent struct {
	Bytes bytes
	Key   int
	Char  string
}

const (
	Tilde = iota + 96
	A
	B
	C
	D
	E
	F
	G
	H
	I
	J
	K
	L
	M
	N
	O
	P
	Q
	R
	S
	T
	U
	V
	W
	X
	Y
	Z
)

const (
	Escape   = 27
	Spacebar = 32
)

const (
	Zero = iota + 48
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
)

const (
	ArrowUp = iota + 65
	ArrowDown
	ArrowRight
	ArrowLeft
)

// used to hold the original stty state
var originalState string

func Parse(input bytes) KeyEvent {
	var event = KeyEvent{Bytes: input, Char: string(input[:])}

	var code byte

	// basic input codes
	if input[1] == 0 && input[2] == 0 {
		code = input[0]

		// space bar
		if code == Spacebar {
			event.Key = Spacebar
		}

		// vanilla escape
		if code == Escape {
			event.Key = Escape
		}

		// number keys
		if code >= 48 && code <= 57 {
			event.Key = int(code)
		}

		// alphabet
		if code >= 97 && code <= 122 {
			event.Key = int(code)
		}

		return event
	}

	// arrow keys
	if input[0] == Escape && input[1] == 91 {
		code = input[2]

		if code >= 65 && code <= 68 {
			event.Key = int(code)
			return event
		}
	}

	return event
}

// fetches original state, sets up TTY for raw (unbuffered) input
func configure() (err error) {
	state, err := stty("-g")
	if err != nil {
		return err
	}

	originalState = state

	// -echo: terminal doesn't echo typed characters back to the terminal
	// -icanon: terminal doesn't interpret special characters (like backspace)
	if _, err := stty("-echo", "-icanon"); err != nil {
		return err
	}

	return
}

// restores the TTY to the original state
func restore() (err error) {
	if _, err = stty("echo"); err != nil {
		return
	}

	if _, err = stty(originalState); err != nil {
		return
	}

	return
}

func stty(args ...string) (string, error) {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
