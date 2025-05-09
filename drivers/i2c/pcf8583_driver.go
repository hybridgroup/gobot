package i2c

import (
	"fmt"
	"log"
	"time"
)

const (
	pcf8583Debug = false

	// PCF8583 supports addresses 0x50 and 0x51
	// The default address applies when the address pin is grounded.
	pcf8583DefaultAddress = 0x50

	// default is 0x10, when set to 0 also some free or unused RAM can be accessed
	pcf8583RamOffset = 0x10
)

// PCF8583Control is used to specify control and status register content
type PCF8583Control uint8

const (
	// registers are named according to the datasheet
	pcf8583Reg_CTRL         = iota // 0x00
	pcf8583Reg_SUBSEC_D0D1         // 0x01
	pcf8583Reg_SEC_D2D3            // 0x02
	pcf8583Reg_MIN_D4D5            // 0x03
	pcf8583Reg_HOUR                // 0x04
	pcf8583Reg_YEARDATE            // 0x05
	pcf8583Reg_WEEKDAYMONTH        // 0x06
	pcf8583Reg_TIMER               // 0x07
	pcf8583Reg_ALARMCTRL           // 0x08, offset for all alarm registers 0x09 ... 0xF

	pcf8583CtrlTimerFlag     PCF8583Control = 0x01 // 50% duty factor, seconds flag if alarm enable bit is 0
	pcf8583CtrlAlarmFlag     PCF8583Control = 0x02 // 50% duty factor, minutes flag if alarm enable bit is 0
	pcf8583CtrlAlarmEnable   PCF8583Control = 0x04 // if enabled, memory 08h is alarm control register
	pcf8583CtrlMask          PCF8583Control = 0x08 // 0: read 05h, 06h unmasked, 1: read date and month count directly
	PCF8583CtrlModeClock50   PCF8583Control = 0x10 // clock mode with 50 Hz
	PCF8583CtrlModeCounter   PCF8583Control = 0x20 // event counter mode
	PCF8583CtrlModeTest      PCF8583Control = 0x30 // test mode
	pcf8583CtrlHoldLastCount PCF8583Control = 0x40 // 0: count, 1: store and hold count in capture latches
	pcf8583CtrlStopCounting  PCF8583Control = 0x80 // 0: count, 1: stop counting, reset divider
)

// PCF8583Driver is a driver for the PCF8583 clock and calendar chip & 240 x 8-bit bit RAM with 1 address program pin.
// please refer to data sheet: https://www.nxp.com/docs/en/data-sheet/PCF8583.pdf
//
// 0 1 0 1 0 0 0 A0|rd
// Lowest bit (rd) is mapped to switch between write(0)/read(1), it is not part of the "real" address.
//
// # PCF8583 is mainly compatible to PCF8593, so this driver should also work for PCF8593 except RAM calls
//
// This driver was tested with Tinkerboard.
type PCF8583Driver struct {
	*Driver
	mode       PCF8583Control // clock 32.768kHz (default), clock 50Hz, event counter
	yearOffset int
	ramOffset  byte
}

