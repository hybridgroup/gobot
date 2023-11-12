package i2c

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
)

const bmp388Debug = false

// the default address is applicable for SDO to VDD, for SDO to GND it will be 0x76
const bmp388DefaultAddress = 0x77

// BMP388Accuracy accuracy type
type (
	BMP388Accuracy  uint8
	BMP388IIRFilter uint8
)

const (
	bmp388ChipID = 0x50

	bmp388RegChipID       = 0x00
	bmp388RegStatus       = 0x03
	bmp388RegPressureData = 0x04 // XLSB, 0x05 LSByte, 0x06 MSByte
	bmp388RegTempData     = 0x07 // XLSB, 0x08 LSByte, 0x09 MSByte
	bmp388RegPWRCTRL      = 0x1B // enable/disable pressure and temperature measurement, mode
	bmp388RegOSR          = 0x1C // Oversampling Rates
	bmp388RegODR          = 0x1D // Output Data Rates
	bmp388RegConf         = 0x1F // config filter for IIR coefficients
	bmp388RegCalib00      = 0x31
	bmp388RegCMD          = 0x7E

	// bits 0, 1 of control register
	bmp388PWRCTRLPressEnableBit = 0x01
	bmp388PWRCTRLTempEnableBit  = 0x02

	// bits 4, 5 of control register (will be shifted on write)
	bmp388PWRCTRLSleep  = 0x00
	bmp388PWRCTRLForced = 0x01 // same as 0x02
	bmp388PWRCTRLNormal = 0x03

	// bits 1, 2 ,3 of config filter IIR filter coefficients (will be shifted on write)
	bmp388ConfFilterCoef0   BMP388IIRFilter = 0 // bypass-mode
	bmp388ConfFilterCoef1   BMP388IIRFilter = 1
	bmp388ConfFilterCoef3   BMP388IIRFilter = 2
	bmp388ConfFilterCoef7   BMP388IIRFilter = 3
	bmp388ConfFilterCoef15  BMP388IIRFilter = 4
	bmp388ConfFilterCoef31  BMP388IIRFilter = 5
	bmp388ConfFilterCoef63  BMP388IIRFilter = 6
	bmp388ConfFilterCoef127 BMP388IIRFilter = 7

	// oversampling rate, a single value is used (could be different for pressure and temperature)
	BMP388AccuracyUltraLow  BMP388Accuracy = 0 // x1 sample
	BMP388AccuracyLow       BMP388Accuracy = 1 // x2 samples
	BMP388AccuracyStandard  BMP388Accuracy = 2 // x4 samples
	BMP388AccuracyHigh      BMP388Accuracy = 3 // x8 samples
	BMP388AccuracyUltraHigh BMP388Accuracy = 4 // x16 samples
	BMP388AccuracyHighest   BMP388Accuracy = 5 // x32 samples

	bmp388CMDReserved        = 0x00 // reserved, no command
	bmp388CMDExtModeEnMiddle = 0x34
	bmp388CMDFifoFlush       = 0xB0 // clears all data in the FIFO, does not change FIFO_CONFIG registers
	bmp388CMDSoftReset       = 0xB6 // triggers a reset, all user configuration settings are overwritten with defaults

	bmp388SeaLevelPressure = 1013.25
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
	*Driver
	calCoeffs   *bmp388CalibrationCoefficients
	ctrlPwrMode uint8
	confFilter  BMP388IIRFilter
}

// NewBMP388Driver creates a new driver with specified i2c interface.
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewBMP388Driver(c Connector, options ...func(Config)) *BMP388Driver {
	d := &BMP388Driver{
		Driver:      NewDriver(c, "BMP388", bmp388DefaultAddress),
		calCoeffs:   &bmp388CalibrationCoefficients{},
		ctrlPwrMode: bmp388PWRCTRLForced,
		confFilter:  bmp388ConfFilterCoef0,
	}
	d.afterStart = d.initialization

	for _, option := range options {
		option(d)
	}

	// TODO: expose commands to API
	return d
}

