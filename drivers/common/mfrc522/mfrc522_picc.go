//nolint:lll // ok here
package mfrc522

import (
	"fmt"
)

const piccDebug = false

// Commands sent to the PICC, used by the PCD to for communication with several PICCs (ISO 14443-3, Type A, section 6.4)
const (
	// Activation
	piccCommandRequestA = 0x26 // REQuest command type A, 7 bit frame, invites PICCs in state IDLE to go to READY
	piccCommandWakeUpA  = 0x52 // Wake-UP command type A, 7 bit frame, invites PICCs in state IDLE and HALT to go to READY
	// Anticollision and SAK
	piccCommandCascadeLevel1 = 0x93 // Select cascade level 1
	piccCommandCascadeLevel2 = 0x95 // Select cascade Level 2
	piccCommandCascadeLevel3 = 0x97 // Select cascade Level 3
	piccCascadeTag           = 0x88 // Cascade tag is used during anti collision
	piccUIDNotComplete       = 0x04 // used on SAK call
	// Halt
	piccCommandHLTA = 0x50 // Halt command, Type A. Instructs an active PICC to go to state HALT.
	piccCommandRATS = 0xE0 // Request command for Answer To Reset.
	// The commands used for MIFARE Classic (from http://www.mouser.com/ds/2/302/MF1S503x-89574.pdf, Section 9)
	// Use MFAuthent to authenticate access to a sector, then use these commands to read/write/modify the blocks on
	// the sector. The read/write commands can also be used for MIFARE Ultralight.
	piccCommandMFRegAUTHRegKEYRegA = 0x60 // Perform authentication with Key A
	piccCommandMFRegAUTHRegKEYRegB = 0x61 // Perform authentication with Key B
	// Reads one 16 byte block from the authenticated sector of the PICC. Also used for MIFARE Ultralight.
	piccCommandMFRegREAD = 0x30
	// Writes one 16 byte block to the authenticated sector of the PICC. Called "COMPATIBILITY WRITE" for MIFARE Ultralight.
	piccCommandMFRegWRITE     = 0xA0
	piccWriteAck              = 0x0A // MIFARE Classic: 4 bit ACK, we use any other value as NAK (data sheet: 0h to 9h, Bh to Fh)
	piccCommandMFRegDECREMENT = 0xC0 // Decrements the contents of a block and stores the result in the internal data register.
	piccCommandMFRegINCREMENT = 0xC1 // Increments the contents of a block and stores the result in the internal data register.
	piccCommandMFRegRESTORE   = 0xC2 // Reads the contents of a block into the internal data register.
	piccCommandMFRegTRANSFER  = 0xB0 // Writes the contents of the internal data register to a block.
	// The commands used for MIFARE Ultralight (from http://www.nxp.com/documents/dataRegsheet/MF0ICU1.pdf, Section 8.6)
	// The piccCommandMFRegREAD and piccCommandMFRegWRITE can also be used for MIFARE Ultralight.
	// piccCommandULRegWRITE = 0xA2 // Writes one 4 byte page to the PICC.
)

const piccReadWriteAuthBlock = uint8(11)

var (
	piccKey                = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	piccUserBlockAddresses = []byte{8, 9, 10}
)

var piccCardFromSak = map[uint8]string{
	0x08: "Classic 1K, Plus 2K-SE-1K(SL1)", 0x18: "Classic 4K, Plus 4K(SL1)",
	0x10: "Plus 2K(SL2)", 0x11: "Plus 4K(SL2)", 0x20: "Plus 2K-SE-1K(SL3), Plus 4K(SL3)",
}

