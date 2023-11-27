package gpio

// GroveRelayDriver represents a Relay with a Grove connector
type GroveRelayDriver struct {
	*RelayDriver
}

// NewGroveRelayDriver return a new GroveRelayDriver given a DigitalWriter and pin.
//
// Supported options:
//
//	"WithName"
//
// Adds the following API Commands:
//
//	"Toggle" - See RelayDriver.Toggle
//	"On" - See RelayDriver.On
//	"Off" - See RelayDriver.Off
//
// Deprecated: Please use [gpio.NewRelayDriver] instead. Development will be discontinued.
func NewGroveRelayDriver(a DigitalWriter, pin string) *GroveRelayDriver {
	return &GroveRelayDriver{
		RelayDriver: NewRelayDriver(a, pin),
	}
}

// GroveLedDriver represents an LED with a Grove connector
type GroveLedDriver struct {
	*LedDriver
}

// NewGroveLedDriver return a new driver for Grove Led given a DigitalWriter and pin.
//
// Supported options:
//
//	"WithName"
//
// Adds the following API Commands:
//
//	"Brightness" - See LedDriver.Brightness
//	"Toggle" - See LedDriver.Toggle
//	"On" - See LedDriver.On
//	"Off" - See LedDriver.Off
//
// Deprecated: Please use [gpio.NewLedDriver] instead. Development will be discontinued.
func NewGroveLedDriver(a DigitalWriter, pin string) *GroveLedDriver {
	return &GroveLedDriver{
		LedDriver: NewLedDriver(a, pin),
	}
}

// GroveBuzzerDriver represents a buzzer with a Grove connector
type GroveBuzzerDriver struct {
	*BuzzerDriver
}

// NewGroveBuzzerDriver return a new driver for Grove buzzer given a DigitalWriter and pin.
//
// Supported options:
//
//	"WithName"
//
// Deprecated: Please use [gpio.NewBuzzerDriver] instead. Development will be discontinued.
func NewGroveBuzzerDriver(a DigitalWriter, pin string, opts ...interface{}) *GroveBuzzerDriver {
	return &GroveBuzzerDriver{
		BuzzerDriver: NewBuzzerDriver(a, pin, opts...),
	}
}

// GroveButtonDriver represents a button sensor with a Grove connector
type GroveButtonDriver struct {
	*ButtonDriver
}

// NewGroveButtonDriver returns a new driver for Grove button with a polling interval of 10 milliseconds given
// a DigitalReader and pin.
//
// Supported options:
//
//	"WithName"
//	"WithButtonPollInterval"
//
// Deprecated: Please use [gpio.NewButtonDriver] instead. Development will be discontinued.
func NewGroveButtonDriver(a DigitalReader, pin string, opts ...interface{}) *GroveButtonDriver {
	return &GroveButtonDriver{
		ButtonDriver: NewButtonDriver(a, pin, opts...),
	}
}

// GroveTouchDriver represents a touch button sensor
// with a Grove connector
type GroveTouchDriver struct {
	*ButtonDriver
}

// NewGroveTouchDriver returns a new driver for Grove touch sensor with a polling interval of 10 milliseconds given
// a DigitalReader and pin.
//
// Supported options:
//
//	"WithName"
//	"WithButtonPollInterval"
//
// Deprecated: Please use [gpio.NewButtonDriver] instead. Development will be discontinued.
func NewGroveTouchDriver(a DigitalReader, pin string, opts ...interface{}) *GroveTouchDriver {
	return &GroveTouchDriver{
		ButtonDriver: NewButtonDriver(a, pin, opts...),
	}
}

// GroveMagneticSwitchDriver represent a magnetic switch sensor with a Grove connector
type GroveMagneticSwitchDriver struct {
	*ButtonDriver
}

// NewGroveMagneticSwitchDriver returns a new driver for Grove magnetic switch sensor with a polling interval of
// 10 milliseconds given a DigitalReader, name and pin.
//
// Supported options:
//
//	"WithName"
//	"WithButtonPollInterval"
//
// Deprecated: Please use [gpio.NewButtonDriver] instead. Development will be discontinued.
func NewGroveMagneticSwitchDriver(a DigitalReader, pin string, opts ...interface{}) *GroveMagneticSwitchDriver {
	return &GroveMagneticSwitchDriver{
		ButtonDriver: NewButtonDriver(a, pin, opts...),
	}
}
