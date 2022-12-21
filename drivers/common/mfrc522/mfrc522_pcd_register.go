package mfrc522

// Page 0: Command and status
const (
	//              0x00     // reserved for future use
	RegCommand = 0x01 // starts and stops command execution
	// ------------ values --------------------
	// commands, see chapter 10.3, table 149 of the datasheet (only 4 lower bits are used for writing)
	Command_Idle             = 0x00 // no action, cancels current command execution
	Command_Mem              = 0x01 // stores 25 bytes into the internal buffer
	Command_GenerateRandomID = 0x02 // generates a 10-byte random ID number
	Command_CalcCRC          = 0x03 // activates the CRC coprocessor or performs a self-test
	Command_Transmit         = 0x04 // transmits data from the FIFO buffer
	// 0x05, 0x06 not used
	Command_NoCmdChange = 0x07 // no command change, can be used to modify the Command register bits without
	Command_Receive     = 0x08 // activates the receiver circuits
	// 0x09..0x0B not used
	Command_Transceive = 0x0C // transmits data from FIFO buffer to antenna and automatically activates the receiver after transmission
	// 0x0D reserved
	Command_MFAuthent = 0x0E // performs the MIFARE standard authentication as a reader
	Command_SoftReset = 0x0F // resets the MFRC522
	// starts the wake up procedure during which this bit is read as a logic 1; it is read as a logic 0 when the
	// is ready; Remark: The PowerDown bit cannot be set when the SoftReset command is activated
	Command_PowerDownBit = 0x10 // Soft power-down mode entered, if 1
	Command_RcvOffBit    = 0x20 // analog part of the receiver is switched off, if 1
	// Command_Reserved67 = 0xC0
)

const (
	RegComIEn = 0x02 // enable and disable the passing of interrupt requests to IRQ pin
	// ------------ values --------------------
	ComIEn_Reset         = 0x80 // see table 25 of data sheet
	ComIEn_TimerIEnBit   = 0x01 // bit 0: allows the timer interrupt request (TimerIRq bit) to be propagated
	ComIEn_ErrIEnBit     = 0x02 // bit 1: allows the error interrupt request (ErrIRq bit) to be propagated
	ComIEn_LoAlertIEnBit = 0x04 // bit 2: allows the low alert interrupt request (LoAlertIRq bit) to be propagated
	ComIEn_HiAlertIEnBit = 0x08 // bit 3: allows the high alert interrupt request (HiAlertIRq bit) to be propagated
	ComIEn_IdleIEnBit    = 0x10 // bit 4: allows the idle interrupt request (IdleIRq bit) to be propagated
	ComIEn_RxIEnBit      = 0x20 // bit 5: allows the receiver interrupt request (RxIRq bit) to be propagated
	ComIEn_TxIEnBit      = 0x40 // bit 6: allows the transmitter interrupt request (TxIRq bit) to be propagated
	// 1: signal on pin IRQ is inverted with respect to the Status1 register’s IRq bit
	// 0: signal on pin IRQ is equal to the IRq bit; in combination with the DivIEn register’s IRqPushPull bit, the
	// default value of logic 1 ensures that the output level on pin IRQ is 3-state
	ComIEn_IRqInv = 0x80 // bit 7: see above
)

const (
	// ------------ unused commands --------------------
	RegDivIEn = 0x03 // enable and disable the passing of interrupt requests to IRQ pin
)

