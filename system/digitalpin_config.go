package system

import (
	"gobot.io/x/gobot"
)

const (
	// IN gpio direction
	IN = "in"
	// OUT gpio direction
	OUT = "out"
	// HIGH gpio level
	HIGH = 1
	// LOW gpio level
	LOW = 0
)

type digitalPinConfig struct {
	label           string
	direction       string
	outInitialState int
	activeLow       bool
}

func newDigitalPinConfig(label string, options ...func(gobot.DigitalPinOptioner) bool) *digitalPinConfig {
	cfg := &digitalPinConfig{
		label:     label,
		direction: IN,
	}
	for _, option := range options {
		option(cfg)
	}
	return cfg
}

// WithPinLabel use a pin label, which will replace the default label "gobotio#".
func WithPinLabel(label string) func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetLabel(label) }
}

// WithPinDirectionOutput initializes the pin as output instead of the default "input".
func WithPinDirectionOutput(initial int) func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetDirectionOutput(initial) }
}

// WithPinDirectionInput initializes the pin as input.
func WithPinDirectionInput() func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetDirectionInput() }
}

// WithPinActiveLow initializes the pin with inverse reaction (applies on input and output).
func WithPinActiveLow() func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetActiveLow() }
}

// SetLabel sets the label to use for next reconfigure. The function is intended to use by WithPinLabel().
func (d *digitalPinConfig) SetLabel(label string) bool {
	if d.label == label {
		return false
	}
	d.label = label
	return true
}

// SetDirectionOutput sets the direction to output for next reconfigure. The function is intended to use
// by WithPinDirectionOutput().
func (d *digitalPinConfig) SetDirectionOutput(initial int) bool {
	if d.direction == OUT {
		// in this case also the initial value will not be written
		return false
	}
	d.direction = OUT
	d.outInitialState = initial
	return true
}

// SetDirectionInput sets the direction to input for next reconfigure. The function is intended to use
// by WithPinDirectionInput().
func (d *digitalPinConfig) SetDirectionInput() bool {
	if d.direction == IN {
		return false
	}
	d.direction = IN
	return true
}

// SetActiveLow sets the pin with inverse reaction (applies on input and output) for next reconfigure. The function
// is intended to use by WithPinActiveLow().
func (d *digitalPinConfig) SetActiveLow() bool {
	if d.activeLow {
		return false
	}
	d.activeLow = true
	return true
}
