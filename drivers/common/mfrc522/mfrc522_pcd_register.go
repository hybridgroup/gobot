//nolint:lll // ok here
package mfrc522

// Page 0: Command and status
const (
	//              0x00     // reserved for future use
	regCommand = 0x01 // starts and stops command execution
	// ------------ values --------------------
	// commands, see chapter 10.3, table 149 of the datasheet (only 4 lower bits are used for writing)
	commandRegIdle             = 0x00 // no action, cancels current command execution
	commandRegMem              = 0x01 // stores 25 bytes into the internal buffer
	commandRegGenerateRandomID = 0x02 // generates a 10-byte random ID number
	commandRegCalcCRC          = 0x03 // activates the CRC coprocessor or performs a self-test
	commandRegTransmit         = 0x04 // transmits data from the FIFO buffer
	// 0x05, 0x06 not used
	// commandRegNoCmdChange = 0x07 // no command change, can be used to modify the Command register bits without
	// commandRegReceive     = 0x08 // activates the receiver circuits
	// 0x09..0x0B not used
	commandRegTransceive = 0x0C // transmits data from FIFO buffer to antenna and automatically activates the receiver after transmission
	// 0x0D reserved
	commandRegMFAuthent = 0x0E // performs the MIFARE standard authentication as a reader
	commandRegSoftReset = 0x0F // resets the MFRC522
	// starts the wake up procedure during which this bit is read as a logic 1; it is read as a logic 0 when the
	// is ready; Remark: The PowerDown bit cannot be set when the SoftReset command is activated
	commandRegPowerDownBit = 0x10 // Soft power-down mode entered, if 1
	commandRegRcvOffBit    = 0x20 // analog part of the receiver is switched off, if 1
	// commandRegReserved67 = 0xC0
)

const (
	regComIEn = 0x02 // enable and disable the passing of interrupt requests to IRQ pin
	// ------------ values --------------------
	comIEnRegReset         = 0x80 // see table 25 of data sheet
	comIEnRegTimerIEnBit   = 0x01 // bit 0: allows the timer interrupt request (TimerIRq bit) to be propagated
	comIEnRegErrIEnBit     = 0x02 // bit 1: allows the error interrupt request (ErrIRq bit) to be propagated
	comIEnRegLoAlertIEnBit = 0x04 // bit 2: allows the low alert interrupt request (LoAlertIRq bit) to be propagated
	comIEnRegHiAlertIEnBit = 0x08 // bit 3: allows the high alert interrupt request (HiAlertIRq bit) to be propagated
	comIEnRegIdleIEnBit    = 0x10 // bit 4: allows the idle interrupt request (IdleIRq bit) to be propagated
	comIEnRegRxIEnBit      = 0x20 // bit 5: allows the receiver interrupt request (RxIRq bit) to be propagated
	comIEnRegTxIEnBit      = 0x40 // bit 6: allows the transmitter interrupt request (TxIRq bit) to be propagated
	// 1: signal on pin IRQ is inverted with respect to the Status1 register’s IRq bit
	// 0: signal on pin IRQ is equal to the IRq bit; in combination with the DivIEn register’s IRqPushPull bit, the
	// default value of logic 1 ensures that the output level on pin IRQ is 3-state
	comIEnRegIRqInv = 0x80 // bit 7: see above
)

const (
// ------------ unused commands --------------------
// regDivIEn = 0x03 // enable and disable the passing of interrupt requests to IRQ pin
)

