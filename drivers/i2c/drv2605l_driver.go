package i2c

import (
	"gobot.io/x/gobot"
)

// DRV2605Mode - operating mode
type DRV2605Mode uint8

// Operating modes, for use in SetMode()
const (
	DRV2605ModeIntTrig     DRV2605Mode = 0x00
	DRV2605ModeExtTrigEdge             = 0x01
	DRV2605ModeExtTrigLvl              = 0x02
	DRV2605ModePWMAnalog               = 0x03
	DRV2605ModeAudioVibe               = 0x04
	DRV2605ModeRealtime                = 0x05
	DRV2605ModeDiagnose                = 0x06
	DRV2605ModeAutocal                 = 0x07
)

const (
	drv2605Address = 0x5A

	drv2605RegStatus = 0x00
	drv2605RegMode   = 0x01

	drv2605Standby = 0x40

	drv2605RegRTPin    = 0x02
	drv2605RegLibrary  = 0x03
	drv2605RegWaveSeq1 = 0x04
	drv2605RegWaveSeq2 = 0x05
	drv2605RegWaveSeq3 = 0x06
	drv2605RegWaveSeq4 = 0x07
	drv2605RegWaveSeq5 = 0x08
	drv2605RegWaveSeq6 = 0x09
	drv2605RegWaveSeq7 = 0x0A
	drv2605RegWaveSeq8 = 0x0B

	drv2605RegGo            = 0x0C
	drv2605RegOverdrive     = 0x0D
	drv2605RegSustainPos    = 0x0E
	drv2605RegSustainNeg    = 0x0F
	drv2605RegBreak         = 0x10
	drv2605RegAudioCtrl     = 0x11
	drv2605RegAudioMinLevel = 0x12
	drv2605RegAudioMaxLevel = 0x13
	drv2605RegAudioMinDrive = 0x14
	drv2605RegAudioMaxDrive = 0x15
	drv2605RegRatedV        = 0x16
	drv2605RegClampV        = 0x17
	drv2605RegAutocalComp   = 0x18
	drv2605RegAutocalEmp    = 0x19
	drv2605RegFeedback      = 0x1A
	drv2605RegControl1      = 0x1B
	drv2605RegControl2      = 0x1C
	drv2605RegControl3      = 0x1D
	drv2605RegControl4      = 0x1E
	drv2605RegVBat          = 0x21
	drv2605RegLRAResoPeriod = 0x22
)

// DRV2605LDriver is the gobot driver for the TI/Adafruit DRV2605L Haptic Controller
//
// Device datasheet: http://www.ti.com/lit/ds/symlink/drv2605l.pdf
//
// Inspired by the Adafruit Python driver by Sean Mealin.
//
// Basic use:
//
//  haptic := i2c.NewDRV2605Driver(adaptor)
//  haptic.SetSequence([]byte{1, 13})
//  haptic.Go()
//
type DRV2605LDriver struct {
	name       string
	connector  Connector
	connection Connection
	Config
}

// NewDRV2605LDriver creates a new driver for the DRV2605L device.
//
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewDRV2605LDriver(conn Connector, options ...func(Config)) *DRV2605LDriver {
	driver := &DRV2605LDriver{
		name:      gobot.DefaultName("DRV2605L"),
		connector: conn,
		Config:    NewConfig(),
	}

	for _, option := range options {
		option(driver)
	}

	return driver
}

// Name returns the name of the device.
func (d *DRV2605LDriver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *DRV2605LDriver) SetName(name string) {
	d.name = name
}

