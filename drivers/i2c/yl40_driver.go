package i2c

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
)

// All address pins are connected to ground.
const yl40DefaultAddress = 0x48

const yl40Debug = false

type YL40Pin string

const (
	// brightness sensor, high brightness - low raw value, scaled to 0..1000 (high brightness - high value)
	YL40Bri YL40Pin = "brightness"
	// temperature sensor, high temperature - low raw value, scaled to °C
	YL40Temp YL40Pin = "temperature"
	// wired to AOUT, scaled to voltage 3.3V
	YL40AIN2 YL40Pin = "analog input AIN2"
	// adjustable resistor, turn clockwise will lower the raw value, scaled to -100..+100% (clockwise)
	YL40Poti YL40Pin = "potentiometer"
	YL40AOUT YL40Pin = "analog output"
)

const (
	// the LED light is visible above ~1.7V
	yl40LedDefaultVal = 1.7
	// default refresh rate, set to zero (cyclic reading deactivated)
	yl40DefaultRefresh = 0
)

type yl40Sensor struct {
	interval time.Duration
	scaler   func(input int) (value float64)
}

type yl40Config struct {
	sensors    map[YL40Pin]*yl40Sensor
	aOutScaler func(input float64) (value int)
}

var yl40Pins = map[YL40Pin]string{
	YL40Bri:  "s.0",
	YL40Temp: "s.1",
	YL40AIN2: "s.2",
	YL40Poti: "s.3",
	YL40AOUT: "aOut.0",
}

// YL40Driver is a Gobot i2c bus driver for the YL-40 module with light dependent resistor (LDR), thermistor (NTC)
// and an potentiometer, one additional analog input and one analog output with an connected LED.
// The module is based on PCF8591 with 4xADC, 1xDAC. For detailed documentation refer to PCF8591Driver.
//
// All values are linear scaled to 3.3V by default. This can be changed, see example "tinkerboard_yl40.go".
//
// This driver was tested with Tinkerboard and this board with temperature & brightness sensor:
// https://www.makershop.de/download/YL_40_yl40.pdf
//
type YL40Driver struct {
	*PCF8591Driver
	conf yl40Config

	aBri  *aio.AnalogSensorDriver
	aTemp *aio.TemperatureSensorDriver
	aAIN2 *aio.AnalogSensorDriver
	aPoti *aio.AnalogSensorDriver
	aOut  *aio.AnalogActuatorDriver
}

// NewYL40Driver creates a new driver with specified i2c interface
// Params:
//    conn Connector - the Adaptor to use with this Driver
//
// Optional parameters:
//  refer to PCF8591Driver for i2c specific options
// 	refer to TemperatureSensorDriver for temperature sensor specific options
// 	refer to AnalogSensorDriver for analog input specific options
//  refer to AnalogActuatorDriver for analog output specific options
//
func NewYL40Driver(a Connector, options ...func(Config)) *YL40Driver {
	options = append(options, WithAddress(yl40DefaultAddress))
	pcf := NewPCF8591Driver(a, options...)

	ntc := aio.TemperatureSensorNtcConf{TC0: 25, R0: 10000.0, B: 3950} //Ohm, R25=10k, B=3950
	defTempScaler := aio.TemperatureSensorNtcScaler(255, 1000, true, ntc)

	defConf := yl40Config{
		sensors: map[YL40Pin]*yl40Sensor{
			YL40Bri: {
				interval: yl40DefaultRefresh,
				scaler:   aio.AnalogSensorLinearScaler(0, 255, 1000, 0),
			},
			YL40Temp: {
				interval: yl40DefaultRefresh,
				scaler:   defTempScaler,
			},
			YL40AIN2: {
				interval: yl40DefaultRefresh,
				scaler:   aio.AnalogSensorLinearScaler(0, 255, 0, 3.3),
			},
			YL40Poti: {
				interval: yl40DefaultRefresh,
				scaler:   aio.AnalogSensorLinearScaler(0, 255, 100, -100),
			},
		},
		aOutScaler: aio.AnalogActuatorLinearScaler(0, 3.3, 0, 255),
	}

	y := &YL40Driver{
		PCF8591Driver: pcf,
		conf:          defConf,
	}

	y.SetName(gobot.DefaultName("YL-40"))

	for _, option := range options {
		option(y)
	}

	// initialize analog drivers
	y.aBri = aio.NewAnalogSensorDriver(pcf, yl40Pins[YL40Bri], y.conf.sensors[YL40Bri].interval)
	y.aTemp = aio.NewTemperatureSensorDriver(pcf, yl40Pins[YL40Temp], y.conf.sensors[YL40Temp].interval)
	y.aAIN2 = aio.NewAnalogSensorDriver(pcf, yl40Pins[YL40AIN2], y.conf.sensors[YL40AIN2].interval)
	y.aPoti = aio.NewAnalogSensorDriver(pcf, yl40Pins[YL40Poti], y.conf.sensors[YL40Poti].interval)
	y.aOut = aio.NewAnalogActuatorDriver(pcf, yl40Pins[YL40AOUT])

	// set input scalers
	y.aBri.SetScaler(y.conf.sensors[YL40Bri].scaler)
	y.aTemp.SetScaler(y.conf.sensors[YL40Temp].scaler)
	y.aAIN2.SetScaler(y.conf.sensors[YL40AIN2].scaler)
	y.aPoti.SetScaler(y.conf.sensors[YL40Poti].scaler)

	// set output scaler
	y.aOut.SetScaler(y.conf.aOutScaler)

	return y
}

