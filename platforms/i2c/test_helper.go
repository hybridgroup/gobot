package i2c

type TestAdaptor struct {
}

func (t TestAdaptor) I2cStart(byte) {
	return
}

func (t TestAdaptor) I2cRead(uint) []byte {
	return []byte{99, 1}
}

func (t TestAdaptor) I2cWrite([]byte) {
	return
}
