package gobotBeaglebone

import (
	"os"
	"strconv"
)

type digitalPin struct {
	PinNum  string
	Mode    string
	PinFile *os.File
	Status  string
}

const GPIO_PATH = "/sys/class/gpio"
const GPIO_DIRECTION_READ = "in"
const GPIO_DIRECTION_WRITE = "out"
const HIGH = 1
const LOW = 0

func newDigitalPin(pinNum int, mode string) *digitalPin {
	d := new(digitalPin)
	d.PinNum = strconv.Itoa(pinNum)

	fi, err := os.OpenFile(GPIO_PATH+"/export", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	fi.WriteString(d.PinNum)
	fi.Close()

	d.setMode(mode)

	return d
}

func (d *digitalPin) setMode(mode string) {
	d.Mode = mode

	if mode == "w" {
		fi, err := os.OpenFile(GPIO_PATH+"/gpio"+d.PinNum+"/direction", os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		fi.WriteString(GPIO_DIRECTION_WRITE)
		fi.Close()
		d.PinFile, err = os.OpenFile(GPIO_PATH+"/gpio"+d.PinNum+"/value", os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
	} else if mode == "r" {
		fi, err := os.OpenFile(GPIO_PATH+"/gpio"+d.PinNum+"/direction", os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		fi.WriteString(GPIO_DIRECTION_READ)
		fi.Close()
		d.PinFile, err = os.OpenFile(GPIO_PATH+"/gpio"+d.PinNum+"/value", os.O_RDONLY, 0666)
		if err != nil {
			panic(err)
		}
	}
}

func (d *digitalPin) digitalWrite(value string) {
	if d.Mode != "w" {
		d.setMode("w")
	}

	d.PinFile.WriteString(value)
	d.PinFile.Sync()
}

func (d *digitalPin) close() {
	fi, err := os.OpenFile(GPIO_PATH+"/unexport", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	fi.WriteString(d.PinNum)
	fi.Close()
	d.PinFile.Close()
}
