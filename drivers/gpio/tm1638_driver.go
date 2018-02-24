package gpio

import (
	"math"

	"strings"

	"gobot.io/x/gobot"
)

const (
	TM1638None = iota
	TM1638Red
	TM1638Green
)

const (
	TM1638DataCmd  = 0x40
	TM1638DispCtrl = 0x80
	TM1638AddrCmd  = 0xC0

	TM1638WriteDisp = 0x00
	TM1638ReadKeys  = 0x02
	TM1638FixedAddr = 0x04
)

// TM1638Driver is the gobot driver for modules based on the TM1638, which has 8 7-segment displays, 8 LEDs and 8 buttons
// Buttons are not supported
//
// Datasheet EN: https://retrocip.cz/files/tm1638.pdf
// Datasheet CN: http://www.datasheetspdf.com/pdf/775356/TitanMicro/TM1638/1
//
// Ported from the Arduino driver https://github.com/rjbatista/tm1638-library

type TM1638Driver struct {
	pinClock  *DirectPinDriver
	pinData    *DirectPinDriver
	pinStrobe  *DirectPinDriver
	fonts      map[string]byte
	name       string
	connection gobot.Connection
	gobot.Commander
}

// NewTM1638Driver return a new TM1638Driver given a gobot.Connection and the clock, data and strobe pins
func NewTM1638Driver(a gobot.Connection, clockPin string, dataPin string, strobePin string) *TM1638Driver {
	t := &TM1638Driver{
		name:       gobot.DefaultName("TM1638"),
		pinClock:  NewDirectPinDriver(a, clockPin),
		pinData:    NewDirectPinDriver(a, dataPin),
		pinStrobe:  NewDirectPinDriver(a, strobePin),
		fonts:      NewTM1638Fonts(),
		connection: a,
		Commander:  gobot.NewCommander(),
	}

	/* TODO : Add commands */

	return t
}

// Start initializes the tm1638, it uses a SPI-like communication protocol
func (t *TM1638Driver) Start() (err error) {

	t.pinStrobe.On()
	t.pinClock.On()

	t.sendCommand(TM1638DataCmd)
	t.sendCommand(TM1638DispCtrl | 8 | 7)

	t.pinStrobe.Off()
	t.send(TM1638AddrCmd)
	for i := 0; i < 16; i++ {
		t.send(TM1638WriteDisp)
	}
	t.pinStrobe.On()

	return
}

// Halt implements the Driver interface
func (t *TM1638Driver) Halt() (err error) { return }

// Name returns the TM1638Drivers name
func (t *TM1638Driver) Name() string { return t.name }

// SetName sets the TM1638Drivers name
func (t *TM1638Driver) SetName(n string) { t.name = n }

// Connection returns the TM1638Driver Connection
func (t *TM1638Driver) Connection() gobot.Connection {
	return t.connection
}

// sendCommand is an auxiliary function to send commands to the TM1638 module
func (t *TM1638Driver) sendCommand(cmd byte) {
	t.pinStrobe.Off()
	t.send(cmd)
	t.pinStrobe.On()
}

// send writes data on the module
func (t *TM1638Driver) send(data byte) {
	for i := 0; i < 8; i++ {
		t.pinClock.Off()

		if (data & 1) > 0 {
			t.pinData.On()
		} else {
			t.pinData.Off()
		}
		data >>= 1

		t.pinClock.On()
	}
}

// sendData is an auxiliary function to send data to the TM1638 module
func (t *TM1638Driver) sendData(address byte, data byte) {
	t.sendCommand(TM1638DataCmd | TM1638FixedAddr)
	t.pinStrobe.Off()
	t.send(TM1638AddrCmd | address)
	t.send(data)
	t.pinStrobe.On()
}

// SetLED changes the color (TM1638None, TM1638Red, TM1638Green) of the specific LED
func (t *TM1638Driver) SetLED(color byte, pos byte) {
	if pos > 7 {
		return
	}
	t.sendData((pos<<1)+1, color)
}

// SetDisplay cuts and sends a byte array to the display (without dots)
func (t *TM1638Driver) SetDisplay(data []byte) {
	minLength := int(math.Min(8, float64(len(data))))
	for i := 0; i < minLength; i++ {
		t.SendChar(byte(i), data[i], false)
	}
}

