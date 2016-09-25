package i2c

import "github.com/hybridgroup/gobot"

var _ gobot.Driver = (*GroveLcdDriver)(nil)
var _ gobot.Driver = (*GroveAccelerometerDriver)(nil)

// GroveLcdDriver is a driver for the Jhd1313m1 LCD display which has two i2c addreses,
// one belongs to a controller and the other controls solely the backlight.
// This module was tested with the Seed Grove LCD RGB Backlight v2.0 display which requires 5V to operate.
// http://www.seeedstudio.com/wiki/Grove_-_LCD_RGB_Backlight
type GroveLcdDriver struct {
	*JHD1313M1Driver
}

// NewGroveLcdDriver creates a new driver with specified name and i2c interface.
func NewGroveLcdDriver(a I2c, name string) *GroveLcdDriver {
	return &GroveLcdDriver{
		JHD1313M1Driver: NewJHD1313M1Driver(a, name),
	}
}

type GroveAccelerometerDriver struct {
	*MMA7660Driver
}

// NewGroveAccelerometerDriver creates a new driver with specified name and i2c interface
func NewGroveAccelerometerDriver(a I2c, name string) *GroveAccelerometerDriver {
	return &GroveAccelerometerDriver{
		MMA7660Driver: NewMMA7660Driver(a, name),
	}
}