const (
	regComIrq = 0x04 // interrupt request bits for communication
	// ------------ values --------------------
	comIrqRegReset         = 0x14 // see table 29 of data sheet
	comIrqRegClearAll      = 0x7F // all bits are set to clear, except the Set1Bit (0 indicates the reset)
	comIrqRegTimerIRqBit   = 0x01 // bit 0: the timer decrements the timer value in register TCounterVal to zero, if 1
	comIrqRegErrIRq1anyBit = 0x02 // bit 1: error bit in the Error register is set, if 1
	// Status1 register’s LoAlert bit is set in opposition to the LoAlert bit, the LoAlertIRq bit stores this event and
	// can only be reset as indicated by the Set1 bit in this register
	// comIrqRegLoAlertIRqBit = 0x04 // bit 2: if 1, see above
	// the Status1 register’s HiAlert bit is set in opposition to the HiAlert bit, the HiAlertIRq bit stores this event
	// and can only be reset as indicated by the Set1 bit in this register
	// comIrqRegHiAlertIRqBit = 0x08 // bit 3: if 1, see above
	// If a command terminates, for example, when the Command register changes its value from any command to Idle command.
	// If an unknown command is started, the Command register Command[3:0] value changes to the idle state and the
	// IdleIRq bit is set. The microcontroller starting the Idle command does not set the IdleIRq bit.
	comIrqRegIdleIRqBit = 0x10 // bit 4: if 1, see above
	// receiver has detected the end of a valid data stream, if the RxMode register’s RxNoErr bit is set to logic 1,
	// the RxIRq bit is only set to logic 1 when data bytes are available in the FIFO
	comIrqRegRxIRqBit = 0x20 // bit 5: if 1, see above
	comIrqRegTxIRqBit = 0x40 // bit 6: set to 1, immediately after the last bit of the transmitted data was sent out
	// 1: indicates that the marked bits in the register are set
	// 0: indicates that the marked bits in the register are cleared
	// comIrqRegSet1Bit = 0x80 // bit 7: see above
)

const (
	regDivIrq = 0x05 // diverse interrupt request bits
	// ------------ values --------------------
	// divIrqRegReset = 0x00 // see table 31 of data sheet
	// divIrqRegReserved01 = 0x03
	divIrqRegCRCIRqBit = 0x04 // bit 2: the CalcCRC command is active and all data is processed
	// divIrqRegReservedBit3 = 0x08
	// this interrupt is set when either a rising or falling signal edge is detected
	// divIrqRegMfinActIRqBit = 0x10 // bit 4: MFIN is active; see above
	// divIrqRegReserved56 = 0x60
	// 1: indicates that the marked bits in the register are set
	// 0: indicates that the marked bits in the register are cleared
	// divIrqRegSet2Bit = 0x80 // bit 7: see above
)

const (
	regError = 0x06 // error bits showing the error status of the last command executed
	// ------------ values --------------------
	// errorRegReset = 0x00 // see table 33 of data sheet
	// set to logic 1 if the SOF is incorrect automatically cleared during receiver start-up phase bit is only valid for
	// 106 kBd; during the MFAuthent command, the ProtocolErr bit is set to logic 1 if the number of bytes received in one
	// data stream is incorrect
	errorRegProtocolErrBit = 0x01 // bit 0: see above
	// automatically cleared during receiver start-up phase; only valid for ISO/IEC 14443 A/MIFARE communication at 106 kBd
	errorRegParityErrBit = 0x02 // bit 1: parity check failed, see above
	// the RxMode register’s RxCRCEn bit is set and the CRC calculation fails automatically; cleared to logic 0 during
	// receiver start-up phase
	errorRegCRCErrBit = 0x04 // bit 2: see above
	// cleared automatically at receiver start-up phase; only valid during the bitwise anticollision at 106 kBd; always
	// set to logic 0 during communication protocols at 212 kBd, 424 kBd and 848 kBd
	errorRegCollErrBit = 0x08 // bit 3: a bit-collision is detected, see above
	// the host or a MFRC522’s internal state machine (e.g. receiver) tries to write data to the FIFO buffer even though
	// it is already full
	errorRegBufferOvflBit = 0x10 // bit 4: FIFO is full, see above
	// errorRegReservedBit5 = 0x20
	// the antenna drivers are automatically switched off
	errorRegTempErrBit = 0x40 // bit 6: internal temperature sensor detects overheating, see above
	// data is written into the FIFO buffer by the host during the MFAuthent command or if data is written into the FIFO
	// buffer by the host during the time between sending the last bit on the RF interface and receiving the last bit on
	// the RF interface
	errorRegWrErrBit = 0x80 // bit 7: see above
)

const (
// ------------ unused commands --------------------
// regStatus1 = 0x07 // communication status bits
)

