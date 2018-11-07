package i2c

import (
	"fmt"
	"math"
	"time"

	"gobot.io/x/gobot"
)

// CCS811DriveMode type
type CCS811DriveMode uint8

// Operating modes which dictate how often measurements are being made. If 0x00 is used as an operating mode, measurements will be disabled
const (
	CCS811DriveModeIdle  CCS811DriveMode = 0x00
	CCS811DriveMode1Sec                  = 0x01
	CCS811DriveMode10Sec                 = 0x02
	CCS811DriveMode60Sec                 = 0x03
	CCS811DriveMode250MS                 = 0x04
)

const (

	//DefaultAddress is the default I2C address for the ccs811
	ccs811DefaultAddress = 0x5A

	//Registers, all definitions have been taken from the datasheet
	//Single byte read only register which indicates if a device is active, if new data is available or if an error occurred.
	ccs811RegStatus = 0x00
	//This is Single byte register, which is used to enable sensor drive mode and interrupts.
	ccs811RegMeasMode = 0x01
	//This multi-byte read only register contains the calculated eCO2 (ppm) and eTVOC (ppb) values followed by the STATUS register, ERROR_ID register and the RAW_DATA register.
	ccs811RegAlgResultData = 0x02
	//Two byte read only register which contains the latest readings from the sensor.
	ccs811RegRawData = 0x03
	//A multi-byte register that can be written with the current Humidity and Temperature values if known.
	ccs811RegEnvData = 0x05
	//Register that holds the NTC value used for temperature calcualtions
	ccs811RegNtc = 0x06
	//Asserting the SW_RESET will restart the CCS811 in Boot mode to enable new application firmware to be downloaded.
	ccs811RegSwReset = 0xFF
	//Single byte read only register which holds the HW ID which is 0x81 for this family of CCS81x devices.
	ccs811RegHwID = 0x20
	//Single byte read only register that contains the hardware version. The value is 0x1X
	ccs811RegHwVersion = 0x21
	//Two byte read only register which contain the version of the firmware bootloader stored in the CCS811 in the format Major.Minor.Trivial
	ccs811RegFwBootVersion = 0x23
	//Two byte read only register which contain the version of the firmware application stored in the CCS811 in the format Major.Minor.Trivial
	ccs811RegFwAppVersion = 0x24
	//To change the mode of the CCS811 from Boot mode to running the application, a single byte write of 0xF4 is required.
	ccs811RegAppStart = 0xF4

	// Constants
	// The hardware ID code
	ccs811HwIDCode = 0x81
)

var (
	// The sequence of bytes needed to do a software reset
	ccs811SwResetSequence = []byte{0x11, 0xE5, 0x72, 0x8A}
)

// CCS811Status represents the current status of the device defined by the ccs811RegStatus.
// The following definitions were taken from https://ams.com/documents/20143/36005/CCS811_DS000459_6-00.pdf/c7091525-c7e5-37ac-eedb-b6c6828b0dcf#page=15
type CCS811Status struct {
	//There is some sort of error on the i2c bus or there is an error with the internal sensor
	HasError byte
	//A new data sample is ready in ccs811RegAlgResultData
	DataReady byte
	//Valid application firmware loaded
	AppValid byte
	//Firmware is in application mode. CCS811 is ready to take sensor measurements
	FwMode byte
}

//NewCCS811Status returns a new instance of the package ccs811 status definiton
func NewCCS811Status(data uint8) *CCS811Status {
	return &CCS811Status{
		HasError:  data & 0x01,
		DataReady: (data >> 3) & 0x01,
		AppValid:  (data >> 4) & 0x01,
		FwMode:    (data >> 7) & 0x01,
	}
}

//CCS811MeasMode represents the current measurement configuration.
//The following definitions were taken from the bit fields of the ccs811RegMeasMode defined in
//https://ams.com/documents/20143/36005/CCS811_DS000459_6-00.pdf/c7091525-c7e5-37ac-eedb-b6c6828b0dcf#page=16
type CCS811MeasMode struct {
	//If intThresh is 1 a data measurement will only be taken when the sensor value mets the threshold constraint.
	//The threshold value is set in the threshold register (0x10)
	intThresh uint8
	//If intDataRdy is 1, the nINT signal (pin 3 of the device) will be driven low when new data is avaliable.
	intDataRdy uint8
	//driveMode represents the sampling rate of the sensor. If the value is 0, the measurement process is idle.
	driveMode CCS811DriveMode
}

