package i2c

// GroveLcdDriver is a driver for the Jhd1313m1 LCD display which has two i2c addreses,
// one belongs to a controller and the other controls solely the backlight.
// This module was tested with the Seed Grove LCD RGB Backlight v2.0 display which requires 5V to operate.
// http://www.seeedstudio.com/wiki/Grove_-_LCD_RGB_Backlight
type GroveLcdDriver struct {
	*JHD1313M1Driver
}

// GroveAccelerometerDriver is a driver for the MMA7660 accelerometer
type GroveAccelerometerDriver struct {
	*MMA7660Driver
}

// NewGroveLcdDriver creates a new driver with specified i2c interface.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewGroveLcdDriver(a Connector, options ...func(Config)) *GroveLcdDriver {
	lcd := &GroveLcdDriver{
		JHD1313M1Driver: NewJHD1313M1Driver(a),
	}

	for _, option := range options {
		option(lcd)
	}

	return lcd
}

// NewGroveAccelerometerDriver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewGroveAccelerometerDriver(a Connector, options ...func(Config)) *GroveAccelerometerDriver {
	mma := &GroveAccelerometerDriver{
		MMA7660Driver: NewMMA7660Driver(a),
	}

	for _, option := range options {
		option(mma)
	}

	return mma
}
