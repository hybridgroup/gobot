package i2c

// INA3221Driver is a driver for the Texas Instruments INA3221 device. The INA3221 is a three-channel
// current and bus voltage monitor with an I2C and SMBUS compatible interface.
//
// INA3221 data sheet and specifications can be found at http://www.ti.com/product/INA3221
//
// This module was tested with SwitchDoc Labs INA3221 breakout board found at http://www.switchdoc.com/

import (
	"gobot.io/x/gobot"
)

// INA3221Channel type that defines which INA3221 channel to read from.
type INA3221Channel uint8

const (
	ina3221Address            uint8   = 0x40 // 1000000 (A0+A1=GND)
	ina3221Read               uint8   = 0x01
	ina3221RegConfig          uint8   = 0x00   // CONFIG REGISTER (R/W)
	ina3221ConfigReset        uint16  = 0x8000 // Reset Bit
	ina3221ConfigEnableChan1  uint16  = 0x4000 // Enable INA3221 Channel 1
	ina3221ConfigEnableChan2  uint16  = 0x2000 // Enable INA3221 Channel 2
	ina3221ConfigEnableChan3  uint16  = 0x1000 // Enable INA3221 Channel 3
	ina3221ConfigAvg2         uint16  = 0x0800 // AVG Samples Bit 2 - See table 3 spec
	ina3221ConfigAvg1         uint16  = 0x0400 // AVG Samples Bit 1 - See table 3 spec
	ina3221ConfigAvg0         uint16  = 0x0200 // AVG Samples Bit 0 - See table 3 spec
	ina3221ConfigVBusCT2      uint16  = 0x0100 // VBUS bit 2 Conversion time - See table 4 spec
	ina3221ConfigVBusCT1      uint16  = 0x0080 // VBUS bit 1 Conversion time - See table 4 spec
	ina3221ConfigVBusCT0      uint16  = 0x0040 // VBUS bit 0 Conversion time - See table 4 spec
	ina3221ConfigVShCT2       uint16  = 0x0020 // Vshunt bit 2 Conversion time - See table 5 spec
	ina3221ConfigVShCT1       uint16  = 0x0010 // Vshunt bit 1 Conversion time - See table 5 spec
	ina3221ConfigVShCT0       uint16  = 0x0008 // Vshunt bit 0 Conversion time - See table 5 spec
	ina3221ConfigMode2        uint16  = 0x0004 // Operating Mode bit 2 - See table 6 spec
	ina3221ConfigMode1        uint16  = 0x0002 // Operating Mode bit 1 - See table 6 spec
	ina3221ConfigMode0        uint16  = 0x0001 // Operating Mode bit 0 - See table 6 spec
	ina3221RegShuntVoltage1   uint8   = 0x01   // SHUNT VOLTAGE REGISTER (R)
	ina3221RegBusVoltage1     uint8   = 0x02   // BUS VOLTAGE REGISTER (R)
	ina3221ShuntResistorValue float64 = 0.1    // default shunt resistor value of 0.1 Ohm

	INA3221Channel1 INA3221Channel = 1
	INA3221Channel2 INA3221Channel = 2
	INA3221Channel3 INA3221Channel = 3
)

// INA3221Driver is a driver for the INA3221 three-channel current and bus voltage monitoring device.
type INA3221Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	halt chan bool
}

// NewINA3221Driver creates a new driver with the specified i2c interface.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):		bus to use with this driver
//		i2c.WithAddress(int):		address to use with this driver
func NewINA3221Driver(c Connector, options ...func(Config)) *INA3221Driver {
	i := &INA3221Driver{
		name:      gobot.DefaultName("INA3221"),
		connector: c,
		Config:    NewConfig(),
	}

	for _, option := range options {
		option(i)
	}

	return i
}

// Name returns the name of the device.
func (i *INA3221Driver) Name() string {
	return i.name
}

// SetName sets the name of the device.
func (i *INA3221Driver) SetName(name string) {
	i.name = name
}

