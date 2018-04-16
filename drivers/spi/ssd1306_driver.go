package spi

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	// default values
	ssd1306Width        = 128
	ssd1306Height       = 64
	ssd1306DcPin        = "16" // for raspberry pi
	ssd1306RstPin       = "18" // for raspberry pi
	ssd1306ExternalVcc  = false
	ssd1306SetStartLine = 0x40
	// fundamental commands
	ssd1306SetContrast          = 0x81
	ssd1306DisplayOnResumeToRAM = 0xA4
	ssd1306DisplayOnResume      = 0xA5
	ssd1306SetDisplayNormal     = 0xA6
	ssd1306SetDisplayInverse    = 0xA7
	ssd1306SetDisplayOff        = 0xAE
	ssd1306SetDisplayOn         = 0xAF
	// scrolling commands
	ssd1306RightHorizontalScroll            = 0x26
	ssd1306LeftHorizontalScroll             = 0x27
	ssd1306VerticalAndRightHorizontalScroll = 0x29
	ssd1306VerticalAndLeftHorizontalScroll  = 0x2A
	ssd1306DeactivateScroll                 = 0x2E
	ssd1306ActivateScroll                   = 0x2F
	ssd1306SetVerticalScrollArea            = 0xA3
	// addressing settings commands
	ssd1306SetMemoryAddressingMode = 0x20
	ssd1306ColumnAddr              = 0x21
	ssd1306PageAddr                = 0x22
	// hardware configuration commands
	ssd1306SetSegmentRemap0   = 0xA0
	ssd1306SetSegmentRemap127 = 0xA1
	ssd1306SetMultiplexRatio  = 0xA8
	ssd1306ComScanInc         = 0xC0
	ssd1306ComScanDec         = 0xC8
	ssd1306SetDisplayOffset   = 0xD3
	ssd1306SetComPins         = 0xDA
	// timing and driving scheme commands
	ssd1306SetDisplayClock      = 0xD5
	ssd1306SetPrechargePeriod   = 0xD9
	ssd1306SetVComDeselectLevel = 0xDB
	ssd1306NOOP                 = 0xE3
	// charge pump command
	ssd1306ChargePumpSetting = 0x8D
)

// ssd1306Init configures the ssd1306 based on the options passed in when the driver was created
func (s *SSD1306Driver) ssd1306Init() {
	s.command(ssd1306SetDisplayOff)
	s.command(ssd1306SetDisplayClock)
	if s.DisplayHeight == 16 {
		s.command(0x60)
	} else {
		s.command(0x80)
	}
	s.command(ssd1306SetMultiplexRatio)
	s.command(uint8(s.DisplayHeight) - 1)
	s.command(ssd1306SetDisplayOffset)
	s.command(0x0)
	s.command(ssd1306SetStartLine)
	s.command(0x0)
	s.command(ssd1306ChargePumpSetting)
	if s.ExternalVcc {
		s.command(0x10)
	} else {
		s.command(0x14)
	}
	s.command(ssd1306SetMemoryAddressingMode)
	s.command(0x00)
	s.command(ssd1306SetSegmentRemap0)
	s.command(0x01)
	s.command(ssd1306ComScanInc)
	s.command(ssd1306SetComPins)
	if s.DisplayHeight == 64 {
		s.command(0x12)
	} else {
		s.command(0x02)
	}
	s.command(ssd1306SetContrast)
	if s.DisplayHeight == 64 {
		if s.ExternalVcc {
			s.command(0x9F)
		} else {
			s.command(0xCF)
		}
	} else {
		s.command(0x8F)
	}
	s.command(ssd1306SetPrechargePeriod)
	if s.ExternalVcc {
		s.command(0x22)
	} else {
		s.command(0xF1)
	}
	s.command(ssd1306SetVComDeselectLevel)
	s.command(0x40)
	s.command(ssd1306DisplayOnResumeToRAM)
	s.command(ssd1306SetDisplayNormal)
	s.command(ssd1306DeactivateScroll)
	s.command(ssd1306SetDisplayOn)
}

