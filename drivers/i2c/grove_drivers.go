package i2c

// GroveLcdDriver is a driver for the Jhd1313m1 LCD display which has two i2c addreses,
// one belongs to a controller and the other controls solely the backlight.
// This module was tested with the Seed Grove LCD RGB Backlight v2.0 display which requires 5V to operate.
// http://www.seeedstudio.com/wiki/Grove_-_LCD_RGB_Backlight
type GroveLcdDriver struct {
	*JHD1313M1Driver
}

// NewGroveLcdDriver creates a new driver with specified i2c interface.
// Params:
//		conn I2cConnector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.Bus(int):	bus to use with this driver
//		i2c.Address(int):	address to use with this driver
//
func NewGroveLcdDriver(a I2cConnector, options ...func(I2cConfig)) *GroveLcdDriver {
	lcd := &GroveLcdDriver{
		JHD1313M1Driver: NewJHD1313M1Driver(a),
	}

	for _, option := range options {
		option(lcd)
	}

	return lcd
}

type GroveAccelerometerDriver struct {
	*MMA7660Driver
}

// NewGroveAccelerometerDriver creates a new driver with specified i2c interface
func NewGroveAccelerometerDriver(a I2cConnector, options ...func(I2cConfig)) *GroveAccelerometerDriver {
	mma := &GroveAccelerometerDriver{
		MMA7660Driver: NewMMA7660Driver(a),
	}

	for _, option := range options {
		option(mma)
	}

	return mma
}
