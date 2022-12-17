package mfrc522

// Page 0: Command and status
const (
	//              0x00     // reserved for future use
	mfrc522RegCommand = 0x01 // starts and stops command execution
	// ------------ values --------------------
	// commands, see chapter 10.3, table 149 of the datasheet (only 4 lower bits are used for writing)
	mfrc522PCD_Idle             = 0x00 // no action, cancels current command execution
	mfrc522PCD_Mem              = 0x01 // stores 25 bytes into the internal buffer
	mfrc522PCD_GenerateRandomID = 0x02 // generates a 10-byte random ID number
	mfrc522PCD_CalcCRC          = 0x03 // activates the CRC coprocessor or performs a self-test
	mfrc522PCD_Transmit         = 0x04 // transmits data from the FIFO buffer
	mfrc522PCD_NoCmdChange      = 0x07 // no command change, can be used to modify the CommandReg register bits without

	// affecting the command, for example, the PowerDown bit
	mfrc522PCD_Receive    = 0x08 // activates the receiver circuits
	mfrc522PCD_Transceive = 0x0C // transmits data from FIFO buffer to antenna and automatically activates the receiver after transmission
	mfrc522PCD_MFAuthent  = 0x0E // performs the MIFARE standard authentication as a reader
	mfrc522PCD_SoftReset  = 0x0F // resets the MFRC522
	mfrc522PCD_RcvOffBit  = 0x20 // analog part of the receiver is switched off, if 1
	// MFRC522 starts the wake up procedure during which this bit is read as a logic 1; it is read as a logic 0 when the
	// MFRC522 is ready; Remark: The PowerDown bit cannot be set when the SoftReset command is activated
	mfrc522PCD_PowerDownBit = 0x10 // Soft power-down mode entered, if 1
)

const (
	// ------------ unused commands --------------------
	mfrc522RegComIEn     = 0x02 // enable and disable interrupt request control bits
	mfrc522RegDivIEn     = 0x03 // enable and disable interrupt request control bits
	mfrc522RegComIrq     = 0x04 // interrupt request bits
	mfrc522RegDivIrq     = 0x05 // interrupt request bits
	mfrc522RegError      = 0x06 // error bits showing the error status of the last command executed
	mfrc522RegStatus1    = 0x07 // communication status bits
	mfrc522RegStatus2    = 0x08 // receiver and transmitter status bits
	mfrc522RegFIFOData   = 0x09 // input and output of 64 byte FIFO buffer
	mfrc522RegFIFOLevel  = 0x0A // number of bytes stored in the FIFO buffer
	mfrc522RegWaterLevel = 0x0B // level for FIFO underflow and overflow warning
	mfrc522RegControl    = 0x0C // miscellaneous control registers
	mfrc522RegBitFraming = 0x0D // adjustments for bit-oriented frames
)

const (
	mfrc522RegColl = 0x0E // bit position of the first bit-collision detected on the RF interface
	// all received bits will be cleared after a collision only used during bitwise anticollision at 106 kBd, otherwise it
	// is set to logic 1
	mfrc522Coll_ValuesAfterCollBit = 0x80 // bit 7: see above
	//mfrc522Coll_reservedBit6       = 0x00
	// no collision detected or the position of the collision is out of the range of CollPos[4:0], if set to 1
	mfrc522Coll_CollPosNotValidBit // bit 5: read-only, see above
	// 4 to 0 CollPos[4:0]-
	// shows the bit position of the first detected collision in a received frame only data bits are interpreted example:
	// 00: indicates a bit-collision in the 32nd bit
	// 01: indicates a bit-collision in the 1st bit
	// 08: indicates a bit-collision in the 8th bit
	// These bits will only be interpreted if the CollPosNotValid bit is set to logic 0
	mfrc522Coll_CollPos = 0x1F // bit 0..4: read-only, see above
	//              0x0F     // reserved for future use
	// Page 1: Command
	//               0x10     // reserved for future use
)

