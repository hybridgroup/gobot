package arietta

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/i2c"
)

type AriettaAdaptor struct {
	gobot.Adaptor

	digitalPins map[string]*sysfsDigitalPin
	pwms        map[string]*Pwm
	i2cDevice   *i2cDevice
}

var digitalPins = []sysfsDigitalPinDesc{
	// Taken from /sys/kernel/debug/pinctrl/pins
	{name: "PA0", hwnum: 0, label: "pioA0"},
	{name: "PA1", hwnum: 1, label: "pioA1"},
	{name: "PA2", hwnum: 2, label: "pioA2"},
	{name: "PA3", hwnum: 3, label: "pioA3"},
	{name: "PA4", hwnum: 4, label: "pioA4"},
	{name: "PA5", hwnum: 5, label: "pioA5"},
	{name: "PA6", hwnum: 6, label: "pioA6"},
	{name: "PA7", hwnum: 7, label: "pioA7"},
	{name: "PA8", hwnum: 8, label: "pioA8"},
	{name: "PA9", hwnum: 9, label: "pioA9"},
	{name: "PA10", hwnum: 10, label: "pioA10"},
	{name: "PA11", hwnum: 11, label: "pioA11"},
	{name: "PA12", hwnum: 12, label: "pioA12"},
	{name: "PA13", hwnum: 13, label: "pioA13"},
	{name: "PA14", hwnum: 14, label: "pioA14"},
	{name: "PA15", hwnum: 15, label: "pioA15"},
	{name: "PA16", hwnum: 16, label: "pioA16"},
	{name: "PA17", hwnum: 17, label: "pioA17"},
	{name: "PA18", hwnum: 18, label: "pioA18"},
	{name: "PA19", hwnum: 19, label: "pioA19"},
	{name: "PA20", hwnum: 20, label: "pioA20"},
	{name: "PA21", hwnum: 21, label: "pioA21"},
	{name: "PA22", hwnum: 22, label: "pioA22"},
	{name: "PA23", hwnum: 23, label: "pioA23"},
	{name: "PA24", hwnum: 24, label: "pioA24"},
	{name: "PA25", hwnum: 25, label: "pioA25"},
	{name: "PA26", hwnum: 26, label: "pioA26"},
	{name: "PA27", hwnum: 27, label: "pioA27"},
	{name: "PA28", hwnum: 28, label: "pioA28"},
	{name: "PA29", hwnum: 29, label: "pioA29"},
	{name: "PA30", hwnum: 30, label: "pioA30"},
	{name: "PA31", hwnum: 31, label: "pioA31"},
	{name: "PB0", hwnum: 32, label: "pioB0"},
	{name: "PB1", hwnum: 33, label: "pioB1"},
	{name: "PB2", hwnum: 34, label: "pioB2"},
	{name: "PB3", hwnum: 35, label: "pioB3"},
	{name: "PB4", hwnum: 36, label: "pioB4"},
	{name: "PB5", hwnum: 37, label: "pioB5"},
	{name: "PB6", hwnum: 38, label: "pioB6"},
	{name: "PB7", hwnum: 39, label: "pioB7"},
	{name: "PB8", hwnum: 40, label: "pioB8"},
	{name: "PB9", hwnum: 41, label: "pioB9"},
	{name: "PB10", hwnum: 42, label: "pioB10"},
	{name: "PB11", hwnum: 43, label: "pioB11"},
	{name: "PB12", hwnum: 44, label: "pioB12"},
	{name: "PB13", hwnum: 45, label: "pioB13"},
	{name: "PB14", hwnum: 46, label: "pioB14"},
	{name: "PB15", hwnum: 47, label: "pioB15"},
	{name: "PB16", hwnum: 48, label: "pioB16"},
	{name: "PB17", hwnum: 49, label: "pioB17"},
	{name: "PB18", hwnum: 50, label: "pioB18"},
	{name: "PB19", hwnum: 51, label: "pioB19"},
	{name: "PB20", hwnum: 52, label: "pioB20"},
	{name: "PB21", hwnum: 53, label: "pioB21"},
	{name: "PB22", hwnum: 54, label: "pioB22"},
	{name: "PB23", hwnum: 55, label: "pioB23"},
	{name: "PB24", hwnum: 56, label: "pioB24"},
	{name: "PB25", hwnum: 57, label: "pioB25"},
	{name: "PB26", hwnum: 58, label: "pioB26"},
	{name: "PB27", hwnum: 59, label: "pioB27"},
	{name: "PB28", hwnum: 60, label: "pioB28"},
	{name: "PB29", hwnum: 61, label: "pioB29"},
	{name: "PB30", hwnum: 62, label: "pioB30"},
	{name: "PB31", hwnum: 63, label: "pioB31"},
	{name: "PC0", hwnum: 64, label: "pioC0"},
	{name: "PC1", hwnum: 65, label: "pioC1"},
	{name: "PC2", hwnum: 66, label: "pioC2"},
	{name: "PC3", hwnum: 67, label: "pioC3"},
	{name: "PC4", hwnum: 68, label: "pioC4"},
	{name: "PC5", hwnum: 69, label: "pioC5"},
	{name: "PC6", hwnum: 70, label: "pioC6"},
	{name: "PC7", hwnum: 71, label: "pioC7"},
	{name: "PC8", hwnum: 72, label: "pioC8"},
	{name: "PC9", hwnum: 73, label: "pioC9"},
	{name: "PC10", hwnum: 74, label: "pioC10"},
	{name: "PC11", hwnum: 75, label: "pioC11"},
	{name: "PC12", hwnum: 76, label: "pioC12"},
	{name: "PC13", hwnum: 77, label: "pioC13"},
	{name: "PC14", hwnum: 78, label: "pioC14"},
	{name: "PC15", hwnum: 79, label: "pioC15"},
	{name: "PC16", hwnum: 80, label: "pioC16"},
	{name: "PC17", hwnum: 81, label: "pioC17"},
	{name: "PC18", hwnum: 82, label: "pioC18"},
	{name: "PC19", hwnum: 83, label: "pioC19"},
	{name: "PC20", hwnum: 84, label: "pioC20"},
	{name: "PC21", hwnum: 85, label: "pioC21"},
	{name: "PC22", hwnum: 86, label: "pioC22"},
	{name: "PC23", hwnum: 87, label: "pioC23"},
	{name: "PC24", hwnum: 88, label: "pioC24"},
	{name: "PC25", hwnum: 89, label: "pioC25"},
	{name: "PC26", hwnum: 90, label: "pioC26"},
	{name: "PC27", hwnum: 91, label: "pioC27"},
	{name: "PC28", hwnum: 92, label: "pioC28"},
	{name: "PC29", hwnum: 93, label: "pioC29"},
	{name: "PC30", hwnum: 94, label: "pioC30"},
	{name: "PC31", hwnum: 95, label: "pioC31"},
	{name: "PD0", hwnum: 96, label: "pioD0"},
	{name: "PD1", hwnum: 97, label: "pioD1"},
	{name: "PD2", hwnum: 98, label: "pioD2"},
	{name: "PD3", hwnum: 99, label: "pioD3"},
	{name: "PD4", hwnum: 100, label: "pioD4"},
	{name: "PD5", hwnum: 101, label: "pioD5"},
	{name: "PD6", hwnum: 102, label: "pioD6"},
	{name: "PD7", hwnum: 103, label: "pioD7"},
	{name: "PD8", hwnum: 104, label: "pioD8"},
	{name: "PD9", hwnum: 105, label: "pioD9"},
	{name: "PD10", hwnum: 106, label: "pioD10"},
	{name: "PD11", hwnum: 107, label: "pioD11"},
	{name: "PD12", hwnum: 108, label: "pioD12"},
	{name: "PD13", hwnum: 109, label: "pioD13"},
	{name: "PD14", hwnum: 110, label: "pioD14"},
	{name: "PD15", hwnum: 111, label: "pioD15"},
	{name: "PD16", hwnum: 112, label: "pioD16"},
	{name: "PD17", hwnum: 113, label: "pioD17"},
	{name: "PD18", hwnum: 114, label: "pioD18"},
	{name: "PD19", hwnum: 115, label: "pioD19"},
	{name: "PD20", hwnum: 116, label: "pioD20"},
	{name: "PD21", hwnum: 117, label: "pioD21"},
	{name: "PD22", hwnum: 118, label: "pioD22"},
	{name: "PD23", hwnum: 119, label: "pioD23"},
	{name: "PD24", hwnum: 120, label: "pioD24"},
	{name: "PD25", hwnum: 121, label: "pioD25"},
	{name: "PD26", hwnum: 122, label: "pioD26"},
	{name: "PD27", hwnum: 123, label: "pioD27"},
	{name: "PD28", hwnum: 124, label: "pioD28"},
	{name: "PD29", hwnum: 125, label: "pioD29"},
	{name: "PD30", hwnum: 126, label: "pioD30"},
	{name: "PD31", hwnum: 127, label: "pioD31"},
}

