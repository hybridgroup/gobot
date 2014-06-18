package i2c

type I2cInterface interface {
	I2cStart(byte)
	I2cRead(uint) []byte
	I2cWrite([]byte)
}
