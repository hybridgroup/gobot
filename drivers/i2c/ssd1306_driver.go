package i2c

import (
	"fmt"
	"image"
)

const ssd1306DefaultAddress = 0x3c

// register addresses for the ssd1306
const (
	// default values
	ssd1306Width        = 128
	ssd1306Height       = 64
	ssd1306ExternalVCC  = false
	ssd1306SetStartLine = 0x40
	// fundamental commands
	ssd1306SetComOutput0 = 0xC0
	ssd1306SetComOutput1 = 0xC1
	ssd1306SetComOutput2 = 0xC2
	ssd1306SetComOutput3 = 0xC3
	ssd1306SetComOutput4 = 0xC4
	ssd1306SetComOutput5 = 0xC5
	ssd1306SetComOutput6 = 0xC6
	ssd1306SetComOutput7 = 0xC7
	ssd1306SetComOutput8 = 0xC8
	ssd1306SetContrast   = 0x81
	// scrolling commands
	// ssd1306ContinuousHScrollRight  = 0x26
	// ssd1306ContinuousHScrollLeft   = 0x27
	// ssd1306ContinuousVHScrollRight = 0x29
	// ssd1306ContinuousVHScrollLeft  = 0x2A
	// ssd1306StopScroll              = 0x2E
	// ssd1306StartScroll             = 0x2F
	// addressing settings commands
	ssd1306SetMemoryAddressingMode = 0x20
	ssd1306ColumnAddr              = 0x21
	ssd1306PageAddr                = 0x22
	// hardware configuration commands
	ssd1306SetSegmentRemap0     = 0xA0
	ssd1306SetSegmentRemap127   = 0xA1
	ssd1306DisplayOnResumeToRAM = 0xA4
	ssd1306SetDisplayNormal     = 0xA6
	ssd1306SetDisplayInverse    = 0xA7
	ssd1306SetDisplayOff        = 0xAE
	ssd1306SetDisplayOn         = 0xAF
	// timing and driving scheme commands
	ssd1306SetDisplayClock      = 0xD5
	ssd1306SetPrechargePeriod   = 0xD9
	ssd1306SetVComDeselectLevel = 0xDB
	ssd1306SetMultiplexRatio    = 0xA8
	ssd1306SetComPins           = 0xDA
	ssd1306SetDisplayOffset     = 0xD3
	// charge pump command
	ssd1306ChargePumpSetting = 0x8D
)

// SSD1306Init contains the initialization settings for the ssd1306 display.
type SSD1306Init struct {
	displayClock         byte
	multiplexRatio       byte
	displayOffset        byte
	startLine            byte
	chargePumpSetting    byte
	memoryAddressingMode byte
	comPins              byte
	contrast             byte
	prechargePeriod      byte
	vComDeselectLevel    byte
}

// GetSequence returns the initialization sequence for the ssd1306 display.
func (i *SSD1306Init) GetSequence() []byte {
	return []byte{
		ssd1306SetDisplayNormal,
		ssd1306SetDisplayOff,
		ssd1306SetDisplayClock, i.displayClock,
		ssd1306SetMultiplexRatio, i.multiplexRatio,
		ssd1306SetDisplayOffset, i.displayOffset,
		ssd1306SetStartLine | i.startLine,
		ssd1306ChargePumpSetting, i.chargePumpSetting,
		ssd1306SetMemoryAddressingMode, i.memoryAddressingMode,
		ssd1306SetSegmentRemap0,
		ssd1306SetComOutput0,
		ssd1306SetComPins, i.comPins,
		ssd1306SetContrast, i.contrast,
		ssd1306SetPrechargePeriod, i.prechargePeriod,
		ssd1306SetVComDeselectLevel, i.vComDeselectLevel,
		ssd1306DisplayOnResumeToRAM,
		ssd1306SetDisplayNormal,
	}
}