var pwms = []PwmDesc{
	{name: "PB11", chip: 0, hwnum: 0},
	{name: "PB12", chip: 0, hwnum: 1},
	{name: "PB13", chip: 0, hwnum: 2},
	{name: "PB14", chip: 0, hwnum: 3},
}

// TODO(michaelh): make configurable.
const i2cLocation = "/dev/i2c-1"

var _ i2c.I2cInterface = (*AriettaAdaptor)(nil)
var _ gpio.DirectPin = (*AriettaAdaptor)(nil)

func NewAriettaAdaptor(name string) *AriettaAdaptor {
	return &AriettaAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"AriettaAdaptor",
		),
		digitalPins: make(map[string]*sysfsDigitalPin),
		pwms:        make(map[string]*Pwm),
	}
}

func (b *AriettaAdaptor) Connect() bool {
	return true
}

func (b *AriettaAdaptor) Finalize() bool {
	for _, pin := range b.digitalPins {
		pin.Finalize()
	}
	for _, pwm := range b.pwms {
		pwm.Finalize()
	}
	if b.i2cDevice != nil {
		b.i2cDevice.finalize()
	}
	return true
}

func (b *AriettaAdaptor) Reconnect() bool  { return true }
func (b *AriettaAdaptor) Disconnect() bool { return true }

