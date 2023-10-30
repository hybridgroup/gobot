package system

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/warthog618/gpiod"
	"gobot.io/x/gobot/v2"
)

const systemGpiodDebug = false

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

var digitalPinGpiodReconfigure = digitalPinGpiodReconfigureLine // to allow unit testing

var (
	digitalPinGpiodUsed      = map[bool]string{true: "used", false: "unused"}
	digitalPinGpiodActiveLow = map[bool]string{true: "low", false: "high"}
	digitalPinGpiodDebounced = map[bool]string{true: "debounced", false: "not debounced"}
)

var digitalPinGpiodDirection = map[gpiod.LineDirection]string{
	gpiod.LineDirectionUnknown: "unknown direction",
	gpiod.LineDirectionInput:   "input", gpiod.LineDirectionOutput: "output",
}

var digitalPinGpiodDrive = map[gpiod.LineDrive]string{
	gpiod.LineDrivePushPull: "push-pull", gpiod.LineDriveOpenDrain: "open-drain",
	gpiod.LineDriveOpenSource: "open-source",
}

var digitalPinGpiodBias = map[gpiod.LineBias]string{
	gpiod.LineBiasUnknown: "unknown", gpiod.LineBiasDisabled: "disabled",
	gpiod.LineBiasPullUp: "pull-up", gpiod.LineBiasPullDown: "pull-down",
}

var digitalPinGpiodEdgeDetect = map[gpiod.LineEdge]string{
	gpiod.LineEdgeNone: "no", gpiod.LineEdgeRising: "rising",
	gpiod.LineEdgeFalling: "falling", gpiod.LineEdgeBoth: "both",
}

var digitalPinGpiodEventClock = map[gpiod.LineEventClock]string{
	gpiod.LineEventClockMonotonic: "monotonic",
	gpiod.LineEventClockRealtime:  "realtime",
}

// newDigitalPinGpiod returns a digital pin given the pin number, with the label "gobotio" followed by the pin number.
// The pin label can be modified optionally. The pin is handled by the character device Kernel ABI.
func newDigitalPinGpiod(chipName string, pin int, options ...func(gobot.DigitalPinOptioner) bool) *digitalPinGpiod {
	if chipName == "" {
		chipName = "gpiochip0"
	}
	cfg := newDigitalPinConfig("gobotio"+strconv.Itoa(pin), options...)
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
		anyChange = option(d) || anyChange
	}
	if anyChange {
		return digitalPinGpiodReconfigure(d, false)
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
	err := digitalPinGpiodReconfigure(d, false)
	if err != nil {
		return fmt.Errorf("gpiod.Export(): %v", err)
	}
	return nil
}