//NewCCS811MeasMode returns a new instance of the package ccs811 measurement mode configuration. This represents the desired intial
//state of the measurement mode register.
func NewCCS811MeasMode() *CCS811MeasMode {
	return &CCS811MeasMode{
		// Disable this by default as this library does not contain the functionality to use the internal interrupt feature.
		intThresh:  0x00,
		intDataRdy: 0x00,
		driveMode:  CCS811DriveMode1Sec,
	}
}

// GetMeasMode returns the measurement mode
func (mm *CCS811MeasMode) GetMeasMode() byte {
	return (mm.intThresh << 2) | (mm.intDataRdy << 3) | uint8((mm.driveMode << 4))
}

//CCS811Driver is the Gobot driver for the CCS811 (air quality sensor) Adafruit breakout board
type CCS811Driver struct {
	name               string
	connector          Connector
	connection         Connection
	measMode           *CCS811MeasMode
	ntcResistanceValue uint32
	Config
}

//NewCCS811Driver creates a new driver for the CCS811 (air quality sensor)
func NewCCS811Driver(a Connector, options ...func(Config)) *CCS811Driver {
	l := &CCS811Driver{
		name:      gobot.DefaultName("CCS811"),
		connector: a,
		measMode:  NewCCS811MeasMode(),
		//Recommended resistance value is 100,000
		ntcResistanceValue: 100000,
		Config:             NewConfig(),
	}

	for _, option := range options {
		option(l)
	}

	return l
}

//WithCCS811MeasMode sets the sampling rate of the device
func WithCCS811MeasMode(mode CCS811DriveMode) func(Config) {
	return func(c Config) {
		d, _ := c.(*CCS811Driver)
		d.measMode.driveMode = mode
	}
}

//WithCCS811NTCResistance sets reistor value used in the temperature calculations.
//This resistor must be placed between pin 4 and pin 8 of the chip
func WithCCS811NTCResistance(val uint32) func(Config) {
	return func(c Config) {
		d, _ := c.(*CCS811Driver)
		d.ntcResistanceValue = val
	}
}

//Start initializes the sensor
func (d *CCS811Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(ccs811DefaultAddress)

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return err
	}

	return d.initialize()
}

//Name returns the Name for the Driver
func (d *CCS811Driver) Name() string { return d.name }

//SetName sets the Name for the Driver
func (d *CCS811Driver) SetName(n string) { d.name = n }

//Connection returns the connection for the Driver
func (d *CCS811Driver) Connection() gobot.Connection { return d.connector.(gobot.Connection) }

//Halt returns true if devices is halted successfully
func (d *CCS811Driver) Halt() (err error) { return }

//GetHardwareVersion returns the hardware version of the device in the form of 0x1X
func (d *CCS811Driver) GetHardwareVersion() (uint8, error) {
	v, err := d.connection.ReadByteData(ccs811RegHwVersion)
	if err != nil {
		return 0, err
	}

	return v, nil
}

//GetFirmwareBootVersion returns the bootloader version
func (d *CCS811Driver) GetFirmwareBootVersion() (uint16, error) {
	v, err := d.connection.ReadWordData(ccs811RegFwBootVersion)
	if err != nil {
		return 0, err
	}

	return v, nil
}

//GetFirmwareAppVersion returns the app code version
func (d *CCS811Driver) GetFirmwareAppVersion() (uint16, error) {
	v, err := d.connection.ReadWordData(ccs811RegFwAppVersion)
	if err != nil {
		return 0, err
	}

	return v, nil
}

//GetStatus returns the current status of the device
func (d *CCS811Driver) GetStatus() (*CCS811Status, error) {
	s, err := d.connection.ReadByteData(ccs811RegStatus)
	if err != nil {
		return nil, err
	}

	cs := NewCCS811Status(s)
	return cs, nil
}