// WithBMP388IIRFilter option sets count of IIR filter coefficients.
// Valid settings are of type "BMP388IIRFilter"
func WithBMP388IIRFilter(val BMP388IIRFilter) func(Config) {
	return func(c Config) {
		if d, ok := c.(*BMP388Driver); ok {
			d.confFilter = val
		} else if bmp388Debug {
			log.Printf("Trying to set IIR filter for non-BMP388Driver %v", c)
		}
	}
}

// Temperature returns the current temperature, in celsius degrees.
func (d *BMP388Driver) Temperature(accuracy BMP388Accuracy) (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	mode := d.ctrlPwrMode<<4 | bmp388PWRCTRLPressEnableBit | bmp388PWRCTRLTempEnableBit
	if err := d.connection.WriteByteData(bmp388RegPWRCTRL, mode); err != nil {
		return 0, err
	}

	// Set Accuracy for temperature
	if err := d.connection.WriteByteData(bmp388RegOSR, uint8(accuracy<<3)); err != nil {
		return 0, err
	}

	rawT, err := d.rawTemp()
	if err != nil {
		return 0.0, err
	}

	temp := d.calculateTemp(rawT)

	return temp, nil
}

// Pressure returns the current barometric pressure, in Pa
func (d *BMP388Driver) Pressure(accuracy BMP388Accuracy) (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	mode := d.ctrlPwrMode<<4 | bmp388PWRCTRLPressEnableBit | bmp388PWRCTRLTempEnableBit
	if err := d.connection.WriteByteData(bmp388RegPWRCTRL, mode); err != nil {
		return 0, err
	}

	// Set Standard Accuracy for pressure
	if err := d.connection.WriteByteData(bmp388RegOSR, uint8(accuracy)); err != nil {
		return 0, err
	}

	rawT, err := d.rawTemp()
	if err != nil {
		return 0.0, err
	}

	rawP, err := d.rawPressure()
	if err != nil {
		return 0.0, err
	}
	tLin := d.calculateTemp(rawT)

	return d.calculatePress(rawP, float64(tLin)), nil
}

// Altitude returns the current altitude in meters based on the
// current barometric pressure and estimated pressure at sea level.
// https://www.weather.gov/media/epz/wxcalc/pressureAltitude.pdf
func (d *BMP388Driver) Altitude(accuracy BMP388Accuracy) (float32, error) {
	atmP, err := d.Pressure(accuracy)
	if err != nil {
		return 0, err
	}
	atmP /= 100.0
	alt := float32(44307.0 * (1.0 - math.Pow(float64(atmP/bmp388SeaLevelPressure), 0.190284)))

	return alt, nil
}

// initialization reads the calibration coefficients.
func (d *BMP388Driver) initialization() error {
	chipID, err := d.connection.ReadByteData(bmp388RegChipID)
	if err != nil {
		return err
	}

	if bmp388ChipID != chipID {
		return fmt.Errorf("Incorrect BMP388 chip ID '0%x' Expected 0x%x", chipID, bmp388ChipID)
	}

	var (
		t1  uint16
		t2  uint16
		t3  int8
		p1  int16
		p2  int16
		p3  int8
		p4  int8
		p5  uint16
		p6  uint16
		p7  int8
		p8  int8
		p9  int16
		p10 int8
		p11 int8
	)

	coefficients := make([]byte, 24)
	if err = d.connection.ReadBlockData(bmp388RegCalib00, coefficients); err != nil {
		return err
	}
	buf := bytes.NewBuffer(coefficients)

	if err := binary.Read(buf, binary.LittleEndian, &t1); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &t2); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &t3); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p1); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p2); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p3); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p4); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p5); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p6); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p7); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p8); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p9); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p10); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p11); err != nil {
		return err
	}

	d.calCoeffs.t1 = float32(float64(t1) / math.Pow(2, -8))
	d.calCoeffs.t2 = float32(float64(t2) / math.Pow(2, 30))
	d.calCoeffs.t3 = float32(float64(t3) / math.Pow(2, 48))
	d.calCoeffs.p1 = float32((float64(p1) - math.Pow(2, 14)) / math.Pow(2, 20))
	d.calCoeffs.p2 = float32((float64(p2) - math.Pow(2, 14)) / math.Pow(2, 29))
	d.calCoeffs.p3 = float32(float64(p3) / math.Pow(2, 32))
	d.calCoeffs.p4 = float32(float64(p4) / math.Pow(2, 37))
	d.calCoeffs.p5 = float32(float64(p5) / math.Pow(2, -3))
	d.calCoeffs.p6 = float32(float64(p6) / math.Pow(2, 6))
	d.calCoeffs.p7 = float32(float64(p7) / math.Pow(2, 8))
	d.calCoeffs.p8 = float32(float64(p8) / math.Pow(2, 15))
	d.calCoeffs.p9 = float32(float64(p9) / math.Pow(2, 48))
	d.calCoeffs.p10 = float32(float64(p10) / math.Pow(2, 48))
	d.calCoeffs.p11 = float32(float64(p11) / math.Pow(2, 65))

	if err := d.connection.WriteByteData(bmp388RegCMD, bmp388CMDSoftReset); err != nil {
		return err
	}

	return d.connection.WriteByteData(bmp388RegConf, uint8(d.confFilter)<<1)
}

