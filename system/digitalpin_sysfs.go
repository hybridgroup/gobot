package system

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gobot.io/x/gobot/v2"
)

const (
	systemSysfsDebug = false
	// gpioPath default linux sysfs gpio path
	gpioPath = "/sys/class/gpio"
)

var errNotExported = errors.New("pin has not been exported")

// digitalPin represents a digital pin
type digitalPinSysfs struct {
	pin string
	*digitalPinConfig
	sfa *sysfsFileAccess

	dirFile       *sysfsFile
	valFile       *sysfsFile
	activeLowFile *sysfsFile
}

// newDigitalPinSysfs returns a digital pin using for the given number. The name of the sysfs file will prepend "gpio"
// to the pin number, eg. a pin number of 10 will have a name of "gpio10". The pin is handled by the sysfs Kernel ABI.
func newDigitalPinSysfs(
	sfa *sysfsFileAccess,
	pin string,
	options ...func(gobot.DigitalPinOptioner) bool,
) *digitalPinSysfs {
	cfg := newDigitalPinConfig("gpio"+pin, options...)
	d := &digitalPinSysfs{
		pin:              pin,
		digitalPinConfig: cfg,
		sfa:              sfa,
	}
	return d
}

// ApplyOptions apply all given options to the pin immediately. Implements interface gobot.DigitalPinOptionApplier.
func (d *digitalPinSysfs) ApplyOptions(options ...func(gobot.DigitalPinOptioner) bool) error {
	anyChange := false
	for _, option := range options {
		anyChange = option(d) || anyChange
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
	return d.reconfigure()
}

// Unexport release the pin
func (d *digitalPinSysfs) Unexport() error {
	unexport, err := d.sfa.openWrite(gpioPath + "/unexport")
	if err != nil {
		return err
	}
	defer unexport.close()

	if d.dirFile != nil {
		d.dirFile.close()
		d.dirFile = nil
	}
	if d.valFile != nil {
		d.valFile.close()
		d.valFile = nil
	}
	if d.activeLowFile != nil {
		d.activeLowFile.close()
		d.activeLowFile = nil
	}

	err = unexport.write([]byte(d.pin))
	if err != nil {
		// If EINVAL then the pin is reserved in the system and can't be unexported, we suppress the error
		var pathError *os.PathError
		if !(errors.As(err, &pathError) && errors.Is(err, Syscall_EINVAL)) {
			return err
		}
	}

	return nil
}

// Write writes the given value to the character device
func (d *digitalPinSysfs) Write(b int) error {
	if d.valFile == nil {
		return errNotExported
	}
	err := d.valFile.write([]byte(strconv.Itoa(b)))
	return err
}

// Read reads a value from character device
func (d *digitalPinSysfs) Read() (int, error) {
	if d.valFile == nil {
		return 0, errNotExported
	}
	buf, err := d.valFile.read()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(buf[0]))
}

func (d *digitalPinSysfs) reconfigure() error {
	exportFile, err := d.sfa.openWrite(gpioPath + "/export")
	if err != nil {
		return err
	}
	defer exportFile.close()

	err = exportFile.write([]byte(d.pin))
	if err != nil {
		// If EBUSY then the pin has already been exported, we suppress the error
		var pathError *os.PathError
		if !(errors.As(err, &pathError) && errors.Is(err, Syscall_EBUSY)) {
			return err
		}
	}

	if d.dirFile != nil {
		d.dirFile.close()
	}

	attempt := 0
	for {
		attempt++
		d.dirFile, err = d.sfa.openReadWrite(fmt.Sprintf("%s/%s/direction", gpioPath, d.label))
		if err == nil {
			break
		}
		if attempt > 10 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if d.valFile != nil {
		d.valFile.close()
	}
	if err == nil {
		d.valFile, err = d.sfa.openReadWrite(fmt.Sprintf("%s/%s/value", gpioPath, d.label))
	}

	// configure direction
	if err == nil {
		err = d.writeDirectionWithInitialOutput()
	}

	// configure inverse logic
	if err == nil {
		if d.activeLow {
			d.activeLowFile, err = d.sfa.openReadWrite(fmt.Sprintf("%s/%s/active_low", gpioPath, d.label))
			if err == nil {
				err = d.activeLowFile.write([]byte("1"))
			}
		}
	}

	// configure bias (inputs and outputs, unsupported)
	if err == nil {
		if d.bias != digitalPinBiasDefault && systemSysfsDebug {
			log.Printf("bias options (%d) are not supported by sysfs, please use hardware resistors instead\n", d.bias)
		}
	}

	// configure debounce period (inputs only), edge detection (inputs only) and drive (outputs only)
	if d.direction == IN {
		// configure debounce (unsupported)
		if d.debouncePeriod != 0 && systemSysfsDebug {
			log.Printf("debounce period option (%d) is not supported by sysfs\n", d.debouncePeriod)
		}

		// configure edge detection
		if err == nil {
			if d.edge != 0 && d.pollInterval <= 0 {
				err = fmt.Errorf("edge detect option (%d) is not implemented for sysfs without discrete polling", d.edge)
			}
		}

		// start discrete polling function and wait for first read is done
		if err == nil {
			if d.pollInterval > 0 {
				err = startEdgePolling(d.label, d.Read, d.pollInterval, d.edge, d.edgeEventHandler, d.pollQuitChan)
			}
		}
	} else if d.drive != digitalPinDrivePushPull && systemSysfsDebug {
		// configure drive (unsupported)
		log.Printf("drive options (%d) are not supported by sysfs\n", d.drive)
	}

	if err != nil {
		if e := d.Unexport(); e != nil {
			err = fmt.Errorf("unexport error '%v' after '%v'", e, err)
		}
	}

	return err
}

func (d *digitalPinSysfs) writeDirectionWithInitialOutput() error {
	if d.dirFile == nil {
		return errNotExported
	}
	if err := d.dirFile.write([]byte(d.direction)); err != nil || d.direction == IN {
		return err
	}

	if d.valFile == nil {
		return errNotExported
	}

	return d.valFile.write([]byte(strconv.Itoa(d.outInitialState)))
}
