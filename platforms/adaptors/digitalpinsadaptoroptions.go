package adaptors

import (
	"time"
)

// DigitalPinsOptionApplier is the interface for digital adaptors options. This provides the possibility for change the
// platform behavior by the user when creating the platform, e.g. by "NewAdaptor()".
// The interface needs to be implemented by each configurable option type.
type DigitalPinsOptionApplier interface {
	apply(cfg *digitalPinsConfiguration)
}

// digitalPinInitializeOption is the type for applying another than the default initializer
type digitalPinsInitializeOption digitalPinInitializer

// digitalPinsSystemSysfsOption is the type to change the default character device implementation to the legacy
// sysfs Kernel ABI
type digitalPinsSystemSysfsOption bool

// digitalPinsForSystemSpiOption is the type to switch the default SPI implementation to GPIO usage
type digitalPinsForSystemSpiOption digitalPinGpiosForSPI

// digitalPinsActiveLowOption is the type to prepare the given pins for inverse reaction on next initialize
type digitalPinsActiveLowOption []string

// digitalPinsPullDownOption is the type to prepare the given pins to be pulled down (high impedance to GND)
// on next initialize
type digitalPinsPullDownOption []string

// digitalPinsPullUpOption is the type to prepare the given pins to be pulled up (high impedance to VDD)
// on next initialize
type digitalPinsPullUpOption []string

// digitalPinsOpenDrainOption is the type to prepare the given output pins to be driven with open drain/collector
// on next initialize
type digitalPinsOpenDrainOption []string

// digitalPinsOpenSourceOption is the type to prepares the given output pins to be driven with open source/emitter
// on next initialize
type digitalPinsOpenSourceOption []string

// digitalPinsDebounceOption is the type to prepare the given input pin to be debounced on next initialize
type digitalPinsDebounceOption struct {
	id     string
	period time.Duration
}

// digitalPinsEventOnFallingEdgeOption is the type to prepare the given input pin to be generate an event
// on falling edge
type digitalPinsEventOnFallingEdgeOption struct {
	id      string
	handler func(int, time.Duration, string, uint32, uint32)
}

// digitalPinsEventOnRisingEdgeOption is the type to prepare the given input pin to be generate an event
// on rising edge
type digitalPinsEventOnRisingEdgeOption struct {
	id      string
	handler func(int, time.Duration, string, uint32, uint32)
}

// digitalPinsEventOnBothEdgesOption is the type to prepare the given input pin to be generate an event
// on rising and falling edges
type digitalPinsEventOnBothEdgesOption struct {
	id      string
	handler func(int, time.Duration, string, uint32, uint32)
}

// digitalPinsPollForEdgeDetectionOption is the type to prepare the given input pin to use a discrete input pin polling
// function together with edge detection.
type digitalPinsPollForEdgeDetectionOption struct {
	id           string
	pollInterval time.Duration
	pollQuitChan chan struct{}
}

func (o digitalPinsInitializeOption) String() string {
	return "pin initializer option for digital IO's"
}

func (o digitalPinsSystemSysfsOption) String() string {
	return "use sysfs vs. cdev system access option for digital pins"
}

func (o digitalPinsForSystemSpiOption) String() string {
	return "use digital pins for SPI option"
}

func (o digitalPinsActiveLowOption) String() string {
	return "digital pins set to active low option"
}

func (o digitalPinsPullDownOption) String() string {
	return "digital pins set to pull down option"
}

func (o digitalPinsPullUpOption) String() string {
	return "digital pins set to pull up option"
}

func (o digitalPinsOpenDrainOption) String() string {
	return "digital pins set to open drain option"
}

func (o digitalPinsOpenSourceOption) String() string {
	return "digital pins set to open source option"
}

func (o digitalPinsDebounceOption) String() string {
	return "set debounce time for digital pin option"
}

func (o digitalPinsEventOnFallingEdgeOption) String() string {
	return "set event on falling edge for digital pin option"
}

func (o digitalPinsEventOnRisingEdgeOption) String() string {
	return "set event on rising edge for digital pin option"
}

func (o digitalPinsEventOnBothEdgesOption) String() string {
	return "set event on rising and falling edge for digital pin option"
}

func (o digitalPinsPollForEdgeDetectionOption) String() string {
	return "discrete polling function for edge detection on digital pin option"
}

func (o digitalPinsInitializeOption) apply(cfg *digitalPinsConfiguration) {
	cfg.initialize = digitalPinInitializer(o)
}

func (o digitalPinsSystemSysfsOption) apply(cfg *digitalPinsConfiguration) {
	c := bool(o)
	cfg.useSysfs = &c
}

func (o digitalPinsForSystemSpiOption) apply(cfg *digitalPinsConfiguration) {
	c := digitalPinGpiosForSPI(o)
	cfg.gpiosForSPI = &c
}

func (o digitalPinsActiveLowOption) apply(cfg *digitalPinsConfiguration) {
	cfg.activeLowPins = append(cfg.activeLowPins, o...)
}

func (o digitalPinsPullDownOption) apply(cfg *digitalPinsConfiguration) {
	cfg.pullDownPins = append(cfg.pullDownPins, o...)
}

func (o digitalPinsPullUpOption) apply(cfg *digitalPinsConfiguration) {
	cfg.pullUpPins = append(cfg.pullUpPins, o...)
}

func (o digitalPinsOpenDrainOption) apply(cfg *digitalPinsConfiguration) {
	cfg.openDrainPins = append(cfg.openDrainPins, o...)
}

func (o digitalPinsOpenSourceOption) apply(cfg *digitalPinsConfiguration) {
	cfg.openSourcePins = append(cfg.openSourcePins, o...)
}

func (o digitalPinsDebounceOption) apply(cfg *digitalPinsConfiguration) {
	pinCfg := digitalPinsDebouncePin(o)
	cfg.debouncePins = append(cfg.debouncePins, pinCfg)
}

func (o digitalPinsEventOnFallingEdgeOption) apply(cfg *digitalPinsConfiguration) {
	pinCfg := digitalPinsEventOnEdgePin(o)
	cfg.eventOnFallingEdgePins = append(cfg.eventOnFallingEdgePins, pinCfg)
}

func (o digitalPinsEventOnRisingEdgeOption) apply(cfg *digitalPinsConfiguration) {
	pinCfg := digitalPinsEventOnEdgePin(o)
	cfg.eventOnRisingEdgePins = append(cfg.eventOnRisingEdgePins, pinCfg)
}

func (o digitalPinsEventOnBothEdgesOption) apply(cfg *digitalPinsConfiguration) {
	pinCfg := digitalPinsEventOnEdgePin(o)
	cfg.eventOnBothEdgesPins = append(cfg.eventOnBothEdgesPins, pinCfg)
}

func (o digitalPinsPollForEdgeDetectionOption) apply(cfg *digitalPinsConfiguration) {
	pinCfg := digitalPinsPollForEdgeDetectionPin(o)
	cfg.pollForEdgeDetectionPins = append(cfg.pollForEdgeDetectionPins, pinCfg)
}