func (d *BMP388Driver) rawTemp() (int32, error) {
	var tp0, tp1, tp2 byte

	data := make([]byte, 3)
	if err := d.connection.ReadBlockData(bmp388RegTempData, data); err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(data)
	// XLSB
	if err := binary.Read(buf, binary.LittleEndian, &tp0); err != nil {
		return 0, err
	}
	// LSB
	if err := binary.Read(buf, binary.LittleEndian, &tp1); err != nil {
		return 0, err
	}
	// MSB
	if err := binary.Read(buf, binary.LittleEndian, &tp2); err != nil {
		return 0, err
	}

	return ((int32(tp2) << 16) | (int32(tp1) << 8) | int32(tp0)), nil
}

func (d *BMP388Driver) rawPressure() (int32, error) {
	var tp0, tp1, tp2 byte

	data := make([]byte, 3)
	if err := d.connection.ReadBlockData(bmp388RegPressureData, data); err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(data)
	// XLSB
	if err := binary.Read(buf, binary.LittleEndian, &tp0); err != nil {
		return 0, err
	}
	// LSB
	if err := binary.Read(buf, binary.LittleEndian, &tp1); err != nil {
		return 0, err
	}
	// MSB
	if err := binary.Read(buf, binary.LittleEndian, &tp2); err != nil {
		return 0, err
	}

	return ((int32(tp2) << 16) | (int32(tp1) << 8) | int32(tp0)), nil
}

func (d *BMP388Driver) calculateTemp(rawTemp int32) float32 {
	// datasheet, sec 9.2 Temperature compensation
	pd1 := float32(rawTemp) - d.calCoeffs.t1
	pd2 := pd1 * d.calCoeffs.t2

	temperatureComp := pd2 + (pd1*pd1)*d.calCoeffs.t3

	return temperatureComp
}

func (d *BMP388Driver) calculatePress(rawPress int32, tLin float64) float32 {
	pd1 := float64(d.calCoeffs.p6) * tLin
	pd2 := float64(d.calCoeffs.p7) * math.Pow(tLin, 2)
	pd3 := float64(d.calCoeffs.p8) * math.Pow(tLin, 3)
	po1 := float64(d.calCoeffs.p5) + pd1 + pd2 + pd3

	pd1 = float64(d.calCoeffs.p2) * tLin
	pd2 = float64(d.calCoeffs.p3) * math.Pow(tLin, 2)
	pd3 = float64(d.calCoeffs.p4) * math.Pow(tLin, 3)
	po2 := float64(rawPress) * (float64(d.calCoeffs.p1) + pd1 + pd2 + pd3)

	pd1 = math.Pow(float64(rawPress), 2)
	pd2 = float64(d.calCoeffs.p9) + float64(d.calCoeffs.p10)*tLin
	pd3 = pd1 * pd2
	pd4 := pd3 + math.Pow(float64(rawPress), 3)*float64(d.calCoeffs.p11)

	pressure := po1 + po2 + pd4

	return float32(pressure)
}
