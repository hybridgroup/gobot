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

var versions = map[uint8]string{0x12: "Counterfeit", 0x88: "FM17522", 0x89: "FM17522E",
	0x90: "MFRC522 0.0", 0x91: "MFRC522 1.0", 0x92: "MFRC522 2.0", 0xB2: "FM17522 1"}

// datasheet:
// https://www.nxp.com/docs/en/data-sheet/MFRC522.pdf
//
// reference implementations:
// * https://github.com/OSSLibraries/Arduino_MFRC522v2
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
//      c BusConnection - the bus connection to use with this driver
func NewMFRC522Common() *MFRC522Common {
	d := &MFRC522Common{}
	return d
}

func (d *MFRC522Common) Initialize(c busConnection) error {
	d.connection = c

	if err := d.softReset(); err != nil {
		return err
	}

	initSequence := [][]byte{
		{RegTxMode, RxTxMode_Reset},
		{RegRxMode, RxTxMode_Reset},
		{RegModWidth, ModWidth_Reset},
		{RegTMode, TMode_TAutoBit | 0x0F}, // timer starts automatically at the end of the transmission
		{RegTPrescaler, 0xFF},
		{RegTReloadL, TReload_Value25ms & 0xFF},
		{RegTReloadH, TReload_Value25ms >> 8},
		{RegTxASK, TxASK_Force100ASKBit},
		{RegMode, Mode_TxWaitRFBit | Mode_PolMFinBit | Mode_CRCPreset6363},
	}
	for _, init := range initSequence {
		d.writeByteData(init[0], init[1])
	}

	if err := d.switchAntenna(true); err != nil {
		return err
	}
	if err := d.setAntennaGain(RFCfg_RxGain38dB); err != nil {
		return err
	}

	return nil
}

func (d *MFRC522Common) PrintReaderVersion() error {
	version, err := d.getVersion()
	if err != nil {
		return err
	}
	fmt.Printf("PCD version: '%s' (0x%X)\n", versions[version], version)
	return nil
}

func (d *MFRC522Common) getVersion() (uint8, error) {
	return d.readByteData(RegVersion)
}

func (d *MFRC522Common) switchAntenna(targetState bool) error {
	val, err := d.readByteData(RegTxControl)
	if err != nil {
		return err
	}
	maskForOn := uint8(TxControl_Tx2RFEn1outputBit | TxControl_Tx1RFEn1outputBit)
	currentState := val&maskForOn == maskForOn

	if targetState == currentState {
		return nil
	}

	if targetState {
		val = val | maskForOn
	} else {
		val = val & maskForOn
	}

	if err := d.writeByteData(RegTxControl, val); err != nil {
		return err
	}

	if targetState {
		time.Sleep(antennaOnTime)
	}

	return nil
}

func (d *MFRC522Common) setAntennaGain(val uint8) error {
	return d.writeByteData(RegRFCfg, val)
}

func (d *MFRC522Common) softReset() error {
	if err := d.writeByteData(RegCommand, Command_SoftReset); err != nil {
		return err
	}

	// The datasheet does not mention how long the SoftReset command takes to complete. According to section 8.8.2 of the
	// datasheet the oscillator start-up time is the start up time of the crystal + 37.74 us.
	// TODO: this can be done better by wait until the power down bit is cleared
	time.Sleep(initTime)
	val, err := d.readByteData(RegCommand)
	if err != nil {
		return err
	}

	if val&Command_PowerDownBit > 1 {
		return fmt.Errorf("initialization takes longer than %s", initTime)
	}
	return nil
}

func (d *MFRC522Common) stopCrypto1() error {
	return d.clearRegisterBitMask(RegStatus2, Status2_MFCrypto1OnBit)
}

