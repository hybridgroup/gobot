package i2c

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"gobot.io/x/gobot"
)

const (
	bmp388ChipID = 0x50

	bmp388RegisterChipID       = 0x00
	bmp388RegisterStatus       = 0x03
	bmp388RegisterConfig       = 0x1F
	bmp388RegisterPressureData = 0x04
	bmp388RegisterTempData     = 0x07
	bmp388RegisterCalib00      = 0x31
	bmp388RegisterCMD          = 0x7E
	// CMD 	: 0x00 nop (reserved. No command.)
	//		: 0x34 extmode_en_middle
	//		: 0xB0 fifo_flush (Clears all data in the FIFO, does not change FIFO_CONFIG registers)
	//		: 0xB6 softreset (Triggers a reset, all user configuration settings are overwritten with their default state)
	bmp388RegisterODR = 0x1D // Output Data Rates
	bmp388RegisterOSR = 0x1C // Oversampling Rates

	bmp388RegisterPWRCTRL = 0x1B
	bmp388PWRCTRLSleep    = 0
	bmp388PWRCTRLForced   = 1
	bmp388PWRCTRLNormal   = 3

	bmp388SeaLevelPressure = 1013.25

	// IIR filter coefficients
	bmp388IIRFIlterCoef0   = 0 // bypass-mode
	bmp388IIRFIlterCoef1   = 1
	bmp388IIRFIlterCoef3   = 2
	bmp388IIRFIlterCoef7   = 3
	bmp388IIRFIlterCoef15  = 4
	bmp388IIRFIlterCoef31  = 5
	bmp388IIRFIlterCoef63  = 6
	bmp388IIRFIlterCoef127 = 7
)

// BMP388Accuracy accuracy type
type BMP388Accuracy uint8

// BMP388Accuracy accuracy modes
const (
	BMP388AccuracyUltraLow  BMP388Accuracy = 0 // x1 sample
	BMP388AccuracyLow       BMP388Accuracy = 1 // x2 samples
	BMP388AccuracyStandard  BMP388Accuracy = 2 // x4 samples
	BMP388AccuracyHigh      BMP388Accuracy = 3 // x8 samples
	BMP388AccuracyUltraHigh BMP388Accuracy = 4 // x16 samples
	BMP388AccuracyHighest   BMP388Accuracy = 5 // x32 samples
)

type bmp388CalibrationCoefficients struct {
	t1  float32
	t2  float32
	t3  float32
	p1  float32
	p2  float32
	p3  float32
	p4  float32
	p5  float32
	p6  float32
	p7  float32
	p8  float32
	p9  float32
	p10 float32
	p11 float32
}

// BMP388Driver is a driver for the BMP388 temperature/pressure sensor
type BMP388Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config

	tpc *bmp388CalibrationCoefficients
}

// NewBMP388Driver creates a new driver with specified i2c interface.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewBMP388Driver(c Connector, options ...func(Config)) *BMP388Driver {
	b := &BMP388Driver{
		name:      gobot.DefaultName("BMP388"),
		connector: c,
		Config:    NewConfig(),
		tpc:       &bmp388CalibrationCoefficients{},
	}

	for _, option := range options {
		option(b)
	}

	// TODO: expose commands to API
	return b
}

// Name returns the name of the device.
func (d *BMP388Driver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *BMP388Driver) SetName(n string) {
	d.name = n
}

// Connection returns the connection of the device.
func (d *BMP388Driver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the BMP388 and loads the calibration coefficients.
func (d *BMP388Driver) Start() (err error) {
	var chipID uint8

	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(bmp180Address)

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return err
	}

	if chipID, err = d.connection.ReadByteData(bmp388RegisterChipID); err != nil {
		return err
	}

	if bmp388ChipID != chipID {
		return fmt.Errorf("Incorrect BMP388 chip ID '0%x' Expected 0x%x", chipID, bmp388ChipID)
	}

	if err := d.initialization(); err != nil {
		return err
	}

	return nil
}

// Halt halts the device.
func (d *BMP388Driver) Halt() (err error) {
	return nil
}

