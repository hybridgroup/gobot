package mfrc522

import (
	"fmt"
	"log"
	"time"
)

// PCD: proximity coupling device (reader unit)
// PICC: proximity integrated circuit card (the card or chip)

const (
	initTime      = 50 * time.Millisecond
	antennaOnTime = 4 * time.Millisecond
)

type busConnection interface {
	WriteByteData(reg byte, data byte) error
	ReadByteData(reg byte) (byte, error)
}

var versions = map[uint8]string{0x12: "Counterfeit", 0x88: "FM17522", 0x89: "FM17522E",
	0x90: "MFRC522 0.0", 0x91: "MFRC522 1.0", 0x92: "MFRC522 2.0",
	0xA3: "RFID-RC522 3.0", 0xB2: "FM17522 1",
}

// datasheet:
// https://www.nxp.com/docs/en/data-sheet/MFRC522.pdf
//
// reference implementations:
// * https://github.com/OSSLibraries/Arduino_MFRC522v2
// * https://periph.io/device/mf-rc522/
type MFRC522Common struct {
	connection busConnection
}

// NewMFRC522Common creates a new Gobot Driver for MFRC522 RFID with specified bus connection
// The device supports SPI, I2C and UART (not implemented yet at gobot system level).
//
// Params:
//      c BusConnection - the bus connection to use with this driver
func NewMFRC522Common() *MFRC522Common {
	d := &MFRC522Common{}
	return d
}

func (d *MFRC522Common) Connect(c busConnection) {
	d.connection = c
}

func (d *MFRC522Common) Check() error {
	log.Println("nice to read you")
	return nil
}

func (d *MFRC522Common) Initialize() error {
	if d.connection == nil {
		return fmt.Errorf("not connected")
	}
	if err := d.softReset(); err != nil {
		return err
	}

	initSequence := [][]byte{
		{mfrc522RegTxMode, mfrc522RxTxMode_Reset},
		{mfrc522RegRxMode, mfrc522RxTxMode_Reset},
		{mfrc522RegModWidth, mfrc522ModWidth_Reset},
		{mfrc522RegTMode, mfrc522TMode_TAutoBit | 0x0F}, // timer starts automatically at the end of the transmission
		{mfrc522RegTPrescaler, 0xFF},                    // together with the 0x0F above, we have 604us
		{mfrc522RegTReloadL, mfrc522RegTReload_Value25ms & 0xFF},
		{mfrc522RegTReloadH, mfrc522RegTReload_Value25ms >> 8},
		{mfrc522RegTxASK, mfrc522TxASK_Force100ASKBit},
		{mfrc522RegMode, mfrc522Mode_TxWaitRFBit | mfrc522Mode_PolMFinBit | mfrc522Mode_CRCPreset6363},
	}
	for _, init := range initSequence {
		d.connection.WriteByteData(init[0], init[1])
	}

	if err := d.switchAntenna(true); err != nil {
		return err
	}

	version, err := d.getVersion()
	if err != nil {
		return err
	}

	fmt.Printf("Initialized version: '%s' (0x%X)\n", versions[version], version)

	return nil
}

func (d *MFRC522Common) getVersion() (uint8, error) {
	return d.connection.ReadByteData(mfrc522RegVersion)
}

func (d *MFRC522Common) switchAntenna(targetState bool) error {
	val, err := d.connection.ReadByteData(mfrc522RegTxControl)
	if err != nil {
		return err
	}
	maskForOn := uint8(mfrc522TxControl_Tx2RFEn1outputBit | mfrc522TxControl_Tx1RFEn1outputBit)
	currentState := val&maskForOn == maskForOn

	if targetState == currentState {
		return nil
	}

	if targetState {
		val = val | maskForOn
	} else {
		val = val & maskForOn
	}

	if err := d.connection.WriteByteData(mfrc522RegTxControl, val); err != nil {
		return err
	}

	if targetState {
		time.Sleep(antennaOnTime)
	}

	return nil
}

func (d *MFRC522Common) softReset() error {
	if err := d.connection.WriteByteData(mfrc522RegCommand, mfrc522PCD_SoftReset); err != nil {
		return err
	}

	// The datasheet does not mention how long the SoftReset command takes to complete. According to section 8.8.2 of the
	// datasheet the oscillator start-up time is the start up time of the crystal + 37.74 us.
	// TODO: this can be done better by wait until the power down bit is cleared
	time.Sleep(initTime)
	val, err := d.connection.ReadByteData(mfrc522RegCommand)
	if err != nil {
		return err
	}

	if val&mfrc522PCD_PowerDownBit > 1 {
		return fmt.Errorf("initialization takes longer than %s", initTime)
	}
	return nil
}
