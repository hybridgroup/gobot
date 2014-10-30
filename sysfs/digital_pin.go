package sysfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	IN       = "in"
	OUT      = "out"
	HIGH     = 1
	LOW      = 0
	GPIOPATH = "/sys/class/gpio"
)

type DigitalPin struct {
	pin       string
	label     string
	direction string
}

// NewDigitalPin returns a DigitalPin given the pin number and sysfs pin label
func NewDigitalPin(pin int, v ...string) *DigitalPin {
	d := &DigitalPin{pin: strconv.Itoa(pin)}
	if len(v) > 0 {
		d.label = v[0]
	} else {
		d.label = "gpio" + d.pin
	}

	return d
}

// Direction returns the current direction of the pin
func (d *DigitalPin) Direction() string {
	return d.direction
}

// SetDirection sets the current direction for specified pin
func (d *DigitalPin) SetDirection(dir string) error {
	d.direction = dir
	_, err := writeFile(fmt.Sprintf("%v/%v/direction", GPIOPATH, d.label), []byte(d.direction))
	return err
}

// Write writes specified value to the pin
func (d *DigitalPin) Write(b int) error {
	_, err := writeFile(fmt.Sprintf("%v/%v/value", GPIOPATH, d.label), []byte(strconv.Itoa(b)))
	return err
}

// Read reads the current value of the pin
func (d *DigitalPin) Read() (n int, err error) {
	buf, err := readFile(fmt.Sprintf("%v/%v/value", GPIOPATH, d.label))
	if err != nil {
		return
	}
	return strconv.Atoi(string(buf[0]))
}

// Export exports the pin for use by the operating system
func (d *DigitalPin) Export() error {
	_, err := writeFile(GPIOPATH+"/export", []byte(d.pin))
	return err
}

// Unexport unexports the pin and releases the pin from the operating system
func (d *DigitalPin) Unexport() error {
	_, err := writeFile(GPIOPATH+"/unexport", []byte(d.pin))
	return err
}

var writeFile = func(path string, data []byte) (i int, err error) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

var readFile = func(path string) (b []byte, err error) {
	return ioutil.ReadFile(path)
}