// SetDisplayText cuts and sends a string to the display (without dots)
func (t *TM1638Driver) SetDisplayText(text string) {
	data := t.fromStringToByteArray(text)
	minLength := int(math.Min(8, float64(len(data))))
	for i := 0; i < minLength; i++ {
		t.SendChar(byte(i), data[i], false)
	}
}

// SendChar sends one byte to the specific position in the display
func (t *TM1638Driver) SendChar(pos byte, data byte, dot bool) {
	if pos > 7 {
		return
	}
	var dotData byte
	if dot {
		dotData = TM1638DispCtrl
	}
	t.sendData(pos<<1, data|(dotData))
}

// fromStringToByteArray translates a string to a byte array with the corresponding representation for the 7-segment LCD, return and empty character if the font is not available
func (t *TM1638Driver) fromStringToByteArray(str string) []byte {
	chars := strings.Split(str, "")
	data := make([]byte, len(chars))

	for index, char := range chars {
		if val, ok := t.fonts[char]; ok {
			data[index] = val
		}
	}
	return data
}

// AddFonts adds new custom fonts or modify the representation of existing ones
func (t *TM1638Driver) AddFonts(fonts map[string]byte) {
	for k, v := range fonts {
		t.fonts[k] = v
	}
}

// ClearFonts removes all the fonts from the driver
func (t *TM1638Driver) ClearFonts() {
	t.fonts = make(map[string]byte)
}

// NewTM1638Fonts returns a map with fonts and their corresponding byte for proper representation on the 7-segment LCD
func NewTM1638Fonts() map[string]byte {
	return map[string]byte{
		" ":  0x00,
		"!":  0x86,
		"'":  0x22,
		"#":  0x7E,
		"$":  0x6D,
		"%":  0x00,
		"&":  0x00,
		"\"": 0x02,
		"(":  0x30,
		")":  0x06,
		"*":  0x63,
		"+":  0x00,
		",":  0x04,
		"-":  0x40,
		".":  0x80,
		"/":  0x52,
		"0":  0x3F,
		"1":  0x06,
		"2":  0x5B,
		"3":  0x4F,
		"4":  0x66,
		"5":  0x6D,
		"6":  0x7D,
		"7":  0x27,
		"8":  0x7F,
		"9":  0x6F,
		":":  0x00,
		";":  0x00,
		"<":  0x00,
		"=":  0x48,
		">":  0x00,
		"?":  0x53,
		"@":  0x5F,
		"A":  0x77,
		"B":  0x7F,
		"C":  0x39,
		"D":  0x3F,
		"E":  0x79,
		"F":  0x71,
		"G":  0x3D,
		"H":  0x76,
		"I":  0x06,
		"J":  0x1F,
		"K":  0x69,
		"L":  0x38,
		"M":  0x15,
		"N":  0x37,
		"O":  0x3F,
		"P":  0x73,
		"Q":  0x67,
		"R":  0x31,
		"S":  0x6D,
		"T":  0x78,
		"U":  0x3E,
		"V":  0x2A,
		"W":  0x1D,
		"X":  0x76,
		"Y":  0x6E,
		"Z":  0x5B,
		"[":  0x39,
		"\\": 0x64, // (this can't be the last char on a line, even in comment or it'll concat)
		"]":  0x0F,
		"^":  0x00,
		"_":  0x08,
		"`":  0x20,
		"a":  0x5F,
		"b":  0x7C,
		"c":  0x58,
		"d":  0x5E,
		"e":  0x7B,
		"f":  0x31,
		"g":  0x6F,
		"h":  0x74,
		"i":  0x04,
		"j":  0x0E,
		"k":  0x75,
		"l":  0x30,
		"m":  0x55,
		"n":  0x54,
		"o":  0x5C,
		"p":  0x73,
		"q":  0x67,
		"r":  0x50,
		"s":  0x6D,
		"t":  0x78,
		"u":  0x1C,
		"v":  0x2A,
		"w":  0x1D,
		"x":  0x76,
		"y":  0x6E,
		"z":  0x47,
		"{":  0x46,
		"|":  0x06,
		"}":  0x70,
		"~":  0x01,
	}
}
