package mfrc522

import (
	"fmt"
	"time"
)

// PCD: proximity coupling device (reader unit)
// PICC: proximity integrated circuit card (the card or chip)

const (
	pcdDebug = false
	initTime = 50 * time.Millisecond
	// at least 5 ms are needed after switch on, see AN10834
	antennaOnTime = 10 * time.Millisecond
)

type busConnection interface {
	ReadByteData(reg byte) (byte, error)
	WriteByteData(reg byte, data byte) error
}

var versions = map[uint8]string{
	0x12: "Counterfeit", 0x88: "FM17522", 0x89: "FM17522E",
	0x90: "MFRC522 0.0", 0x91: "MFRC522 1.0", 0x92: "MFRC522 2.0", 0xB2: "FM17522 1",
}

// MFRC522Common is the Gobot Driver for MFRC522 RFID.
// datasheet:
// https://www.nxp.com/docs/en/data-sheet/MFRC522.pdf
//
// reference implementations:
// * https://github.com/OSSLibraries/ArduinoRegMFRC522v2
// * https://github.com/jdevelop/golang-rpi-extras
// * https://github.com/pimylifeup/MFRC522-python
// * https://periph.io/device/mf-rc522/
type MFRC522Common struct {
	connection      busConnection
	firstCardAccess bool
}

// NewMFRC522Common creates a new Gobot Driver for MFRC522 RFID with specified bus connection
// The device supports SPI, I2C and UART (not implemented yet at gobot system level).
//
// Params:
//
//	c BusConnection - the bus connection to use with this driver
func NewMFRC522Common() *MFRC522Common {
	d := &MFRC522Common{}
	return d
}

// Initialize sets the connection and initializes the driver.
func (d *MFRC522Common) Initialize(c busConnection) error {
	d.connection = c

	if err := d.softReset(); err != nil {
		return err
	}

	initSequence := [][]byte{
		{regTxMode, rxtxModeRegReset},
		{regRxMode, rxtxModeRegReset},
		{regModWidth, modWidthRegReset},
		{regTMode, tModeRegTAutoBit | 0x0F}, // timer starts automatically at the end of the transmission
		{regTPrescaler, 0xFF},
		{regTReloadL, tReloadRegValue25ms & 0xFF},
		{regTReloadH, tReloadRegValue25ms >> 8},
		{regTxASK, txASKRegForce100ASKBit},
		{regMode, modeRegTxWaitRFBit | modeRegPolMFinBit | modeRegCRCPreset6363},
	}
	for _, init := range initSequence {
		if err := d.writeByteData(init[0], init[1]); err != nil {
			return err
		}
	}

	if err := d.switchAntenna(true); err != nil {
		return err
	}
	if err := d.setAntennaGain(rfcCfgRegRxGain38dB); err != nil {
		return err
	}

	return nil
}

// PrintReaderVersion gets and prints the reader (pcd) version.
func (d *MFRC522Common) PrintReaderVersion() error {
	version, err := d.getVersion()
	if err != nil {
		return err
	}
	fmt.Printf("PCD version: '%s' (0x%X)\n", versions[version], version)
	return nil
}

func (d *MFRC522Common) getVersion() (uint8, error) {
	return d.readByteData(regVersion)
}

func (d *MFRC522Common) switchAntenna(targetState bool) error {
	val, err := d.readByteData(regTxControl)
	if err != nil {
		return err
	}
	maskForOn := uint8(txControlRegTx2RFEn1outputBit | txControlRegTx1RFEn1outputBit)
	currentState := val&maskForOn == maskForOn

	if targetState == currentState {
		return nil
	}

	if targetState {
		val = val | maskForOn
	} else {
		val = val & ^maskForOn
	}

	if err := d.writeByteData(regTxControl, val); err != nil {
		return err
	}

	if targetState {
		time.Sleep(antennaOnTime)
	}

	return nil
}

func (d *MFRC522Common) setAntennaGain(val uint8) error {
	return d.writeByteData(regRFCfg, val)
}

