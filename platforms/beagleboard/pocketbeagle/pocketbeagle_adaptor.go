package pocketbeagle

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/platforms/beagleboard/beaglebone"
	"gobot.io/x/gobot/v2/system"
)

const pwmPeriodDefault = 500000 // 0.5 ms = 2 kHz

// PocketBeagleAdaptor is the Gobot Adaptor for the PocketBeagle
// For more information check out:
//
//	http://beagleboard.org/pocket
type PocketBeagleAdaptor struct {
	*beaglebone.Adaptor
	sys *system.Accesser
}

// NewAdaptor creates a new Adaptor for the PocketBeagle
// tested with sysfs and cdev on "Linux BeagleBone 5.10.168-ti-r71"
//
// Optional parameters:
//
//	adaptors.WithGpioCdevAccess():	use character device driver instead of sysfs
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//	adaptors.WithGpiosPullUp(pins): will be silently ignored (for some pins)
//	adaptors.WithGpioDebounce(inPinNum, debounceTime): is only supported for debounceTime < 8ms
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
func NewAdaptor(opts ...interface{}) *PocketBeagleAdaptor {
	sys := system.NewAccesser()
	a := PocketBeagleAdaptor{
		Adaptor: beaglebone.NewAdaptor(opts...),
		sys:     sys,
	}

	a.SetName(gobot.DefaultName("PocketBeagle"))

	var digitalPinsOpts []adaptors.DigitalPinsOptionApplier
	pwmPinsOpts := []adaptors.PwmPinsOptionApplier{adaptors.WithPWMDefaultPeriod(pwmPeriodDefault)}
	for _, opt := range opts {
		switch o := opt.(type) {
		case adaptors.DigitalPinsOptionApplier:
			digitalPinsOpts = append(digitalPinsOpts, o)
		case adaptors.PwmPinsOptionApplier:
			pwmPinsOpts = append(pwmPinsOpts, o)
		case func(system.Optioner):
			// do nothing, already applied on NewAdaptor(opts)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.Name()))
		}
	}

	analogPinTranslator := adaptors.NewAnalogPinTranslator(sys, analogPinMap)
	digitalPinTranslator := adaptors.NewDigitalPinTranslator(sys, gpioPinDefinitions)
	pwmPinTranslator := adaptors.NewPWMPinTranslator(sys, pwmPinMap)

	a.AnalogPinsAdaptor = adaptors.NewAnalogPinsAdaptor(sys, analogPinTranslator.Translate)
	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys,
		a.getTranslateAndMuxDigitalPinFunc(digitalPinTranslator.Translate), digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, a.getTranslateAndMuxPWMPinFunc(pwmPinTranslator.Translate),
		pwmPinsOpts...)

	return &a
}

func (a *PocketBeagleAdaptor) getTranslateAndMuxDigitalPinFunc(
	digitalPinTranslate func(id string) (string, int, error),
) func(id string) (string, int, error) {
	return func(id string) (string, int, error) {
		if a.sys.IsSysfsDigitalPinAccess() {
			// mux is done by id, not by line
			if err := a.muxPin(id, "gpio"); err != nil {
				return "", -1, err
			}
		}

		return digitalPinTranslate(id)
	}
}

func (a *PocketBeagleAdaptor) getTranslateAndMuxPWMPinFunc(
	pwmPinTranslate func(id string) (string, int, error),
) func(id string) (string, int, error) {
	return func(id string) (string, int, error) {
		path, channel, err := pwmPinTranslate(id)
		if err != nil {
			return path, channel, err
		}

		if err := a.muxPin(id, "pwm"); err != nil {
			return "", -1, err
		}

		return path, channel, nil
	}
}

func (a *PocketBeagleAdaptor) muxPin(pin, cmd string) error {
	path := fmt.Sprintf("/sys/devices/platform/ocp/ocp:%s_pinmux/state", pin)
	fi, e := a.sys.OpenFile(path, os.O_WRONLY, 0o666)
	defer fi.Close() //nolint:staticcheck // for historical reasons
	if e != nil {
		return e
	}
	_, e = fi.WriteString(cmd)
	return e
}