// IsCardPresent is used to poll for a card in range. After an successful request, the card is halted.
func (d *MFRC522Common) IsCardPresent() error {
	d.firstCardAccess = true

	if err := d.writeByteData(regTxMode, rxtxModeRegReset); err != nil {
		return err
	}
	if err := d.writeByteData(regRxMode, rxtxModeRegReset); err != nil {
		return err
	}
	if err := d.writeByteData(regModWidth, modWidthRegReset); err != nil {
		return err
	}

	answer := []byte{0x00, 0x00} // also called ATQA
	if err := d.piccRequest(piccCommandWakeUpA, answer); err != nil {
		return err
	}

	if piccDebug {
		fmt.Printf("Card found: %v\n\n", answer)
	}
	if err := d.piccHalt(); err != nil {
		return err
	}
	return nil
}

// ReadText reads a card with the dedicated workflow: REQA, Activate, Perform Transaction, Halt/Deselect.
// see "Card Polling" in https://www.nxp.com/docs/en/application-note/AN10834.pdf.
// and return the result as text string.
// TODO: make this more usable, e.g. by given length of text
func (d *MFRC522Common) ReadText() (string, error) {
	answer := []byte{0x00, 0x00}
	if err := d.piccRequest(piccCommandWakeUpA, answer); err != nil {
		return "", err
	}

	uid, err := d.piccActivate()
	if err != nil {
		return "", err
	}

	if piccDebug {
		fmt.Printf("uid: %v\n", uid)
	}

	if err := d.piccAuthenticate(piccReadWriteAuthBlock, piccKey, uid); err != nil {
		if piccDebug {
			fmt.Println("authenticate failed for address", piccReadWriteAuthBlock)
		}
		return "", err
	}

	var content []byte
	for _, block := range piccUserBlockAddresses {
		blockData, err := d.piccRead(block)
		if err != nil {
			if piccDebug {
				fmt.Println("read failed at block", block)
			}
			return "", err
		}
		content = append(content, blockData...)
	}
	if piccDebug {
		fmt.Println("content:", string(content), content)
	}

	if err := d.piccHalt(); err != nil {
		return "", err
	}

	return string(content), d.stopCrypto1()
}

// WriteText writes the given string to the card. All old values will be overwritten.
func (d *MFRC522Common) WriteText(text string) error {
	answer := []byte{0x00, 0x00}
	if err := d.piccRequest(piccCommandWakeUpA, answer); err != nil {
		return err
	}

	uid, err := d.piccActivate()
	if err != nil {
		return err
	}

	if piccDebug {
		fmt.Printf("uid: %v\n", uid)
	}

	if err := d.piccAuthenticate(piccReadWriteAuthBlock, piccKey, uid); err != nil {
		if piccDebug {
			fmt.Println("authenticate failed for address", piccReadWriteAuthBlock)
		}
		return err
	}

	// prepare data with text and trailing zero's
	textData := append([]byte(text), make([]byte, len(piccUserBlockAddresses)*16)...)

	for i, blockNum := range piccUserBlockAddresses {
		blockData := textData[i*16 : (i+1)*16]
		err := d.piccWrite(blockNum, blockData)
		if err != nil {
			if piccDebug {
				fmt.Println("write failed at block", blockNum)
			}
			return err
		}
	}

	if err := d.piccHalt(); err != nil {
		return err
	}

	return d.stopCrypto1()
}

func (d *MFRC522Common) piccHalt() error {
	if piccDebug {
		fmt.Println("-halt-")
	}
	haltCommand := []byte{piccCommandHLTA, 0x00}
	crcResult := []byte{0x00, 0x00}
	if err := d.calculateCRC(haltCommand, crcResult); err != nil {
		return err
	}
	haltCommand = append(haltCommand, crcResult...)

	txLastBits := uint8(0x00) // we use all 8 bits
	if err := d.communicateWithPICC(commandRegTransceive, haltCommand, []byte{}, txLastBits, false); err != nil {
		// an error is the sign for successful halt
		if piccDebug {
			fmt.Println("this is not treated as error:", err)
		}
		return nil
	}

	return fmt.Errorf("something went wrong with halt")
}

