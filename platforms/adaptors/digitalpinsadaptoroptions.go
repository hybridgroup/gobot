package adaptors

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

// DigitalPinsOptioner is the interface for digital adaptors options. This provides the possibility for change the
// platform behavior by the user when creating the platform, e.g. by "NewAdaptor()".
// TODO: change to applier-architecture, see options of pwmpinsadaptor.go
type DigitalPinsOptioner interface {
	setDigitalPinInitializer(initializer digitalPinInitializer)
	setDigitalPinsForSystemGpiod()
	setDigitalPinsForSystemSpi(sclkPin, ncsPin, sdoPin, sdiPin string)
	prepareDigitalPinsActiveLow(pin string, otherPins ...string)
	prepareDigitalPinsPullDown(pin string, otherPins ...string)
	prepareDigitalPinsPullUp(pin string, otherPins ...string)
	prepareDigitalPinsOpenDrain(pin string, otherPins ...string)
	prepareDigitalPinsOpenSource(pin string, otherPins ...string)
	prepareDigitalPinDebounce(pin string, period time.Duration)
	prepareDigitalPinEventOnFallingEdge(pin string, handler func(lineOffset int, timestamp time.Duration,
		detectedEdge string, seqno uint32, lseqno uint32))
	prepareDigitalPinEventOnRisingEdge(pin string, handler func(lineOffset int, timestamp time.Duration,
		detectedEdge string, seqno uint32, lseqno uint32))
	prepareDigitalPinEventOnBothEdges(pin string, handler func(lineOffset int, timestamp time.Duration,
		detectedEdge string, seqno uint32, lseqno uint32))
	prepareDigitalPinPollForEdgeDetection(pin string, pollInterval time.Duration, pollQuitChan chan struct{})
}

func (a *DigitalPinsAdaptor) setDigitalPinInitializer(pinInit digitalPinInitializer) {
	a.initialize = pinInit
}

func (a *DigitalPinsAdaptor) setDigitalPinsForSystemGpiod() {
	system.WithDigitalPinGpiodAccess()(a.sys)
}

func (a *DigitalPinsAdaptor) setDigitalPinsForSystemSpi(sclkPin, ncsPin, sdoPin, sdiPin string) {
	system.WithSpiGpioAccess(a, sclkPin, ncsPin, sdoPin, sdiPin)(a.sys)
}

func (a *DigitalPinsAdaptor) prepareDigitalPinsActiveLow(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinActiveLow())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinsPullDown(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinPullDown())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinsPullUp(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinPullUp())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinsOpenDrain(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinOpenDrain())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinsOpenSource(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinOpenSource())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinDebounce(id string, period time.Duration) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinDebounce(period))
}

func (a *DigitalPinsAdaptor) prepareDigitalPinEventOnFallingEdge(id string, handler func(int, time.Duration, string,
	uint32, uint32),
) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinEventOnFallingEdge(handler))
}

func (a *DigitalPinsAdaptor) prepareDigitalPinEventOnRisingEdge(id string, handler func(int, time.Duration, string,
	uint32, uint32),
) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinEventOnRisingEdge(handler))
}

func (a *DigitalPinsAdaptor) prepareDigitalPinEventOnBothEdges(id string, handler func(int, time.Duration, string,
	uint32, uint32),
) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinEventOnBothEdges(handler))
}

func (a *DigitalPinsAdaptor) prepareDigitalPinPollForEdgeDetection(
	id string,
	pollInterval time.Duration,
	pollQuitChan chan struct{},
) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinPollForEdgeDetection(pollInterval, pollQuitChan))
}