// Connection returns the connection of the device.
func (i *INA3221Driver) Connection() gobot.Connection {
	return i.connector.(gobot.Connection)
}

// Start initializes the INA3221
func (i *INA3221Driver) Start() error {
	var err error
	bus := i.GetBusOrDefault(i.connector.GetDefaultBus())
	address := i.GetAddressOrDefault(int(ina3221Address))

	if i.connection, err = i.connector.GetConnection(address, bus); err != nil {
		return err
	}

	if err := i.initialize(); err != nil {
		return err
	}

	return nil
}

// Halt halts the device.
func (i *INA3221Driver) Halt() error {
	return nil
}

// GetBusVoltage gets the bus voltage in Volts
func (i *INA3221Driver) GetBusVoltage(channel INA3221Channel) (float64, error) {
	value, err := i.getBusVoltageRaw(channel)
	if err != nil {
		return 0, err
	}

	return float64(value) * .001, nil
}

// GetShuntVoltage Gets the shunt voltage in mV
func (i *INA3221Driver) GetShuntVoltage(channel INA3221Channel) (float64, error) {
	value, err := i.getShuntVoltageRaw(channel)
	if err != nil {
		return 0, err
	}

	return float64(value) * float64(.005), nil
}

// GetCurrent gets the current value in mA, taking into account the config settings and current LSB
func (i *INA3221Driver) GetCurrent(channel INA3221Channel) (float64, error) {
	value, err := i.GetShuntVoltage(channel)
	if err != nil {
		return 0, err
	}

	ma := value / ina3221ShuntResistorValue
	return ma, nil
}

// GetLoadVoltage gets the load voltage in mV
func (i *INA3221Driver) GetLoadVoltage(channel INA3221Channel) (float64, error) {
	bv, err := i.GetBusVoltage(channel)
	if err != nil {
		return 0, err
	}

	sv, err := i.GetShuntVoltage(channel)
	if err != nil {
		return 0, err
	}

	return bv + (sv / 1000.0), nil
}

// getBusVoltageRaw gets the raw bus voltage (16-bit signed integer, so +-32767)
func (i *INA3221Driver) getBusVoltageRaw(channel INA3221Channel) (int16, error) {
	val, err := i.readWordFromRegister(ina3221RegBusVoltage1 + (uint8(channel)-1)*2)
	if err != nil {
		return 0, err
	}

	value := int32(val)
	if value > 0x7FFF {
		value -= 0x10000
	}

	return int16(value), nil
}

// getShuntVoltageRaw gets the raw shunt voltage (16-bit signed integer, so +-32767)
func (i *INA3221Driver) getShuntVoltageRaw(channel INA3221Channel) (int16, error) {
	val, err := i.readWordFromRegister(ina3221RegShuntVoltage1 + (uint8(channel)-1)*2)
	if err != nil {
		return 0, err
	}

	value := int32(val)
	if value > 0x7FFF {
		value -= 0x10000
	}

	return int16(value), nil
}

// reads word from supplied register address
func (i *INA3221Driver) readWordFromRegister(reg uint8) (uint16, error) {
	val, err := i.connection.ReadWordData(reg)
	if err != nil {
		return 0, err
	}

	return uint16(((val & 0x00FF) << 8) | ((val & 0xFF00) >> 8)), nil
}

// initialize initializes the INA3221 device
func (i *INA3221Driver) initialize() error {
	config := ina3221ConfigEnableChan1 |
		ina3221ConfigEnableChan2 |
		ina3221ConfigEnableChan3 |
		ina3221ConfigAvg1 |
		ina3221ConfigVBusCT2 |
		ina3221ConfigVShCT2 |
		ina3221ConfigMode2 |
		ina3221ConfigMode1 |
		ina3221ConfigMode0

	return i.connection.WriteBlockData(ina3221RegConfig, []byte{byte(config >> 8), byte(config & 0x00FF)})
}
