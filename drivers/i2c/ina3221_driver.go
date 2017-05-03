package i2c

import (
	"gobot.io/x/gobot"
)

type Channel uint8

const (
	// I2C ADDRESS/BITS
	INA3221_ADDRESS             uint8   = 0x40 //1000000 (A0+A1=GND)
	INA3221_READ                uint8   = 0x01
	INA3221_REG_CONFIG          uint8   = 0x00   // CONFIG REGISTER (R/W)
	INA3221_CONFIG_RESET        uint16  = 0x8000 //Reset Bit
	INA3221_CONFIG_ENABLE_CHAN1 uint16  = 0x4000 //Enable Channel 1
	INA3221_CONFIG_ENABLE_CHAN2 uint16  = 0x2000 //Enable Channel 2
	INA3221_CONFIG_ENABLE_CHAN3 uint16  = 0x1000 //Enable Channel 3
	INA3221_CONFIG_AVG2         uint16  = 0x0800 //AVG Samples Bit 2 - See table 3 spec
	INA3221_CONFIG_AVG1         uint16  = 0x0400 //AVG Samples Bit 1 - See table 3 spec
	INA3221_CONFIG_AVG0         uint16  = 0x0200 //AVG Samples Bit 0 - See table 3 spec
	INA3221_CONFIG_VBUS_CT2     uint16  = 0x0100 //VBUS bit 2 Conversion time - See table 4 spec
	INA3221_CONFIG_VBUS_CT1     uint16  = 0x0080 //VBUS bit 1 Conversion time - See table 4 spec
	INA3221_CONFIG_VBUS_CT0     uint16  = 0x0040 //VBUS bit 0 Conversion time - See table 4 spec
	INA3221_CONFIG_VSH_CT2      uint16  = 0x0020 //Vshunt bit 2 Conversion time - See table 5 spec
	INA3221_CONFIG_VSH_CT1      uint16  = 0x0010 //Vshunt bit 1 Conversion time - See table 5 spec
	INA3221_CONFIG_VSH_CT0      uint16  = 0x0008 //Vshunt bit 0 Conversion time - See table 5 spec
	INA3221_CONFIG_MODE_2       uint16  = 0x0004 //Operating Mode bit 2 - See table 6 spec
	INA3221_CONFIG_MODE_1       uint16  = 0x0002 //Operating Mode bit 1 - See table 6 spec
	INA3221_CONFIG_MODE_0       uint16  = 0x0001 //Operating Mode bit 0 - See table 6 spec
	INA3221_REG_SHUNTVOLTAGE_1  uint8   = 0x01   // SHUNT VOLTAGE REGISTER (R)
	INA3221_REG_BUSVOLTAGE_1    uint8   = 0x02   // BUS VOLTAGE REGISTER (R)
	SHUNT_RESISTOR_VALUE        float64 = 0.1    //default shunt resistor value of 0.1 Ohm

	INA3221Channel1 Channel = 1
	INA3221Channel2 Channel = 2
	INA3221Channel3 Channel = 3
)

type INA3221Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	halt chan bool
}

func NewINA3221Driver(c Connector, options ...func(Config)) *INA3221Driver {
	i := &INA3221Driver{
		name:      "Ina3221",
		connector: c,
		Config: NewConfig(),
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
	address := i.GetAddressOrDefault(int(INA3221_ADDRESS))

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
func (i *INA3221Driver) GetBusVoltage(channel Channel) (float64, error) {
	value, err := i.getBusVoltageRaw(channel)
	if err != nil {
		return 0, err
	}

	return float64(value) * .001, nil
}

// GetShuntVoltage Gets the shunt voltage in mV
func (i *INA3221Driver) GetShuntVoltage(channel Channel) (float64, error) {
	value, err := i.getShuntVoltageRaw(channel)
	if err != nil {
		return 0, err
	}

	return float64(value) * float64(.005), nil
}

// GetCurrent gets the current value in mA, taking into account the config settings and current LSB
func (i *INA3221Driver) GetCurrent(channel Channel) (float64, error) {
	value, err := i.GetShuntVoltage(channel)
	if err != nil {
		return 0, err
	}

	ma := value / SHUNT_RESISTOR_VALUE
	return ma, nil
}

// GetLoadVoltage gets the load voltage in mV
func (i *INA3221Driver) GetLoadVoltage(channel Channel) (float64, error) {
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
func (i *INA3221Driver) getBusVoltageRaw(channel Channel) (int16, error) {
	val, err := i.connection.ReadWordData(INA3221_REG_BUSVOLTAGE_1 + (uint8(channel)-1)*2)
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
func (i *INA3221Driver) getShuntVoltageRaw(channel Channel) (int16, error) {
	v, err := i.connection.ReadWordData(INA3221_REG_SHUNTVOLTAGE_1 + (uint8(channel)-1)*2)
	if err != nil {
		return 0, err
	}

	value := int32(v)
	if value > 0x7FFF {
		value -= 0x10000
	}

	return int16(value), nil
}

// initialize initializes the device
func (i *INA3221Driver) initialize() error {
	config := INA3221_CONFIG_ENABLE_CHAN1 |
		INA3221_CONFIG_ENABLE_CHAN2 |
		INA3221_CONFIG_ENABLE_CHAN3 |
		INA3221_CONFIG_AVG1 |
		INA3221_CONFIG_VBUS_CT2 |
		INA3221_CONFIG_VSH_CT2 |
		INA3221_CONFIG_MODE_2 |
		INA3221_CONFIG_MODE_1 |
		INA3221_CONFIG_MODE_0

	return i.connection.WriteWordData(INA3221_REG_CONFIG, config)
}
