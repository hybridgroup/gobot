package system

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	gpiocdev "github.com/warthog618/go-gpiocdev"

	"gobot.io/x/gobot/v2"
)

const systemCdevDebug = false

type cdevLine interface {
	SetValue(value int) error
	Value() (int, error)
	Close() error
}

type digitalPinCdev struct {
	chipName string
	pin      int
	*digitalPinConfig
	line cdevLine
}

var digitalPinCdevReconfigure = digitalPinCdevReconfigureLine // to allow unit testing

var (
	digitalPinCdevUsed      = map[bool]string{true: "used", false: "unused"}
	digitalPinCdevActiveLow = map[bool]string{true: "low", false: "high"}
	digitalPinCdevDebounced = map[bool]string{true: "debounced", false: "not debounced"}
)

var digitalPinCdevDirection = map[gpiocdev.LineDirection]string{
	gpiocdev.LineDirectionUnknown: "unknown direction",
	gpiocdev.LineDirectionInput:   "input", gpiocdev.LineDirectionOutput: "output",
}

var digitalPinCdevDrive = map[gpiocdev.LineDrive]string{
	gpiocdev.LineDrivePushPull: "push-pull", gpiocdev.LineDriveOpenDrain: "open-drain",
	gpiocdev.LineDriveOpenSource: "open-source",
}

var digitalPinCdevBias = map[gpiocdev.LineBias]string{
	gpiocdev.LineBiasUnknown: "unknown", gpiocdev.LineBiasDisabled: "disabled",
	gpiocdev.LineBiasPullUp: "pull-up", gpiocdev.LineBiasPullDown: "pull-down",
}

var digitalPinCdevEdgeDetect = map[gpiocdev.LineEdge]string{
	gpiocdev.LineEdgeNone: "no", gpiocdev.LineEdgeRising: "rising",
	gpiocdev.LineEdgeFalling: "falling", gpiocdev.LineEdgeBoth: "both",
}

var digitalPinCdevEventClock = map[gpiocdev.LineEventClock]string{
	gpiocdev.LineEventClockMonotonic: "monotonic",
	gpiocdev.LineEventClockRealtime:  "realtime",
}

// newDigitalPinCdev returns a digital pin given the pin number, with the label "gobotio" followed by the pin number.
// The pin label can be modified optionally. The pin is handled by the character device Kernel ABI.
func newDigitalPinCdev(chipName string, pin int, options ...func(gobot.DigitalPinOptioner) bool) *digitalPinCdev {
	if chipName == "" {
		chipName = "gpiochip0"
	}
	cfg := newDigitalPinConfig("gobotio"+strconv.Itoa(pin), options...)
	d := &digitalPinCdev{
		chipName:         chipName,
		pin:              pin,
		digitalPinConfig: cfg,
	}
	return d
}

// ApplyOptions apply all given options to the pin immediately. Implements interface gobot.DigitalPinOptionApplier.
func (d *digitalPinCdev) ApplyOptions(options ...func(gobot.DigitalPinOptioner) bool) error {
	anyChange := false
	for _, option := range options {
		anyChange = option(d) || anyChange
	}
	if anyChange {
		return digitalPinCdevReconfigure(d, false)
	}
	return nil
}

// DirectionBehavior gets the direction behavior when the pin is used the next time. This means its possibly not in
// this direction type at the moment. Implements the interface gobot.DigitalPinValuer, but should be rarely used.
func (d *digitalPinCdev) DirectionBehavior() string {
	return d.direction
}

// Export sets the pin as used by this driver. Implements the interface gobot.DigitalPinner.
func (d *digitalPinCdev) Export() error {
	err := digitalPinCdevReconfigure(d, false)
	if err != nil {
		return fmt.Errorf("cdev.Export(): %v", err)
	}
	return nil
}

// Unexport releases the pin as input. Implements the interface gobot.DigitalPinner.
func (d *digitalPinCdev) Unexport() error {
	var errs []string
	if d.line != nil {
		if err := digitalPinCdevReconfigure(d, true); err != nil {
			errs = append(errs, err.Error())
		}
		if err := d.line.Close(); err != nil {
			err = fmt.Errorf("cdev.Unexport()-line.Close(): %v", err)
			errs = append(errs, err.Error())
		}
	}
	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ","))
}

// Write writes the given value to the character device. Implements the interface gobot.DigitalPinner.
func (d *digitalPinCdev) Write(val int) error {
	if val < 0 {
		val = 0
	}
	if val > 1 {
		val = 1
	}

	err := d.line.SetValue(val)
	if err != nil {
		return fmt.Errorf("cdev.Write(): %v", err)
	}
	return nil
}

// Read reads the given value from character device. Implements the interface gobot.DigitalPinner.
func (d *digitalPinCdev) Read() (int, error) {
	val, err := d.line.Value()
	if err != nil {
		return 0, fmt.Errorf("cdev.Read(): %v", err)
	}
	return val, err
}

// ListLines is used for development purposes.
func (d *digitalPinCdev) ListLines() error {
	c, err := gpiocdev.NewChip(d.chipName, gpiocdev.WithConsumer(d.label))
	if err != nil {
		return err
	}
	for i := 0; i < c.Lines(); i++ {
		li, err := c.LineInfo(i)
		if err != nil {
			return err
		}
		fmt.Println(digitalPinCdevFmtLine(li))
	}

	return nil
}

// List is used for development purposes.
func (d *digitalPinCdev) List() error {
	c, err := gpiocdev.NewChip(d.chipName)
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
	fmt.Println(digitalPinCdevFmtLine(li))

	return nil
}