// NewPCF8583Driver creates a new driver with specified i2c interface
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
//	i2c.WithPCF8583Mode(PCF8583Control): mode of this driver
func NewPCF8583Driver(c Connector, options ...func(Config)) *PCF8583Driver {
	d := &PCF8583Driver{
		Driver:    NewDriver(c, "PCF8583", pcf8583DefaultAddress),
		ramOffset: pcf8583RamOffset,
	}
	d.afterStart = d.initialize

	for _, option := range options {
		option(d)
	}

	// API commands
	//nolint:forcetypeassert // ok here
	d.AddCommand("WriteTime", func(params map[string]interface{}) interface{} {
		val := params["val"].(time.Time)
		err := d.WriteTime(val)
		return map[string]interface{}{"err": err}
	})
	d.AddCommand("ReadTime", func(_ map[string]interface{}) interface{} {
		val, err := d.ReadTime()
		return map[string]interface{}{"val": val, "err": err}
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("WriteCounter", func(params map[string]interface{}) interface{} {
		val := params["val"].(int32)
		err := d.WriteCounter(val)
		return map[string]interface{}{"err": err}
	})
	d.AddCommand("ReadCounter", func(_ map[string]interface{}) interface{} {
		val, err := d.ReadCounter()
		return map[string]interface{}{"val": val, "err": err}
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("WriteRAM", func(params map[string]interface{}) interface{} {
		address := params["address"].(uint8)
		val := params["val"].(uint8)
		err := d.WriteRAM(address, val)
		return map[string]interface{}{"err": err}
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("ReadRAM", func(params map[string]interface{}) interface{} {
		address := params["address"].(uint8)
		val, err := d.ReadRAM(address)
		return map[string]interface{}{"val": val, "err": err}
	})
	return d
}

// WithPCF8583Mode is used to change the mode between 32.678kHz clock, 50Hz clock, event counter
// Valid settings are of type "PCF8583Control"
func WithPCF8583Mode(mode PCF8583Control) func(Config) {
	return func(c Config) {
		d, ok := c.(*PCF8583Driver)
		if ok {
			if !mode.isClockMode() && !mode.isCounterMode() {
				panic(fmt.Sprintf("%s: mode 0x%02x is not supported", d.name, mode))
			}
			d.mode = mode
		} else if pcf8583Debug {
			log.Printf("trying to set mode for non-PCF8583Driver %v", c)
		}
	}
}

// WriteTime setup the clock registers with the given time
func (d *PCF8583Driver) WriteTime(val time.Time) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// according to chapter 7.11 of the product data sheet, the stop counting flag of the control/status register
	// must be set before, so we read the control byte before and only set/reset the stop
	ctrlRegVal, err := d.connection.ReadByteData(uint8(pcf8583Reg_CTRL))
	if err != nil {
		return err
	}
	if !PCF8583Control(ctrlRegVal).isClockMode() {
		return fmt.Errorf("%s: can't write time because the device is in wrong mode 0x%02x", d.name, ctrlRegVal)
	}
	year, month, day := val.Date()
	err = d.connection.WriteBlockData(uint8(pcf8583Reg_CTRL),
		[]byte{
			ctrlRegVal | uint8(pcf8583CtrlStopCounting),
			// sub seconds in 1/10th seconds
			pcf8583encodeBcd(uint8(val.Nanosecond() / 1000000 / 10)), //nolint:gosec // TODO: fix later
			pcf8583encodeBcd(uint8(val.Second())),                    //nolint:gosec // TODO: fix later
			pcf8583encodeBcd(uint8(val.Minute())),                    //nolint:gosec // TODO: fix later
			pcf8583encodeBcd(uint8(val.Hour())),                      //nolint:gosec // TODO: fix later
			// year, date (we keep the year counter zero and set the offset)
			pcf8583encodeBcd(uint8(day)), //nolint:gosec // TODO: fix later
			// month, weekday (not BCD): Sunday = 0, Monday = 1 ...
			uint8(val.Weekday())<<5 | pcf8583encodeBcd(uint8(month)), //nolint:gosec // TODO: fix later
		})
	if err != nil {
		return err
	}
	d.yearOffset = year
	return d.run(ctrlRegVal)
}

// ReadTime reads the clock and returns the value
func (d *PCF8583Driver) ReadTime() (time.Time, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// according to chapter 7.1 of the product data sheet, the setting of "hold last count" flag
	// is not needed when reading with auto increment
	ctrlRegVal, err := d.connection.ReadByteData(uint8(pcf8583Reg_CTRL))
	if err != nil {
		return time.Time{}, err
	}
	if !PCF8583Control(ctrlRegVal).isClockMode() {
		return time.Time{}, fmt.Errorf("%s: can't read time because the device is in wrong mode 0x%02x", d.name, ctrlRegVal)
	}
	// auto increment feature is used
	clockDataSize := 6
	data := make([]byte, clockDataSize)
	read, err := d.connection.Read(data)
	if err != nil {
		return time.Time{}, err
	}
	if read != clockDataSize {
		return time.Time{}, fmt.Errorf("%s: %d bytes read, but %d expected", d.name, read, clockDataSize)
	}
	nanos := int(pcf8583decodeBcd(data[0])) * 1000000 * 10 // sub seconds in 1/10th seconds
	seconds := int(pcf8583decodeBcd(data[1]))
	minutes := int(pcf8583decodeBcd(data[2]))
	hours := int(pcf8583decodeBcd(data[3]))
	// year, date (the device can only count 4 years)
	year := int(data[4]>>6) + d.yearOffset        // use the first two bits, no BCD
	date := int(pcf8583decodeBcd(data[4] & 0x3F)) // remove the year-bits for date
	// weekday (not used here), month
	month := time.Month(pcf8583decodeBcd(data[5] & 0x1F)) // remove the weekday-bits
	return time.Date(year, month, date, hours, minutes, seconds, nanos, time.UTC), nil
}

// WriteCounter writes the counter registers
func (d *PCF8583Driver) WriteCounter(val int32) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// we don't care of negative values here
	// according to chapter 7.11 of the product data sheet, the stop counting flag of the control/status register
	// must be set before, so we read the control byte before and only set/reset the stop
	ctrlRegVal, err := d.connection.ReadByteData(uint8(pcf8583Reg_CTRL))
	if err != nil {
		return err
	}
	if !PCF8583Control(ctrlRegVal).isCounterMode() {
		return fmt.Errorf("%s: can't write counter because the device is in wrong mode 0x%02x", d.name, ctrlRegVal)
	}
	err = d.connection.WriteBlockData(uint8(pcf8583Reg_CTRL),
		[]byte{
			ctrlRegVal | uint8(pcf8583CtrlStopCounting), // stop
			//nolint:gosec // TODO: fix later
			pcf8583encodeBcd(uint8(val % 100)), // 2 lowest digits
			//nolint:gosec // TODO: fix later
			pcf8583encodeBcd(uint8((val / 100) % 100)), // 2 middle digits
			//nolint:gosec // TODO: fix later
			pcf8583encodeBcd(uint8((val / 10000) % 100)), // 2 highest digits
		})
	if err != nil {
		return err
	}
	return d.run(ctrlRegVal)
}

// ReadCounter reads the counter registers
func (d *PCF8583Driver) ReadCounter() (int32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// according to chapter 7.1 of the product data sheet, the setting of "hold last count" flag
	// is not needed when reading with auto increment
	ctrlRegVal, err := d.connection.ReadByteData(uint8(pcf8583Reg_CTRL))
	if err != nil {
		return 0, err
	}
	if !PCF8583Control(ctrlRegVal).isCounterMode() {
		return 0, fmt.Errorf("%s: can't read counter because the device is in wrong mode 0x%02x", d.name, ctrlRegVal)
	}
	// auto increment feature is used
	counterDataSize := 3
	data := make([]byte, counterDataSize)
	read, err := d.connection.Read(data)
	if err != nil {
		return 0, err
	}
	if read != counterDataSize {
		return 0, fmt.Errorf("%s: %d bytes read, but %d expected", d.name, read, counterDataSize)
	}
	return int32(pcf8583decodeBcd(data[0])) +
		int32(pcf8583decodeBcd(data[1]))*100 +
		int32(pcf8583decodeBcd(data[2]))*10000, nil
}

// WriteRAM writes a value to a given address in memory (0x00-0xFF)
func (d *PCF8583Driver) WriteRAM(address uint8, val uint8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	realAddress := uint16(address) + uint16(d.ramOffset)
	if realAddress > 0xFF {
		return fmt.Errorf("%s: RAM address overflow %d", d.name, realAddress)
	}
	return d.connection.WriteByteData(uint8(realAddress), val)
}

// ReadRAM reads a value from a given address (0x00-0xFF)
func (d *PCF8583Driver) ReadRAM(address uint8) (uint8, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	realAddress := uint16(address) + uint16(d.ramOffset)
	if realAddress > 0xFF {
		return 0, fmt.Errorf("%s: RAM address overflow %d", d.name, realAddress)
	}
	return d.connection.ReadByteData(uint8(realAddress))
}

func (d *PCF8583Driver) run(ctrlRegVal uint8) error {
	ctrlRegVal = ctrlRegVal & ^uint8(pcf8583CtrlStopCounting) // reset stop bit
	return d.connection.WriteByteData(uint8(pcf8583Reg_CTRL), ctrlRegVal)
}

func (d *PCF8583Driver) initialize() error {
	// switch to configured mode
	ctrlRegVal, err := d.connection.ReadByteData(uint8(pcf8583Reg_CTRL))
	if err != nil {
		return err
	}
	if d.mode.isModeDiffer(PCF8583Control(ctrlRegVal)) {
		ctrlRegVal = ctrlRegVal&^uint8(PCF8583CtrlModeTest) | uint8(d.mode)
		if err = d.connection.WriteByteData(uint8(pcf8583Reg_CTRL), ctrlRegVal); err != nil {
			return err
		}
		if pcf8583Debug {
			if PCF8583Control(ctrlRegVal).isCounterMode() {
				log.Printf("%s switched to counter mode 0x%02x", d.name, ctrlRegVal)
			} else {
				log.Printf("%s switched to clock mode 0x%02x", d.name, ctrlRegVal)
			}
		}
	}
	return nil
}

func (c PCF8583Control) isClockMode() bool {
	return uint8(c)&uint8(PCF8583CtrlModeCounter) == 0
}

func (c PCF8583Control) isCounterMode() bool {
	counterModeSet := (uint8(c) & uint8(PCF8583CtrlModeCounter)) != 0
	clockMode50Set := (uint8(c) & uint8(PCF8583CtrlModeClock50)) != 0
	return counterModeSet && !clockMode50Set
}

func (c PCF8583Control) isModeDiffer(mode PCF8583Control) bool {
	return uint8(c)&uint8(PCF8583CtrlModeTest) != uint8(mode)&uint8(PCF8583CtrlModeTest)
}

func pcf8583encodeBcd(val byte) byte {
	// decimal 12 => 0x12
	if val > 99 {
		val = 99
		if pcf8583Debug {
			log.Printf("PCF8583 BCD value (%d) exceeds limit of 99, now limited.", val)
		}
	}
	hi, lo := val/10, val%10
	return hi<<4 | lo
}

func pcf8583decodeBcd(bcd byte) byte {
	// 0x12 => decimal 12
	hi, lo := bcd>>4, bcd&0x0f
	if hi > 9 {
		hi = 9
		if pcf8583Debug {
			log.Printf("PCF8583 BCD value (%02x) exceeds limit 0x99 on most significant digit, now limited", bcd)
		}
	}
	if lo > 9 {
		lo = 9
		if pcf8583Debug {
			log.Printf("PCF8583 BCD value (%02x) exceeds limit 0x99 on least significant digit, now limited", bcd)
		}
	}
	return 10*hi + lo
}
