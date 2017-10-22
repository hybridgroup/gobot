package i2c

import (
	"gobot.io/x/gobot"
)

const ssd1306I2CAddress = 0x3c

const ssd1306Width = 128
const ssd1306Height = 64

const ssd1306PageSize = 8

const ssd1306SetMemoryAddressingMode = 0x20
const ssd1306SetComOutput0 = 0xC0
const ssd1306SetComOutput1 = 0xC1
const ssd1306SetComOutput2 = 0xC2
const ssd1306SetComOutput3 = 0xC3
const ssd1306SetComOutput4 = 0xC4
const ssd1306SetComOutput5 = 0xC5
const ssd1306SetComOutput6 = 0xC6
const ssd1306SetComOutput7 = 0xC7
const ssd1306SetComOutput8 = 0xC8
const ssd1306ColumnAddr = 0x21
const ssd1306PageAddr = 0x22
const ssd1306SetContrast = 0x81
const ssd1306SetSegmentRemap0 = 0xA0
const ssd1306SetSegmentRemap127 = 0xA1
const ssd1306DisplayOnResumeToRAM = 0xA4
const ssd1306SetDisplayNormal = 0xA6
const ssd1306SetDisplayInverse = 0xA7
const ssd1306SetDisplayOff = 0xAE
const ssd1306SetDisplayOn = 0xAF
const ssd1306ContinuousHScrollRight = 0x26
const ssd1306ContinuousHScrollLeft = 0x27
const ssd1306ContinuousVHScrollRight = 0x29
const ssd1306ContinuousVHScrollLeft = 0x2A
const ssd1306StopScroll = 0x2E
const ssd1306StartScroll = 0x2F
const ssd1306SetStartLine = 0x40
const ssd1306ChargePumpSetting = 0x8D
const ssd1306SetDisplayClock = 0xD5
const ssd1306SetMultiplexRatio = 0xA8
const ssd1306SetComPins = 0xDA
const ssd1306SetDisplayOffset = 0xD3
const ssd1306SetPrechargePeriod = 0xD9
const ssd1306SetVComDeselectLevel = 0xDB

var ssd1306InitSequence []byte = []byte{
	ssd1306SetDisplayNormal,
	ssd1306SetDisplayOff,
	ssd1306SetDisplayClock, 0x80, // the suggested ratio 0x80
	ssd1306SetMultiplexRatio, 0x3F,
	ssd1306SetDisplayOffset, 0x0, //no offset
	ssd1306SetStartLine | 0x0, //SETSTARTLINE
	ssd1306ChargePumpSetting, 0x14,
	ssd1306SetMemoryAddressingMode, 0x00, //0x0 act like ks0108
	ssd1306SetSegmentRemap0,
	ssd1306SetComOutput0,
	ssd1306SetComPins, 0x12, //COMSCANDEC
	ssd1306SetContrast, 0xCF,
	ssd1306SetPrechargePeriod, 0xF1,
	ssd1306SetVComDeselectLevel, 0x40,
	ssd1306DisplayOnResumeToRAM,
	ssd1306SetDisplayNormal,
	ssd1306StopScroll,
	ssd1306SetSegmentRemap0,
	ssd1306SetSegmentRemap127,
	ssd1306SetComOutput8,
	ssd1306SetMemoryAddressingMode, 0x00,
	ssd1306SetContrast, 0xff,
}

// DisplayBuffer represents the display buffer intermediate memory
type DisplayBuffer struct {
	Width, Height int
	buffer        []byte
}

// NewDisplayBuffer creates a new DisplayBuffer
func NewDisplayBuffer(Width, Height int) *DisplayBuffer {
	s := &DisplayBuffer{
		Width:  Width,
		Height: Height,
	}
	s.buffer = make([]byte, s.Size())
	return s
}

// Size returns the memory size of the display buffer
func (s *DisplayBuffer) Size() int {
	return (s.Width * s.Height) / ssd1306PageSize
}

// Clear the contents of the display buffer
func (s *DisplayBuffer) Clear() {
	s.buffer = make([]byte, s.Size())
}

// Set sets the x, y pixel with c color
func (s *DisplayBuffer) Set(x, y, c int) {
	idx := x + (y/ssd1306PageSize)*s.Width
	bit := uint(y) % ssd1306PageSize

	if c == 0 {
		s.buffer[idx] &= ^(1 << bit)
	} else {
		s.buffer[idx] |= (1 << bit)
	}
}

// SSD1306Driver is a Gobot Driver for a SSD1306 Display
type SSD1306Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	gobot.Commander

	DisplayWidth  int
	DisplayHeight int
	Buffer        *DisplayBuffer
}