const (
	regStatus2 = 0x08 // receiver and transmitter status bits
	// ------------ values --------------------
	// status2RegReset = 0x00 // see table 37 of data sheet
	// bit 0..2 shows the state of the transmitter and receiver state machines
	// status2RegModemStateIdle = 0x00 // idle
	// status2RegModemStateWait = 0x01 // wait for the BitFraming register’s StartSend bit
	// the minimum time for TxWait is defined by the TxWait register
	// status2RegModemStateTxWait       = 0x02 // wait until RF field is present if the TMode register’s TxWaitRF bit is set to logic 1
	// status2RegModemStateTransmitting = 0x03
	// the minimum time for RxWait is defined by the RxWait register
	// status2RegModemStateRxWait      = 0x04 // wait until RF field is present if the TMode register’s TxWaitRF bit is set to logic 1
	// status2RegModemStateWaitForData = 0x05
	// status2RegModemStateReceiving   = 0x06
	// all data communication with the card is encrypted; can only be set to logic 1 by a successful execution of the
	// MFAuthent command; only valid in Read/Write mode for MIFARE standard cards; this bit is cleared by software
	status2RegMFCrypto1OnBit = 0x08 // bit 3: indicates that the MIFARE Crypto1 unit is switched on and, see above
	// status2RegReserved45 = 0x30
	// 1: the I2C-bus input filter is set to the High-speed mode independent of the I2C-bus protocol
	// 0: the I2C-bus input filter is set to the I2C-bus protocol used
	// status2RegI2cForceHSBit     = 0x40 // I2C-bus input filter settings, see above
	// status2RegTempSensClear1Bit = 0x80 // clears the temperature error if the temperature is below the alarm limit of 125C
)

const (
	regFIFOData = 0x09 // input and output of 64 byte FIFO buffer
)

const (
	regFIFOLevel = 0x0A // number of bytes stored in the FIFO buffer
	// ------------ values --------------------
	// fifoLevelRegReset = 0x00 // see table 41 of data sheet
	// indicates the number of bytes stored in the FIFO buffer writing to the FIFOData register increments and reading
	// decrements the FIFOLevel value
	// fifoLevelRegValue = 0x7F // bit 0..6: see above
	// immediately clears the internal FIFO buffer’s read and write pointer and Error register’s BufferOvfl bit reading
	// this bit always returns 0
	fifoLevelRegFlushBufferBit = 0x80 // bit 7: see above
)

const (
// ------------ unused commands --------------------
// regWaterLevel = 0x0B // level for FIFO underflow and overflow warning
)

const (
	regControl = 0x0C // miscellaneous control registers
	// ------------ values --------------------
	// controlRegReset = 0x10 // see table 45 of data sheet
	// indicates the number of valid bits in the last received byte
	// if this value is 000b, the whole byte is valid
	controlRegRxLastBits = 0x07 // bit 0..2: see above
	// controlRegReserved3to5 = 0x38
	// controlRegTStartNowBit = 0x40 // bit 6: timer starts immediately, if 1; reading always returns logic 0
	// controlRegTStopNow     = 0x80 // bit 7: timer stops immediately, if 1; reading always returns logic 0
)

const (
	regBitFraming = 0x0D // adjustments for bit-oriented frames
	// ------------ values --------------------
	bitFramingRegReset = 0x00 // see table 47 of data sheet
	// used for transmission of bit oriented frames: defines the number of bits of the last byte that will be transmitted
	// 000b indicates that all bits of the last byte will be transmitted
	bitFramingRegTxLastBits = 0x07 // bit 0..2: see above
	// bitFramingRegReservedBit3 = 0x08
	// used for reception of bit-oriented frames: defines the bit position for the first bit received to be stored in the
	// FIFO buffer,  example:
	// 0: LSB of the received bit is stored at bit position 0, the second received bit is stored at bit position 1
	// 1: LSB of the received bit is stored at bit position 1, the second received bit is stored at bit position 2
	// 7: LSB of the received bit is stored at bit position 7, the second received bit is stored in the next byte that
	// follows at bit position 0
	// These bits are only to be used for bitwise anticollision at 106 kBd, for all other modes they are set to 0
	// bitFramingRegRxAlign = 0x70 // bit 4..6: see above
	// starts the transmission of data, only valid in combination with the Transceive command
	bitFramingRegStartSendBit = 0x80 // bit 7: see above
)

