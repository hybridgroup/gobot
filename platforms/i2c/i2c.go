package i2c

type I2cInterface interface {
	I2cStart(address byte) (err error)
	I2cRead(len uint) (data []byte, err error)
	I2cWrite(buf []byte) (err error)
}
