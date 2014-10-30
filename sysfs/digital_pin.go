package sysfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	IN  = "in"
	OUT = "out"
)

var HIGH = []byte("1")
var LOW = []byte("0")

// writeFile validates file existence and writes data into it
func writeFile(name string, data []byte) (i int, err error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

type DigitalPin struct {
	pin       string
	direction string
}

// newDigitalPin returns an exported digital pin
func NewDigitalPin(pin int) *DigitalPin {
	d := &DigitalPin{pin: strconv.Itoa(pin)}
	d.Export()
	return d
}
func (d *DigitalPin) Direction() string {
	return d.direction
}

// setDir sets writes a directory using direction path for specified pin.
func (d *DigitalPin) SetDirection(dir string) error {
	d.direction = dir
	_, err := writeFile(fmt.Sprintf("/sys/class/gpio/%v/direction", d.pin), []byte(d.direction))
	return err
}

// Write writes specified value to gpio value path
func (d *DigitalPin) Write(p []byte) (n int, err error) {
	return writeFile(fmt.Sprintf("/sys/class/gpio/%v/value", d.pin), p)
}

// Read reads from gpio value path
func (d *DigitalPin) Read(p []byte) (n int, err error) {
	p, err = ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/%v/value", d.pin))
	return len(p), err
}

// export writes directory for gpio export path
func (d *DigitalPin) Export() error {
	_, err := writeFile("/sys/class/gpio/export", []byte(d.pin))
	return err
}

// unexport writes directory for gpio unexport path
func (d *DigitalPin) Unexport() error {
	_, err := writeFile("/sys/class/gpio/unexport", []byte(d.pin))
	return err
}
