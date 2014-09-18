package edison

import (
	"io/ioutil"
	"strconv"
)

func gpioPath() string {
	return "/sys/class/gpio"
}
func gpioExportPath() string {
	return gpioPath() + "/export"
}
func gpioUnExportPath() string {
	return gpioPath() + "/unexport"
}
func gpioDirectionPath(pin string) string {
	return gpioPath() + "/gpio" + pin + "/direction"
}
func gpioValuePath(pin string) string {
	return gpioPath() + "/gpio" + pin + "/value"
}

type digitalPin struct {
	pin string
	dir string
}

func newDigitalPin(pin int) *digitalPin {
	d := &digitalPin{pin: strconv.Itoa(pin)}
	d.export()
	return d
}

func (d *digitalPin) setDir(dir string) {
	d.dir = dir
	err := writeFile(gpioDirectionPath(d.pin), dir)
	if err != nil {
		panic(err)
	}
}

func (d *digitalPin) digitalWrite(value string) {
	if d.dir != "out" {
		d.setDir("out")
	}
	err := writeFile(gpioValuePath(d.pin), value)
	if err != nil {
		panic(err)
	}
}

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

func (d *digitalPin) export() {
	writeFile(gpioExportPath(), d.pin)
}

func (d *digitalPin) unexport() {
	writeFile(gpioUnExportPath(), d.pin)
}

func (d *digitalPin) close() {
	d.unexport()
}