func (d *MFRC522Common) communicateWithPICC(command uint8, sendData []byte, backData []byte, txLastBits uint8,
	checkCRC bool) error {
	irqEn := 0x00
	waitIRq := uint8(0x00)
	switch command {
	case Command_MFAuthent:
		irqEn = ComIEn_IdleIEnBit | ComIEn_ErrIEnBit
		waitIRq = ComIrq_IdleIRqBit
	case Command_Transceive:
		irqEn = ComIEn_TimerIEnBit | ComIEn_ErrIEnBit | ComIEn_LoAlertIEnBit
		irqEn = irqEn | ComIEn_IdleIEnBit | ComIEn_RxIEnBit | ComIEn_TxIEnBit
		waitIRq = uint8(ComIrq_IdleIRqBit | ComIrq_RxIRqBit)
	}

	// TODO: this is not used at the moment (propagation of IRQ pin)
	if err := d.writeByteData(RegComIEn, uint8(irqEn|ComIEn_IRqInv)); err != nil {
		return err
	}
	if err := d.writeByteData(RegComIrq, ComIrq_ClearAll); err != nil {
		return err
	}
	if err := d.writeByteData(RegFIFOLevel, FIFOLevel_FlushBufferBit); err != nil {
		return err
	}
	// stop any active command
	if err := d.writeByteData(RegCommand, Command_Idle); err != nil {
		return err
	}
	// prepare and start communication
	if err := d.writeFifo(sendData); err != nil {
		return err
	}
	if err := d.writeByteData(RegBitFraming, txLastBits); err != nil {
		return err
	}
	if err := d.writeByteData(RegCommand, command); err != nil {
		return err
	}
	if command == Command_Transceive {
		if err := d.setRegisterBitMask(RegBitFraming, BitFraming_StartSendBit); err != nil {
			return err
		}
	}

	// Wait for the command to complete. On initialization the TAuto flag in TMode register is set. This means the timer
	// automatically starts when the PCD stops transmitting.
	const maxTries = 5
	i := 0
	for ; i < maxTries; i++ {
		irqs, err := d.readByteData(RegComIrq)
		if err != nil {
			return err
		}
		if irqs&waitIRq > 0 {
			// One of the interrupts that signal success has been set.
			break
		}
		if irqs&ComIrq_TimerIRqBit == ComIrq_TimerIRqBit {
			return fmt.Errorf("the timer interrupt occurred")
		}
		time.Sleep(time.Millisecond)
	}

	if err := d.clearRegisterBitMask(RegBitFraming, BitFraming_StartSendBit); err != nil {
		return err
	}

	if i >= maxTries {
		return fmt.Errorf("no data available after %d tries", maxTries)
	}

	errorRegValue, err := d.readByteData(RegError)
	if err != nil {
		return err
	}
	// stop if any errors except collisions were detected.
	if err := d.getFirstError(errorRegValue &^ Error_CollErrBit); err != nil {
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

	if err := d.getFirstError(errorRegValue & Error_CollErrBit); err != nil {
		return err
	}

	if backLen > 2 && checkCRC {
		// the last 2 bytes of data contains the CRC
		if backLen == 3 && rxLastBits == 0x04 {
			return fmt.Errorf("CRC: MIFARE Classic NAK is not OK")
		}
		if backLen < 4 || backLen == 4 && rxLastBits != 0 {
			return fmt.Errorf("CRC: at least the 2 byte CRC_A value and all 8 bits of the last byte must be received")
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
	if err := d.writeByteData(RegCommand, Command_Idle); err != nil {
		return err
	}
	if err := d.writeByteData(RegDivIrq, DivIrq_CRCIRqBit); err != nil {
		return err
	}
	if err := d.writeByteData(RegFIFOLevel, FIFOLevel_FlushBufferBit); err != nil {
		return err
	}
	if err := d.writeFifo(data); err != nil {
		return err
	}
	if err := d.writeByteData(RegCommand, Command_CalcCRC); err != nil {
		return err
	}

	const maxTries = 3
	for i := 0; i < maxTries; i++ {
		irqs, err := d.readByteData(RegDivIrq)
		if err != nil {
			return err
		}
		if irqs&DivIrq_CRCIRqBit == DivIrq_CRCIRqBit {
			if err := d.writeByteData(RegCommand, Command_Idle); err != nil {
				return err
			}
			result[0], err = d.readByteData(RegCRCResultL)
			if err != nil {
				return err
			}
			result[1], err = d.readByteData(RegCRCResultH)
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
		if err := d.writeByteData(RegFIFOData, b); err != nil {
			return err
		}
	}
	return nil
}

func (d *MFRC522Common) readFifo(backData []byte) (uint8, error) {
	n, err := d.readByteData(RegFIFOLevel) // Number of bytes in the FIFO
	if n > uint8(len(backData)) {
		return 0, fmt.Errorf("more data in FIFO (%d) than expected (%d)", n, len(backData))
	}

	if n < uint8(len(backData)) {
		return 0, fmt.Errorf("less data in FIFO (%d) than expected (%d)", n, len(backData))
	}

	// the register command is always the same, the pointer in FIFO is incremented automatically after each read
	for i := 0; i < int(n); i++ {
		byteVal, err := d.readByteData(RegFIFOData)
		if err != nil {
			return 0, err
		}
		backData[i] = byteVal
	}

	rxLastBits, err := d.readByteData(RegControl)
	if err != nil {
		return 0, err
	}
	return rxLastBits & Control_RxLastBits, nil
}

func (d *MFRC522Common) getFirstError(errorRegValue uint8) error {
	if errorRegValue == 0 {
		return nil
	}

	if errorRegValue&Error_ProtocolErrBit == Error_ProtocolErrBit {
		return fmt.Errorf("a protocol error occurred")
	}
	if errorRegValue&Error_ParityErrBit == Error_ParityErrBit {
		return fmt.Errorf("a parity error occurred")
	}
	if errorRegValue&Error_CRCErrBit == Error_CRCErrBit {
		return fmt.Errorf("a CRC error occurred")
	}
	if errorRegValue&Error_CollErrBit == Error_CollErrBit {
		return fmt.Errorf("a collision error occurred")
	}
	if errorRegValue&Error_BufferOvflBit == Error_BufferOvflBit {
		return fmt.Errorf("a buffer overflow error occurred")
	}
	if errorRegValue&Error_TempErrBit == Error_TempErrBit {
		return fmt.Errorf("a temperature error occurred")
	}
	if errorRegValue&Error_WrErrBit == Error_WrErrBit {
		return fmt.Errorf("a temperature error occurred")
	}
	return fmt.Errorf("an unknown error occurred")
}