const (
	RegComIrq = 0x04 // interrupt request bits for communication
	// ------------ values --------------------
	ComIrq_Reset         = 0x14 // see table 29 of data sheet
	ComIrq_ClearAll      = 0x7F // all bits are set to clear, except the Set1Bit (0 indicates the reset)
	ComIrq_TimerIRqBit   = 0x01 // bit 0: the timer decrements the timer value in register TCounterVal to zero, if 1
	ComIrq_ErrIRq1anyBit = 0x02 // bit 1: error bit in the Error register is set, if 1
	// Status1 register’s LoAlert bit is set in opposition to the LoAlert bit, the LoAlertIRq bit stores this event and
	// can only be reset as indicated by the Set1 bit in this register
	ComIrq_LoAlertIRqBit = 0x04 // bit 2: if 1, see above
	// the Status1 register’s HiAlert bit is set in opposition to the HiAlert bit, the HiAlertIRq bit stores this event
	// and can only be reset as indicated by the Set1 bit in this register
	ComIrq_HiAlertIRqBit = 0x08 // bit 3: if 1, see above
	// If a command terminates, for example, when the Command register changes its value from any command to Idle command.
	// If an unknown command is started, the Command register Command[3:0] value changes to the idle state and the
	// IdleIRq bit is set. The microcontroller starting the Idle command does not set the IdleIRq bit.
	ComIrq_IdleIRqBit = 0x10 // bit 4: if 1, see above
	// receiver has detected the end of a valid data stream, if the RxMode register’s RxNoErr bit is set to logic 1,
	// the RxIRq bit is only set to logic 1 when data bytes are available in the FIFO
	ComIrq_RxIRqBit = 0x20 // bit 5: if 1, see above
	ComIrq_TxIRqBit = 0x40 // bit 6: set to 1, immediately after the last bit of the transmitted data was sent out
	// 1: indicates that the marked bits in the register are set
	// 0: indicates that the marked bits in the register are cleared
	ComIrq_Set1Bit = 0x80 // bit 7: see above
)

const (
	RegDivIrq = 0x05 // diverse interrupt request bits
	// ------------ values --------------------
	DivIrq_Reset = 0x00 // see table 31 of data sheet
	//DivIrq_Reserved01 = 0x03
	DivIrq_CRCIRqBit = 0x04 // bit 2: the CalcCRC command is active and all data is processed
	//DivIrq_ReservedBit3 = 0x08
	// this interrupt is set when either a rising or falling signal edge is detected
	DivIrq_MfinActIRqBit = 0x10 // bit 4: MFIN is active; see above
	//DivIrq_Reserved56 = 0x60
	// 1: indicates that the marked bits in the register are set
	// 0: indicates that the marked bits in the register are cleared
	DivIrq_Set2Bit = 0x80 // bit 7: see above
)

const (
	RegError = 0x06 // error bits showing the error status of the last command executed
	// ------------ values --------------------
	Error_Reset = 0x00 // see table 33 of data sheet
	// set to logic 1 if the SOF is incorrect automatically cleared during receiver start-up phase bit is only valid for
	// 106 kBd; during the MFAuthent command, the ProtocolErr bit is set to logic 1 if the number of bytes received in one
	// data stream is incorrect
	Error_ProtocolErrBit = 0x01 // bit 0: see above
	// automatically cleared during receiver start-up phase; only valid for ISO/IEC 14443 A/MIFARE communication at 106 kBd
	Error_ParityErrBit = 0x02 // bit 1: parity check failed, see above
	// the RxMode register’s RxCRCEn bit is set and the CRC calculation fails automatically; cleared to logic 0 during
	// receiver start-up phase
	Error_CRCErrBit = 0x04 // bit 2: see above
	// cleared automatically at receiver start-up phase; only valid during the bitwise anticollision at 106 kBd; always
	// set to logic 0 during communication protocols at 212 kBd, 424 kBd and 848 kBd
	Error_CollErrBit = 0x08 // bit 3: a bit-collision is detected, see above
	//the host or a MFRC522’s internal state machine (e.g. receiver) tries to write data to the FIFO buffer even though
	// it is already full
	Error_BufferOvflBit = 0x10 // bit 4: FIFO is full, see above
	//Error_ReservedBit5 = 0x20
	// the antenna drivers are automatically switched off
	Error_TempErrBit = 0x40 // bit 6: internal temperature sensor detects overheating, see above
	// data is written into the FIFO buffer by the host during the MFAuthent command or if data is written into the FIFO
	// buffer by the host during the time between sending the last bit on the RF interface and receiving the last bit on
	// the RF interface
	Error_WrErrBit = 0x80 // bit 7: see above
)

const (
	// ------------ unused commands --------------------
	RegStatus1 = 0x07 // communication status bits
)

