package i2c

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*MPR121Driver)(nil)

const (
	// The device has a configurable I2C address by connecting the ADDR pin to the VSS, VDD, SDA or SCL
	// lines This results in I2C addresses of 0x5A, 0x5B, 0x5C and 0x5D
	// The grove i2c touch sensor is setup to use 0x5A
	mpr121address = 0x5A

	// MPR121 Registers
	// Registers 0x2B ~ 0x7F are control and configuration registers which need to be correctly configured before any capacitance
	// measurement and touch detection.

	// MHD_R MHD Rising
	MHD_R = 0x2B
	// NHD_R NHD Amount Rising
	NHD_R = 0x2C
	// NCL_R NCL Rising
	NCL_R = 0x2D
	//FDL_R FDL Rising
	FDL_R = 0x2E
	// MHD_F MHD Falling
	MHD_F = 0x2F
	// NHD_F NHD Amount Falling
	NHD_F       = 0x30
	NCL_F       = 0x31
	FDL_F       = 0x32
	ELE0_T      = 0x41
	ELE0_R      = 0x42
	ELE1_T      = 0x43
	ELE1_R      = 0x44
	ELE2_T      = 0x45
	ELE2_R      = 0x46
	ELE3_T      = 0x47
	ELE3_R      = 0x48
	ELE4_T      = 0x49
	ELE4_R      = 0x4A
	ELE5_T      = 0x4B
	ELE5_R      = 0x4C
	ELE6_T      = 0x4D
	ELE6_R      = 0x4E
	ELE7_T      = 0x4F
	ELE7_R      = 0x50
	ELE8_T      = 0x51
	ELE8_R      = 0x52
	ELE9_T      = 0x53
	ELE9_R      = 0x54
	ELE10_T     = 0x55
	ELE10_R     = 0x56
	ELE11_T     = 0x57
	ELE11_R     = 0x58
	FIL_CFG     = 0x5D
	ELE_CFG     = 0x5E
	GPIO_CTRL0  = 0x73
	GPIO_CTRL1  = 0x74
	GPIO_DATA   = 0x75
	GPIO_DIR    = 0x76
	GPIO_EN     = 0x77
	GPIO_SET    = 0x78
	GPIO_CLEAR  = 0x79
	GPIO_TOGGLE = 0x7A
	ATO_CFG0    = 0x7B
	ATO_CFGU    = 0x7D
	ATO_CFGL    = 0x7E
	ATO_CFGT    = 0x7F

	// Global Constants
	TOU_THRESH = 0x0F
	REL_THRESH = 0x0A
)

// MPR121Driver is a driver for a proximity capacity touch sensor controller
// http://www.seeedstudio.com/wiki/images/b/b9/Freescale_Semiconductor%3BMPR121QR2.pdf
type MPR121Driver struct {
	name       string
	connection I2c
	interval   time.Duration
	gobot.Eventer
	touchStates []bool
}

// NewMPR121Driver creates a new driver with specified name and i2c interface.
// Note that the MPR driver needs to have its interrupt port connected to another
// pin and monitored to trigger the reading of its values.
func NewMPR121Driver(a I2c, name string, v ...time.Duration) *MPR121Driver {
	m := &MPR121Driver{
		name:        name,
		connection:  a,
		Eventer:     gobot.NewEventer(),
		interval:    10 * time.Millisecond,
		touchStates: make([]bool, 12),
	}

	if len(v) > 0 {
		m.interval = v[0]
	}
	m.AddEvent(Error)
	return m
}

func (h *MPR121Driver) Name() string                 { return h.name }
func (h *MPR121Driver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start initalizes and start the driver
func (h *MPR121Driver) Start() (errs []error) {
	if err := h.connection.I2cStart(mpr121address); err != nil {
		return []error{err}
	}

	registers := [][2]byte{
		// Section A - Controls filtering when data is > baseline.
		[2]byte{MHD_R, 0x01},
		[2]byte{NHD_R, 0x01},
		[2]byte{NCL_R, 0x00},
		[2]byte{FDL_R, 0x00},

		// Section B - Controls filtering when data is < baseline.
		[2]byte{MHD_F, 0x01},
		[2]byte{NHD_F, 0x01},
		[2]byte{NCL_F, 0xFF},
		[2]byte{FDL_F, 0x02},

		// Section C - Sets touch and release thresholds for each electrode
		[2]byte{ELE0_T, TOU_THRESH},
		[2]byte{ELE0_R, REL_THRESH},
		[2]byte{ELE1_T, TOU_THRESH},
		[2]byte{ELE1_R, REL_THRESH},
		[2]byte{ELE2_T, TOU_THRESH},
		[2]byte{ELE2_R, REL_THRESH},
		[2]byte{ELE3_T, TOU_THRESH},
		[2]byte{ELE3_R, REL_THRESH},
		[2]byte{ELE4_T, TOU_THRESH},
		[2]byte{ELE4_R, REL_THRESH},
		[2]byte{ELE5_T, TOU_THRESH},
		[2]byte{ELE5_R, REL_THRESH},
		[2]byte{ELE6_T, TOU_THRESH},
		[2]byte{ELE6_R, REL_THRESH},
		[2]byte{ELE7_T, TOU_THRESH},
		[2]byte{ELE7_R, REL_THRESH},
		[2]byte{ELE8_T, TOU_THRESH},
		[2]byte{ELE8_R, REL_THRESH},
		[2]byte{ELE9_T, TOU_THRESH},
		[2]byte{ELE9_R, REL_THRESH},
		[2]byte{ELE10_T, TOU_THRESH},
		[2]byte{ELE10_R, REL_THRESH},
		[2]byte{ELE11_T, TOU_THRESH},
		[2]byte{ELE11_R, REL_THRESH},

		// Section D
		// Set the Filter Configuration
		// Set ESI2
		[2]byte{FIL_CFG, 0x04},

		// Section E
		// Electrode Configuration
		// Set ELE_CFG to 0x00 to return to standby mode
		[2]byte{ELE_CFG, 0x0C}, // Enables all 12 Electrodes
	}

	for _, reg := range registers {
		if err := h.connection.I2cWrite(mpr121address, []byte{reg[0], reg[1]}); err != nil {
			return []error{err}
		}
	}

	return nil
}

func (h *MPR121Driver) TouchValues() ([]bool, error) {
	var lsb, msb byte
	var touched uint16

	vals := make([]bool, 12)
	ret, err := h.connection.I2cRead(mpr121address, 2)
	if err != nil {
		gobot.Publish(h.Event(Error), err)
		return nil, err
	}
	if len(ret) != 2 {
		return vals, fmt.Errorf("unexpected i2c value returned")
	}
	buf := bytes.NewBuffer(ret)
	binary.Read(buf, binary.BigEndian, &lsb)
	binary.Read(buf, binary.BigEndian, &msb)
	if lsb != 0x0 {
		fmt.Printf("%X %X", lsb, msb)
	}
	touched = (uint16(msb) << 8) | uint16(lsb)
	// Check what electrodes were pressed
	for i := 0; i < 12; i++ {
		vals[i] = (touched&(1<<uint(i)) == 1)
	}

	return vals, err
}

// Halt returns true if devices is halted successfully
func (h *MPR121Driver) Halt() (err []error) { return }