const (
	regColl = 0x0E // bit position of the first bit-collision detected on the RF interface
	// 4 to 0 CollPos[4:0]-
	// shows the bit position of the first detected collision in a received frame only data bits are interpreted example:
	// 00: indicates a bit-collision in the 32nd bit
	// 01: indicates a bit-collision in the 1st bit
	// 08: indicates a bit-collision in the 8th bit
	// These bits will only be interpreted if the CollPosNotValid bit is set to logic 0
	// collRegCollPos = 0x1F // bit 0..4: read-only, see above
	// no collision detected or the position of the collision is out of the range of CollPos[4:0], if set to 1
	// collRegCollPosNotValidBit = 0x20 // bit 5: read-only, see above
	// collRegReservedBit6       = 0x40
	// all received bits will be cleared after a collision only used during bitwise anticollision at 106 kBd, otherwise it
	// is set to logic 1
	collRegValuesAfterCollBit = 0x80 // bit 7: see above
)

//              0x0F     // reserved for future use
// Page 1: Command
//               0x10     // reserved for future use

const (
	regMode = 0x11 // defines general modes for transmitting and receiving
	// ------------ values --------------------
	// modeRegReset = 0x3F // see table 55 of data sheet
	// bit 0..1: defines the preset value for the CRC coprocessor for the CalcCRC command; Remark: during any
	// communication, the preset values are selected automatically according to the definition of bits in the rxModeReg
	// and TxMode registers
	modeRegCRCPreset0000 = 0x00 // 0x0000
	modeRegCRCPreset6363 = 0x01 // 0x6363
	modeRegCRCPresetA671 = 0x10 // 0xA671
	modeRegCRCPresetFFFF = 0x11 // 0xFFFF
	// modeRegReservedBit2 = 0x04
	// defines the polarity of pin MFIN; Remark: the internal envelope signal is encoded active LOW, changing this bit
	// generates a MFinActIRq event
	modeRegPolMFinBit = 0x08 // bit 3: polarity of pin MFIN is active HIGH, is set to 1
	// modeRegReservedBit4 = 0x10
	modeRegTxWaitRFBit = 0x20 // bit 5: transmitter can only be started if an RF field is generated, if set to 1
	// modeRegReservedBit6 = 0x40
	// CRC coprocessor calculates the CRC with MSB first 0 in the CRCResult register the values for the CRCResultMSB[7:0]
	// bits and the CRCResultLSB[7:0] bits are bit reversed; Remark: during RF communication this bit is ignored
	// modeRegMSBFirstBit = 0x80 // bit 7: see above, if set to 1
)

const (
	regTxMode = 0x12 // defines transmission data rate and framing
	regRxMode = 0x13 // defines reception data rate and framing
	// ------------ values --------------------
	rxtxModeRegReset = 0x00
	// txModeRegReserved  = 0x07 // bit 0..2 reserved for TX
	// rxModeRegReserved  = 0x03 // bit 0,1 reserved for RX
	// 0: receiver is deactivated after receiving a data frame
	// 1: able to receive more than one data frame; only valid for data rates above 106 kBd in order to handle the
	// polling command; after setting this bit the Receive and Transceive commands will not terminate automatically.
	// Multiple reception can only be deactivated by writing any command (except the Receive command) to the commandReg
	// register, or by the host clearing the bit if set to logic 1, an error byte is added to the FIFO buffer at the
	// end of a received data stream which is a copy of the Error register value. For the  version 2.0 the
	// CRC status is reflected in the signal CRCOk, which indicates the actual status of the CRC coprocessor. For the
	//  version 1.0 the CRC status is reflected in the signal CRCErr.
//	rxModeRegRxMultipleBit = 0x04
// an invalid received data stream (less than 4 bits received) will be ignored and the receiver remains active
// rxModeRegRxNoErrBit = 0x08 // bit 3
// txModeRegInvModBit  = 0x08 // bit 3: modulation of transmitted data is inverted, if 1
// bit 4..6: defines the bit rate during data transmission; the  handles transfer speeds up to 848 kBd
// rxtxModeRegSpeed106kBd = 0x00 //106 kBd
// rxtxModeRegSpeed212kBd = 0x10 //212 kBd
// rxtxModeRegSpeed424kBd = 0x20 //424 kBd
// rxtxModeRegSpeed848kBd = 0x30 //848 kBd
// rxtxModeRegSpeedRes1   = 0x40 //reserved
// rxtxModeRegSpeedRes2   = 0x50 //reserved
// rxtxModeRegSpeedRes3   = 0x60 //reserved
// rxtxModeRegSpeedRes4   = 0x70 //reserved
// RX: enables the CRC calculation during reception
// TX: enables CRC generation during data transmission
// rxtxModeRegTxCRCEnBit = 0x80 // bit 7: can only be set to logic 0 at 106 kBd
)

