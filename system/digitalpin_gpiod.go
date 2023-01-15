package system

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/warthog618/gpiod"
	"gobot.io/x/gobot"
)

const systemGpiodDebug = true

type cdevLine interface {
	SetValue(value int) error
	Value() (int, error)
	Close() error
}

type digitalPinGpiod struct {
	chipName string
	pin      int
	*digitalPinConfig
	line cdevLine
}

var used = map[bool]string{true: "used", false: "unused"}
var activeLow = map[bool]string{true: "low", false: "high"}
var debounced = map[bool]string{true: "debounced", false: "not debounced"}

var direction = map[gpiod.LineDirection]string{gpiod.LineDirectionUnknown: "unknown direction",
	gpiod.LineDirectionInput: "input", gpiod.LineDirectionOutput: "output"}

var drive = map[gpiod.LineDrive]string{gpiod.LineDrivePushPull: "push-pull", gpiod.LineDriveOpenDrain: "open-drain",
	gpiod.LineDriveOpenSource: "open-source"}

var bias = map[gpiod.LineBias]string{gpiod.LineBiasUnknown: "unknown", gpiod.LineBiasDisabled: "disabled",
	gpiod.LineBiasPullUp: "pull-up", gpiod.LineBiasPullDown: "pull-down"}

var edgeDetect = map[gpiod.LineEdge]string{gpiod.LineEdgeNone: "no", gpiod.LineEdgeRising: "rising",
	gpiod.LineEdgeFalling: "falling", gpiod.LineEdgeBoth: "both"}

var eventClock = map[gpiod.LineEventClock]string{gpiod.LineEventClockMonotonic: "monotonic",
	gpiod.LineEventClockRealtime: "realtime"}

// newDigitalPinGpiod returns a digital pin given the pin number, with the label "gobotio" followed by the pin number.
// The pin label can be modified optionally. The pin is handled by the character device Kernel ABI.
func newDigitalPinGpiod(chipName string, pin int, options ...func(gobot.DigitalPinOptioner) bool) *digitalPinGpiod {
	if chipName == "" {
		chipName = "gpiochip0"
	}
	cfg := newDigitalPinConfig("gobotio"+strconv.Itoa(int(pin)), options...)
	d := &digitalPinGpiod{
		chipName:         chipName,
		pin:              pin,
		digitalPinConfig: cfg,
	}
	return d
}

// ApplyOptions apply all given options to the pin immediately. Implements interface gobot.DigitalPinOptionApplier.
func (d *digitalPinGpiod) ApplyOptions(options ...func(gobot.DigitalPinOptioner) bool) error {
	anyChange := false
	for _, option := range options {
		anyChange = anyChange || option(d)
	}
	if anyChange {
		return d.reconfigure(false)
	}
	return nil
}

// DirectionBehavior gets the direction behavior when the pin is used the next time. This means its possibly not in
// this direction type at the moment. Implements the interface gobot.DigitalPinValuer, but should be rarely used.
func (d *digitalPinGpiod) DirectionBehavior() string {
	return d.direction
}

// Export sets the pin as used by this driver. Implements the interface gobot.DigitalPinner.
func (d *digitalPinGpiod) Export() error {
	err := d.reconfigure(false)
	if err != nil {
		return fmt.Errorf("gpiod.Export(): %v", err)
	}
	return nil
}