// DisplayBuffer represents the display buffer intermediate memory
type DisplayBuffer struct {
	width, height, pageSize int
	buffer                  []byte
}

// NewDisplayBuffer creates a new DisplayBuffer
func NewDisplayBuffer(width, height, pageSize int) *DisplayBuffer {
	s := &DisplayBuffer{
		width:    width,
		height:   height,
		pageSize: pageSize,
	}
	s.buffer = make([]byte, s.Size())
	return s
}

// Size returns the memory size of the display buffer
func (d *DisplayBuffer) Size() int {
	return (d.width * d.height) / d.pageSize
}

// Clear the contents of the display buffer
func (d *DisplayBuffer) Clear() {
	d.buffer = make([]byte, d.Size())
}

// SetPixel sets the x, y pixel with c color
func (d *DisplayBuffer) SetPixel(x, y, c int) {
	idx := x + (y/d.pageSize)*d.width
	bit := uint(y) % uint(d.pageSize)
	if c == 0 {
		d.buffer[idx] &= ^(1 << bit)
	} else {
		d.buffer[idx] |= (1 << bit)
	}
}

// Set sets the display buffer with the given buffer
func (d *DisplayBuffer) Set(buf []byte) {
	d.buffer = buf
}

// SSD1306Driver is a Gobot Driver for a SSD1306 Display
type SSD1306Driver struct {
	name          string
	connector     Connector
	connection    Connection
	dcDriver      *gpio.DirectPinDriver
	rstDriver     *gpio.DirectPinDriver
	pageSize      int
	DisplayWidth  int
	DisplayHeight int
	DCPin         string
	RSTPin        string
	ExternalVcc   bool
	buffer        *DisplayBuffer
	Config
	gobot.Commander
}

// NewSSD1306Driver creates a new SSD1306Driver.
//
// Params:
//      conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//      spi.WithBus(int):    		bus to use with this driver
//     	spi.WithChip(int):    		chip to use with this driver
//      spi.WithMode(int):    		mode to use with this driver
//      spi.WithBits(int):    		number of bits to use with this driver
//      spi.WithSpeed(int64):   	speed in Hz to use with this driver
//      spi.WithDisplayWidth(int): 	width of display (defaults to 128)
//      spi.WithDisplayHeight(int): height of display (defaults to 64)
//      spi.WithDCPin(string): 		gpio pin number connected to dc pin on display (defaults to "16")
//      spi.WithRstPin(string): 	gpio pin number connected to rst pin on display (defaults to "18")
//      spi.WithExternalVCC(bool): 	set to true if using external vcc (defaults to false)
//
func NewSSD1306Driver(a gobot.Adaptor, options ...func(Config)) *SSD1306Driver {
	// cast adaptor to spi connector since we also need the adaptor for gpio
	b, ok := a.(Connector)
	if !ok {
		panic("unable to get gobot connector for ssd1306")
	}
	s := &SSD1306Driver{
		name:          gobot.DefaultName("SSD1306"),
		Commander:     gobot.NewCommander(),
		connector:     b,
		DisplayWidth:  ssd1306Width,
		DisplayHeight: ssd1306Height,
		DCPin:         ssd1306DcPin,
		RSTPin:        ssd1306RstPin,
		ExternalVcc:   ssd1306ExternalVcc,
		Config:        NewConfig(),
	}
	for _, option := range options {
		option(s)
	}
	s.dcDriver = gpio.NewDirectPinDriver(a, s.DCPin)
	s.rstDriver = gpio.NewDirectPinDriver(a, s.RSTPin)
	s.pageSize = s.DisplayHeight / 8
	s.buffer = NewDisplayBuffer(s.DisplayWidth, s.DisplayHeight, s.pageSize)
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

// Start sets up the needed connection, and initialized the device.
func (s *SSD1306Driver) Start() (err error) {
	bus := s.GetBusOrDefault(s.connector.GetSpiDefaultBus())
	chip := s.GetChipOrDefault(s.connector.GetSpiDefaultChip())
	mode := s.GetModeOrDefault(s.connector.GetSpiDefaultMode())
	bits := s.GetBitsOrDefault(s.connector.GetSpiDefaultBits())
	maxSpeed := s.GetSpeedOrDefault(s.connector.GetSpiDefaultMaxSpeed())

	s.connection, err = s.connector.GetSpiConnection(bus, chip, mode, bits, maxSpeed)
	if err != nil {
		return err
	}
	s.ssd1306Init()
	return
}

// Halt returns true if device is halted successfully.
func (s *SSD1306Driver) Halt() (err error) {
	s.Reset()
	s.Off()
	return nil
}

// WithDisplayWidth option sets the SSD1306Driver DisplayWidth option.
func WithDisplayWidth(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*SSD1306Driver)
		if ok {
			d.DisplayWidth = val
		} else {
			panic("unable to set display width for ssd1306")
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
			panic("unable to set display height for ssd1306")
		}
	}
}