const (
	mfrc522RegMode = 0x11 // defines general modes for transmitting and receiving
	// ------------ values --------------------
	mfrc522Mode_Reset = 0x3F // see table 55 of data sheet
	// CRC coprocessor calculates the CRC with MSB first 0 in the CRCResultReg register the values for the CRCResultMSB[7:0]
	// bits and the CRCResultLSB[7:0] bits are bit reversed; Remark: during RF communication this bit is ignored
	mfrc522Mode_MSBFirstBit = 0x80 // bit 7: see above, if set to 1
	//mfrc522Mode_ReservedBit6 = 0x40
	mfrc522Mode_TxWaitRFBit = 0x20 // bit 5: transmitter can only be started if an RF field is generated, if set to 1
	//mfrc522Mode_ReservedBit4 = 0x10
	// defines the polarity of pin MFIN; Remark: the internal envelope signal is encoded active LOW, changing this bit
	// generates a MFinActIRq event
	mfrc522Mode_PolMFinBit = 0x80 // bit 3: polarity of pin MFIN is active HIGH, is set to 1
	//mfrc522Mode_ReservedBit2 = 0x40
	// bit 0..2: defines the preset value for the CRC coprocessor for the CalcCRC command; Remark: during any
	// communication, the preset values are selected automatically according to the definition of bits in the RxModeReg
	// and TxModeReg registers
	mfrc522Mode_CRCPreset0000 = 0x00 // 0x0000
	mfrc522Mode_CRCPreset6363 = 0x01 // 0x6363
	mfrc522Mode_CRCPresetA671 = 0x10 // 0xA671
	mfrc522Mode_CRCPresetFFFF = 0x11 // 0xFFFF
)

const (
	mfrc522RegTxMode = 0x12 // defines transmission data rate and framing
	mfrc522RegRxMode = 0x13 // defines reception data rate and framing
	// ------------ values --------------------
	mfrc522RxTxMode_Reset = 0x00
	//TxMode_Reserved  = 0x07 // bit 0..2 reserved for TX
	//RxMode_Reserved  = 0x03 // bit 0,1 reserved for RX
	// 0: receiver is deactivated after receiving a data frame
	// 1: able to receive more than one data frame; only valid for data rates above 106 kBd in order to handle the
	// polling command; after setting this bit the Receive and Transceive commands will not terminate automatically.
	// Multiple reception can only be deactivated by writing any command (except the Receive command) to the CommandReg
	// register, or by the host clearing the bit if set to logic 1, an error byte is added to the FIFO buffer at the
	// end of a received data stream which is a copy of the ErrorReg register value. For the MFRC522 version 2.0 the
	// CRC status is reflected in the signal CRCOk, which indicates the actual status of the CRC coprocessor. For the
	// MFRC522 version 1.0 the CRC status is reflected in the signal CRCErr.
	mfrc522RxMode_RxMultipleBit = 0x04
	// an invalid received data stream (less than 4 bits received) will be ignored and the receiver remains active
	mfrc522RxMode_RxNoErrBit = 0x08 // bit 3
	mfrc522TxMode_InvModBit  = 0x08 // bit 3: modulation of transmitted data is inverted, if 1
	// bit 4..6: defines the bit rate during data transmission; the MFRC522 handles transfer speeds up to 848 kBd
	mfrc522RxTxMode_Speed106kBd = 0x00 //106 kBd
	mfrc522RxTxMode_Speed212kBd = 0x10 //212 kBd
	mfrc522RxTxMode_Speed424kBd = 0x20 //424 kBd
	mfrc522RxTxMode_Speed848kBd = 0x30 //848 kBd
	mfrc522RxTxMode_SpeedRes1   = 0x40 //reserved
	mfrc522RxTxMode_SpeedRes2   = 0x50 //reserved
	mfrc522RxTxMode_SpeedRes3   = 0x60 //reserved
	mfrc522RxTxMode_SpeedRes4   = 0x70 //reserved
	// RX: enables the CRC calculation during reception
	// TX: enables CRC generation during data transmission
	mfrc522RxTxMode_TxCRCEnBit = 0x80 // bit 7: can only be set to logic 0 at 106 kBd
)

