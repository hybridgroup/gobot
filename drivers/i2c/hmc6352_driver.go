package i2c

const hmc6352DefaultAddress = 0x21

// HMC6352Driver is a Driver for a HMC6352 digital compass
type HMC6352Driver struct {
	*Driver
}

// NewHMC6352Driver creates a new driver with specified i2c interface
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewHMC6352Driver(c Connector, options ...func(Config)) *HMC6352Driver {
	d := &HMC6352Driver{
		Driver: NewDriver(c, "HMC6352", hmc6352DefaultAddress),
	}
	d.afterStart = d.initialize

	for _, option := range options {
		option(d)
	}

	return d
}

// Heading returns the current heading
func (d *HMC6352Driver) Heading() (uint16, error) {
	if _, err := d.write([]byte("A")); err != nil {
		return 0, err
	}
	buf := []byte{0, 0}
	bytesRead, err := d.read(buf)
	if err != nil {
		return 0, err
	}
	if bytesRead == 2 {
		heading := (uint16(buf[1]) + uint16(buf[0])*256) / 10
		return heading, nil
	}

	return 0, ErrNotEnoughBytes
}

func (d *HMC6352Driver) initialize() error {
	if _, err := d.write([]byte("A")); err != nil {
		return err
	}
	return nil
}
