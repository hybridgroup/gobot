package system

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"gobot.io/x/gobot"
)

const (
	// gpioPath default linux sysfs gpio path
	gpioPath = "/sys/class/gpio"
)

var errNotExported = errors.New("pin has not been exported")

// digitalPin represents a digital pin
type digitalPinSysfs struct {
	pin string
	*digitalPinConfig
	fs filesystem

	dirFile File
	valFile File
}

// newDigitalPinSysfs returns a digital pin using for the given number. The name of the sysfs file will prepend "gpio"
// to the pin number, eg. a pin number of 10 will have a name of "gpio10". The pin is handled by the sysfs Kernel ABI.
func newDigitalPinSysfs(fs filesystem, pin string, options ...func(gobot.DigitalPinOptioner) bool) *digitalPinSysfs {
	cfg := newDigitalPinConfig("gpio"+pin, options...)
	d := &digitalPinSysfs{
		pin:              pin,
		digitalPinConfig: cfg,
		fs:               fs,
	}
	return d
}

// ApplyOptions apply all given options to the pin immediately. Implements interface gobot.DigitalPinOptionApplier.
func (d *digitalPinSysfs) ApplyOptions(options ...func(gobot.DigitalPinOptioner) bool) error {
	anyChange := false
	for _, option := range options {
		anyChange = anyChange || option(d)
	}
	if anyChange {
		return d.reconfigure()
	}
	return nil
}

// DirectionBehavior gets the direction behavior when the pin is used the next time. This means its possibly not in
// this direction type at the moment. Implements the interface gobot.DigitalPinValuer, but should be rarely used.
func (d *digitalPinSysfs) DirectionBehavior() string {
	return d.direction
}

// Export sets the pin as exported with the configured direction
func (d *digitalPinSysfs) Export() error {
	err := d.reconfigure()
	return err
}

// Unexport release the pin
func (d *digitalPinSysfs) Unexport() error {
	unexport, err := d.fs.openFile(gpioPath+"/unexport", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer unexport.Close()

	if d.dirFile != nil {
		d.dirFile.Close()
		d.dirFile = nil
	}
	if d.valFile != nil {
		d.valFile.Close()
		d.valFile = nil
	}

	_, err = writeFile(unexport, []byte(d.pin))
	if err != nil {
		// If EINVAL then the pin is reserved in the system and can't be unexported
		e, ok := err.(*os.PathError)
		if !ok || e.Err != syscall.EINVAL {
			return err
		}
	}

	return nil
}

// Write writes the given value to the character device
func (d *digitalPinSysfs) Write(b int) error {
	_, err := writeFile(d.valFile, []byte(strconv.Itoa(b)))
	return err
}

// Read reads the given value from character device
func (d *digitalPinSysfs) Read() (int, error) {
	buf, err := readFile(d.valFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(buf[0]))
}

func (d *digitalPinSysfs) reconfigure() error {
	exportFile, err := d.fs.openFile(gpioPath+"/export", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer exportFile.Close()

	_, err = writeFile(exportFile, []byte(d.pin))
	if err != nil {
		// If EBUSY then the pin has already been exported
		e, ok := err.(*os.PathError)
		if !ok || e.Err != syscall.EBUSY {
			return err
		}
	}

	if d.dirFile != nil {
		d.dirFile.Close()
	}

	attempt := 0
	for {
		attempt++
		d.dirFile, err = d.fs.openFile(fmt.Sprintf("%s/%s/direction", gpioPath, d.label), os.O_RDWR, 0644)
		if err == nil {
			break
		}
		if attempt > 10 {
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}

	if d.valFile != nil {
		d.valFile.Close()
	}
	if err == nil {
		d.valFile, err = d.fs.openFile(fmt.Sprintf("%s/%s/value", gpioPath, d.label), os.O_RDWR, 0644)
	}

	// configure line
	if err == nil {
		err = d.writeDirectionWithInitialOutput()
	}

	if err != nil {
		d.Unexport()
	}

	return err
}

func (d *digitalPinSysfs) writeDirectionWithInitialOutput() error {
	if _, err := writeFile(d.dirFile, []byte(d.direction)); err != nil || d.direction == IN {
		return err
	}
	_, err := writeFile(d.valFile, []byte(strconv.Itoa(d.outInitialState)))
	return err
}

// Linux sysfs / GPIO specific sysfs docs.
//  https://www.kernel.org/doc/Documentation/filesystems/sysfs.txt
//  https://www.kernel.org/doc/Documentation/gpio/sysfs.txt

var writeFile = func(f File, data []byte) (i int, err error) {
	if f == nil {
		return 0, errNotExported
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
		return nil, errNotExported
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