const (
	// ------------ unused commands --------------------
	mfrc522RegTxControl = 0x14 // controls the logical behavior of the antenna driver pins TX1 and TX2
	// ------------ values --------------------
	mfrc522RegTxControl_Reset       = 0x80 // see table 61 of data sheet
	mfrc522TxControl_InvTx2RFOnBit  = 0x80 // bit 7: output signal on pin TX2 inverted if driver TX2 is enabled, if 1
	mfrc522TxControl_InvTx1RFOnBit  = 0x40 // bit 6: output signal on pin TX1 inverted if driver TX1 is enabled, if 1
	mfrc522TxControl_InvTx2RFOffBit = 0x20 // bit 5: output signal on pin TX2 inverted if driver TX2 is disabled, if 1
	mfrc522TxControl_InvTx1RFOffBit = 0x10 // bit 4: output signal on pin TX1 inverted if driver TX1 is disabled, if 1
	// signal on pin TX2 continuously delivers the unmodulated 13.56 MHz energy carrier0Tx2CW bit is enabled to modulate
	// the 13.56 MHz energy carrier
	mfrc522TxControl_Tx2CW1outputBit = 0x08 // bit 3: see above
	// mfrc522TxControl_ReservedBit2 = 0x04
	// signal on pin TX2 delivers the 13.56 MHz energy carrier modulated by the transmission data
	mfrc522TxControl_Tx2RFEn1outputBit = 0x02 // bit 1: see above
	// signal on pin TX1 delivers the 13.56 MHz energy carrier modulated by the transmission data
	mfrc522TxControl_Tx1RFEn1outputBit = 0x01 // bit 0: see above
)

const (
	mfrc522RegTxASK = 0x15 // controls the setting of the transmission modulation
	// ------------ values --------------------
	mfrc522TxASK_Reset = 0x00 // see table 63 of data sheet
	// mfrc522TxASK_Reserved = 0x3F // bit 0..5
	mfrc522TxASK_Force100ASKBit = 0x40 // bit 6: forces a 100 % ASK modulation independent of the ModGsPReg register
	// mfrc522TxASK_ReservedBit7 = 0x80
)

const (
	mfrc522RegTxSel       = 0x16 // selects the internal sources for the antenna driver
	mfrc522RegRxSel       = 0x17 // selects internal receiver settings
	mfrc522RegRxThreshold = 0x18 // selects thresholds for the bit decoder
	mfrc522RegDemod       = 0x19 // defines demodulator settings
	//               0x1A     // reserved for future use
	//               0x1B     // reserved for future use
	mfrc522RegMfTx = 0x1C // controls some MIFARE communication transmit parameters
	mfrc522RegMfRx = 0x1D // controls some MIFARE communication receive parameters
	//               0x1E     // reserved for future use
	mfrc522RegSerialSpeed = 0x1F // selects the speed of the serial UART interface

	// Page 2: Configuration
	//               0x20        // reserved for future use
	mfrc522RegCRCResultH = 0x21 // shows the MSB and LSB values of the CRC calculation
	mfrc522RegCRCResultL = 0x22
	//               0x23        // reserved for future use
)

const (
	mfrc522RegModWidth = 0x24 // controls the ModWidth setting?
	// ------------ values --------------------
	mfrc522ModWidth_Reset = 0x26 // see table 93 of data sheet
)

