package system

import (
	"time"

	"gobot.io/x/gobot/v2"
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

const (
	digitalPinBiasDefault  = 0 // GPIO uses the hardware default
	digitalPinBiasDisable  = 1 // GPIO has pull disabled
	digitalPinBiasPullDown = 2 // GPIO has pull up enabled
	digitalPinBiasPullUp   = 3 // GPIO has pull down enabled

	// open drain and open source allows the connection of output ports with the same mode (OR logic)
	// * for open drain/collector pull up the ports with an external resistor/load
	// * for open source/emitter pull down the ports with an external resistor/load
	digitalPinDrivePushPull   = 0 // the pin will be driven actively high and	low (default)
	digitalPinDriveOpenDrain  = 1 // the pin will be driven active to low only
	digitalPinDriveOpenSource = 2 // the pin will be driven active to high only

	digitalPinEventNone          = 0 // no event will be triggered on any pin change (default)
	digitalPinEventOnFallingEdge = 1 // an event will be triggered on changes from high to low state
	digitalPinEventOnRisingEdge  = 2 // an event will be triggered on changes from low to high state
	digitalPinEventOnBothEdges   = 3 // an event will be triggered on all changes
)

const (
	// DigitalPinEventRisingEdge indicates an inactive to active event.
	DigitalPinEventRisingEdge = "rising edge"
	// DigitalPinEventFallingEdge indicates an active to inactive event.
	DigitalPinEventFallingEdge = "falling edge"
)

type digitalPinConfig struct {
	label            string
	direction        string
	outInitialState  int
	activeLow        bool
	bias             int
	drive            int
	debouncePeriod   time.Duration
	edge             int
	edgeEventHandler func(lineOffset int, timestamp time.Duration, detectedEdge string, seqno uint32, lseqno uint32)
	pollInterval     time.Duration
	pollQuitChan     chan struct{}
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

// WithPinPullDown initializes the pin to be pulled down (high impedance to GND, applies on input and output).
// This is working since Kernel 5.5.
func WithPinPullDown() func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetBias(digitalPinBiasPullDown) }
}

// WithPinPullUp initializes the pin to be pulled up (high impedance to VDD, applies on input and output).
// This is working since Kernel 5.5.
func WithPinPullUp() func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetBias(digitalPinBiasPullUp) }
}

// WithPinOpenDrain initializes the output pin to be driven with open drain/collector.
func WithPinOpenDrain() func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetDrive(digitalPinDriveOpenDrain) }
}

// WithPinOpenSource initializes the output pin to be driven with open source/emitter.
func WithPinOpenSource() func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetDrive(digitalPinDriveOpenSource) }
}

// WithPinDebounce initializes the input pin to be debounced.
func WithPinDebounce(period time.Duration) func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool { return d.SetDebounce(period) }
}

// WithPinEventOnFallingEdge initializes the input pin for edge detection and call the event handler on falling edges.
func WithPinEventOnFallingEdge(handler func(lineOffset int, timestamp time.Duration, detectedEdge string, seqno uint32,
	lseqno uint32),
) func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool {
		return d.SetEventHandlerForEdge(handler, digitalPinEventOnFallingEdge)
	}
}

// WithPinEventOnRisingEdge initializes the input pin for edge detection and call the event handler on rising edges.
func WithPinEventOnRisingEdge(handler func(lineOffset int, timestamp time.Duration, detectedEdge string, seqno uint32,
	lseqno uint32),
) func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool {
		return d.SetEventHandlerForEdge(handler, digitalPinEventOnRisingEdge)
	}
}

// WithPinEventOnBothEdges initializes the input pin for edge detection and call the event handler on all edges.
func WithPinEventOnBothEdges(handler func(lineOffset int, timestamp time.Duration, detectedEdge string, seqno uint32,
	lseqno uint32),
) func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool {
		return d.SetEventHandlerForEdge(handler, digitalPinEventOnBothEdges)
	}
}

// WithPinPollForEdgeDetection initializes a discrete input pin polling function to use for edge detection.
func WithPinPollForEdgeDetection(
	pollInterval time.Duration,
	pollQuitChan chan struct{},
) func(gobot.DigitalPinOptioner) bool {
	return func(d gobot.DigitalPinOptioner) bool {
		return d.SetPollForEdgeDetection(pollInterval, pollQuitChan)
	}
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

// SetBias sets the pin bias (applies on input and output) for next reconfigure. The function
// is intended to use by WithPinPullUp() and WithPinPullDown().
func (d *digitalPinConfig) SetBias(bias int) bool {
	if d.bias == bias {
		return false
	}
	d.bias = bias
	return true
}

// SetDrive sets the pin drive mode (applies on output only) for next reconfigure. The function
// is intended to use by WithPinOpenDrain(), WithPinOpenSource() and WithPinPushPull().
func (d *digitalPinConfig) SetDrive(drive int) bool {
	if d.drive == drive {
		return false
	}
	d.drive = drive
	return true
}

// SetDebounce sets the input pin with the given debounce period for next reconfigure. The function
// is intended to use by WithPinDebounce().
func (d *digitalPinConfig) SetDebounce(period time.Duration) bool {
	if d.debouncePeriod == period {
		return false
	}
	d.debouncePeriod = period
	return true
}

// SetEventHandlerForEdge sets the input pin to edge detection to call the event handler on specified edge. The
// function is intended to use by WithPinEventOnFallingEdge(), WithPinEventOnRisingEdge() and WithPinEventOnBothEdges().
func (d *digitalPinConfig) SetEventHandlerForEdge(
	handler func(int, time.Duration, string, uint32, uint32),
	edge int,
) bool {
	if d.edge == edge {
		return false
	}
	d.edge = edge
	d.edgeEventHandler = handler
	return true
}

// SetPollForEdgeDetection use a discrete input polling method to detect edges. A poll interval of zero or smaller
// will deactivate this function. Please note: Using this feature is CPU consuming and less accurate than using cdev
// event handler (go-gpiocdev package) and should be done only if the former is not implemented or not working for
// the adaptor. E.g. sysfs driver in gobot has not implemented edge detection yet. The function is only useful
// together with SetEventHandlerForEdge() and its corresponding With*() functions.
// The function is intended to use by WithPinPollForEdgeDetection().
//
//nolint:nonamedreturns // useful here
func (d *digitalPinConfig) SetPollForEdgeDetection(
	pollInterval time.Duration,
	pollQuitChan chan struct{},
) (changed bool) {
	if d.pollInterval == pollInterval {
		return false
	}
	d.pollInterval = pollInterval
	d.pollQuitChan = pollQuitChan
	return true
}