const (
	RegStatus2 = 0x08 // receiver and transmitter status bits
	// ------------ values --------------------
	Status2_Reset = 0x00 // see table 37 of data sheet
	// bit 0..2 shows the state of the transmitter and receiver state machines
	Status2_ModemStateIdle = 0x00 // idle
	Status2_ModemStateWait = 0x01 // wait for the BitFraming register’s StartSend bit
	// the minimum time for TxWait is defined by the TxWait register
	Status2_ModemStateTxWait       = 0x02 // wait until RF field is present if the TMode register’s TxWaitRF bit is set to logic 1
	Status2_ModemStateTransmitting = 0x03
	// the minimum time for RxWait is defined by the RxWait register
	Status2_ModemStateRxWait      = 0x04 // wait until RF field is present if the TMode register’s TxWaitRF bit is set to logic 1
	Status2_ModemStateWaitForData = 0x05
	Status2_ModemStateReceiving   = 0x06
	// all data communication with the card is encrypted; can only be set to logic 1 by a successful execution of the
	// MFAuthent command; only valid in Read/Write mode for MIFARE standard cards; this bit is cleared by software
	Status2_MFCrypto1OnBit = 0x08 // bit 3: indicates that the MIFARE Crypto1 unit is switched on and, see above
	//Status2_Reserved45 = 0x30
	// 1: the I2C-bus input filter is set to the High-speed mode independent of the I2C-bus protocol
	// 0: the I2C-bus input filter is set to the I2C-bus protocol used
	Status2_I2cForceHSBit     = 0x40 // I2C-bus input filter settings, see above
	Status2_TempSensClear1Bit = 0x80 // clears the temperature error if the temperature is below the alarm limit of 125C
)

const (
	RegFIFOData = 0x09 // input and output of 64 byte FIFO buffer
)

const (
	RegFIFOLevel = 0x0A // number of bytes stored in the FIFO buffer
	// ------------ values --------------------
	FIFOLevel_Reset = 0x00 // see table 41 of data sheet
	// indicates the number of bytes stored in the FIFO buffer writing to the FIFOData register increments and reading
	// decrements the FIFOLevel value
	FIFOLevel_Value = 0x7F // bit 0..6: see above
	// immediately clears the internal FIFO buffer’s read and write pointer and Error register’s BufferOvfl bit reading
	// this bit always returns 0
	FIFOLevel_FlushBufferBit = 0x80 // bit 7: see above
)

const (
	// ------------ unused commands --------------------
	RegWaterLevel = 0x0B // level for FIFO underflow and overflow warning
)

const (
	RegControl = 0x0C // miscellaneous control registers
	// ------------ values --------------------
	Control_Reset = 0x10 // see table 45 of data sheet
	// indicates the number of valid bits in the last received byte
	// if this value is 000b, the whole byte is valid
	Control_RxLastBits = 0x07 // bit 0..2: see above
	//Control_Reserved3to5 = 0x38
	Control_TStartNowBit = 0x40 // bit 6: timer starts immediately, if 1; reading always returns logic 0
	Control_TStopNow     = 0x80 // bit 7: timer stops immediately, if 1; reading always returns logic 0
)

const (
	RegBitFraming = 0x0D // adjustments for bit-oriented frames
	// ------------ values --------------------
	BitFraming_Reset = 0x00 // see table 47 of data sheet
	// used for transmission of bit oriented frames: defines the number of bits of the last byte that will be transmitted
	// 000b indicates that all bits of the last byte will be transmitted
	BitFraming_TxLastBits = 0x07 // bit 0..2: see above
	//BitFraming_ReservedBit3 = 0x08
	// used for reception of bit-oriented frames: defines the bit position for the first bit received to be stored in the
	// FIFO buffer,  example:
	// 0: LSB of the received bit is stored at bit position 0, the second received bit is stored at bit position 1
	// 1: LSB of the received bit is stored at bit position 1, the second received bit is stored at bit position 2
	// 7: LSB of the received bit is stored at bit position 7, the second received bit is stored in the next byte that
	// follows at bit position 0
	// These bits are only to be used for bitwise anticollision at 106 kBd, for all other modes they are set to 0
	BitFraming_RxAlign = 0x70 // bit 4..6: see above
	//starts the transmission of data, only valid in combination with the Transceive command
	BitFraming_StartSendBit = 0x80 // bit 7: see above
)

