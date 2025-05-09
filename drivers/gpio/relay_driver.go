package gpio

import (
	"fmt"

	"gobot.io/x/gobot/v2"
)

// relayOptionApplier needs to be implemented by each configurable option type
type relayOptionApplier interface {
	apply(cfg *relayConfiguration)
}

// relayConfiguration contains all changeable attributes of the driver.
type relayConfiguration struct {
	inverted bool
}

// relayInvertedOption is the type for applying inverted behavior to the configuration
type relayInvertedOption bool

// RelayDriver represents a digital relay
type RelayDriver struct {
	*driver
	relayCfg *relayConfiguration
	high     bool
}

// NewRelayDriver return a new RelayDriver given a DigitalWriter and pin.
//
// Supported options:
//
//	"WithName"
//	"WithRelayInverted"
//
// Adds the following API Commands:
//
//	"Toggle" - See RelayDriver.Toggle
//	"On" - See RelayDriver.On
//	"Off" - See RelayDriver.Off
func NewRelayDriver(a DigitalWriter, pin string, opts ...interface{}) *RelayDriver {
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &RelayDriver{
		driver:   newDriver(a.(gobot.Connection), "Relay", withPin(pin)),
		relayCfg: &relayConfiguration{},
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case relayOptionApplier:
			o.apply(d.relayCfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	d.AddCommand("Toggle", func(_ map[string]interface{}) interface{} {
		return d.Toggle()
	})

	d.AddCommand("On", func(_ map[string]interface{}) interface{} {
		return d.On()
	})

	d.AddCommand("Off", func(_ map[string]interface{}) interface{} {
		return d.Off()
	})

	return d
}

// WithRelayInverted change the relay action to inverted.
func WithRelayInverted() relayOptionApplier {
	return relayInvertedOption(true)
}

// State return true if the relay is On and false if the relay is Off
func (d *RelayDriver) State() bool {
	if d.relayCfg.inverted {
		return !d.high
	}
	return d.high
}

// On sets the relay to a high state.
func (d *RelayDriver) On() error {
	newValue := byte(1)
	if d.relayCfg.inverted {
		newValue = 0
	}
	if err := d.digitalWrite(d.driverCfg.pin, newValue); err != nil {
		return err
	}

	d.high = !d.relayCfg.inverted

	return nil
}

// Off sets the relay to a low state.
func (d *RelayDriver) Off() error {
	newValue := byte(0)
	if d.relayCfg.inverted {
		newValue = 1
	}
	if err := d.digitalWrite(d.driverCfg.pin, newValue); err != nil {
		return err
	}

	d.high = d.relayCfg.inverted

	return nil
}

// Toggle sets the relay to the opposite of it's current state
func (d *RelayDriver) Toggle() error {
	if d.State() {
		return d.Off()
	}

	return d.On()
}

// IsInverted returns true if the relay acts inverted
func (d *RelayDriver) IsInverted() bool {
	return d.relayCfg.inverted
}

func (o relayInvertedOption) String() string {
	return "relay acts inverted option"
}

func (o relayInvertedOption) apply(cfg *relayConfiguration) {
	cfg.inverted = bool(o)
}
