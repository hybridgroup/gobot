package gpio

import (
	"gobot.io/x/gobot"
)

const (
	MAX7219Digit0 = 0x01
	MAX7219Digit1 = 0x02
	MAX7219Digit2 = 0x03
	MAX7219Digit3 = 0x04
	MAX7219Digit4 = 0x05
	MAX7219Digit5 = 0x06
	MAX7219Digit6 = 0x07
	MAX7219Digit7 = 0x08

	MAX7219DecodeMode  = 0x09
	MAX7219Intensity   = 0x0a
	MAX7219ScanLimit   = 0x0b
	MAX7219Shutdown    = 0x0c
	MAX7219DisplayTest = 0x0f
)

// MAX7219Driver is the gobot driver for the MAX7219 LED driver
//
// Datasheet: https://datasheets.maximintegrated.com/en/ds/MAX7219-MAX7221.pdf
type MAX7219Driver struct {
	pinClock   *DirectPinDriver
	pinData    *DirectPinDriver
	pinCS      *DirectPinDriver
	name       string
	count      uint
	connection gobot.Connection
	gobot.Commander
}

// NewMAX7219Driver return a new MAX7219Driver given a gobot.Connection, pins and how many chips are chained
func NewMAX7219Driver(a gobot.Connection, clockPin string, dataPin string, csPin string, count uint) *MAX7219Driver {
	t := &MAX7219Driver{
		name:       gobot.DefaultName("MAX7219Driver"),
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

// Start initializes the max7219, it uses a SPI-like communication protocol
func (a *MAX7219Driver) Start() (err error) {
	a.pinData.On()
	a.pinClock.On()
	a.pinCS.On()

	a.All(MAX7219ScanLimit, 0x07)
	a.All(MAX7219DecodeMode, 0x00)
	a.All(MAX7219Shutdown, 0x01)
	a.All(MAX7219DisplayTest, 0x00)
	a.ClearAll()
	a.All(MAX7219Intensity, 0x0f&0x0f)

	return
}

// Halt implements the Driver interface
func (a *MAX7219Driver) Halt() (err error) { return }

// Name returns the MAX7219Drivers name
func (a *MAX7219Driver) Name() string { return a.name }

// SetName sets the MAX7219Drivers name
func (a *MAX7219Driver) SetName(n string) { a.name = n }

// Connection returns the MAX7219Driver Connection
func (a *MAX7219Driver) Connection() gobot.Connection {
	return a.connection
}

// SetIntensity changes the intensity (from 1 to 7) of the display
func (a *MAX7219Driver) SetIntensity(level byte) {
	if level > 15 {
		level = 15
	}
	a.All(MAX7219Intensity, level&level)
}

// ClearAll turns off all LEDs of all modules
func (a *MAX7219Driver) ClearAll() {
	for i := 1; i <= 8; i++ {
		a.All(byte(i), 0)
	}
}

// ClearAll turns off all LEDs of the given module
func (a *MAX7219Driver) ClearOne(which uint) {
	for i := 1; i <= 8; i++ {
		a.One(which, byte(i), 0)
	}
}

// sendData is an auxiliary function to send data to the MAX7219Driver module
func (a *MAX7219Driver) sendData(address byte, data byte) {
	a.pinCS.Off()
	a.send(address)
	a.send(data)
	a.pinCS.On()
}

// send writes data on the module
func (a *MAX7219Driver) send(data byte) {
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
func (a *MAX7219Driver) All(address byte, data byte) {
	a.pinCS.Off()
	var c uint
	for c = 0; c < a.count; c++ {
		a.send(address)
		a.send(data)
	}
	a.pinCS.On()
}

// One sends data to a specific module
func (a *MAX7219Driver) One(which uint, address byte, data byte) {
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
