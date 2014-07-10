package gpio

import (
	"strconv"

	"github.com/hybridgroup/gobot"
)

type DirectPinDriver struct {
	gobot.Driver
}

func NewDirectPinDriver(a DirectPin, name string, pin string) *DirectPinDriver {
	d := &DirectPinDriver{
		Driver: *gobot.NewDriver(
			name,
			"DirectPinDriver",
			a.(gobot.AdaptorInterface),
			pin,
		),
	}

	d.AddCommand("DigitalRead", func(params map[string]interface{}) interface{} {
		return d.DigitalRead()
	})
	d.AddCommand("DigitalWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.DigitalWrite(byte(level))
		return nil
	})
	d.AddCommand("AnalogRead", func(params map[string]interface{}) interface{} {
		return d.AnalogRead()
	})
	d.AddCommand("AnalogWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.AnalogWrite(byte(level))
		return nil
	})
	d.AddCommand("PwmWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.PwmWrite(byte(level))
		return nil
	})
	d.AddCommand("ServoWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.ServoWrite(byte(level))
		return nil
	})

	return d
}

func (d *DirectPinDriver) adaptor() DirectPin {
	return d.Adaptor().(DirectPin)
}
func (d *DirectPinDriver) Start() bool { return true }
func (d *DirectPinDriver) Halt() bool  { return true }
func (d *DirectPinDriver) Init() bool  { return true }

func (d *DirectPinDriver) DigitalRead() int {
	return d.adaptor().DigitalRead(d.Pin())
}

func (d *DirectPinDriver) DigitalWrite(level byte) {
	d.adaptor().DigitalWrite(d.Pin(), level)
}

func (d *DirectPinDriver) AnalogRead() int {
	return d.adaptor().AnalogRead(d.Pin())
}

func (d *DirectPinDriver) AnalogWrite(level byte) {
	d.adaptor().AnalogWrite(d.Pin(), level)
}

func (d *DirectPinDriver) PwmWrite(level byte) {
	d.adaptor().PwmWrite(d.Pin(), level)
}

func (d *DirectPinDriver) ServoWrite(level byte) {
	d.adaptor().ServoWrite(d.Pin(), level)
}
