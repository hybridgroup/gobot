package gpio

import (
	"gobot.io/x/gobot"
)

const (
	MAX72xxDigit0 = 0x01
	MAX72xxDigit1 = 0x02
	MAX72xxDigit2 = 0x03
	MAX72xxDigit3 = 0x04
	MAX72xxDigit4 = 0x05
	MAX72xxDigit5 = 0x06
	MAX72xxDigit6 = 0x07
	MAX72xxDigit7 = 0x08

	MAX72xxDecodeMode  = 0x09
	MAX72xxIntensity   = 0x0a
	MAX72xxScanLimit   = 0x0b
	MAX72xxShutdown    = 0x0c
	MAX72xxDisplayTest = 0x0f
)

// MAX72xxDriver is the gobot driver for the MAX7219 & MAX7221 LED drivers
//
// Datasheet: https://datasheets.maximintegrated.com/en/ds/MAX7219-MAX7221.pdf
type MAX72xxDriver struct {
	pinClock   *DirectPinDriver
	pinData    *DirectPinDriver
	pinCS      *DirectPinDriver
	name       string
	count      uint
	connection gobot.Connection
	gobot.Commander
}

// NewMAX72xxDriver return a new MAX72xxDriver given a gobot.Connection, pins and how many chips are chained
func NewMAX72xxDriver(a gobot.Connection, clockPin string, dataPin string, csPin string, count uint) *MAX72xxDriver {
	t := &MAX72xxDriver{
		name:       gobot.DefaultName("MAX72xxDriver"),
		pinClock:   NewDirectPinDriver(a, clockPin),
		pinData:    NewDirectPinDriver(a, dataPin),
		pinCS:      NewDirectPinDriver(a, csPin),
		count:      count,
		connection: a,
		Commander:  gobot.NewCommander(),
	}

	/* TODO : Add commands */

	return t
}

// Start initializes the max72xx, it uses a SPI-like communication protocol
func (a *MAX72xxDriver) Start() (err error) {
	a.pinData.On()
	a.pinClock.On()
	a.pinCS.On()

	a.All(MAX72xxScanLimit, 0x07)
	a.All(MAX72xxDecodeMode, 0x00)
	a.All(MAX72xxShutdown, 0x01)
	a.All(MAX72xxDisplayTest, 0x00)
	a.ClearAll()
	a.All(MAX72xxIntensity, 0x0f&0x0f)

	return
}

// Halt implements the Driver interface
func (a *MAX72xxDriver) Halt() (err error) { return }

// Name returns the MAX72xxDrivers name
func (a *MAX72xxDriver) Name() string { return a.name }

// SetName sets the MAX72xxDrivers name
func (a *MAX72xxDriver) SetName(n string) { a.name = n }

// Connection returns the MAX72xxDriver Connection
func (a *MAX72xxDriver) Connection() gobot.Connection {
	return a.connection
}

// SetIntensity changes the intensity (from 1 to 7) of the display
func (a *MAX72xxDriver) SetIntensity(level byte) {
	if level > 15 {
		level = 15
	}
	a.All(MAX72xxIntensity, level&level)
}

// ClearAll turns off all LEDs of all modules
func (a *MAX72xxDriver) ClearAll() {
	for i := 1; i <= 8; i++ {
		a.All(byte(i), 0)
	}
}

// ClearAll turns off all LEDs of the given module
func (a *MAX72xxDriver) ClearOne(which uint) {
	for i := 1; i <= 8; i++ {
		a.One(which, byte(i), 0)
	}
}

// sendData is an auxiliary function to send data to the MAX72xxDriver module
func (a *MAX72xxDriver) sendData(address byte, data byte) {
	a.pinCS.Off()
	a.send(address)
	a.send(data)
	a.pinCS.On()
}

// send writes data on the module
func (a *MAX72xxDriver) send(data byte) {
	var i byte
	for i = 8; i > 0; i-- {
		mask := byte(0x01 << (i - 1))

		a.pinClock.Off()
		if data&mask > 0 {
			a.pinData.On()
		} else {
			a.pinData.Off()
		}
		a.pinClock.On()
	}
}

// All sends the same data to all the modules
func (a *MAX72xxDriver) All(address byte, data byte) {
	a.pinCS.Off()
	var c uint
	for c = 0; c < a.count; c++ {
		a.send(address)
		a.send(data)
	}
	a.pinCS.On()
}

// One sends data to a specific module
func (a *MAX72xxDriver) One(which uint, address byte, data byte) {
	a.pinCS.Off()
	var c uint
	for c = 0; c < a.count; c++ {
		if c == which {
			a.send(address)
			a.send(data)
		} else {
			a.send(0)
			a.send(0)
		}
	}
	a.pinCS.On()
}
