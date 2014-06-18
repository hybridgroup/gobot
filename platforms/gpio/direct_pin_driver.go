package gpio

import (
	"github.com/hybridgroup/gobot"
	"strconv"
)

type DirectPinDriver struct {
	gobot.Driver
}

func NewDirectPinDriver(a DirectPin, name string, pin string) *DirectPinDriver {
	d := &DirectPinDriver{
		Driver: gobot.Driver{
			Name:     name,
			Pin:      pin,
			Commands: make(map[string]func(map[string]interface{}) interface{}),
			Adaptor:  a.(gobot.AdaptorInterface),
		},
	}

	d.Driver.AddCommand("DigitalRead", func(params map[string]interface{}) interface{} {
		return d.DigitalRead()
	})
	d.Driver.AddCommand("DigitalWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.DigitalWrite(byte(level))
		return nil
	})
	d.Driver.AddCommand("AnalogRead", func(params map[string]interface{}) interface{} {
		return d.AnalogRead()
	})
	d.Driver.AddCommand("AnalogWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.AnalogWrite(byte(level))
		return nil
	})
	d.Driver.AddCommand("PwmWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.PwmWrite(byte(level))
		return nil
	})
	d.Driver.AddCommand("ServoWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.ServoWrite(byte(level))
		return nil
	})

	return d
}

func (d *DirectPinDriver) adaptor() DirectPin {
	return d.Driver.Adaptor.(DirectPin)
}
func (d *DirectPinDriver) Start() bool { return true }
func (d *DirectPinDriver) Halt() bool  { return true }
func (d *DirectPinDriver) Init() bool  { return true }

func (d *DirectPinDriver) DigitalRead() int {
	return d.adaptor().DigitalRead(d.Pin)
}

func (d *DirectPinDriver) DigitalWrite(level byte) {
	d.adaptor().DigitalWrite(d.Pin, level)
}

func (d *DirectPinDriver) AnalogRead() int {
	return d.adaptor().AnalogRead(d.Pin)
}

func (d *DirectPinDriver) AnalogWrite(level byte) {
	d.adaptor().AnalogWrite(d.Pin, level)
}

func (d *DirectPinDriver) PwmWrite(level byte) {
	d.adaptor().PwmWrite(d.Pin, level)
}

func (d *DirectPinDriver) ServoWrite(level byte) {
	d.adaptor().ServoWrite(d.Pin, level)
}
