package i2c

const mma7660DefaultAddress = 0x4c

const (
	MMA7660_X              = 0x00
	MMA7660_Y              = 0x01
	MMA7660_Z              = 0x02
	MMA7660_TILT           = 0x03
	MMA7660_SRST           = 0x04
	MMA7660_SPCNT          = 0x05
	MMA7660_INTSU          = 0x06
	MMA7660_MODE           = 0x07
	MMA7660_STAND_BY       = 0x00
	MMA7660_ACTIVE         = 0x01
	MMA7660_SR             = 0x08
	MMA7660_AUTO_SLEEP_120 = 0x00
	MMA7660_AUTO_SLEEP_64  = 0x01
	MMA7660_AUTO_SLEEP_32  = 0x02
	MMA7660_AUTO_SLEEP_16  = 0x03
	MMA7660_AUTO_SLEEP_8   = 0x04
	MMA7660_AUTO_SLEEP_4   = 0x05
	MMA7660_AUTO_SLEEP_2   = 0x06
	MMA7660_AUTO_SLEEP_1   = 0x07
	MMA7660_PDET           = 0x09
	MMA7660_PD             = 0x0A
)

type MMA7660Driver struct {
	*Driver
}

// NewMMA7660Driver creates a new driver with specified i2c interface
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewMMA7660Driver(c Connector, options ...func(Config)) *MMA7660Driver {
	d := &MMA7660Driver{
		Driver: NewDriver(c, "MMA7660", mma7660DefaultAddress),
	}
	d.afterStart = d.initialize

	for _, option := range options {
		option(d)
	}

	// TODO: add commands for API
	return d
}

// Acceleration returns the acceleration of the provided x, y, z
//
//nolint:nonamedreturns // is sufficient here
func (d *MMA7660Driver) Acceleration(x, y, z float64) (ax, ay, az float64) {
	return x / 21.0, y / 21.0, z / 21.0
}

// XYZ returns the raw x,y and z axis from the mma7660
//
//nolint:nonamedreturns // is sufficient here
func (d *MMA7660Driver) XYZ() (x float64, y float64, z float64, err error) {
	buf := []byte{0, 0, 0}
	bytesRead, err := d.connection.Read(buf)
	if err != nil {
		return
	}

	if bytesRead != 3 {
		err = ErrNotEnoughBytes
		return
	}

	for _, val := range buf {
		if ((val >> 6) & 0x01) == 1 {
			err = ErrNotReady
			return
		}
	}

	x = float64((int8(buf[0]) << 2)) / 4.0
	y = float64((int8(buf[1]) << 2)) / 4.0
	z = float64((int8(buf[2]) << 2)) / 4.0

	return
}

func (d *MMA7660Driver) initialize() error {
	if _, err := d.connection.Write([]byte{MMA7660_MODE, MMA7660_STAND_BY}); err != nil {
		return err
	}

	if _, err := d.connection.Write([]byte{MMA7660_SR, MMA7660_AUTO_SLEEP_32}); err != nil {
		return err
	}

	if _, err := d.connection.Write([]byte{MMA7660_MODE, MMA7660_ACTIVE}); err != nil {
		return err
	}

	return nil
}