//GetTemperature returns the device temperature in celcius.
//If you do not have an NTC resistor installed, this function should not be called
func (d *CCS811Driver) GetTemperature() (float32, error) {

	buf, err := d.read(ccs811RegNtc, 4)
	if err != nil {
		return 0, err
	}

	vref := ((uint16(buf[0]) << 8) | uint16(buf[1]))
	vrntc := ((uint16(buf[2]) << 8) | uint16(buf[3]))
	rntc := (float32(vrntc) * float32(d.ntcResistanceValue) / float32(vref))

	ntcTemp := float32(math.Log(float64(rntc / 10000.0)))
	ntcTemp /= 3380.0
	ntcTemp += 1.0 / (25 + 273.15)
	ntcTemp = 1.0 / ntcTemp
	ntcTemp -= 273.15

	return ntcTemp, nil
}

//GetGasData returns the data for the gas sensor.
//eco2 is returned in ppm and tvoc is returned in ppb
func (d *CCS811Driver) GetGasData() (uint16, uint16, error) {

	data, err := d.read(ccs811RegAlgResultData, 4)
	if err != nil {
		return 0, 0, err
	}

	// Bit masks defined by https://ams.com/documents/20143/36005/CCS811_AN000369_2-00.pdf/25d0db9a-92b9-fa7f-362c-a7a4d1e292be#page=14
	eco2 := (uint16(data[0]) << 8) | uint16(data[1])
	tvoC := (uint16(data[2]) << 8) | uint16(data[3])

	return eco2, tvoC, nil
}

//HasData returns true if the device has not errored and temperature/gas data is avaliable
func (d *CCS811Driver) HasData() (bool, error) {
	s, err := d.GetStatus()
	if err != nil {
		return false, err
	}

	if !(s.DataReady == 0x01) || (s.HasError == 0x01) {
		return false, nil
	}

	return true, nil
}

//EnableExternalInterrupt enables the external output hardware interrupt pin 3.
func (d *CCS811Driver) EnableExternalInterrupt() error {
	d.measMode.intDataRdy = 1
	return d.connection.WriteByteData(ccs811RegMeasMode, d.measMode.GetMeasMode())
}

//DisableExternalInterrupt disables the external output hardware interrupt pin 3.
func (d *CCS811Driver) DisableExternalInterrupt() error {
	d.measMode.intDataRdy = 0
	return d.connection.WriteByteData(ccs811RegMeasMode, d.measMode.GetMeasMode())
}

//updateMeasMode writes the current value of measMode to the measurement mode register.
func (d *CCS811Driver) updateMeasMode() error {
	return d.connection.WriteByteData(ccs811RegMeasMode, d.measMode.GetMeasMode())
}

//ResetDevice does a software reset of the device. After this operation is done,
//the user must start the app code before the sensor can take any measurements
func (d *CCS811Driver) resetDevice() error {
	return d.connection.WriteBlockData(ccs811RegSwReset, ccs811SwResetSequence)
}

//startApp starts the app code in the device. This operation has to be done after a
//software reset to start taking sensor measurements.
func (d *CCS811Driver) startApp() error {
	//Write without data is needed to start the app code
	_, err := d.connection.Write([]byte{ccs811RegAppStart})
	return err
}

func (d *CCS811Driver) initialize() error {
	deviceID, err := d.connection.ReadByteData(ccs811RegHwID)
	if err != nil {
		return fmt.Errorf("Failed to get the device id from ccs811RegHwID with error: %s", err.Error())
	}

	// Verify that the connected device is the CCS811 sensor
	if deviceID != ccs811HwIDCode {
		return fmt.Errorf("The fetched device id %d is not the known id %d with error", deviceID, ccs811HwIDCode)
	}

	if err := d.resetDevice(); err != nil {
		return fmt.Errorf("Was not able to reset the device with error: %s", err.Error())
	}

	// Required sleep to allow device to switch states
	time.Sleep(100 * time.Millisecond)

	if err := d.startApp(); err != nil {
		return fmt.Errorf("Failed to start app code with error: %s", err.Error())
	}

	if err := d.updateMeasMode(); err != nil {
		return fmt.Errorf("Failed to update the measMode register with error: %s", err.Error())
	}

	return nil
}

// An implementation of the ReadBlockData i2c operation. This code was copied from the BMP280Driver code
func (d *CCS811Driver) read(reg byte, n int) ([]byte, error) {
	if _, err := d.connection.Write([]byte{reg}); err != nil {
		return nil, err
	}
	buf := make([]byte, n)
	bytesRead, err := d.connection.Read(buf)
	if bytesRead != n || err != nil {
		return nil, err
	}
	return buf, nil
}