// Unexport releases the pin as input. Implements the interface gobot.DigitalPinner.
func (d *digitalPinGpiod) Unexport() error {
	var errs []string
	if d.line != nil {
		if err := d.reconfigure(true); err != nil {
			errs = append(errs, err.Error())
		}
		if err := d.line.Close(); err != nil {
			err = fmt.Errorf("gpiod.Unexport()-line.Close(): %v", err)
			errs = append(errs, err.Error())
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(errs, ","))
}

// Write writes the given value to the character device. Implements the interface gobot.DigitalPinner.
func (d *digitalPinGpiod) Write(val int) error {
	if val < 0 {
		val = 0
	}
	if val > 1 {
		val = 1
	}

	err := d.line.SetValue(val)
	if err != nil {
		return fmt.Errorf("gpiod.Write(): %v", err)
	}
	return nil
}

// Read reads the given value from character device. Implements the interface gobot.DigitalPinner.
func (d *digitalPinGpiod) Read() (int, error) {
	val, err := d.line.Value()
	if err != nil {
		return 0, fmt.Errorf("gpiod.Read(): %v", err)
	}
	return val, err
}

// ListLines is used for development purposes.
func (d *digitalPinGpiod) ListLines() error {
	c, err := gpiod.NewChip(d.chipName, gpiod.WithConsumer(d.label))
	if err != nil {
		return err
	}
	for i := 0; i < c.Lines(); i++ {
		li, err := c.LineInfo(i)
		if err != nil {
			return err
		}
		fmt.Println(fmtLine(li))
	}

	return nil
}

// List is used for development purposes.
func (d *digitalPinGpiod) List() error {
	c, err := gpiod.NewChip(d.chipName)
	if err != nil {
		return err
	}
	defer c.Close()
	l, err := c.RequestLine(d.pin)
	if err != nil && l != nil {
		l.Close()
		l = nil
	}
	li, err := l.Info()
	if err != nil {
		return err
	}
	fmt.Println(fmtLine(li))

	return nil
}

func (d *digitalPinGpiod) reconfigure(forceInput bool) error {
	// cleanup old line
	if d.line != nil {
		d.line.Close()
	}
	d.line = nil

	// acquire chip, temporary
	// the given label is applied to all lines, which are requested on the chip
	gpiodChip, err := gpiod.NewChip(d.chipName, gpiod.WithConsumer(d.label))
	id := fmt.Sprintf("%s-%d", d.chipName, d.pin)
	if err != nil {
		return fmt.Errorf("gpiod.reconfigure(%s)-lib.NewChip(%s): %v", id, d.chipName, err)
	}
	defer gpiodChip.Close()

	// acquire line
	gpiodLine, err := gpiodChip.RequestLine(d.pin)
	if err != nil {
		if gpiodLine != nil {
			gpiodLine.Close()
		}
		d.line = nil

		return fmt.Errorf("gpiod.reconfigure(%s)-c.RequestLine(%d): %v", id, d.pin, err)
	}
	d.line = gpiodLine

	// configure direction
	if d.direction == IN || forceInput {
		if err := gpiodLine.Reconfigure(gpiod.AsInput); err != nil {
			return fmt.Errorf("gpiod.reconfigure(%s)-l.Reconfigure(gpiod.AsInput): %v", id, err)
		}
	} else {
		if err := gpiodLine.Reconfigure(gpiod.AsOutput(d.outInitialState)); err != nil {
			return fmt.Errorf("gpiod.reconfigure(%s)-l.Reconfigure(gpiod.AsOutput(%d)): %v", id, d.outInitialState, err)
		}
	}

	// configure inverse logic
	if d.activeLow {
		if err := gpiodLine.Reconfigure(gpiod.AsActiveLow); err != nil {
			return fmt.Errorf("gpiod.reconfigure(%s)-l.Reconfigure(gpiod.AsActiveLow): %v", id, err)
		}
	}

	// configure bias
	switch d.bias {
	case digitalPinBiasPullDown:
		if err := gpiodLine.Reconfigure(gpiod.LineBiasPullDown); err != nil {
			return fmt.Errorf("gpiod.reconfigure(%s)-l.Reconfigure(gpiod.LineBiasPullDown): %v", id, err)
		}
	case digitalPinBiasPullUp:
		if err := gpiodLine.Reconfigure(gpiod.LineBiasPullUp); err != nil {
			return fmt.Errorf("gpiod.reconfigure(%s)-l.Reconfigure(gpiod.LineBiasPullUp): %v", id, err)
		}
	default:
		if err := gpiodLine.Reconfigure(gpiod.LineBiasUnknown); err != nil {
			return fmt.Errorf("gpiod.reconfigure(%s)-l.Reconfigure(gpiod.LineBiasUnknown): %v", id, err)
		}
	}

	return nil
}

func fmtLine(li gpiod.LineInfo) string {
	var consumer string
	if li.Consumer != "" {
		consumer = fmt.Sprintf(" by '%s'", li.Consumer)
	}
	return fmt.Sprintf("++ Info line %d '%s', %s%s ++\n Config: %s\n",
		li.Offset, li.Name, used[li.Used], consumer, fmtLineConfig(li.Config))
}

func fmtLineConfig(cfg gpiod.LineConfig) string {
	t := "active-%s, %s, %s, %s bias, %s edge detect, %s, debounce-period: %v, %s event clock"
	return fmt.Sprintf(t, activeLow[cfg.ActiveLow], direction[cfg.Direction], drive[cfg.Drive], bias[cfg.Bias],
		edgeDetect[cfg.EdgeDetection], debounced[cfg.Debounced], cfg.DebouncePeriod, eventClock[cfg.EventClock])
}