func (b *AriettaAdaptor) PwmWrite(pin string, val byte) {
	b.findPwm(pin).PwmWrite(val)
}

func (b *AriettaAdaptor) DigitalRead(pin string) int {
	return b.findDigitalPin(pin).DigitalRead()
}

func (b *AriettaAdaptor) DigitalWrite(pin string, val byte) {
	b.findDigitalPin(pin).DigitalWrite(val)
}

// TODO(michaelh): implement.  Stubbed out so the adapter implements
// gpio.DigitalPin.
func (b *AriettaAdaptor) AnalogRead(pin string) int {
	panic("Not implemented.")
}

func (b *AriettaAdaptor) AnalogWrite(pin string, val byte) {
	panic("Not implemented.")
}

func (b *AriettaAdaptor) InitServo() {
	panic("Not implemented.")
}

func (b *AriettaAdaptor) ServoWrite(pin string, val byte) {
	panic("Not implemented.")
}

func (b *AriettaAdaptor) I2cStart(address byte) {
	if b.i2cDevice == nil {
		b.i2cDevice = newI2cDevice(i2cLocation, address)
	}
	b.i2cDevice.start()
}

func (b *AriettaAdaptor) I2cWrite(data []byte) {
	b.i2cDevice.write(data)
}

func (b *AriettaAdaptor) I2cRead(size uint) []byte {
	return b.i2cDevice.read(size)
}

// TODO(michaelh): make multithread safe.  Pins are created lazily.
func (b *AriettaAdaptor) findDigitalPin(pin string) *sysfsDigitalPin {
	d, ok := b.digitalPins[pin]

	if ok {
		return d
	}

	for _, m := range digitalPins {
		if m.name == pin {
			d := newSysfsDigitalPin(&m)
			b.digitalPins[pin] = d
			return d
		}
	}

	panic(fmt.Sprintf("No such digital pin %v.", pin))
}

func (b *AriettaAdaptor) findPwm(pin string) *Pwm {
	d, ok := b.pwms[pin]

	if ok {
		return d
	}

	for _, m := range pwms {
		if m.name == pin {
			d := newPwm(&m)
			b.pwms[pin] = d
			return d
		}
	}

	panic(fmt.Sprintf("No such PWM channel %v.", pin))
}