// Temperature returns the current temperature, in celsius degrees.
func (d *BMP388Driver) Temperature(accuracy BMP388Accuracy) (temp float32, err error) {
	var rawT int32

	// Enable Pressure and Temperature measurement, set FORCED operating mode
	var mode byte = (bmp388PWRCTRLForced << 4) | 3 // 1100|1|1 == mode|T|P
	if err = d.connection.WriteByteData(bmp388RegisterPWRCTRL, mode); err != nil {
		return 0, err
	}

	// Set Accuracy for temperature
	if err = d.connection.WriteByteData(bmp388RegisterOSR, uint8(accuracy<<3)); err != nil {
		return 0, err
	}

	if rawT, err = d.rawTemp(); err != nil {
		return 0.0, err
	}
	temp = d.calculateTemp(rawT)
	return
}

// Pressure returns the current barometric pressure, in Pa
func (d *BMP388Driver) Pressure(accuracy BMP388Accuracy) (press float32, err error) {
	var rawT, rawP int32

	// Enable Pressure and Temperature measurement, set FORCED operating mode
	var mode byte = (bmp388PWRCTRLForced << 4) | 3 // 1100|1|1 == mode|T|P
	if err = d.connection.WriteByteData(bmp388RegisterPWRCTRL, mode); err != nil {
		return 0, err
	}

	// Set Standard Accuracy for pressure
	if err = d.connection.WriteByteData(bmp388RegisterOSR, uint8(accuracy)); err != nil {
		return 0, err
	}

	if rawT, err = d.rawTemp(); err != nil {
		return 0.0, err
	}

	if rawP, err = d.rawPressure(); err != nil {
		return 0.0, err
	}
	tLin := d.calculateTemp(rawT)
	return d.calculatePress(rawP, float64(tLin)), nil
}

// Altitude returns the current altitude in meters based on the
// current barometric pressure and estimated pressure at sea level.
// https://www.weather.gov/media/epz/wxcalc/pressureAltitude.pdf
func (d *BMP388Driver) Altitude(accuracy BMP388Accuracy) (alt float32, err error) {
	atmP, _ := d.Pressure(accuracy)
	atmP /= 100.0
	alt = float32(44307.0 * (1.0 - math.Pow(float64(atmP/bmp388SeaLevelPressure), 0.190284)))

	return
}

// initialization reads the calibration coefficients.
func (d *BMP388Driver) initialization() (err error) {
	var (
		coefficients []byte
		t1           uint16
		t2           uint16
		t3           int8
		p1           int16
		p2           int16
		p3           int8
		p4           int8
		p5           uint16
		p6           uint16
		p7           int8
		p8           int8
		p9           int16
		p10          int8
		p11          int8
	)

	if coefficients, err = d.read(bmp388RegisterCalib00, 24); err != nil {
		return err
	}
	buf := bytes.NewBuffer(coefficients)

	binary.Read(buf, binary.LittleEndian, &t1)
	binary.Read(buf, binary.LittleEndian, &t2)
	binary.Read(buf, binary.LittleEndian, &t3)
	binary.Read(buf, binary.LittleEndian, &p1)
	binary.Read(buf, binary.LittleEndian, &p2)
	binary.Read(buf, binary.LittleEndian, &p3)
	binary.Read(buf, binary.LittleEndian, &p4)
	binary.Read(buf, binary.LittleEndian, &p5)
	binary.Read(buf, binary.LittleEndian, &p6)
	binary.Read(buf, binary.LittleEndian, &p7)
	binary.Read(buf, binary.LittleEndian, &p8)
	binary.Read(buf, binary.LittleEndian, &p9)
	binary.Read(buf, binary.LittleEndian, &p10)
	binary.Read(buf, binary.LittleEndian, &p11)

	d.tpc.t1 = float32(float64(t1) / math.Pow(2, -8))
	d.tpc.t2 = float32(float64(t2) / math.Pow(2, 30))
	d.tpc.t3 = float32(float64(t3) / math.Pow(2, 48))
	d.tpc.p1 = float32((float64(p1) - math.Pow(2, 14)) / math.Pow(2, 20))
	d.tpc.p2 = float32((float64(p2) - math.Pow(2, 14)) / math.Pow(2, 29))
	d.tpc.p3 = float32(float64(p3) / math.Pow(2, 32))
	d.tpc.p4 = float32(float64(p4) / math.Pow(2, 37))
	d.tpc.p5 = float32(float64(p5) / math.Pow(2, -3))
	d.tpc.p6 = float32(float64(p6) / math.Pow(2, 6))
	d.tpc.p7 = float32(float64(p7) / math.Pow(2, 8))
	d.tpc.p8 = float32(float64(p8) / math.Pow(2, 15))
	d.tpc.p9 = float32(float64(p9) / math.Pow(2, 48))
	d.tpc.p10 = float32(float64(p10) / math.Pow(2, 48))
	d.tpc.p11 = float32(float64(p11) / math.Pow(2, 65))

	// Perform a power on reset. All user configuration settings are overwritten
	// with their default state.
	if err = d.connection.WriteByteData(bmp388RegisterCMD, 0xB6); err != nil {
		return err
	}

	//  set IIR filter to off
	if err = d.connection.WriteByteData(bmp388RegisterConfig, bmp388IIRFIlterCoef0<<1); err != nil {
		return err
	}

	return nil
}

