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

// WithLabel use a pin label, which will replace the default label "gobotio#".
func WithLabel(label string) func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetLabel(label) }
}

// WithDirectionOutput initializes the pin as output instead of the default "input".
func WithDirectionOutput(initial int) func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetDirectionOutput(initial) }
}

// WithDirectionInput initializes the pin as input.
func WithDirectionInput() func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetDirectionInput() }
}

// SetLabel sets the label to use for next reconfigure. The function is intended to use by WithLabel().
func (d *digitalPinConfig) SetLabel(label string) bool {
	if d.label == label {
		return false
	}
	d.label = label
	return true
}

// SetDirectionOutput sets the direction to output for next reconfigure. The function is intended to use by WithLabel().
func (d *digitalPinConfig) SetDirectionOutput(initial int) bool {
	if d.direction == OUT {
		// in this case also the initial value will not be written
		return false
	}
	d.direction = OUT
	d.outInitialState = initial
	return true
}

// SetDirectionInput sets the direction to input for next reconfigure. The function is intended to use by WithLabel().
func (d *digitalPinConfig) SetDirectionInput() bool {
	if d.direction == IN {
		return false
	}
	d.direction = IN
	return true
}