func (d *MFRC522Common) softReset() error {
	if err := d.writeByteData(regCommand, commandRegSoftReset); err != nil {
		return err
	}

	// The datasheet does not mention how long the SoftReset command takes to complete. According to section 8.8.2 of the
	// datasheet the oscillator start-up time is the start up time of the crystal + 37.74 us.
	// TODO: this can be done better by wait until the power down bit is cleared
	time.Sleep(initTime)
	val, err := d.readByteData(regCommand)
	if err != nil {
		return err
	}

	if val&commandRegPowerDownBit > 1 {
		return fmt.Errorf("initialization takes longer than %s", initTime)
	}
	return nil
}

func (d *MFRC522Common) stopCrypto1() error {
	return d.clearRegisterBitMask(regStatus2, status2RegMFCrypto1OnBit)
}

func (d *MFRC522Common) communicateWithPICC(command uint8, sendData []byte, backData []byte, txLastBits uint8,
	checkCRC bool,
) error {
	irqEn := 0x00
	waitIRq := uint8(0x00)
	switch command {
	case commandRegMFAuthent:
		irqEn = comIEnRegIdleIEnBit | comIEnRegErrIEnBit
		waitIRq = comIrqRegIdleIRqBit
	case commandRegTransceive:
		irqEn = comIEnRegTimerIEnBit | comIEnRegErrIEnBit | comIEnRegLoAlertIEnBit
		irqEn = irqEn | comIEnRegIdleIEnBit | comIEnRegRxIEnBit | comIEnRegTxIEnBit
		waitIRq = uint8(comIrqRegIdleIRqBit | comIrqRegRxIRqBit)
	}

	// TODO: this is not used at the moment (propagation of IRQ pin)
	//nolint:gosec // TODO: fix later
	if err := d.writeByteData(regComIEn, uint8(irqEn|comIEnRegIRqInv)); err != nil {
		return err
	}
	if err := d.writeByteData(regComIrq, comIrqRegClearAll); err != nil {
		return err
	}
	if err := d.writeByteData(regFIFOLevel, fifoLevelRegFlushBufferBit); err != nil {
		return err
	}
	// stop any active command
	if err := d.writeByteData(regCommand, commandRegIdle); err != nil {
		return err
	}
	// prepare and start communication
	if err := d.writeFifo(sendData); err != nil {
		return err
	}
	if err := d.writeByteData(regBitFraming, txLastBits); err != nil {
		return err
	}
	if err := d.writeByteData(regCommand, command); err != nil {
		return err
	}
	if command == commandRegTransceive {
		if err := d.setRegisterBitMask(regBitFraming, bitFramingRegStartSendBit); err != nil {
			return err
		}
	}

	// Wait for the command to complete. On initialization the TAuto flag in TMode register is set. This means the timer
	// automatically starts when the PCD stops transmitting.
	const maxTries = 5
	i := 0
	for ; i < maxTries; i++ {
		irqs, err := d.readByteData(regComIrq)
		if err != nil {
			return err
		}
		if irqs&waitIRq > 0 {
			// One of the interrupts that signal success has been set.
			break
		}
		if irqs&comIrqRegTimerIRqBit == comIrqRegTimerIRqBit {
			return fmt.Errorf("the timer interrupt occurred")
		}
		time.Sleep(time.Millisecond)
	}

	if err := d.clearRegisterBitMask(regBitFraming, bitFramingRegStartSendBit); err != nil {
		return err
	}

	if i >= maxTries {
		return fmt.Errorf("no data available after %d tries", maxTries)
	}

	errorRegValue, err := d.readByteData(regError)
	if err != nil {
		return err
	}
	// stop if any errors except collisions were detected.
	if err := d.getFirstError(errorRegValue &^ errorRegCollErrBit); err != nil {
		return err
	}

	backLen := len(backData)
	var rxLastBits uint8
	if backLen > 0 {
		rxLastBits, err = d.readFifo(backData)
		if err != nil {
			return err
		}
		if pcdDebug {
			fmt.Printf("rxLastBits: 0x%02x\n", rxLastBits)
		}
	}

	if err := d.getFirstError(errorRegValue & errorRegCollErrBit); err != nil {
		return err
	}

	if backLen > 2 && checkCRC {
		// the last 2 bytes of data contains the CRC
		if backLen == 3 && rxLastBits == 0x04 {
			return fmt.Errorf("CRC: MIFARE Classic NAK is not OK")
		}
		if backLen < 4 || backLen == 4 && rxLastBits != 0 {
			return fmt.Errorf("CRC: at least the 2 byte CRCRegA value and all 8 bits of the last byte must be received")
		}

		crcResult := []byte{0x00, 0x00}
		crcData := backData[:backLen-2]
		dataCrc := backData[backLen-2:]
		err := d.calculateCRC(crcData, crcResult)
		if err != nil {
			return err
		}
		if dataCrc[0] != crcResult[0] || dataCrc[1] != crcResult[1] {
			return fmt.Errorf("CRC: values not match %v - %v", crcResult, dataCrc)
		}
	}

	return nil
}

