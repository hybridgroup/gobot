//go:build !linux

package system

import (
	"fmt"
	"strconv"

	"gobot.io/x/gobot/v2"
)

var errGPIOCdevNotSupported = fmt.Errorf("GPIO character device is not supported on this platform")

type cdevLine interface {
	SetValue(value int) error
	Value() (int, error)
	Close() error
}

type digitalPinCdev struct {
	*digitalPinConfig

	chipName string
	pin      int
	line     cdevLine
}

var digitalPinCdevReconfigure = func(*digitalPinCdev, bool) error {
	return errGPIOCdevNotSupported
}

func newDigitalPinCdev(chipName string, pin int, options ...func(gobot.DigitalPinOptioner) bool) *digitalPinCdev {
	if chipName == "" {
		chipName = "gpiochip0"
	}
	cfg := newDigitalPinConfig("gobotio"+strconv.Itoa(pin), options...)
	return &digitalPinCdev{
		chipName:         chipName,
		pin:              pin,
		digitalPinConfig: cfg,
	}
}

func (d *digitalPinCdev) ApplyOptions(_ ...func(gobot.DigitalPinOptioner) bool) error {
	return errGPIOCdevNotSupported
}

func (d *digitalPinCdev) DirectionBehavior() string {
	return d.direction
}

func (d *digitalPinCdev) Export() error {
	return errGPIOCdevNotSupported
}

func (d *digitalPinCdev) Unexport() error {
	return errGPIOCdevNotSupported
}

func (d *digitalPinCdev) Write(_ int) error {
	return errGPIOCdevNotSupported
}

func (d *digitalPinCdev) Read() (int, error) {
	return 0, errGPIOCdevNotSupported
}

func (d *digitalPinCdev) ListLines() error {
	return errGPIOCdevNotSupported
}

func (d *digitalPinCdev) List() error {
	return errGPIOCdevNotSupported
}