const (
	RegColl = 0x0E // bit position of the first bit-collision detected on the RF interface
	// 4 to 0 CollPos[4:0]-
	// shows the bit position of the first detected collision in a received frame only data bits are interpreted example:
	// 00: indicates a bit-collision in the 32nd bit
	// 01: indicates a bit-collision in the 1st bit
	// 08: indicates a bit-collision in the 8th bit
	// These bits will only be interpreted if the CollPosNotValid bit is set to logic 0
	Coll_CollPos = 0x1F // bit 0..4: read-only, see above
	// no collision detected or the position of the collision is out of the range of CollPos[4:0], if set to 1
	Coll_CollPosNotValidBit = 0x20 // bit 5: read-only, see above
	//Coll_ReservedBit6       = 0x40
	// all received bits will be cleared after a collision only used during bitwise anticollision at 106 kBd, otherwise it
	// is set to logic 1
	Coll_ValuesAfterCollBit = 0x80 // bit 7: see above
)

//              0x0F     // reserved for future use
// Page 1: Command
//               0x10     // reserved for future use

const (
	RegMode = 0x11 // defines general modes for transmitting and receiving
	// ------------ values --------------------
	Mode_Reset = 0x3F // see table 55 of data sheet
	// bit 0..1: defines the preset value for the CRC coprocessor for the CalcCRC command; Remark: during any
	// communication, the preset values are selected automatically according to the definition of bits in the RxModeReg
	// and TxMode registers
	Mode_CRCPreset0000 = 0x00 // 0x0000
	Mode_CRCPreset6363 = 0x01 // 0x6363
	Mode_CRCPresetA671 = 0x10 // 0xA671
	Mode_CRCPresetFFFF = 0x11 // 0xFFFF
	//Mode_ReservedBit2 = 0x04
	// defines the polarity of pin MFIN; Remark: the internal envelope signal is encoded active LOW, changing this bit
	// generates a MFinActIRq event
	Mode_PolMFinBit = 0x08 // bit 3: polarity of pin MFIN is active HIGH, is set to 1
	//Mode_ReservedBit4 = 0x10
	Mode_TxWaitRFBit = 0x20 // bit 5: transmitter can only be started if an RF field is generated, if set to 1
	//Mode_ReservedBit6 = 0x40
	// CRC coprocessor calculates the CRC with MSB first 0 in the CRCResult register the values for the CRCResultMSB[7:0]
	// bits and the CRCResultLSB[7:0] bits are bit reversed; Remark: during RF communication this bit is ignored
	Mode_MSBFirstBit = 0x80 // bit 7: see above, if set to 1
)

const (
	RegTxMode = 0x12 // defines transmission data rate and framing
	RegRxMode = 0x13 // defines reception data rate and framing
	// ------------ values --------------------
	RxTxMode_Reset = 0x00
	//TxMode_Reserved  = 0x07 // bit 0..2 reserved for TX
	//RxMode_Reserved  = 0x03 // bit 0,1 reserved for RX
	// 0: receiver is deactivated after receiving a data frame
	// 1: able to receive more than one data frame; only valid for data rates above 106 kBd in order to handle the
	// polling command; after setting this bit the Receive and Transceive commands will not terminate automatically.
	// Multiple reception can only be deactivated by writing any command (except the Receive command) to the CommandReg
	// register, or by the host clearing the bit if set to logic 1, an error byte is added to the FIFO buffer at the
	// end of a received data stream which is a copy of the Error register value. For the  version 2.0 the
	// CRC status is reflected in the signal CRCOk, which indicates the actual status of the CRC coprocessor. For the
	//  version 1.0 the CRC status is reflected in the signal CRCErr.
	RxMode_RxMultipleBit = 0x04
	// an invalid received data stream (less than 4 bits received) will be ignored and the receiver remains active
	RxMode_RxNoErrBit = 0x08 // bit 3
	TxMode_InvModBit  = 0x08 // bit 3: modulation of transmitted data is inverted, if 1
	// bit 4..6: defines the bit rate during data transmission; the  handles transfer speeds up to 848 kBd
	RxTxMode_Speed106kBd = 0x00 //106 kBd
	RxTxMode_Speed212kBd = 0x10 //212 kBd
	RxTxMode_Speed424kBd = 0x20 //424 kBd
	RxTxMode_Speed848kBd = 0x30 //848 kBd
	RxTxMode_SpeedRes1   = 0x40 //reserved
	RxTxMode_SpeedRes2   = 0x50 //reserved
	RxTxMode_SpeedRes3   = 0x60 //reserved
	RxTxMode_SpeedRes4   = 0x70 //reserved
	// RX: enables the CRC calculation during reception
	// TX: enables CRC generation during data transmission
	RxTxMode_TxCRCEnBit = 0x80 // bit 7: can only be set to logic 0 at 106 kBd
)

