package beaglebone

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type analogPin struct {
	pinNum string
}

// newAnalogPing returns a new analog pin with specified pin number
func newAnalogPin(pinNum string) *analogPin {
	d := new(analogPin)
	d.pinNum = pinNum

	return d
}

// analogRead returns the value from the specified analog pin.
// returns an int value between 0 and 1023.
func (a *analogPin) analogRead() int {
	var err error
	var fi *os.File

	ocp, err := filepath.Glob(Ocp)
	if err != nil {
		panic(err)
	}

	helper, err := filepath.Glob(fmt.Sprintf("%v/helper.*", ocp[0]))
	if err != nil {
		panic(err)
	}

	fi, err = os.Open(fmt.Sprintf("%v/%v", helper[0], a.pinNum))
	if err != nil {
		panic(err)
	}

	var buf = make([]byte, 1024)
	fi.Read(buf)
	fi.Close()

	i, _ := strconv.Atoi(strings.Split(string(buf), "\n")[0])
	return i
}