// 16 bit CRC will be calculated for the given data
func (d *MFRC522Common) calculateCRC(data []byte, result []byte) error {
	// Stop any active command.
	if err := d.writeByteData(regCommand, commandRegIdle); err != nil {
		return err
	}
	if err := d.writeByteData(regDivIrq, divIrqRegCRCIRqBit); err != nil {
		return err
	}
	if err := d.writeByteData(regFIFOLevel, fifoLevelRegFlushBufferBit); err != nil {
		return err
	}
	if err := d.writeFifo(data); err != nil {
		return err
	}
	if err := d.writeByteData(regCommand, commandRegCalcCRC); err != nil {
		return err
	}

	const maxTries = 3
	for i := 0; i < maxTries; i++ {
		irqs, err := d.readByteData(regDivIrq)
		if err != nil {
			return err
		}
		if irqs&divIrqRegCRCIRqBit == divIrqRegCRCIRqBit {
			if err := d.writeByteData(regCommand, commandRegIdle); err != nil {
				return err
			}
			result[0], err = d.readByteData(regCRCResultL)
			if err != nil {
				return err
			}
			result[1], err = d.readByteData(regCRCResultH)
			if err != nil {
				return err
			}
			return nil
		}
		time.Sleep(time.Millisecond)
	}

	return fmt.Errorf("no CRC available after %d tries", maxTries)
}

func (d *MFRC522Common) writeFifo(fifoData []byte) error {
	// the register command is always the same, the pointer in FIFO is incremented automatically after each write
	for _, b := range fifoData {
		if err := d.writeByteData(regFIFOData, b); err != nil {
			return err
		}
	}
	return nil
}

func (d *MFRC522Common) readFifo(backData []byte) (uint8, error) {
	n, err := d.readByteData(regFIFOLevel) // Number of bytes in the FIFO
	if err != nil {
		return 0, err
	}
	//nolint:gosec // TODO: fix later
	if n > uint8(len(backData)) {
		return 0, fmt.Errorf("more data in FIFO (%d) than expected (%d)", n, len(backData))
	}
	//nolint:gosec // TODO: fix later
	if n < uint8(len(backData)) {
		return 0, fmt.Errorf("less data in FIFO (%d) than expected (%d)", n, len(backData))
	}

	// the register command is always the same, the pointer in FIFO is incremented automatically after each read
	for i := 0; i < int(n); i++ {
		byteVal, err := d.readByteData(regFIFOData)
		if err != nil {
			return 0, err
		}
		backData[i] = byteVal
	}

	rxLastBits, err := d.readByteData(regControl)
	if err != nil {
		return 0, err
	}
	return rxLastBits & controlRegRxLastBits, nil
}

func (d *MFRC522Common) getFirstError(errorRegValue uint8) error {
	if errorRegValue == 0 {
		return nil
	}

	if errorRegValue&errorRegProtocolErrBit == errorRegProtocolErrBit {
		return fmt.Errorf("a protocol error occurred")
	}
	if errorRegValue&errorRegParityErrBit == errorRegParityErrBit {
		return fmt.Errorf("a parity error occurred")
	}
	if errorRegValue&errorRegCRCErrBit == errorRegCRCErrBit {
		return fmt.Errorf("a CRC error occurred")
	}
	if errorRegValue&errorRegCollErrBit == errorRegCollErrBit {
		return fmt.Errorf("a collision error occurred")
	}
	if errorRegValue&errorRegBufferOvflBit == errorRegBufferOvflBit {
		return fmt.Errorf("a buffer overflow error occurred")
	}
	if errorRegValue&errorRegTempErrBit == errorRegTempErrBit {
		return fmt.Errorf("a temperature error occurred")
	}
	if errorRegValue&errorRegWrErrBit == errorRegWrErrBit {
		return fmt.Errorf("a temperature error occurred")
	}
	return fmt.Errorf("an unknown error occurred")
}