const (
	RegTxControl = 0x14 // controls the logical behavior of the antenna driver pins TX1 and TX2
	// ------------ values --------------------
	RegTxControl_Reset = 0x80 // see table 61 of data sheet
	// signal on pin TX1 delivers the 13.56 MHz energy carrier modulated by the transmission data
	TxControl_Tx1RFEn1outputBit = 0x01 // bit 0: see above
	// signal on pin TX2 delivers the 13.56 MHz energy carrier modulated by the transmission data
	TxControl_Tx2RFEn1outputBit = 0x02 // bit 1: see above
	//TxControl_ReservedBit2 = 0x04
	// signal on pin TX2 continuously delivers the unmodulated 13.56 MHz energy carrier0Tx2CW bit is enabled to modulate
	// the 13.56 MHz energy carrier
	TxControl_Tx2CW1outputBit = 0x08 // bit 3: see above
	TxControl_InvTx1RFOffBit  = 0x10 // bit 4: output signal on pin TX1 inverted if driver TX1 is disabled, if 1
	TxControl_InvTx2RFOffBit  = 0x20 // bit 5: output signal on pin TX2 inverted if driver TX2 is disabled, if 1
	TxControl_InvTx1RFOnBit   = 0x40 // bit 6: output signal on pin TX1 inverted if driver TX1 is enabled, if 1
	TxControl_InvTx2RFOnBit   = 0x80 // bit 7: output signal on pin TX2 inverted if driver TX2 is enabled, if 1
)

const (
	RegTxASK = 0x15 // controls the setting of the transmission modulation
	// ------------ values --------------------
	TxASK_Reset = 0x00 // see table 63 of data sheet
	//TxASK_Reserved = 0x3F // bit 0..5
	TxASK_Force100ASKBit = 0x40 // bit 6: forces a 100 % ASK modulation independent of the ModGsP register
	//TxASK_ReservedBit7 = 0x80
)

const (
	RegTxSel       = 0x16 // selects the internal sources for the antenna driver
	RegRxSel       = 0x17 // selects internal receiver settings
	RegRxThreshold = 0x18 // selects thresholds for the bit decoder
	RegDemod       = 0x19 // defines demodulator settings
	//               0x1A     // reserved for future use
	//               0x1B     // reserved for future use
	RegMfTx = 0x1C // controls some MIFARE communication transmit parameters
	RegMfRx = 0x1D // controls some MIFARE communication receive parameters
	//               0x1E     // reserved for future use
	RegSerialSpeed = 0x1F // selects the speed of the serial UART interface

	// Page 2: Configuration
	//               0x20 // reserved for future use
	RegCRCResultH = 0x21 // shows the MSB and LSB values of the CRC calculation
	RegCRCResultL = 0x22
	//               0x23 // reserved for future use
)

const (
	RegModWidth = 0x24
	// ------------ values --------------------
	ModWidth_Reset = 0x26 // see table 93 of data sheet
)

//               0x25        // reserved for future use

const (
	RegRFCfg = 0x26 // configures the receiver gain
	// ------------ values --------------------
	RFCfg_Reset = 0x48 // see table 97 of data sheet
	//RFCfg_Reserved03 = 0x07
	// bit 4..6: defines the receiver’s signal voltage gain factor
	RFCfg_RxGain18dB  = 0x00
	RFCfg_RxGain23dB  = 0x10
	RFCfg_RxGain018dB = 0x20
	RFCfg_RxGain023dB = 0x30
	RFCfg_RxGain33dB  = 0x40
	RFCfg_RxGain38dB  = 0x50
	RFCfg_RxGain43dB  = 0x60
	RFCfg_RxGain48dB  = 0x70
	//RFCfg_Reserved7 = 0x80
)