func (d *BMP388Driver) rawTemp() (temp int32, err error) {
	var data []byte
	var tp0, tp1, tp2 byte

	if data, err = d.read(bmp388RegisterTempData, 3); err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(data)

	binary.Read(buf, binary.LittleEndian, &tp0) // XLSB
	binary.Read(buf, binary.LittleEndian, &tp1) // LSB
	binary.Read(buf, binary.LittleEndian, &tp2) // MSB

	temp = ((int32(tp2) << 16) | (int32(tp1) << 8) | int32(tp0))
	return
}

func (d *BMP388Driver) rawPressure() (press int32, err error) {
	var data []byte
	var tp0, tp1, tp2 byte

	if data, err = d.read(bmp388RegisterPressureData, 3); err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(data)

	binary.Read(buf, binary.LittleEndian, &tp0) // XLSB
	binary.Read(buf, binary.LittleEndian, &tp1) // LSB
	binary.Read(buf, binary.LittleEndian, &tp2) // MSB

	press = ((int32(tp2) << 16) | (int32(tp1) << 8) | int32(tp0))

	return
}

func (d *BMP388Driver) calculateTemp(rawTemp int32) float32 {
	// datasheet, sec 9.2 Temperature compensation
	pd1 := float32(rawTemp) - d.tpc.t1
	pd2 := pd1 * d.tpc.t2

	temperatureComp := pd2 + (pd1*pd1)*d.tpc.t3

	return temperatureComp
}

func (d *BMP388Driver) calculatePress(rawPress int32, tLin float64) float32 {
	pd1 := float64(d.tpc.p6) * tLin
	pd2 := float64(d.tpc.p7) * math.Pow(tLin, 2)
	pd3 := float64(d.tpc.p8) * math.Pow(tLin, 3)
	po1 := float64(d.tpc.p5) + pd1 + pd2 + pd3

	pd1 = float64(d.tpc.p2) * tLin
	pd2 = float64(d.tpc.p3) * math.Pow(tLin, 2)
	pd3 = float64(d.tpc.p4) * math.Pow(tLin, 3)
	po2 := float64(rawPress) * (float64(d.tpc.p1) + pd1 + pd2 + pd3)

	pd1 = math.Pow(float64(rawPress), 2)
	pd2 = float64(d.tpc.p9) + float64(d.tpc.p10)*tLin
	pd3 = pd1 * pd2
	pd4 := pd3 + math.Pow(float64(rawPress), 3)*float64(d.tpc.p11)

	pressure := po1 + po2 + pd4

	return float32(pressure)
}

func (d *BMP388Driver) read(address byte, n int) ([]byte, error) {
	if _, err := d.connection.Write([]byte{address}); err != nil {
		return nil, err
	}
	buf := make([]byte, n)
	bytesRead, err := d.connection.Read(buf)
	if bytesRead != n || err != nil {
		return nil, err
	}
	return buf, nil
}
