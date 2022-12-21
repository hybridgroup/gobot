package mfrc522

import "fmt"

func (d *MFRC522Common) readByteData(reg uint8) (uint8, error) {
	if d.connection == nil {
		return 0, fmt.Errorf("not connected")
	}
	return d.connection.ReadByteData(reg)
}

func (d *MFRC522Common) writeByteData(reg uint8, data uint8) error {
	if d.connection == nil {
		return fmt.Errorf("not connected")
	}
	return d.connection.WriteByteData(reg, data)
}

func (d *MFRC522Common) setRegisterBitMask(reg uint8, mask uint8) error {
	val, err := d.readByteData(reg)
	if err != nil {
		return err
	}
	if err := d.writeByteData(reg, val|mask); err != nil {
		return err
	}
	return nil
}

func (d *MFRC522Common) clearRegisterBitMask(reg uint8, mask uint8) error {
	val, err := d.readByteData(reg)
	if err != nil {
		return err
	}
	if err := d.writeByteData(reg, val&^mask); err != nil {
		return err
	}
	return nil
}