// 128x64 init sequence
var ssd1306Init128x64 = &SSD1306Init{
	displayClock:         0x80,
	multiplexRatio:       0x3F,
	displayOffset:        0x00,
	startLine:            0x00,
	chargePumpSetting:    0x14, // 0x10 if external vcc is set
	memoryAddressingMode: 0x00,
	comPins:              0x12,
	contrast:             0xCF, // 0x9F if external vcc is set
	prechargePeriod:      0xF1, // 0x22 if external vcc is set
	vComDeselectLevel:    0x40,
}

// 128x32 init sequence
var ssd1306Init128x32 = &SSD1306Init{
	displayClock:         0x80,
	multiplexRatio:       0x1F,
	displayOffset:        0x00,
	startLine:            0x00,
	chargePumpSetting:    0x14, // 0x10 if external vcc is set
	memoryAddressingMode: 0x00,
	comPins:              0x02,
	contrast:             0x8F, // 0x9F if external vcc is set
	prechargePeriod:      0xF1, // 0x22 if external vcc is set
	vComDeselectLevel:    0x40,
}

// 96x16 init sequence
var ssd1306Init96x16 = &SSD1306Init{
	displayClock:         0x60,
	multiplexRatio:       0x0F,
	displayOffset:        0x00,
	startLine:            0x00,
	chargePumpSetting:    0x14, // 0x10 if external vcc is set
	memoryAddressingMode: 0x00,
	comPins:              0x02,
	contrast:             0x8F, // 0x9F if external vcc is set
	prechargePeriod:      0xF1, // 0x22 if external vcc is set
	vComDeselectLevel:    0x40,
}

// DisplayBuffer represents the display buffer intermediate memory.
type DisplayBuffer struct {
	width, height, pageSize int
	buffer                  []byte
}

// NewDisplayBuffer creates a new DisplayBuffer.
func NewDisplayBuffer(width, height, pageSize int) *DisplayBuffer {
	d := &DisplayBuffer{
		width:    width,
		height:   height,
		pageSize: pageSize,
	}
	d.buffer = make([]byte, d.Size())
	return d
}

// Size returns the memory size of the display buffer.
func (d *DisplayBuffer) Size() int {
	return (d.width * d.height) / d.pageSize
}

// Clear the contents of the display buffer.
func (d *DisplayBuffer) Clear() {
	d.buffer = make([]byte, d.Size())
}

// SetPixel sets the x, y pixel with c color.
func (d *DisplayBuffer) SetPixel(x, y, c int) {
	idx := x + (y/d.pageSize)*d.width
	bit := uint(y) % uint(d.pageSize) //nolint:gosec // TODO: fix later
	if c == 0 {
		d.buffer[idx] &= ^(1 << bit)
	} else {
		d.buffer[idx] |= (1 << bit)
	}
}

// Set sets the display buffer with the given buffer.
func (d *DisplayBuffer) Set(buf []byte) {
	d.buffer = buf
}

// SSD1306Driver is a Gobot Driver for a SSD1306 Display.
type SSD1306Driver struct {
	*Driver
	initSequence  *SSD1306Init
	displayWidth  int
	displayHeight int
	externalVCC   bool
	pageSize      int
	buffer        *DisplayBuffer
}

// NewSSD1306Driver creates a new SSD1306Driver.
//
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	WithBus(int):    			bus to use with this driver
//	WithAddress(int):    		address to use with this driver
//	WithSSD1306DisplayWidth(int): 	width of display (defaults to 128)
//	WithSSD1306DisplayHeight(int): 	height of display (defaults to 64)
//	WithSSD1306ExternalVCC:          set true when using an external OLED supply (defaults to false)
func NewSSD1306Driver(c Connector, options ...func(Config)) *SSD1306Driver {
	s := &SSD1306Driver{
		Driver:        NewDriver(c, "SSD1306", ssd1306DefaultAddress),
		displayHeight: ssd1306Height,
		displayWidth:  ssd1306Width,
		externalVCC:   ssd1306ExternalVCC,
	}
	s.afterStart = s.initialize

	// set options
	for _, option := range options {
		option(s)
	}
	// set page size
	s.pageSize = 8
	// set display buffer
	s.buffer = NewDisplayBuffer(s.displayWidth, s.displayHeight, s.pageSize)
	// add commands
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
		s.Clear()
		return map[string]interface{}{}
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

// WithSSD1306DisplayWidth option sets the SSD1306Driver DisplayWidth option.
func WithSSD1306DisplayWidth(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*SSD1306Driver)
		if ok {
			d.displayWidth = val
		}
	}
}