// Connection returns the connection of the device.
func (d *DRV2605LDriver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the device.
func (d *DRV2605LDriver) Start() (err error) {
	if err := d.initialize(); err != nil {
		return err
	}
	return nil
}

func (d *DRV2605LDriver) writeByteRegisters(regValPairs []struct{ reg, val uint8 }) (err error) {
	for _, rv := range regValPairs {
		if err = d.connection.WriteByteData(rv.reg, rv.val); err != nil {
			break
		}
	}
	return err
}

func (d *DRV2605LDriver) initialize() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(drv2605Address)

	d.connection, err = d.connector.GetConnection(address, bus)
	if err != nil {
		return
	}

	feedback, err := d.connection.ReadByteData(drv2605RegFeedback)
	if err != nil {
		return
	}

	control, err := d.connection.ReadByteData(drv2605RegControl3)
	if err != nil {
		return
	}

	err = d.writeByteRegisters([]struct{ reg, val uint8 }{
		// leave standby, enter "internal trig" mode
		{drv2605RegMode, 0},
		// turn off real-time play
		{drv2605RegRTPin, 0},
		// init wave sequencer with "strong click"
		{drv2605RegWaveSeq1, 1},
		{drv2605RegWaveSeq1, 0},
		// no physical parameter tweaks
		{drv2605RegSustainPos, 0},
		{drv2605RegSustainNeg, 0},
		{drv2605RegBreak, 0},
		// set up ERM open loop
		{drv2605RegFeedback, feedback & 0x7f},
		{drv2605RegControl3, control | 0x20},
	})

	return
}

// SetMode sets the device in one of the eight modes as described in the
// datasheet. Defaults to mode 0, internal trig.
func (d *DRV2605LDriver) SetMode(newMode DRV2605Mode) (err error) {
	mode, err := d.connection.ReadByteData(drv2605RegMode)
	if err != nil {
		return err
	}

	// clear mode bits (lower three bits)
	mode &= 0xf8
	// set new mode bits
	mode |= uint8(newMode)

	err = d.connection.WriteByteData(drv2605RegMode, mode)

	return err
}

// SetStandbyMode controls device low power mode
func (d *DRV2605LDriver) SetStandbyMode(standby bool) (err error) {
	modeVal, err := d.connection.ReadByteData(drv2605RegMode)
	if err != nil {
		return err
	}
	if standby {
		modeVal |= drv2605Standby
	} else {
		modeVal &= 0xFF ^ drv2605Standby
	}

	err = d.connection.WriteByteData(drv2605RegMode, modeVal)

	return err
}

// SelectLibrary selects which waveform library to play from, 1-7.
// See datasheet for more info.
func (d *DRV2605LDriver) SelectLibrary(library uint8) (err error) {
	err = d.connection.WriteByteData(drv2605RegLibrary, library&0x7)
	return err
}

// GetPauseWaveform returns a special waveform ID used in SetSequence() to encode
// pauses between waveforms. Time is specified in tens of milliseconds
// ranging from 0ms (delayTime10MS = 0) to 1270ms (delayTime10MS = 127).
// Times out of range are clipped to fit.
func (d *DRV2605LDriver) GetPauseWaveform(delayTime10MS uint8) (pauseID uint8) {
	if delayTime10MS > 127 {
		delayTime10MS = 127
	}

	return delayTime10MS | 0x80
}

// SetSequence sets the sequence of waveforms to be played by the sequencer,
// specified by waveform id as described in the datasheet.
// The sequencer can play at most 8 waveforms in sequence, longer
// sequences will be truncated.
// A waveform id of zero marks the end of the sequence.
// Pauses can be encoded using GetPauseWaveform().
func (d *DRV2605LDriver) SetSequence(waveforms []uint8) (err error) {
	if len(waveforms) < 8 {
		waveforms = append(waveforms, 0)
	}
	if len(waveforms) > 8 {
		waveforms = waveforms[0:8]
	}
	for i, w := range waveforms {
		if err = d.connection.WriteByteData(uint8(drv2605RegWaveSeq1+i), w); err != nil {
			return err
		}
	}
	return nil
}

// Go plays the current sequence of waveforms.
func (d *DRV2605LDriver) Go() (err error) {
	err = d.connection.WriteByteData(drv2605RegGo, 1)
	return err
}

// Halt halts the device.
func (d *DRV2605LDriver) Halt() (err error) {
	if d.connection != nil {
		// stop playback
		if err = d.connection.WriteByteData(drv2605RegGo, 0); err != nil {
			return err
		}

		// enter standby
		return d.SetStandbyMode(true)
	}
	return
}
