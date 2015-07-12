package i2c

import "github.com/hybridgroup/gobot"

var _ gobot.Driver = (*MMA7660Driver)(nil)

const mma7660Address = 0x4c

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
	name       string
	connection I2c
}

// NewMMA7660Driver creates a new driver with specified name and i2c interface
func NewMMA7660Driver(a I2c, name string) *MMA7660Driver {
	return &MMA7660Driver{
		name:       name,
		connection: a,
	}
}

func (h *MMA7660Driver) Name() string                 { return h.name }
func (h *MMA7660Driver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start initialized the mma7660
func (h *MMA7660Driver) Start() (errs []error) {
	if err := h.connection.I2cStart(mma7660Address); err != nil {
		return []error{err}
	}

	if err := h.connection.I2cWrite(mma7660Address, []byte{MMA7660_MODE, MMA7660_STAND_BY}); err != nil {
		return []error{err}
	}

	if err := h.connection.I2cWrite(mma7660Address, []byte{MMA7660_SR, MMA7660_AUTO_SLEEP_32}); err != nil {
		return []error{err}
	}

	if err := h.connection.I2cWrite(mma7660Address, []byte{MMA7660_MODE, MMA7660_ACTIVE}); err != nil {
		return []error{err}
	}

	return
}

// Halt returns true if devices is halted successfully
func (h *MMA7660Driver) Halt() (errs []error) { return }

// Acceleration returns the acceleration  of the provided x, y, z
func (h *MMA7660Driver) Acceleration(x, y, z float64) (ax, ay, az float64) {
	return x / 21.0, y / 21.0, z / 21.0
}

// XYZ returns the raw x,y and z axis from the  mma7660
func (h *MMA7660Driver) XYZ() (x float64, y float64, z float64, err error) {
	ret, err := h.connection.I2cRead(mma7660Address, 3)
	if err != nil {
		return
	}

	if len(ret) != 3 {
		err = ErrNotEnoughBytes
		return
	}

	for _, val := range ret {
		if ((val >> 6) & 0x01) == 1 {
			err = ErrNotReady
			return
		}
	}

	x = float64((int8(ret[0]) << 2)) / 4.0
	y = float64((int8(ret[1]) << 2)) / 4.0
	z = float64((int8(ret[2]) << 2)) / 4.0

	return
}