const (
	// ------------ unused commands --------------------
	RegGsN    = 0x27 // selects the conductance of the antenna driver pins TX1 and TX2 for modulation
	RegCWGsP  = 0x28 // defines the conductance of the p-driver output during periods of no modulation
	RegModGsP = 0x29 // defines the conductance of the p-driver output during periods of modulation

)
const (
	RegTMode      = 0x2A // defines settings for the internal timer
	RegTPrescaler = 0x2B // the lower 8 bits of the TPrescaler value. The 4 high bits are in TModeReg.
	// ------------ values --------------------
	TMode_Reset      = 0x00 // see table 105 of data sheet
	TPrescaler_Reset = 0x00 // see table 107 of data sheet
	// timer starts automatically at the end of the transmission in all communication modes at all speeds; if the
	// RxMode register’s RxMultiple bit is not set, the timer stops immediately after receiving the 5th bit (1 start
	// bit, 4 data bits); if the RxMultiple bit is set to logic 1 the timer never stops, in which case the timer can be
	// stopped by setting the Control register’s TStopNow bit to logic 1
	TMode_TAutoBit = 0x80 // bit 7: see above
	// bit 6,5: indicates that the timer is not influenced by the protocol; internal timer is running in
	// gated mode; Remark: in gated mode, the Status1 register’s TRunning bit is logic 1 when the timer is enabled by
	// the TMode register’s TGated bits; this bit does not influence the gating signal
	TMode_TGatedNon  = 0x00 // non-gated mode
	TMode_TGatedMFIN = 0x20 // gated by pin MFIN
	TMode_TGatedAUX1 = 0x40 // gated by pin AUX1
	// 1: timer automatically restarts its count-down from the 16-bit timer reload value instead of counting down to zero
	// 0: timer decrements to 0 and the ComIrq register’s TimerIRq bit is set to logic 1
	TMode_TAutoRestartBit = 0x10 // bit 4, see above
	// defines the higher 4 bits of the TPrescaler value; The following formula is used to calculate the timer
	// frequency if the Demod register’s TPrescalEven bit in Demot register’s set to logic 0:
	// ftimer = 13.56 MHz / (2*TPreScaler+1); TPreScaler = [TPrescaler_Hi:TPrescaler_Lo]
	// TPrescaler value on 12 bits) (Default TPrescalEven bit is logic 0)
	// The following formula is used to calculate the timer frequency if the Demod register’s TPrescalEven bit is set
	// to logic 1: ftimer = 13.56 MHz / (2*TPreScaler+2).
	TMode_TPrescaler_Value25us  = 0x0A9 // 169  => f_timer=40kHz, timer period of 25μs.
	TMode_TPrescaler_Value38us  = 0x0FF // 255  => f_timer=26kHz, timer period of 38μs.
	TMode_TPrescaler_Value500us = 0xD3E // 3390 => f_timer= 2kHz, timer period of 500us.
	TMode_TPrescaler_Value604us = 0xFFF // 4095 => f_timer=1.65kHz, timer period of 604us.
)

const (
	// defines the 16-bit timer reload value; on a start event, the timer loads the timer reload value changing this
	// register affects the timer only at the next start event
	RegTReloadH = 0x2C
	RegTReloadL = 0x2D
	// ------------ values --------------------
	TReload_Reset      = 0x0000 // see table 109, 111
	TReload_Value25ms  = 0x03E8 // 1000,  25ms before timeout
	TReload_Value833ms = 0x001E //   30, 833ms before timeout
)

const (
	// ------------ unused commands --------------------
	RegTCounterValueH = 0x2E // shows the 16-bit timer value
	RegTCounterValueL = 0x2F

	// Page 3: Test Registers
	//               0x30      // reserved for future use
	RegTestSel1     = 0x31 // general test signal configuration
	RegTestSel2     = 0x32 // general test signal configuration
	RegTestPinEn    = 0x33 // enables pin output driver on pins D1 to D7
	RegTestPinValue = 0x34 // defines the values for D1 to D7 when it is used as an I/O bus
	RegTestBus      = 0x35 // shows the status of the internal test bus
	RegAutoTest     = 0x36 // controls the digital self-test
	RegVersion      = 0x37 // shows the software version
	RegAnalogTest   = 0x38 // controls the pins AUX1 and AUX2
	RegTestDAC1     = 0x39 // defines the test value for TestDAC1
	RegTestDAC2     = 0x3A // defines the test value for TestDAC2
	RegTestADC      = 0x3B // shows the value of ADC I and Q channels
	//               0x3C      // reserved for production tests
	//               0x3D      // reserved for production tests
	//               0x3E      // reserved for production tests
	//               0x3F      // reserved for production tests
)