// WithYL40Interval option sets the interval for refresh of given pin in YL40 driver
func WithYL40Interval(pin YL40Pin, val time.Duration) func(Config) {
	return func(c Config) {
		y, ok := c.(*YL40Driver)
		if ok {
			if sensor, ok := y.conf.sensors[pin]; ok {
				sensor.interval = val
			}
		} else if yl40Debug {
			log.Printf("trying to set interval for '%s' refresh for non-YL40Driver %v", pin, c)
		}
	}
}

// WithYL40InputScaler option sets the input scaler of given input pin in YL40 driver
func WithYL40InputScaler(pin YL40Pin, scaler func(input int) (value float64)) func(Config) {
	return func(c Config) {
		y, ok := c.(*YL40Driver)
		if ok {
			if sensor, ok := y.conf.sensors[pin]; ok {
				sensor.scaler = scaler
			}
		} else if yl40Debug {
			log.Printf("trying to set input scaler for '%s' for non-YL40Driver %v", pin, c)
		}
	}
}

// WithYL40OutputScaler option sets the output scaler in YL40 driver
func WithYL40OutputScaler(scaler func(input float64) (value int)) func(Config) {
	return func(c Config) {
		y, ok := c.(*YL40Driver)
		if ok {
			y.conf.aOutScaler = scaler
		} else if yl40Debug {
			log.Printf("trying to set output scaler for '%s' for non-YL40Driver %v", YL40AOUT, c)
		}
	}
}

// Start initializes the driver
func (y *YL40Driver) Start() (err error) {
	// must be the first one
	if err := y.PCF8591Driver.Start(); err != nil {
		return err
	}
	if err := y.aBri.Start(); err != nil {
		return err
	}
	if err := y.aTemp.Start(); err != nil {
		return err
	}
	if err := y.aAIN2.Start(); err != nil {
		return err
	}
	if err := y.aPoti.Start(); err != nil {
		return err
	}
	if err := y.aOut.Start(); err != nil {
		return err
	}
	return y.Write(yl40LedDefaultVal)
}

// Halt stops the driver
func (y *YL40Driver) Halt() (err error) {
	// we try halt on each device, not stopping on the first error
	var errors []string
	if err := y.aBri.Halt(); err != nil {
		errors = append(errors, err.Error())
	}
	if err := y.aTemp.Halt(); err != nil {
		errors = append(errors, err.Error())
	}
	if err := y.aAIN2.Halt(); err != nil {
		errors = append(errors, err.Error())
	}
	if err := y.aPoti.Halt(); err != nil {
		errors = append(errors, err.Error())
	}
	if err := y.aOut.Halt(); err != nil {
		errors = append(errors, err.Error())
	}
	// must be the last one
	if err := y.PCF8591Driver.Halt(); err != nil {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		return fmt.Errorf("Halt the driver %s", strings.Join(errors, ", "))
	}
	return nil
}

// Read returns the current reading from the given pin of the driver
// For the analog output pin the last written value is returned
func (y *YL40Driver) Read(pin YL40Pin) (val float64, err error) {
	switch pin {
	case YL40Bri:
		return y.aBri.ReadValue()
	case YL40Temp:
		return y.aTemp.ReadValue()
	case YL40AIN2:
		return y.aAIN2.ReadValue()
	case YL40Poti:
		return y.aPoti.ReadValue()
	case YL40AOUT:
		return y.aOut.Value(), nil
	default:
		return 0, fmt.Errorf("Analog reading from pin '%s' not supported", pin)
	}
}

// ReadBrightness returns the current reading from brightness pin of the driver
func (y *YL40Driver) ReadBrightness() (val float64, err error) {
	return y.Read(YL40Bri)
}

// ReadTemperature returns the current reading from temperature pin of the driver
func (y *YL40Driver) ReadTemperature() (val float64, err error) {
	return y.Read(YL40Temp)
}

// ReadAIN2 returns the current reading from analog input pin 2 pin of the driver
func (y *YL40Driver) ReadAIN2() (val float64, err error) {
	return y.Read(YL40AIN2)
}

// ReadPotentiometer returns the current reading from potentiometer pin of the driver
func (y *YL40Driver) ReadPotentiometer() (val float64, err error) {
	return y.Read(YL40Poti)
}

// Value returns the last read or written value from the given pin of the driver
func (y *YL40Driver) Value(pin YL40Pin) (val float64, err error) {
	switch pin {
	case YL40Bri:
		return y.aBri.Value(), nil
	case YL40Temp:
		return y.aTemp.Value(), nil
	case YL40AIN2:
		return y.aAIN2.Value(), nil
	case YL40Poti:
		return y.aPoti.Value(), nil
	case YL40AOUT:
		return y.aOut.Value(), nil
	default:
		return 0, fmt.Errorf("Get analog value from pin '%s' not supported", pin)
	}
}

// Brightness returns the last read brightness of the driver
func (y *YL40Driver) Brightness() (val float64, err error) {
	return y.Value(YL40Bri)
}

// Temperature returns the last read temperature of the driver
func (y *YL40Driver) Temperature() (val float64, err error) {
	return y.Value(YL40Temp)
}

// AIN2 returns the last read analog input value of the driver
func (y *YL40Driver) AIN2() (val float64, err error) {
	return y.Value(YL40AIN2)
}

// Potentiometer returns the last read potentiometer value of the driver
func (y *YL40Driver) Potentiometer() (val float64, err error) {
	return y.Value(YL40Poti)
}

// AOUT returns the last written value of the driver
func (y *YL40Driver) AOUT() (val float64, err error) {
	return y.Value(YL40AOUT)
}

// Write writes the given value to the analog output
func (y *YL40Driver) Write(val float64) (err error) {
	return y.aOut.Write(val)
}