func (d *MFRC522Common) piccWrite(block uint8, blockData []byte) error {
	if piccDebug {
		fmt.Println("-write-")
		fmt.Println("blockData:", blockData, len(blockData))
	}
	if len(blockData) != 16 {
		return fmt.Errorf("the block to write needs to be exactly 16 bytes long, but has %d bytes", len(blockData))
	}
	// MIFARE Classic protocol requires two steps to perform a write.
	// Step 1: Tell the PICC we want to write to block blockAddr.
	writeDataCommand := []byte{piccCommandMFRegWRITE, block}
	crcResult := []byte{0x00, 0x00}
	if err := d.calculateCRC(writeDataCommand, crcResult); err != nil {
		return err
	}
	writeDataCommand = append(writeDataCommand, crcResult...)

	txLastBits := uint8(0x00) // we use all 8 bits
	backData := make([]byte, 1)
	if err := d.communicateWithPICC(commandRegTransceive, writeDataCommand, backData, txLastBits, false); err != nil {
		return err
	}
	if backData[0]&piccWriteAck != piccWriteAck {
		return fmt.Errorf("preparation of write on MIFARE classic failed (%v)", backData)
	}
	if piccDebug {
		fmt.Println("backData", backData)
	}

	// Step 2: Transfer the data
	if err := d.calculateCRC(blockData, crcResult); err != nil {
		return err
	}

	var writeData []byte
	writeData = append(writeData, blockData...)
	writeData = append(writeData, crcResult...)
	if err := d.communicateWithPICC(commandRegTransceive, writeData, []byte{}, txLastBits, false); err != nil {
		return err
	}

	return nil
}

func (d *MFRC522Common) piccRead(block uint8) ([]byte, error) {
	if piccDebug {
		fmt.Println("-read-")
	}
	readDataCommand := []byte{piccCommandMFRegREAD, block}
	crcResult := []byte{0x00, 0x00}
	if err := d.calculateCRC(readDataCommand, crcResult); err != nil {
		return nil, err
	}
	readDataCommand = append(readDataCommand, crcResult...)

	txLastBits := uint8(0x00)    // we use all 8 bits
	backData := make([]byte, 18) // 16 data byte and 2 byte CRC
	if err := d.communicateWithPICC(commandRegTransceive, readDataCommand, backData, txLastBits, true); err != nil {
		return nil, err
	}

	return backData[:16], nil
}

func (d *MFRC522Common) piccAuthenticate(address uint8, key []byte, uid []byte) error {
	if piccDebug {
		fmt.Println("-authenticate-")
	}

	buf := []byte{piccCommandMFRegAUTHRegKEYRegA, address}
	buf = append(buf, key...)
	buf = append(buf, uid...)

	if err := d.communicateWithPICC(commandRegMFAuthent, buf, []byte{}, 0, false); err != nil {
		return err
	}

	return nil
}

