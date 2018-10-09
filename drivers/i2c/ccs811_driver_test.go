package i2c

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

// The CCS811 Meets the Driver Definition
var _ gobot.Driver = (*CCS811Driver)(nil)

// --------- HELPERS
func initTestCCS811Driver() (driver *CCS811Driver) {
	driver, _ = initTestCCS811DriverWithStubbedAdaptor()
	return
}

func initTestCCS811DriverWithStubbedAdaptor() (*CCS811Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewCCS811Driver(adaptor), adaptor
}

// --------- BASE TESTS
func TestNewCCS811Driver(t *testing.T) {
	// Does it return a pointer to an instance of CCS811Driver?
	var c interface{} = NewCCS811Driver(newI2cTestAdaptor())
	_, ok := c.(*CCS811Driver)
	if !ok {
		t.Errorf("NewCCS811Driver() should have returned a *CCS811Driver")
	}
}

func TestCCS811DriverSetName(t *testing.T) {
	// Does it change the name of the driver
	d := initTestCCS811Driver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestCCS811Connection(t *testing.T) {
	// Does it create an instance of gobot.Connection
	ccs811 := initTestCCS811Driver()
	gobottest.Refute(t, ccs811.Connection(), nil)
}

// // --------- CONFIG OVERIDE TESTS

func TestCCS811DriverWithBus(t *testing.T) {
	// Can it update the bus
	d := NewCCS811Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestCCS811DriverWithAddress(t *testing.T) {
	// Can it update the address
	d := NewCCS811Driver(newI2cTestAdaptor(), WithAddress(0xFF))
	gobottest.Assert(t, d.GetAddressOrDefault(0x5a), 0xFF)
}

func TestCCS811DriverWithCCS811MeasMode(t *testing.T) {
	// Can it update the measurement mode
	d := NewCCS811Driver(newI2cTestAdaptor(), WithCCS811MeasMode(CCS811DriveMode10Sec))
	gobottest.Assert(t, d.measMode.driveMode, CCS811DriveMode(CCS811DriveMode10Sec))
}

func TestCCS811DriverWithCCS811NTCResistance(t *testing.T) {
	// Can it update the ntc resitor value used for temp calcuations
	d := NewCCS811Driver(newI2cTestAdaptor(), WithCCS811NTCResistance(0xFF))
	gobottest.Assert(t, d.ntcResistanceValue, uint32(0xFF))
}

// // --------- DRIVER SPECIFIC TESTS

func TestCCS811DriverGetGasData(t *testing.T) {

	cases := []struct {
		readReturn func([]byte) (int, error)
		eco2       uint16
		tvoc       uint16
		err        error
	}{
		// Can it compute the gas data with ideal values taken from the bus
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{1, 156, 0, 86})
				return 4, nil
			},
			eco2: 412,
			tvoc: 86,
			err:  nil,
		},
		// Can it compute the gas data with the max values possible taken from the bus
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{255, 255, 255, 255})
				return 4, nil
			},
			eco2: 65535,
			tvoc: 65535,
			err:  nil,
		},
		// Does it return an error when the i2c operation fails
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{255, 255, 255, 255})
				return 4, errors.New("Error")
			},
			eco2: 0,
			tvoc: 0,
			err:  errors.New("Error"),
		},
	}

	d, adaptor := initTestCCS811DriverWithStubbedAdaptor()
	// Create stub function as it is needed by read submethod in driver code
	adaptor.i2cWriteImpl = func([]byte) (int, error) { return 0, nil }

	d.Start()
	for _, c := range cases {
		adaptor.i2cReadImpl = c.readReturn
		eco2, tvoc, err := d.GetGasData()
		gobottest.Assert(t, eco2, c.eco2)
		gobottest.Assert(t, tvoc, c.tvoc)
		gobottest.Assert(t, err, c.err)
	}

}

func TestCCS811DriverGetTemperature(t *testing.T) {

	cases := []struct {
		readReturn func([]byte) (int, error)
		temp       float32
		err        error
	}{
		// Can it compute the temperature data with ideal values taken from the bus
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{10, 197, 0, 248})
				return 4, nil
			},
			temp: 27.811005,
			err:  nil,
		},
		// Can it compute the temperature data without bus values overflowing
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{129, 197, 10, 248})
				return 4, nil
			},
			temp: 29.48822,
			err:  nil,
		},
		// Can it compute a negative temperature
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{255, 255, 255, 255})
				return 4, nil
			},
			temp: -25.334152,
			err:  nil,
		},
		// Does it return an error if the i2c bus errors
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{129, 197, 0, 248})
				return 4, errors.New("Error")
			},
			temp: 0,
			err:  errors.New("Error"),
		},
	}

	d, adaptor := initTestCCS811DriverWithStubbedAdaptor()
	// Create stub function as it is needed by read submethod in driver code
	adaptor.i2cWriteImpl = func([]byte) (int, error) { return 0, nil }

	d.Start()
	for _, c := range cases {
		adaptor.i2cReadImpl = c.readReturn
		temp, err := d.GetTemperature()
		gobottest.Assert(t, temp, c.temp)
		gobottest.Assert(t, err, c.err)
	}

}

func TestCCS811DriverHasData(t *testing.T) {

	cases := []struct {
		readReturn func([]byte) (int, error)
		result     bool
		err        error
	}{
		// Does it return true for HasError = 0 and DataRdy = 1
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x08})
				return 1, nil
			},
			result: true,
			err:    nil,
		},
		// Does it return false for HasError = 1 and DataRdy = 1
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x09})
				return 1, nil
			},
			result: false,
			err:    nil,
		},
		// Does it return false for HasError = 1 and DataRdy = 0
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x01})
				return 1, nil
			},
			result: false,
			err:    nil,
		},
		// Does it return false for HasError = 0 and DataRdy = 0
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x00})
				return 1, nil
			},
			result: false,
			err:    nil,
		},
		// Does it return an error when the i2c read operation fails
		{
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x00})
				return 1, errors.New("Error")
			},
			result: false,
			err:    errors.New("Error"),
		},
	}

	d, adaptor := initTestCCS811DriverWithStubbedAdaptor()
	// Create stub function as it is needed by read submethod in driver code
	adaptor.i2cWriteImpl = func([]byte) (int, error) { return 0, nil }

	d.Start()
	for _, c := range cases {
		adaptor.i2cReadImpl = c.readReturn
		result, err := d.HasData()
		gobottest.Assert(t, result, c.result)
		gobottest.Assert(t, err, c.err)
	}

}
