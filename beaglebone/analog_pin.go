package gobotBeaglebone

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

func newAnalogPin(pinNum string) *analogPin {
	var err error
	var fi *os.File

	d := new(analogPin)
	d.pinNum = pinNum

	slot, err := filepath.Glob(SLOTS)
	if err != nil {
		panic(err)
	}
	fi, err = os.OpenFile(fmt.Sprintf("%v/slots", slot[0]), os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	fi.WriteString("cape-bone-iio")
	fi.Close()
	return d
}

func (me *analogPin) analogRead() int {
	var err error
	var fi *os.File

	ocp, err := filepath.Glob(OCP)
	if err != nil {
		panic(err)
	}

	helper, err := filepath.Glob(fmt.Sprintf("%v/helper.*", ocp[0]))
	if err != nil {
		panic(err)
	}

	fi, err = os.Open(fmt.Sprintf("%v/%v", helper[0], me.pinNum))
	if err != nil {
		panic(err)
	}

	var buf []byte = make([]byte, 1024)
	fi.Read(buf)
	fi.Close()

	i, _ := strconv.Atoi(strings.Split(string(buf), "\n")[0])
	return i
}