// activate a card with the dedicated workflow: Anticollision and optional Request for "Answer To Select" (RATS) and
// "Protocol Parameter Selection" (PPS).
// see "Card Activation" in https://www.nxp.com/docs/en/application-note/AN10834.pdf.
// note: the card needs to be in ready state, e.g. by a request or wake up is done before
func (d *MFRC522Common) piccActivate() ([]byte, error) {
	if err := d.clearRegisterBitMask(regColl, collRegValuesAfterCollBit); err != nil {
		return nil, err
	}
	if err := d.writeByteData(regBitFraming, bitFramingRegReset); err != nil {
		return nil, err
	}

	// start cascade level 1 (0x93) for return:
	// * one size UID (4 byte): UID0..3 and one byte BCC or
	// * cascade tag (0x88) and UID0..2 and BCC
	// in the latter case the UID is incomplete and the next cascade level needs to be started.
	// cascade level 2 (0x95) return:
	// * double size UID (7 byte): UID3..6 and one byte BCC or
	// * cascade tag (0x88) and UID3..5 and BCC
	// cascade level 3 (0x97) return:
	// * triple size UID (10 byte): UID6..9
	// after each anticollision check (request of next UID) the SAK needs to be done (same level command)
	// BCC: Block Check Character
	// SAK: Select Acknowledge

	var uid []byte
	var sak uint8

	for cascadeLevel := 1; cascadeLevel < 3; cascadeLevel++ {
		var piccCommand uint8
		switch cascadeLevel {
		case 1:
			piccCommand = piccCommandCascadeLevel1
		case 2:
			piccCommand = piccCommandCascadeLevel2
		case 3:
			piccCommand = piccCommandCascadeLevel3
		default:
			return nil, fmt.Errorf("unknown cascade level %d", cascadeLevel)
		}

		if piccDebug {
			fmt.Println("-anti collision-")
		}

		txLastBits := uint8(0x00) // we use all 8 bits
		numValidBits := uint8(4 * 8)
		sendForAnticol := []byte{piccCommand, numValidBits}
		backData := []byte{0x00, 0x00, 0x00, 0x00, 0x00} // 4 bytes CT/UID and BCC
		if err := d.communicateWithPICC(commandRegTransceive, sendForAnticol, backData, txLastBits, false); err != nil {
			return nil, err
		}

		// TODO: no real anticollision check yet

		// check BCC
		bcc := byte(0)
		for _, v := range backData[:4] {
			bcc = bcc ^ v
		}
		if bcc != backData[4] {
			return nil, fmt.Errorf("BCC mismatch, expected %02x actual %02x", bcc, backData[4])
		}

		if backData[0] == piccCascadeTag {
			uid = append(uid, backData[1:3]...)
			if piccDebug {
				fmt.Printf("next cascade is needed after SAK, uid: %v", uid)
			}
		} else {
			uid = append(uid, backData[:4]...)
			if piccDebug {
				fmt.Printf("backData: %v, uid: %v\n", backData, uid)
			}
		}

		if piccDebug {
			fmt.Println("-select acknowledge-")
		}
		sendCommand := []byte{piccCommand}
		sendCommand = append(sendCommand, 0x70)        // 7 bytes
		sendCommand = append(sendCommand, backData...) // uid including BCC
		crcResult := []byte{0x00, 0x00}
		if err := d.calculateCRC(sendCommand, crcResult); err != nil {
			return uid, err
		}
		sendCommand = append(sendCommand, crcResult...)
		sakData := []byte{0x00, 0x00, 0x00}
		if err := d.communicateWithPICC(commandRegTransceive, sendCommand, sakData, txLastBits, false); err != nil {
			return nil, err
		}
		bcc = byte(0)
		for _, v := range sakData[:2] {
			bcc = bcc ^ v
		}
		if piccDebug {
			fmt.Printf("sak data: %v\n", sakData)
		}
		if sakData[0] != piccUIDNotComplete {
			sak = sakData[0]
			break
		}
		if piccDebug {
			fmt.Printf("next cascade called, SAK: %v\n", sakData[0])
		}
	}

	if piccDebug || d.firstCardAccess {
		d.firstCardAccess = false
		fmt.Printf("card '%s' selected\n", piccCardFromSak[sak])
	}
	return uid, nil
}

func (d *MFRC522Common) piccRequest(reqMode uint8, answer []byte) error {
	if len(answer) < 2 {
		return fmt.Errorf("at least 2 bytes room needed for the answer")
	}

	if err := d.clearRegisterBitMask(regColl, collRegValuesAfterCollBit); err != nil {
		return err
	}

	// for request A and wake up the short frame format is used - transmit only 7 bits of the last (and only) byte.
	txLastBits := uint8(0x07 & bitFramingRegTxLastBits)
	if err := d.communicateWithPICC(commandRegTransceive, []byte{reqMode}, answer, txLastBits, false); err != nil {
		return err
	}

	return nil
}