const (
	regTxControl = 0x14 // controls the logical behavior of the antenna driver pins TX1 and TX2
	// ------------ values --------------------
	// regtxControlRegReset = 0x80 // see table 61 of data sheet
	// signal on pin TX1 delivers the 13.56 MHz energy carrier modulated by the transmission data
	txControlRegTx1RFEn1outputBit = 0x01 // bit 0: see above
	// signal on pin TX2 delivers the 13.56 MHz energy carrier modulated by the transmission data
	txControlRegTx2RFEn1outputBit = 0x02 // bit 1: see above
	// txControlRegReservedBit2 = 0x04
	// signal on pin TX2 continuously delivers the unmodulated 13.56 MHz energy carrier0Tx2CW bit is enabled to modulate
	// the 13.56 MHz energy carrier
	// txControlRegTx2CW1outputBit = 0x08 // bit 3: see above
	// txControlRegInvTx1RFOffBit  = 0x10 // bit 4: output signal on pin TX1 inverted if driver TX1 is disabled, if 1
	// txControlRegInvTx2RFOffBit  = 0x20 // bit 5: output signal on pin TX2 inverted if driver TX2 is disabled, if 1
	// txControlRegInvTx1RFOnBit   = 0x40 // bit 6: output signal on pin TX1 inverted if driver TX1 is enabled, if 1
	// txControlRegInvTx2RFOnBit   = 0x80 // bit 7: output signal on pin TX2 inverted if driver TX2 is enabled, if 1
)

const (
	regTxASK = 0x15 // controls the setting of the transmission modulation
	// ------------ values --------------------
	// txASKRegReset = 0x00 // see table 63 of data sheet
	// txASKRegReserved = 0x3F // bit 0..5
	txASKRegForce100ASKBit = 0x40 // bit 6: forces a 100 % ASK modulation independent of the ModGsP register
	// txASKRegReservedBit7 = 0x80
)

const (
	// regTxSel       = 0x16 // selects the internal sources for the antenna driver
	// regRxSel       = 0x17 // selects internal receiver settings
	// regRxThreshold = 0x18 // selects thresholds for the bit decoder
	// regDemod       = 0x19 // defines demodulator settings
	//               0x1A     // reserved for future use
	//               0x1B     // reserved for future use
	// regMfTx = 0x1C // controls some MIFARE communication transmit parameters
	// regMfRx = 0x1D // controls some MIFARE communication receive parameters
	//               0x1E     // reserved for future use
	// regSerialSpeed = 0x1F // selects the speed of the serial UART interface

	// Page 2: Configuration
	//               0x20 // reserved for future use
	regCRCResultH = 0x21 // shows the MSB and LSB values of the CRC calculation
	regCRCResultL = 0x22
	//               0x23 // reserved for future use
)

const (
	regModWidth = 0x24
	// ------------ values --------------------
	modWidthRegReset = 0x26 // see table 93 of data sheet
)

//               0x25        // reserved for future use

const (
	regRFCfg = 0x26 // configures the receiver gain
	// ------------ values --------------------
	// rfcCfgRegReset = 0x48 // see table 97 of data sheet
	// rfcCfgRegReserved03 = 0x07
	// bit 4..6: defines the receiver’s signal voltage gain factor
	rfcCfgRegRxGain18dB  = 0x00
	rfcCfgRegRxGain23dB  = 0x10
	rfcCfgRegRxGain018dB = 0x20
	rfcCfgRegRxGain023dB = 0x30
	rfcCfgRegRxGain33dB  = 0x40
	rfcCfgRegRxGain38dB  = 0x50
	rfcCfgRegRxGain43dB  = 0x60
	rfcCfgRegRxGain48dB  = 0x70
	// rfcCfgRegReserved7 = 0x80
)

const (
// ------------ unused commands --------------------
// regGsN    = 0x27 // selects the conductance of the antenna driver pins TX1 and TX2 for modulation
// regCWGsP  = 0x28 // defines the conductance of the p-driver output during periods of no modulation
// regModGsP = 0x29 // defines the conductance of the p-driver output during periods of modulation
)

