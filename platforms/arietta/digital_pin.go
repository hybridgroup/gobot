package arietta

import (
	"os"
)

const (
	gpioPath = "/sys/class/gpio"
)

// sysfsDigitalPinDesc describes a digital pin exposed through sysfs.
type sysfsDigitalPinDesc struct {
	// name is the canonical name of the pin.
	name string
	// hwnum is the pin ID used to export the pin.
	hwnum int
	// label is the name the exported pin appears as under
	// /sys/class/gpio.
	label string
}

// sysfsDigitalPin is a digital pin accessed through the Linux sysfs
// API.
type sysfsDigitalPin struct {
	desc *sysfsDigitalPinDesc
	file *os.File
	// mode is the current mode if any.  Used to lazily set the
	// mode.
	mode string
}

func newSysfsDigitalPin(desc *sysfsDigitalPinDesc) *sysfsDigitalPin {
	d := &sysfsDigitalPin{
		desc: desc,
	}

	writeInt(desc.hwnum, gpioPath, "export")
	return d
}

func (d *sysfsDigitalPin) setMode(mode string) {
	if mode != d.mode {
		d.closeFile()
		path := joinPath(gpioPath, d.desc.label)

		switch {
		case mode == "w":
			writeStr("out", path, "direction")
			d.file = openOrDie(os.O_WRONLY, path, "value")
		case mode == "r":
			writeStr("in", path, "direction")
			d.file = openOrDie(os.O_RDONLY, path, "value")
		default:
			panic("Illegal pin mode.")
		}
		d.mode = mode
	}
}

func (d *sysfsDigitalPin) DigitalWrite(value byte) {
	d.setMode("w")
	if value == 0 {
		d.file.WriteString("0")
	} else {
		d.file.WriteString("1")
	}
	// TODO(michaelh): check if needed.
	d.file.Sync()
}

func (d *sysfsDigitalPin) DigitalRead() int {
	d.setMode("r")
	var buf []byte = make([]byte, 1)
	d.file.ReadAt(buf, 0)

	if buf[0] == '0' {
		return 0
	} else {
		return 1
	}
}

func (d *sysfsDigitalPin) closeFile() {
	if d.file != nil {
		d.file.Close()
		d.file = nil
	}
}

func (d *sysfsDigitalPin) Finalize() bool {
	d.closeFile()
	writeInt(d.desc.hwnum, gpioPath, "unexport")
	return true
}
