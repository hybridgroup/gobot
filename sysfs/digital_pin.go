package sysfs

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

const (
	// IN gpio direction
	IN = "in"
	// OUT gpio direction
	OUT = "out"
	// HIGH gpio level
	HIGH = 1
	// LOW gpio level
	LOW = 0
	// GPIOPATH default linux gpio path
	GPIOPATH = "/sys/class/gpio"
)

// DigitalPin is the interface for sysfs gpio interactions
type DigitalPin interface {
	// Unexport unexports the pin and releases the pin from the operating system
	Unexport() error
	// Export exports the pin for use by the operating system
	Export() error
	// Read reads the current value of the pin
	Read() (int, error)
	// Direction sets the direction for the pin
	Direction(string) error
	// Write writes to the pin
	Write(int) error
}

type digitalPin struct {
	pin       string
	label     string
	value     File
	direction File
}

var (
	notExportedErr = errors.New("pin not exported")
)

// NewDigitalPin returns a DigitalPin given the pin number and an optional sysfs pin label.
// If no label is supplied the default label will prepend "gpio" to the pin number,
// eg. a pin number of 10 will have a label of "gpio10"
func NewDigitalPin(pin int, v ...string) DigitalPin {
	d := &digitalPin{pin: strconv.Itoa(pin)}
	if len(v) > 0 {
		d.label = v[0]
	} else {
		d.label = "gpio" + d.pin
	}

	return d
}

func (d *digitalPin) Direction(dir string) error {
	_, err := writeFile(d.direction, []byte(dir))
	return err
}

func (d *digitalPin) Write(b int) error {
	_, err := writeFile(d.value, []byte(strconv.Itoa(b)))
	return err
}

func (d *digitalPin) Read() (n int, err error) {
	buf, err := readFile(d.value)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(buf[0]))
}

func (d *digitalPin) Export() error {
	export, err := OpenFile(GPIOPATH+"/export", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer export.Close()

	if _, err := writeFile(export, []byte(d.pin)); err != nil {
		// If EBUSY then the pin has already been exported
		if err.(*os.PathError).Err != syscall.EBUSY {
			return err
		}
	}
	if d.direction == nil {
		d.direction, err = OpenFile(fmt.Sprintf("%v/%v/direction", GPIOPATH, d.label), os.O_RDWR, 0644)
	}
	if d.value == nil && err == nil {
		d.value, err = OpenFile(fmt.Sprintf("%v/%v/value", GPIOPATH, d.label), os.O_RDWR, 0644)
	}
	return err
}

func (d *digitalPin) Unexport() error {
	unexport, err := OpenFile(GPIOPATH+"/unexport", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer unexport.Close()

	if _, err := writeFile(unexport, []byte(d.pin)); err != nil {
		// If EINVAL then the pin is reserved in the system and can't be unexported
		if err.(*os.PathError).Err != syscall.EINVAL {
			return err
		}
	}

	if d.direction != nil {
		d.direction.Close()
		d.direction = nil
	}

	if d.value != nil {
		d.value.Close()
		d.value = nil
	}

	return nil
}

var writeFile = func(file File, data []byte) (i int, err error) {
	if file == nil {
		return 0, notExportedErr
	}
	return file.Write(data)
}

var readFile = func(file File) ([]byte, error) {
	if file == nil {
		return make([]byte, 0), notExportedErr
	}

	buf := make([]byte, 2)
	_, err := file.Read(buf)
	return buf, err
}
