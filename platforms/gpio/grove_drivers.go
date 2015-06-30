package gpio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*GroveTouchDriver)(nil)
var _ gobot.Driver = (*GroveSoundSensorDriver)(nil)
var _ gobot.Driver = (*GroveButtonDriver)(nil)
var _ gobot.Driver = (*GroveBuzzerDriver)(nil)
var _ gobot.Driver = (*GroveLightSensorDriver)(nil)
var _ gobot.Driver = (*GroveLedDriver)(nil)
var _ gobot.Driver = (*GroveRotaryDriver)(nil)
var _ gobot.Driver = (*GroveRelayDriver)(nil)

type GroveRelayDriver struct {
	*RelayDriver
}

type GroveRotaryDriver struct {
	*AnalogSensorDriver
}

type GroveLedDriver struct {
	*LedDriver
}

type GroveLightSensorDriver struct {
	*AnalogSensorDriver
}

type GroveBuzzerDriver struct {
	*BuzzerDriver
}

type GroveButtonDriver struct {
	*ButtonDriver
}

type GroveSoundSensorDriver struct {
	*AnalogSensorDriver
}

type GroveTouchDriver struct {
	*ButtonDriver
}

func NewGroveTouchDriver(a DigitalReader, name string, pin string, v ...time.Duration) *GroveButtonDriver {
	return &GroveTouchDriver{
		ButtonDriver: NewButtonDriver(a, name, pin, v...),
	}
}

func NewGroveButtonDriver(a DigitalReader, name string, pin string, v ...time.Duration) *GroveButtonDriver {
	return &GroveButtonDriver{
		ButtonDriver: NewButtonDriver(a, name, pin, v...),
	}
}

func NewGroveBuzzerDriver(a DigitalWriter, name string, pin string) *GroveBuzzerDriver {
	return &GroveBuzzerDriver{
		BuzzerDriver: NewBuzzerDriver(a, name, pin),
	}
}

func NewGroveLedDriver(a DigitalWriter, name string, pin string) *GroveLedDriver {
	return &GroveLedDriver{
		LedDriver: NewLedDriver(a, name, pin),
	}
}

func NewGroveRelayDriver(a DigitalWriter, name string, pin string) *GroveRelayDriver {
	return &GroveRelayDriver{
		RelayDriver: NewRelayDriver(a, name, pin),
	}
}

func NewGroveRotaryDriver(a AnalogReader, name string, pin string, v ...time.Duration) *GroveRotaryDriver {
	return &GroveRotaryDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, name, pin, v...),
	}
}

func NewGroveLightSensorDriver(a AnalogReader, name string, pin string, v ...time.Duration) *GroveLightSensorDriver {
	return &GroveLightSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, name, pin, v...),
	}
}

func NewGroveSoundSensorDriver(a AnalogReader, name string, pin string, v ...time.Duration) *GroveSoundSensorDriver {
	return &GroveSoundSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, name, pin, v...),
	}
}
