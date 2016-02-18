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
	pin   string
	label string

	value     File
	direction File
}

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

var notExportedError = errors.New("pin has not been exported")

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
	export, err := fs.OpenFile(GPIOPATH+"/export", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer export.Close()

	_, err = writeFile(export, []byte(d.pin))
	if err != nil {
		// If EBUSY then the pin has already been exported
		if err.(*os.PathError).Err != syscall.EBUSY {
			return err
		}
	}

	if d.direction != nil {
		d.direction.Close()
	}

	d.direction, err = fs.OpenFile(fmt.Sprintf("%v/%v/direction", GPIOPATH, d.label), os.O_RDWR, 0644)

	if d.value != nil {
		d.value.Close()
	}
	if err == nil {
		d.value, err = fs.OpenFile(fmt.Sprintf("%v/%v/value", GPIOPATH, d.label), os.O_RDWR, 0644)
	}

	if err != nil {
		// Should we unexport here?
		// If we don't unexport we should make sure to close d.direction and d.value here
		d.Unexport()
	}

	return err
}

func (d *digitalPin) Unexport() error {
	unexport, err := fs.OpenFile(GPIOPATH+"/unexport", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer unexport.Close()

	if d.direction != nil {
		d.direction.Close()
		d.direction = nil
	}
	if d.value != nil {
		d.value.Close()
		d.value = nil
	}

	_, err = writeFile(unexport, []byte(d.pin))
	if err != nil {
		// If EINVAL then the pin is reserved in the system and can't be unexported
		if err.(*os.PathError).Err != syscall.EINVAL {
			return err
		}
	}

	return nil
}

// Linux sysfs / GPIO specific sysfs docs.
//  https://www.kernel.org/doc/Documentation/filesystems/sysfs.txt
//  https://www.kernel.org/doc/Documentation/gpio/sysfs.txt

var writeFile = func(f File, data []byte) (i int, err error) {
	if f == nil {
		return 0, notExportedError
	}

	// sysfs docs say:
	// > When writing sysfs files, userspace processes should first read the
	// > entire file, modify the values it wishes to change, then write the
	// > entire buffer back.
	// however, this seems outdated/inaccurate (docs are from back in the Kernel BitKeeper days).

	i, err = f.Write(data)
	return i, err
}

var readFile = func(f File) ([]byte, error) {
	if f == nil {
		return nil, notExportedError
	}

	// sysfs docs say:
	// > If userspace seeks back to zero or does a pread(2) with an offset of '0' the [..] method will
	// > be called again, rearmed, to fill the buffer.

	// TODO: Examine if seek is needed if full buffer is read from sysfs file.

	buf := make([]byte, 2)
	_, err := f.Seek(0, os.SEEK_SET)
	if err == nil {
		_, err = f.Read(buf)
	}
	return buf, err
}
