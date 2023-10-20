package i2c

// PCA9501 supports addresses from 0x00 to 0x7F
// 0x00 - 0x3F: GPIO, 0x40 - 0x7F: EEPROM
//
// 0 EE A5 A4 A3 A2 A1 A0|rd
// Lowest bit (rd) is mapped to switch between write(0)/read(1), it is not part of the "real" address.
// Highest bit (EE) is mapped to switch between GPIO(0)/EEPROM(1).
//
// The EEPROM address will be generated from GPIO address in this driver.
const pca9501DefaultAddress = 0x3F // this applies, if all 6 address pins left open (have pull up resistors)

// PCA9501Driver is a Gobot Driver for the PCA9501 8-bit GPIO & 2-kbit EEPROM with 6 address program pins.
// 2-kbit EEPROM has 256 byte, means addresses between 0x00-0xFF
//
// please refer to data sheet: https://www.nxp.com/docs/en/data-sheet/PCA9501.pdf
//
// PCA9501 is the replacement for PCF8574, so this driver should also work for PCF8574 except EEPROM calls
type PCA9501Driver struct {
	connectionMem Connection
	*Driver
}

// NewPCA9501Driver creates a new driver with specified i2c interface
// Params:
//
//	a Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewPCA9501Driver(a Connector, options ...func(Config)) *PCA9501Driver {
	p := &PCA9501Driver{
		Driver: NewDriver(a, "PCA9501", pca9501DefaultAddress, options...),
	}
	p.afterStart = p.initialize

	// API commands
	p.AddCommand("WriteGPIO", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(uint8)
		val := params["val"].(uint8)
		err := p.WriteGPIO(pin, val)
		return map[string]interface{}{"err": err}
	})

	p.AddCommand("ReadGPIO", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(uint8)
		val, err := p.ReadGPIO(pin)
		return map[string]interface{}{"val": val, "err": err}
	})

	p.AddCommand("WriteEEPROM", func(params map[string]interface{}) interface{} {
		address := params["address"].(uint8)
		val := params["val"].(uint8)
		err := p.WriteEEPROM(address, val)
		return map[string]interface{}{"err": err}
	})

	p.AddCommand("ReadEEPROM", func(params map[string]interface{}) interface{} {
		address := params["address"].(uint8)
		val, err := p.ReadEEPROM(address)
		return map[string]interface{}{"val": val, "err": err}
	})
	return p
}

// WriteGPIO writes a value to a gpio pin (0-7)
func (p *PCA9501Driver) WriteGPIO(pin uint8, val uint8) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// read current value of CTRL register, 0 is output, 1 is no output
	iodir, err := p.connection.ReadByte()
	if err != nil {
		return err
	}
	// set pin as output by clearing bit
	iodirVal := clearBit(iodir, pin)
	// write CTRL register
	err = p.connection.WriteByte(iodirVal)
	if err != nil {
		return err
	}
	// read current value of port
	cVal, err := p.connection.ReadByte()
	if err != nil {
		return err
	}
	// set or reset the bit in value
	var nVal uint8
	if val == 0 {
		nVal = clearBit(cVal, pin)
	} else {
		nVal = setBit(cVal, pin)
	}
	// write new value to port
	err = p.connection.WriteByte(nVal)
	if err != nil {
		return err
	}
	return nil
}

// ReadGPIO reads a value from a given gpio pin (0-7)
func (p *PCA9501Driver) ReadGPIO(pin uint8) (uint8, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// read current value of CTRL register, 0 is no input, 1 is an input
	iodir, err := p.connection.ReadByte()
	if err != nil {
		return 0, err
	}
	// set pin as input by setting bit
	iodirVal := setBit(iodir, pin)
	// write CTRL register
	err = p.connection.WriteByte(iodirVal)
	if err != nil {
		return 0, err
	}
	// read port and create return bit
	val, err := p.connection.ReadByte()
	if err != nil {
		return val, err
	}
	val = 1 << pin & val
	if val > 1 {
		val = 1
	}
	return val, nil
}

// ReadEEPROM reads a value from a given address (0x00-0xFF)
// Note: only this sequence for memory read is supported: "STARTW-DATA1-STARTR-DATA2-STOP"
// DATA1: EEPROM address, DATA2: read value
func (p *PCA9501Driver) ReadEEPROM(address uint8) (uint8, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.connectionMem.ReadByteData(address)
}

// WriteEEPROM writes a value to a given address in memory (0x00-0xFF)
func (p *PCA9501Driver) WriteEEPROM(address uint8, val uint8) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.connectionMem.WriteByteData(address, val)
}

func (p *PCA9501Driver) initialize() (err error) {
	// initialize the EEPROM connection
	bus := p.GetBusOrDefault(p.connector.DefaultI2cBus())
	addressMem := p.GetAddressOrDefault(pca9501DefaultAddress) | 0x40
	p.connectionMem, err = p.connector.GetI2cConnection(addressMem, bus)
	return
}