const (
	// ------------ unused commands --------------------
	//               0x25        // reserved for future use
	mfrc522RegRFCfg  = 0x26 // configures the receiver gain
	mfrc522RegGsN    = 0x27 // selects the conductance of the antenna driver pins TX1 and TX2 for modulation
	mfrc522RegCWGsP  = 0x28 // defines the conductance of the p-driver output during periods of no modulation
	mfrc522RegModGsP = 0x29 // defines the conductance of the p-driver output during periods of modulation

)
const (
	mfrc522RegTMode      = 0x2A // defines settings for the internal timer
	mfrc522RegTPrescaler = 0x2B // the lower 8 bits of the TPrescaler value. The 4 high bits are in TModeReg.
	// ------------ values --------------------
	mfrc522TMode_Reset      = 0x00 // see table 105 of data sheet
	mfrc522TPrescaler_Reset = 0x00 // see table 107 of data sheet
	// timer starts automatically at the end of the transmission in all communication modes at all speeds; if the
	// RxModeReg register’s RxMultiple bit is not set, the timer stops immediately after receiving the 5th bit (1 start
	// bit, 4 data bits); if the RxMultiple bit is set to logic 1 the timer never stops, in which case the timer can be
	// stopped by setting the ControlReg register’s TStopNow bit to logic 1
	mfrc522TMode_TAutoBit = 0x80 // bit 7: see above
	// bit 6,5: indicates that the timer is not influenced by the protocol; internal timer is running in
	// gated mode; Remark: in gated mode, the Status1Reg register’s TRunning bit is logic 1 when the timer is enabled by
	// the TModeReg register’s TGated bits; this bit does not influence the gating signal
	mfrc522TMode_TGatedNon  = 0x00 // non-gated mode
	mfrc522TMode_TGatedMFIN = 0x20 // gated by pin MFIN
	mfrc522TMode_TGatedAUX1 = 0x40 // gated by pin AUX1
	// 1: timer automatically restarts its count-down from the 16-bit timer reload value instead of counting down to zero
	// 0: timer decrements to 0 and the ComIrqReg register’s TimerIRq bit is set to logic 1
	mfrc522TMode_TAutoRestartBit = 0x10 // bit 4, see above
	// defines the higher 4 bits of the TPrescaler value; The following formula is used to calculate the timer
	// frequency if the DemodReg register’s TPrescalEven bit in Demot Regis set to logic 0:
	// ftimer = 13.56 MHz / (2*TPreScaler+1); TPreScaler = [TPrescaler_Hi:TPrescaler_Lo]
	// TPrescaler value on 12 bits) (Default TPrescalEven bit is logic 0)
	// The following formula is used to calculate the timer frequency if the DemodReg register’s TPrescalEven bit is set
	// to logic 1: ftimer = 13.56 MHz / (2*TPreScaler+2).
	mfrc522TMode_TPrescaler_Value25us  = 0x0A9 // 169  => f_timer=40kHz, timer period of 25μs.
	mfrc522TMode_TPrescaler_Value38us  = 0x0FF // 255  => f_timer=26kHz, timer period of 38μs.
	mfrc522TMode_TPrescaler_Value500us = 0xD3E // 3390 => f_timer= 2kHz, timer period of 500us.
	mfrc522TMode_TPrescaler_Value604us = 0xFFF // 4095 => f_timer=1.65kHz, timer period of 604us.
)

const (
	// defines the 16-bit timer reload value; on a start event, the timer loads the timer reload value changing this
	// register affects the timer only at the next start event
	mfrc522RegTReloadH = 0x2C
	mfrc522RegTReloadL = 0x2D
	// ------------ values --------------------
	mfrc522RegTReload_Reset      = 0x0000 // see table 109, 111
	mfrc522RegTReload_Value25ms  = 0x03E8 // 1000,  25ms before timeout
	mfrc522RegTReload_Value833ms = 0x001E //   30, 833ms before timeout
)

const (
	// ------------ unused commands --------------------
	mfrc522RegTCounterValueH = 0x2E // shows the 16-bit timer value
	mfrc522RegTCounterValueL = 0x2F

	// Page 3: Test Registers
	//               0x30      // reserved for future use
	mfrc522RegTestSel1     = 0x31 // general test signal configuration
	mfrc522RegTestSel2     = 0x32 // general test signal configuration
	mfrc522RegTestPinEn    = 0x33 // enables pin output driver on pins D1 to D7
	mfrc522RegTestPinValue = 0x34 // defines the values for D1 to D7 when it is used as an I/O bus
	mfrc522RegTestBus      = 0x35 // shows the status of the internal test bus
	mfrc522RegAutoTest     = 0x36 // controls the digital self-test
	mfrc522RegVersion      = 0x37 // shows the software version
	mfrc522RegAnalogTest   = 0x38 // controls the pins AUX1 and AUX2
	mfrc522RegTestDAC1     = 0x39 // defines the test value for TestDAC1
	mfrc522RegTestDAC2     = 0x3A // defines the test value for TestDAC2
	mfrc522RegTestADC      = 0x3B // shows the value of ADC I and Q channels
	//               0x3C      // reserved for production tests
	//               0x3D      // reserved for production tests
	//               0x3E      // reserved for production tests
	//               0x3F      // reserved for production tests
)