// WithDCPin option sets the SSD1306Driver DC Pin option.
func WithDCPin(val string) func(Config) {
	return func(c Config) {
		d, ok := c.(*SSD1306Driver)
		if ok {
			d.DCPin = val
		} else {
			panic("unable to set dc pin for ssd1306")
		}
	}
}

// WithRstPin option sets the SSD1306Driver RST pin option.
func WithRstPin(val string) func(Config) {
	return func(c Config) {
		d, ok := c.(*SSD1306Driver)
		if ok {
			d.RSTPin = val
		} else {
			panic("unable to set rst pin for ssd1306")
		}
	}
}

// WithExternalVCC option sets the SSD1306Driver external vcc option.
func WithExternalVCC(val bool) func(Config) {
	return func(c Config) {
		d, ok := c.(*SSD1306Driver)
		if ok {
			d.ExternalVcc = val
		} else {
			panic("unable to set rst pin for ssd1306")
		}
	}
}

// On turns on the display.
func (s *SSD1306Driver) On() (err error) {
	return s.command(ssd1306SetDisplayOn)
}

// Off turns off the display.
func (s *SSD1306Driver) Off() (err error) {
	return s.command(ssd1306SetDisplayOff)
}

// Clear clears the display buffer.
func (s *SSD1306Driver) Clear() (err error) {
	s.buffer.Clear()
	return nil
}

// Set sets a pixel in the display buffer.
func (s *SSD1306Driver) Set(x, y, c int) {
	s.buffer.SetPixel(x, y, c)
}

// Reset re-initializes the device to a clean state.
func (s *SSD1306Driver) Reset() (err error) {
	s.rstDriver.DigitalWrite(1)
	time.Sleep(10 * time.Millisecond)
	s.rstDriver.DigitalWrite(0)
	time.Sleep(10 * time.Millisecond)
	s.rstDriver.DigitalWrite(1)
	return nil
}

// SetBufferAndDisplay sets the display buffer with the given buffer and displays the image.
func (s *SSD1306Driver) SetBufferAndDisplay(buf []byte) (err error) {
	s.buffer.Set(buf)
	return s.Display()
}

// SetContrast sets the display contrast (0-255).
func (s *SSD1306Driver) SetContrast(contrast byte) (err error) {
	if contrast < 0 || contrast > 255 {
		return fmt.Errorf("contrast value must be between 0-255")
	}
	if err = s.command(ssd1306SetContrast); err != nil {
		return err
	}
	return s.command(contrast)
}

// Display sends the memory buffer to the display.
func (s *SSD1306Driver) Display() (err error) {
	s.command(ssd1306ColumnAddr)
	s.command(0)
	s.command(uint8(s.DisplayWidth) - 1)
	s.command(ssd1306PageAddr)
	s.command(0)
	s.command(uint8(s.pageSize) - 1)
	if err = s.dcDriver.DigitalWrite(1); err != nil {
		return err
	}
	return s.connection.Tx(append([]byte{0x40}, s.buffer.buffer...), nil)
}

// command sends a unique command
func (s *SSD1306Driver) command(b byte) (err error) {
	if err = s.dcDriver.DigitalWrite(0); err != nil {
		return err
	}
	err = s.connection.Tx([]byte{b}, nil)
	return err
}