// NewSSD1306Driver creates a new SSD1306Driver.
//
// Params:
//        conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//        WithBus(int):    			bus to use with this driver
//        WithAddress(int):    		address to use with this driver
//        WithDisplayWidth(int): 	width of display (defaults to 128)
//        WithDisplayHeight(int): 	height of display (defaults to 64)
//
func NewSSD1306Driver(a Connector, options ...func(Config)) *SSD1306Driver {
	s := &SSD1306Driver{
		name:          gobot.DefaultName("SSD1306"),
		Commander:     gobot.NewCommander(),
		connector:     a,
		Config:        NewConfig(),
		DisplayHeight: ssd1306Height,
		DisplayWidth:  ssd1306Width,
	}

	for _, option := range options {
		option(s)
	}

	s.Buffer = NewDisplayBuffer(s.DisplayWidth, s.DisplayHeight)

	s.AddCommand("Display", func(params map[string]interface{}) interface{} {
		err := s.Display()
		return map[string]interface{}{"err": err}
	})

	s.AddCommand("On", func(params map[string]interface{}) interface{} {
		err := s.On()
		return map[string]interface{}{"err": err}
	})

	s.AddCommand("Off", func(params map[string]interface{}) interface{} {
		err := s.Off()
		return map[string]interface{}{"err": err}
	})

	s.AddCommand("Clear", func(params map[string]interface{}) interface{} {
		err := s.Clear()
		return map[string]interface{}{"err": err}
	})

	s.AddCommand("SetContrast", func(params map[string]interface{}) interface{} {
		contrast := byte(params["contrast"].(byte))
		err := s.SetContrast(contrast)
		return map[string]interface{}{"err": err}
	})

	s.AddCommand("Set", func(params map[string]interface{}) interface{} {
		x := int(params["x"].(int))
		y := int(params["y"].(int))
		c := int(params["c"].(int))

		s.Set(x, y, c)
		return nil
	})

	return s
}

// Name returns the Name for the Driver
func (s *SSD1306Driver) Name() string { return s.name }

// SetName sets the Name for the Driver
func (s *SSD1306Driver) SetName(n string) { s.name = n }

// Connection returns the connection for the Driver
func (s *SSD1306Driver) Connection() gobot.Connection { return s.connector.(gobot.Connection) }

// Start starts the Driver up, and writes start command
func (s *SSD1306Driver) Start() (err error) {
	bus := s.GetBusOrDefault(s.connector.GetDefaultBus())
	address := s.GetAddressOrDefault(ssd1306I2CAddress)

	s.connection, err = s.connector.GetConnection(address, bus)
	if err != nil {
		return
	}

	s.Init()
	s.On()

	return
}

// Halt returns true if device is halted successfully
func (s *SSD1306Driver) Halt() (err error) { return nil }

// WithDisplayWidth option sets the SSD1306Driver DisplayWidth option.
func WithDisplayWidth(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*SSD1306Driver)
		if ok {
			d.DisplayWidth = val
		} else {
			// TODO: return error for trying to set DisplayWidth for non-SSD1306Driver
			return
		}
	}
}

// WithDisplayHeight option sets the SSD1306Driver DisplayHeight option.
func WithDisplayHeight(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*SSD1306Driver)
		if ok {
			d.DisplayHeight = val
		} else {
			// TODO: return error for trying to set DisplayHeight for non-SSD1306Driver
			return
		}
	}
}

// Init turns display on
func (s *SSD1306Driver) Init() (err error) {
	s.Off()
	s.commands(ssd1306InitSequence)

	s.commands([]byte{ssd1306ColumnAddr, 0, // Start at 0,
		byte(s.Buffer.Width) - 1, // End at last column (127?)
	})

	s.commands([]byte{ssd1306PageAddr, 0, // Start at 0,
		(byte(s.Buffer.Height) / ssd1306PageSize) - 1, // End at page 7
	})

	return nil
}

// On turns display on
func (s *SSD1306Driver) On() (err error) {
	return s.command(ssd1306SetDisplayOn)
}

// Off turns display off
func (s *SSD1306Driver) Off() (err error) {
	return s.command(ssd1306SetDisplayOff)
}

// Clear clears
func (s *SSD1306Driver) Clear() (err error) {
	s.Buffer.Clear()
	return nil
}

// Set sets a pixel
func (s *SSD1306Driver) Set(x, y, c int) {
	s.Buffer.Set(x, y, c)
}

// Reset sends the memory buffer to the display
func (s *SSD1306Driver) Reset() (err error) {
	s.Off()
	s.Clear()
	s.On()
	return nil
}

// SetContrast sets the display contrast
func (s *SSD1306Driver) SetContrast(contrast byte) (err error) {
	err = s.commands([]byte{ssd1306SetContrast, contrast})
	return
}

// Display sends the memory buffer to the display
func (s *SSD1306Driver) Display() (err error) {
	// Write the buffer
	_, err = s.connection.Write(append([]byte{0x40}, s.Buffer.buffer...))
	return err
}

// command sends a unique command
func (s *SSD1306Driver) command(b byte) (err error) {
	_, err = s.connection.Write([]byte{0x80, b})
	return
}

// commands sends a command sequence
func (s *SSD1306Driver) commands(commands []byte) (err error) {
	var command []byte
	for _, d := range commands {
		command = append(command, []byte{0x80, d}...)
	}
	_, err = s.connection.Write(command)
	return
}
