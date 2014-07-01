package beaglebone

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

const GPIOPath = "/sys/class/gpio"
const GPIODirectionRead = "in"
const GPIODirectionWrite = "out"
const HIGH = 1
const LOW = 0

func newDigitalPin(pinNum int, mode string) *digitalPin {
	d := new(digitalPin)
	d.PinNum = strconv.Itoa(pinNum)

	fi, err := os.OpenFile(GPIOPath+"/export", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	fi.WriteString(d.PinNum)

	d.setMode(mode)

	return d
}

func (d *digitalPin) setMode(mode string) {
	d.Mode = mode

	if mode == "w" {
		fi, err := os.OpenFile(GPIOPath+"/gpio"+d.PinNum+"/direction", os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		fi.WriteString(GPIODirectionWrite)
		fi.Close()
		d.PinFile, err = os.OpenFile(GPIOPath+"/gpio"+d.PinNum+"/value", os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
	} else if mode == "r" {
		fi, err := os.OpenFile(GPIOPath+"/gpio"+d.PinNum+"/direction", os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		fi.WriteString(GPIODirectionRead)
		fi.Close()
		d.PinFile, err = os.OpenFile(GPIOPath+"/gpio"+d.PinNum+"/value", os.O_RDONLY, 0666)
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

func (d *digitalPin) digitalRead() int {
	if d.Mode != "r" {
		d.setMode("r")
	}

	var buf []byte = make([]byte, 1)
	d.PinFile.ReadAt(buf, 0)

	i, _ := strconv.Atoi(string(buf[0]))
	return i
}

func (d *digitalPin) close() {
	fi, err := os.OpenFile(GPIOPath+"/unexport", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	fi.WriteString(d.PinNum)
	fi.Close()
	d.PinFile.Close()
}
