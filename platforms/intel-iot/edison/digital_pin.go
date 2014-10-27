package edison

import (
	"io/ioutil"
	"strconv"
)

// gpioPath returns base gpio path
func gpioPath() string {
	return "/sys/class/gpio"
}

// gpioExportPath returns export path
func gpioExportPath() string {
	return gpioPath() + "/export"
}

// gpioUnExportPath returns unexport path
func gpioUnExportPath() string {
	return gpioPath() + "/unexport"
}

// gpioDirectionPath returns direction path for specified pin
func gpioDirectionPath(pin string) string {
	return gpioPath() + "/gpio" + pin + "/direction"
}

// gpioValuePath returns value path for specified pin
func gpioValuePath(pin string) string {
	return gpioPath() + "/gpio" + pin + "/value"
}

type digitalPin struct {
	pin string
	dir string
}

// newDigitalPin returns an exported digital pin
func newDigitalPin(pin int) *digitalPin {
	d := &digitalPin{pin: strconv.Itoa(pin)}
	d.export()
	return d
}

// setDir sets writes a directory using direction path for specified pin.
// It panics on error
func (d *digitalPin) setDir(dir string) {
	d.dir = dir
	err := writeFile(gpioDirectionPath(d.pin), dir)
	if err != nil {
		panic(err)
	}
}

// digitalWrite writes specified value to gpio value path
// It panics on error
func (d *digitalPin) digitalWrite(value string) {
	if d.dir != "out" {
		d.setDir("out")
	}
	err := writeFile(gpioValuePath(d.pin), value)
	if err != nil {
		panic(err)
	}
}

// digitalRead reads from gpio value path
func (d *digitalPin) digitalRead() int {
	if d.dir != "in" {
		d.setDir("in")
	}

	buf, err := ioutil.ReadFile(gpioValuePath(d.pin))
	if err != nil {
		panic(err)
	}

	i, err := strconv.Atoi(string(buf[0]))
	if err != nil {
		panic(err)
	}
	return i
}

// export writes directory for gpio export path
func (d *digitalPin) export() {
	writeFile(gpioExportPath(), d.pin)
}

// unexport writes directory for gpio unexport path
func (d *digitalPin) unexport() {
	writeFile(gpioUnExportPath(), d.pin)
}

// close unexports digital pin
func (d *digitalPin) close() {
	d.unexport()
}