// Unexport releases the pin as input. Implements the interface gobot.DigitalPinner.
func (d *digitalPinGpiod) Unexport() error {
	var errs []string
	if d.line != nil {
		if err := digitalPinGpiodReconfigure(d, true); err != nil {
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
		fmt.Println(digitalPinGpiodFmtLine(li))
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
	fmt.Println(digitalPinGpiodFmtLine(li))

	return nil
}

func digitalPinGpiodReconfigureLine(d *digitalPinGpiod, forceInput bool) error {
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

	// collect line configuration options
	var opts []gpiod.LineReqOption

	// configure direction, debounce period (inputs only), edge detection (inputs only) and drive (outputs only)
	if d.direction == IN || forceInput {
		if systemGpiodDebug {
			log.Printf("input (%s): debounce %s, edge %d, handler %t, inverse %t, bias %d",
				id, d.debouncePeriod, d.edge, d.edgeEventHandler != nil, d.activeLow, d.bias)
		}
		opts = append(opts, gpiod.AsInput)
		if !forceInput && d.drive != digitalPinDrivePushPull && systemGpiodDebug {
			log.Printf("\n++ drive option (%d) is dropped for input++\n", d.drive)
		}
		if d.debouncePeriod != 0 {
			opts = append(opts, gpiod.WithDebounce(d.debouncePeriod))
		}
		// edge detection
		if d.edgeEventHandler != nil && d.pollInterval <= 0 {
			// use edge detection provided by gpiod
			wrappedHandler := digitalPinGpiodGetWrappedEventHandler(d.edgeEventHandler)
			switch d.edge {
			case digitalPinEventOnFallingEdge:
				opts = append(opts, gpiod.WithEventHandler(wrappedHandler), gpiod.WithFallingEdge)
			case digitalPinEventOnRisingEdge:
				opts = append(opts, gpiod.WithEventHandler(wrappedHandler), gpiod.WithRisingEdge)
			case digitalPinEventOnBothEdges:
				opts = append(opts, gpiod.WithEventHandler(wrappedHandler), gpiod.WithBothEdges)
			default:
				opts = append(opts, gpiod.WithoutEdges)
			}
		}
	} else {
		if systemGpiodDebug {
			log.Printf("output (%s): ini-state %d, drive %d, inverse %t, bias %d",
				id, d.outInitialState, d.drive, d.activeLow, d.bias)
		}
		opts = append(opts, gpiod.AsOutput(d.outInitialState))
		switch d.drive {
		case digitalPinDriveOpenDrain:
			opts = append(opts, gpiod.AsOpenDrain)
		case digitalPinDriveOpenSource:
			opts = append(opts, gpiod.AsOpenSource)
		default:
			opts = append(opts, gpiod.AsPushPull)
		}
		if d.debouncePeriod != 0 && systemGpiodDebug {
			log.Printf("\n++debounce option (%d) is dropped for output++\n", d.drive)
		}
		if d.edgeEventHandler != nil || d.edge != digitalPinEventNone && systemGpiodDebug {
			log.Printf("\n++edge detection is dropped for output++\n")
		}
	}

	// configure inverse logic (inputs and outputs)
	if d.activeLow {
		opts = append(opts, gpiod.AsActiveLow)
	}

	// configure bias (inputs and outputs)
	switch d.bias {
	case digitalPinBiasPullDown:
		opts = append(opts, gpiod.WithPullDown)
	case digitalPinBiasPullUp:
		opts = append(opts, gpiod.WithPullUp)
	default:
		opts = append(opts, gpiod.WithBiasAsIs)
	}

	// acquire line with collected options
	gpiodLine, err := gpiodChip.RequestLine(d.pin, opts...)
	if err != nil {
		if gpiodLine != nil {
			gpiodLine.Close()
		}
		d.line = nil

		return fmt.Errorf("gpiod.reconfigure(%s)-c.RequestLine(%d, %v): %v", id, d.pin, opts, err)
	}
	d.line = gpiodLine

	// start discrete polling function and wait for first read is done
	if (d.direction == IN || forceInput) && d.pollInterval > 0 {
		if err := startEdgePolling(d.label, d.Read, d.pollInterval, d.edge, d.edgeEventHandler,
			d.pollQuitChan); err != nil {
			return err
		}
	}

	return nil
}

func digitalPinGpiodGetWrappedEventHandler(
	handler func(int, time.Duration, string, uint32, uint32),
) func(gpiod.LineEvent) {
	return func(evt gpiod.LineEvent) {
		detectedEdge := "none"
		switch evt.Type {
		case gpiod.LineEventRisingEdge:
			detectedEdge = DigitalPinEventRisingEdge
		case gpiod.LineEventFallingEdge:
			detectedEdge = DigitalPinEventFallingEdge
		}
		handler(evt.Offset, evt.Timestamp, detectedEdge, evt.Seqno, evt.LineSeqno)
	}
}

func digitalPinGpiodFmtLine(li gpiod.LineInfo) string {
	var consumer string
	if li.Consumer != "" {
		consumer = fmt.Sprintf(" by '%s'", li.Consumer)
	}
	return fmt.Sprintf("++ Info line %d '%s', %s%s ++\n Config: %s\n",
		li.Offset, li.Name, digitalPinGpiodUsed[li.Used], consumer, digitalPinGpiodFmtLineConfig(li.Config))
}

func digitalPinGpiodFmtLineConfig(cfg gpiod.LineConfig) string {
	t := "active-%s, %s, %s, %s bias, %s edge detect, %s, debounce-period: %v, %s event clock"
	return fmt.Sprintf(t, digitalPinGpiodActiveLow[cfg.ActiveLow], digitalPinGpiodDirection[cfg.Direction],
		digitalPinGpiodDrive[cfg.Drive], digitalPinGpiodBias[cfg.Bias], digitalPinGpiodEdgeDetect[cfg.EdgeDetection],
		digitalPinGpiodDebounced[cfg.Debounced], cfg.DebouncePeriod, digitalPinGpiodEventClock[cfg.EventClock])
}