// WithSSD1306DisplayHeight option sets the SSD1306Driver DisplayHeight option.
func WithSSD1306DisplayHeight(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*SSD1306Driver)
		if ok {
			d.displayHeight = val
		}
	}
}

// WithSSD1306ExternalVCC option sets the SSD1306Driver ExternalVCC option.
func WithSSD1306ExternalVCC(val bool) func(Config) {
	return func(c Config) {
		d, ok := c.(*SSD1306Driver)
		if ok {
			d.externalVCC = val
		}
	}
}

// Init initializes the ssd1306 display.
func (s *SSD1306Driver) Init() error {
	// turn off screen
	if err := s.Off(); err != nil {
		return err
	}
	// run through initialization commands
	if err := s.commands(s.initSequence.GetSequence()); err != nil {
		return err
	}
	if err := s.commands([]byte{ssd1306ColumnAddr, 0, byte(s.buffer.width) - 1}); err != nil {
		return err
	}

	return s.commands([]byte{ssd1306PageAddr, 0, (byte(s.buffer.height / s.pageSize)) - 1})
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
func (s *SSD1306Driver) Clear() {
	s.buffer.Clear()
}

// Set sets a pixel in the buffer.
func (s *SSD1306Driver) Set(x, y, c int) {
	s.buffer.SetPixel(x, y, c)
}

// Reset clears display.
func (s *SSD1306Driver) Reset() error {
	if err := s.Off(); err != nil {
		return err
	}
	s.Clear()

	return s.On()
}

// SetContrast sets the display contrast.
func (s *SSD1306Driver) SetContrast(contrast byte) error {
	return s.commands([]byte{ssd1306SetContrast, contrast})
}

// Display sends the memory buffer to the display.
func (s *SSD1306Driver) Display() error {
	_, err := s.connection.Write(append([]byte{0x40}, s.buffer.buffer...))
	return err
}

// ShowImage takes a standard Go image and displays it in monochrome.
func (s *SSD1306Driver) ShowImage(img image.Image) error {
	if img.Bounds().Dx() != s.displayWidth || img.Bounds().Dy() != s.displayHeight {
		return fmt.Errorf("image must match display width and height: %dx%d", s.displayWidth, s.displayHeight)
	}
	s.Clear()
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

// command sends a command to the ssd1306
func (s *SSD1306Driver) command(b byte) error {
	_, err := s.connection.Write([]byte{0x80, b})
	return err
}

// commands sends a command sequence to the ssd1306
func (s *SSD1306Driver) commands(commands []byte) error {
	var command []byte
	for _, d := range commands {
		command = append(command, []byte{0x80, d}...)
	}
	_, err := s.connection.Write(command)
	return err
}

func (s *SSD1306Driver) initialize() error {
	// check device size for supported resolutions
	switch {
	case s.displayWidth == 128 && s.displayHeight == 64:
		s.initSequence = ssd1306Init128x64
	case s.displayWidth == 128 && s.displayHeight == 32:
		s.initSequence = ssd1306Init128x32
	case s.displayWidth == 96 && s.displayHeight == 16:
		s.initSequence = ssd1306Init96x16
	default:
		return fmt.Errorf("%dx%d resolution is unsupported, supported resolutions: 128x64, 128x32, 96x16",
			s.displayWidth, s.displayHeight)
	}
	// check for external vcc
	if s.externalVCC {
		s.initSequence.chargePumpSetting = 0x10
		s.initSequence.contrast = 0x9F
		s.initSequence.prechargePeriod = 0x22
	}
	if err := s.Init(); err != nil {
		return err
	}

	return s.On()
}