func digitalPinCdevReconfigureLine(d *digitalPinCdev, forceInput bool) error {
	// cleanup old line
	if d.line != nil {
		d.line.Close()
	}
	d.line = nil

	// acquire chip, temporary
	// the given label is applied to all lines, which are requested on the chip
	gpiodChip, err := gpiocdev.NewChip(d.chipName, gpiocdev.WithConsumer(d.label))
	id := fmt.Sprintf("%s-%d", d.chipName, d.pin)
	if err != nil {
		return fmt.Errorf("cdev.reconfigure(%s)-lib.NewChip(%s): %v", id, d.chipName, err)
	}
	defer gpiodChip.Close()

	// collect line configuration options
	var opts []gpiocdev.LineReqOption

	// configure direction, debounce period (inputs only), edge detection (inputs only) and drive (outputs only)
	if d.direction == IN || forceInput {
		if systemCdevDebug {
			log.Printf("input (%s): debounce %s, edge %d, handler %t, inverse %t, bias %d",
				id, d.debouncePeriod, d.edge, d.edgeEventHandler != nil, d.activeLow, d.bias)
		}
		opts = append(opts, gpiocdev.AsInput)
		if !forceInput && d.drive != digitalPinDrivePushPull && systemCdevDebug {
			log.Printf("\n++ drive option (%d) is dropped for input++\n", d.drive)
		}
		if d.debouncePeriod != 0 {
			opts = append(opts, gpiocdev.WithDebounce(d.debouncePeriod))
		}
		// edge detection
		if d.edgeEventHandler != nil && d.pollInterval <= 0 {
			// use edge detection provided by gpiocdev
			wrappedHandler := digitalPinCdevGetWrappedEventHandler(d.edgeEventHandler)
			switch d.edge {
			case digitalPinEventOnFallingEdge:
				opts = append(opts, gpiocdev.WithEventHandler(wrappedHandler), gpiocdev.WithFallingEdge)
			case digitalPinEventOnRisingEdge:
				opts = append(opts, gpiocdev.WithEventHandler(wrappedHandler), gpiocdev.WithRisingEdge)
			case digitalPinEventOnBothEdges:
				opts = append(opts, gpiocdev.WithEventHandler(wrappedHandler), gpiocdev.WithBothEdges)
			default:
				opts = append(opts, gpiocdev.WithoutEdges)
			}
		}
	} else {
		if systemCdevDebug {
			log.Printf("output (%s): ini-state %d, drive %d, inverse %t, bias %d",
				id, d.outInitialState, d.drive, d.activeLow, d.bias)
		}
		opts = append(opts, gpiocdev.AsOutput(d.outInitialState))
		switch d.drive {
		case digitalPinDriveOpenDrain:
			opts = append(opts, gpiocdev.AsOpenDrain)
		case digitalPinDriveOpenSource:
			opts = append(opts, gpiocdev.AsOpenSource)
		default:
			opts = append(opts, gpiocdev.AsPushPull)
		}
		if d.debouncePeriod != 0 && systemCdevDebug {
			log.Printf("\n++debounce option (%d) is dropped for output++\n", d.drive)
		}
		if d.edgeEventHandler != nil || d.edge != digitalPinEventNone && systemCdevDebug {
			log.Printf("\n++edge detection is dropped for output++\n")
		}
	}

	// configure inverse logic (inputs and outputs)
	if d.activeLow {
		opts = append(opts, gpiocdev.AsActiveLow)
	}

	// configure bias (inputs and outputs)
	switch d.bias {
	case digitalPinBiasPullDown:
		opts = append(opts, gpiocdev.WithPullDown)
	case digitalPinBiasPullUp:
		opts = append(opts, gpiocdev.WithPullUp)
	default:
		opts = append(opts, gpiocdev.WithBiasAsIs)
	}

	// acquire line with collected options
	gpiodLine, err := gpiodChip.RequestLine(d.pin, opts...)
	if err != nil {
		if gpiodLine != nil {
			gpiodLine.Close()
		}
		d.line = nil

		return fmt.Errorf("cdev.reconfigure(%s)-c.RequestLine(%d, %v): %v", id, d.pin, opts, err)
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

func digitalPinCdevGetWrappedEventHandler(
	handler func(int, time.Duration, string, uint32, uint32),
) func(gpiocdev.LineEvent) {
	return func(evt gpiocdev.LineEvent) {
		detectedEdge := "none"
		switch evt.Type {
		case gpiocdev.LineEventRisingEdge:
			detectedEdge = DigitalPinEventRisingEdge
		case gpiocdev.LineEventFallingEdge:
			detectedEdge = DigitalPinEventFallingEdge
		}
		handler(evt.Offset, evt.Timestamp, detectedEdge, evt.Seqno, evt.LineSeqno)
	}
}

func digitalPinCdevFmtLine(li gpiocdev.LineInfo) string {
	var consumer string
	if li.Consumer != "" {
		consumer = fmt.Sprintf(" by '%s'", li.Consumer)
	}
	return fmt.Sprintf("++ Info line %d '%s', %s%s ++\n Config: %s\n",
		li.Offset, li.Name, digitalPinCdevUsed[li.Used], consumer, digitalPinCdevFmtLineConfig(li.Config))
}

func digitalPinCdevFmtLineConfig(cfg gpiocdev.LineConfig) string {
	t := "active-%s, %s, %s, %s bias, %s edge detect, %s, debounce-period: %v, %s event clock"
	return fmt.Sprintf(t, digitalPinCdevActiveLow[cfg.ActiveLow], digitalPinCdevDirection[cfg.Direction],
		digitalPinCdevDrive[cfg.Drive], digitalPinCdevBias[cfg.Bias], digitalPinCdevEdgeDetect[cfg.EdgeDetection],
		digitalPinCdevDebounced[cfg.Debounced], cfg.DebouncePeriod, digitalPinCdevEventClock[cfg.EventClock])
}