const (
	regTMode      = 0x2A // defines settings for the internal timer
	regTPrescaler = 0x2B // the lower 8 bits of the TPrescaler value. The 4 high bits are in tModeReg.
	// ------------ values --------------------
	// tModeRegReset      = 0x00 // see table 105 of data sheet
	// tPrescalerRegReset = 0x00 // see table 107 of data sheet
	// timer starts automatically at the end of the transmission in all communication modes at all speeds; if the
	// RxMode register’s RxMultiple bit is not set, the timer stops immediately after receiving the 5th bit (1 start
	// bit, 4 data bits); if the RxMultiple bit is set to logic 1 the timer never stops, in which case the timer can be
	// stopped by setting the Control register’s TStopNow bit to logic 1
	tModeRegTAutoBit = 0x80 // bit 7: see above
	// bit 6,5: indicates that the timer is not influenced by the protocol; internal timer is running in
	// gated mode; Remark: in gated mode, the Status1 register’s TRunning bit is logic 1 when the timer is enabled by
	// the TMode register’s TGated bits; this bit does not influence the gating signal
	// tModeRegTGatedNon  = 0x00 // non-gated mode
	// tModeRegTGatedMFIN = 0x20 // gated by pin MFIN
	// tModeRegTGatedAUX1 = 0x40 // gated by pin AUX1
	// 1: timer automatically restarts its count-down from the 16-bit timer reload value instead of counting down to zero
	// 0: timer decrements to 0 and the ComIrq register’s TimerIRq bit is set to logic 1
	// tModeRegTAutoRestartBit = 0x10 // bit 4, see above
	// defines the higher 4 bits of the TPrescaler value; The following formula is used to calculate the timer
	// frequency if the Demod register’s TPrescalEven bit in Demot register’s set to logic 0:
	// ftimer = 13.56 MHz / (2*TPreScaler+1); TPreScaler = [tPrescalerRegHi:tPrescalerRegLo]
	// TPrescaler value on 12 bits) (Default TPrescalEven bit is logic 0)
	// The following formula is used to calculate the timer frequency if the Demod register’s TPrescalEven bit is set
	// to logic 1: ftimer = 13.56 MHz / (2*TPreScaler+2).
	// tModeRegtPrescalerRegValue25us  = 0x0A9 // 169  => fRegtimer=40kHz, timer period of 25μs.
	// tModeRegtPrescalerRegValue38us  = 0x0FF // 255  => fRegtimer=26kHz, timer period of 38μs.
	// tModeRegtPrescalerRegValue500us = 0xD3E // 3390 => fRegtimer= 2kHz, timer period of 500us.
	// tModeRegtPrescalerRegValue604us = 0xFFF // 4095 => fRegtimer=1.65kHz, timer period of 604us.
)

const (
	// defines the 16-bit timer reload value; on a start event, the timer loads the timer reload value changing this
	// register affects the timer only at the next start event
	regTReloadH = 0x2C
	regTReloadL = 0x2D
	// ------------ values --------------------
	tReloadRegReset      = 0x0000 // see table 109, 111
	tReloadRegValue25ms  = 0x03E8 // 1000,  25ms before timeout
	tReloadRegValue833ms = 0x001E //   30, 833ms before timeout
)

const (
	// ------------ unused commands --------------------
	// regTCounterValueH = 0x2E // shows the 16-bit timer value
	// regTCounterValueL = 0x2F

	// Page 3: Test Registers
	//               0x30      // reserved for future use
	regTestSel1     = 0x31 // general test signal configuration
	regTestSel2     = 0x32 // general test signal configuration
	regTestPinEn    = 0x33 // enables pin output driver on pins D1 to D7
	regTestPinValue = 0x34 // defines the values for D1 to D7 when it is used as an I/O bus
	regTestBus      = 0x35 // shows the status of the internal test bus
	regAutoTest     = 0x36 // controls the digital self-test
	regVersion      = 0x37 // shows the software version
	regAnalogTest   = 0x38 // controls the pins AUX1 and AUX2
	regTestDAC1     = 0x39 // defines the test value for TestDAC1
	regTestDAC2     = 0x3A // defines the test value for TestDAC2
	regTestADC      = 0x3B // shows the value of ADC I and Q channels
	//               0x3C      // reserved for production tests
	//               0x3D      // reserved for production tests
	//               0x3E      // reserved for production tests
	//               0x3F      // reserved for production tests
)
