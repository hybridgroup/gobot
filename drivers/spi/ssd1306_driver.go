package spi

import (
	"fmt"
	"image"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
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
	bit := uint(y) % uint(d.pageSize) //nolint:gosec // TODO: fix later
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
	*Driver
	dcDriver      *gpio.DirectPinDriver
	rstDriver     *gpio.DirectPinDriver
	pageSize      int
	DisplayWidth  int
	DisplayHeight int
	DCPin         string
	RSTPin        string
	ExternalVcc   bool
	buffer        *DisplayBuffer
}

// NewSSD1306Driver creates a new SSD1306Driver.
//
// Params:
//
//	conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	 spi.WithBusNumber(int):  bus to use with this driver
//		spi.WithChipNumber(int): chip to use with this driver
//	 spi.WithMode(int):    	 mode to use with this driver
//	 spi.WithBitCount(int):   number of bits to use with this driver
//	 spi.WithSpeed(int64):    speed in Hz to use with this driver
//	 spi.WithDisplayWidth(int): 	width of display (defaults to 128)
//	 spi.WithDisplayHeight(int): height of display (defaults to 64)
//	 spi.WithDCPin(string): 		gpio pin number connected to dc pin on display (defaults to "16")
//	 spi.WithRstPin(string): 	gpio pin number connected to rst pin on display (defaults to "18")
//	 spi.WithExternalVCC(bool): 	set to true if using external vcc (defaults to false)
func NewSSD1306Driver(a gobot.Adaptor, options ...func(Config)) *SSD1306Driver {
	// cast adaptor to spi connector since we also need the adaptor for gpio
	b, ok := a.(Connector)
	if !ok {
		panic("unable to get gobot connector for ssd1306")
	}
	s := &SSD1306Driver{
		Driver:        NewDriver(b, "SSD1306"),
		DisplayWidth:  ssd1306Width,
		DisplayHeight: ssd1306Height,
		DCPin:         ssd1306DcPin,
		RSTPin:        ssd1306RstPin,
		ExternalVcc:   ssd1306ExternalVcc,
	}
	s.afterStart = s.initialize
	s.beforeHalt = s.shutdown

	for _, option := range options {
		option(s)
	}
	s.dcDriver = gpio.NewDirectPinDriver(a, s.DCPin)
	s.rstDriver = gpio.NewDirectPinDriver(a, s.RSTPin)
	s.pageSize = s.DisplayHeight / 8
	s.buffer = NewDisplayBuffer(s.DisplayWidth, s.DisplayHeight, s.pageSize)
	s.AddCommand("Display", func(_ map[string]interface{}) interface{} {
		err := s.Display()
		return map[string]interface{}{"err": err}
	})
	s.AddCommand("On", func(_ map[string]interface{}) interface{} {
		err := s.On()
		return map[string]interface{}{"err": err}
	})
	s.AddCommand("Off", func(_ map[string]interface{}) interface{} {
		err := s.Off()
		return map[string]interface{}{"err": err}
	})
	s.AddCommand("Clear", func(_ map[string]interface{}) interface{} {
		err := s.Clear()
		return map[string]interface{}{"err": err}
	})
	//nolint:forcetypeassert // ok here
	s.AddCommand("SetContrast", func(params map[string]interface{}) interface{} {
		contrast := params["contrast"].(byte)
		err := s.SetContrast(contrast)
		return map[string]interface{}{"err": err}
	})
	//nolint:forcetypeassert // ok here
	s.AddCommand("Set", func(params map[string]interface{}) interface{} {
		x := params["x"].(int)
		y := params["y"].(int)
		c := params["c"].(int)
		s.Set(x, y, c)
		return nil
	})
	return s
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
func (s *SSD1306Driver) On() error {
	return s.command(ssd1306SetDisplayOn)
}

// Off turns off the display.
func (s *SSD1306Driver) Off() error {
	return s.command(ssd1306SetDisplayOff)
}

// Clear clears the display buffer.
func (s *SSD1306Driver) Clear() error {
	s.buffer.Clear()
	return nil
}

// Set sets a pixel in the display buffer.
func (s *SSD1306Driver) Set(x, y, c int) {
	s.buffer.SetPixel(x, y, c)
}

// Reset re-initializes the device to a clean state.
func (s *SSD1306Driver) Reset() error {
	if err := s.rstDriver.DigitalWrite(1); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := s.rstDriver.DigitalWrite(0); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := s.rstDriver.DigitalWrite(1); err != nil {
		return err
	}
	return nil
}

// SetBufferAndDisplay sets the display buffer with the given buffer and displays the image.
func (s *SSD1306Driver) SetBufferAndDisplay(buf []byte) error {
	s.buffer.Set(buf)
	return s.Display()
}

// SetContrast sets the display contrast (0-255).
func (s *SSD1306Driver) SetContrast(contrast byte) error {
	if err := s.command(ssd1306SetContrast); err != nil {
		return err
	}
	return s.command(contrast)
}

// Display sends the memory buffer to the display.
func (s *SSD1306Driver) Display() error {
	if err := s.command(ssd1306ColumnAddr); err != nil {
		return err
	}
	if err := s.command(0); err != nil {
		return err
	}
	//nolint:gosec // TODO: fix later
	if err := s.command(uint8(s.DisplayWidth) - 1); err != nil {
		return err
	}
	if err := s.command(ssd1306PageAddr); err != nil {
		return err
	}
	if err := s.command(0); err != nil {
		return err
	}
	//nolint:gosec // TODO: fix later
	if err := s.command(uint8(s.pageSize) - 1); err != nil {
		return err
	}
	if err := s.dcDriver.DigitalWrite(1); err != nil {
		return err
	}
	return s.connection.WriteBlockData(0x40, s.buffer.buffer)
}

// ShowImage takes a standard Go image and shows it on the display in monochrome.
func (s *SSD1306Driver) ShowImage(img image.Image) error {
	if img.Bounds().Dx() != s.DisplayWidth || img.Bounds().Dy() != s.DisplayHeight {
		return fmt.Errorf("Image must match the display width and height")
	}

	if err := s.Clear(); err != nil {
		return err
	}
	for y, w, h := 0, img.Bounds().Dx(), img.Bounds().Dy(); y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(x, y)
			if r, g, b, _ := c.RGBA(); r > 0 || g > 0 || b > 0 {
				s.Set(x, y, 1)
			}
		}
	}
	return s.Display()
}

