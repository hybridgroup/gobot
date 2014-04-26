package gobotI2C

type TestAdaptor struct {
}

func (t TestAdaptor) I2cStart(byte) {
	return
}

func (t TestAdaptor) I2cRead(uint16) []uint16 {
	return []uint16{99, 1}
}

func (t TestAdaptor) I2cWrite([]uint16) {
	return
}
