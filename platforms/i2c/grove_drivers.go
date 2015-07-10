package i2c

import "github.com/hybridgroup/gobot"

var _ gobot.Driver = (*GroveLcdDriver)(nil)
var _ gobot.Driver = (*GroveAccelerometerDriver)(nil)

type GroveLcdDriver struct {
	*JHD1313M1Driver
}

type GroveAccelerometerDriver struct {
	*MMA7660Driver
}

func NewGroveLcdDriver(a I2c, name string) *GroveLcdDriver {
	return &GroveLcdDriver{
		JHD1313M1Driver: NewJHD1313M1Driver(a, name),
	}
}

func NewGroveAccelerometerDriver(a I2c, name string) *GroveAccelerometerDriver {
	return &GroveAccelerometerDriver{
		MMA7660Driver: NewMMA7660Driver(a, name),
	}
}