// command sends a unique command
func (s *SSD1306Driver) command(b byte) error {
	if err := s.dcDriver.DigitalWrite(0); err != nil {
		return err
	}
	return s.connection.WriteByte(b)
}

// initialize configures the ssd1306 based on the options passed in when the driver was created
func (s *SSD1306Driver) initialize() error {
	if err := s.command(ssd1306SetDisplayOff); err != nil {
		return err
	}
	if err := s.command(ssd1306SetDisplayClock); err != nil {
		return err
	}
	if s.DisplayHeight == 16 {
		if err := s.command(0x60); err != nil {
			return err
		}
	} else {
		if err := s.command(0x80); err != nil {
			return err
		}
	}
	if err := s.command(ssd1306SetMultiplexRatio); err != nil {
		return err
	}
	//nolint:gosec // TODO: fix later
	if err := s.command(uint8(s.DisplayHeight) - 1); err != nil {
		return err
	}
	if err := s.command(ssd1306SetDisplayOffset); err != nil {
		return err
	}
	if err := s.command(0x0); err != nil {
		return err
	}
	if err := s.command(ssd1306SetStartLine); err != nil {
		return err
	}
	if err := s.command(0x0); err != nil {
		return err
	}
	if err := s.command(ssd1306ChargePumpSetting); err != nil {
		return err
	}
	if s.ExternalVcc {
		if err := s.command(0x10); err != nil {
			return err
		}
	} else {
		if err := s.command(0x14); err != nil {
			return err
		}
	}
	if err := s.command(ssd1306SetMemoryAddressingMode); err != nil {
		return err
	}
	if err := s.command(0x00); err != nil {
		return err
	}
	if err := s.command(ssd1306SetSegmentRemap0); err != nil {
		return err
	}
	if err := s.command(0x01); err != nil {
		return err
	}
	if err := s.command(ssd1306ComScanInc); err != nil {
		return err
	}
	if err := s.command(ssd1306SetComPins); err != nil {
		return err
	}
	if s.DisplayHeight == 64 {
		if err := s.command(0x12); err != nil {
			return err
		}
	} else {
		if err := s.command(0x02); err != nil {
			return err
		}
	}
	if err := s.command(ssd1306SetContrast); err != nil {
		return err
	}
	if s.DisplayHeight == 64 {
		if s.ExternalVcc {
			if err := s.command(0x9F); err != nil {
				return err
			}
		} else {
			if err := s.command(0xCF); err != nil {
				return err
			}
		}
	} else {
		if err := s.command(0x8F); err != nil {
			return err
		}
	}
	if err := s.command(ssd1306SetPrechargePeriod); err != nil {
		return err
	}
	if s.ExternalVcc {
		if err := s.command(0x22); err != nil {
			return err
		}
	} else {
		if err := s.command(0xF1); err != nil {
			return err
		}
	}
	if err := s.command(ssd1306SetVComDeselectLevel); err != nil {
		return err
	}
	if err := s.command(0x40); err != nil {
		return err
	}
	if err := s.command(ssd1306DisplayOnResumeToRAM); err != nil {
		return err
	}
	if err := s.command(ssd1306SetDisplayNormal); err != nil {
		return err
	}
	if err := s.command(ssd1306DeactivateScroll); err != nil {
		return err
	}
	return s.command(ssd1306SetDisplayOn)
}

func (s *SSD1306Driver) shutdown() error {
	if err := s.Reset(); err != nil {
		return err
	}
	return s.Off()
}
