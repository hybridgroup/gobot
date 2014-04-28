package i2c

type I2cInterface interface {
	I2cStart(byte)
	I2cRead(uint16) []uint16
	I2cWrite([]uint16)
}
