package i2c

import (
	"fmt"
	"log"
	"math"
	"sort"
)

const (
	hmc5883lDebug          = false
	hmc5883lDefaultAddress = 0x1E // default I2C Address
)

const (
	hmc5883lRegA      = 0x00 // Address of Configuration register A
	hmc5883lRegB      = 0x01 // Address of Configuration register B
	hmc5883lRegMode   = 0x02 // Address of node register
	hmc5883lAxisX     = 0x03 // Address of X-axis MSB data register
	hmc5883lAxisZ     = 0x05 // Address of Z-axis MSB data register
	hmc5883lAxisY     = 0x07 // Address of Y-axis MSB data register
	hmc5883lRegStatus = 0x09 // Address of status register
	hmc5883lRegIdA    = 0x0A // Address of identification register A
	hmc5883lRegIdB    = 0x0B // Address of identification register B
	hmc5883lRegIdC    = 0x0C // Address of identification register C

	hmc5883lRegA_SamplesAvg1      = 0x00 // no samples averaged
	hmc5883lRegA_SamplesAvg2      = 0x01 // 2 samples averaged
	hmc5883lRegA_SamplesAvg4      = 0x02 // 4 samples averaged
	hmc5883lRegA_SamplesAvg8      = 0x03 // 8 samples averaged
	hmc5883lRegA_OutputRate750    = 0x00 // data output rate 0.75 Hz
	hmc5883lRegA_OutputRate1500   = 0x01 // data output rate 1.5 Hz
	hmc5883lRegA_OutputRate3000   = 0x02 // data output rate 3.0 Hz
	hmc5883lRegA_OutputRate7500   = 0x03 // data output rate 7.5 Hz
	hmc5883lRegA_OutputRate15000  = 0x04 // data output rate 15.0 Hz
	hmc5883lRegA_OutputRate30000  = 0x05 // data output rate 30.0 Hz
	hmc5883lRegA_OutputRate75000  = 0x06 // data output rate 75.0 Hz
	hmc5883lRegA_MeasNormal       = 0x00 // normal measurement configuration
	hmc5883lRegA_MeasPositiveBias = 0x01 // positive bias for X, Y, Z
	hmc5883lRegA_MeasNegativeBias = 0x02 // negative bias for X, Y, Z

	hmc5883lRegB_Gain1370 = 0x00 // gain is 1370 Gauss
	hmc5883lRegB_Gain1090 = 0x01 // gain is 1090 Gauss
	hmc5883lRegB_Gain820  = 0x02 // gain is 820 Gauss
	hmc5883lRegB_Gain660  = 0x03 // gain is 660 Gauss
	hmc5883lRegB_Gain440  = 0x04 // gain is 440 Gauss
	hmc5883lRegB_Gain390  = 0x05 // gain is 390 Gauss
	hmc5883lRegB_Gain330  = 0x06 // gain is 330 Gauss
	hmc5883lRegB_Gain230  = 0x07 // gain is 230 Gauss

	hmc5883lRegM_Continuous = 0x00 // continuous measurement mode
	hmc5883lRegM_Single     = 0x01 // return to idle after a single measurement
	hmc5883lRegM_Idle       = 0x10 // idle mode
)

// HMC5883LDriver is a Gobot Driver for a HMC5883 I2C 3 axis digital compass.
//
// This driver was tested with Tinkerboard & Digispark adaptor and a HMC5883L breakout board GY-273,
// available from various distributors.
//
// datasheet:
// http://www.adafruit.com/datasheets/HMC5883L_3-Axis_Digital_Compass_IC.pdf
//
// reference implementations:
// * https://github.com/gvalkov/micropython-esp8266-hmc5883l
// * https://github.com/adafruit/Adafruit_HMC5883_Unified
type HMC5883LDriver struct {
	*Driver
	samplesAvg      uint8
	outputRate      uint32 // in mHz
	applyBias       int8
	measurementMode int
	gain            float64 // in 1/Gauss
}

var hmc5883lSamplesAvgBits = map[uint8]int{
	1: hmc5883lRegA_SamplesAvg1,
	2: hmc5883lRegA_SamplesAvg2,
	4: hmc5883lRegA_SamplesAvg4,
	8: hmc5883lRegA_SamplesAvg8,
}

var hmc5883lOutputRateBits = map[uint32]int{
	750:   hmc5883lRegA_OutputRate750,
	1500:  hmc5883lRegA_OutputRate1500,
	3000:  hmc5883lRegA_OutputRate3000,
	7500:  hmc5883lRegA_OutputRate7500,
	15000: hmc5883lRegA_OutputRate15000,
	30000: hmc5883lRegA_OutputRate30000,
	75000: hmc5883lRegA_OutputRate75000,
}

var hmc5883lMeasurementFlowBits = map[int8]int{
	0:  hmc5883lRegA_MeasNormal,
	1:  hmc5883lRegA_MeasPositiveBias,
	-1: hmc5883lRegA_MeasNegativeBias,
}

var hmc5883lGainBits = map[float64]int{
	1370.0: hmc5883lRegB_Gain1370,
	1090.0: hmc5883lRegB_Gain1090,
	820.0:  hmc5883lRegB_Gain820,
	660.0:  hmc5883lRegB_Gain660,
	440.0:  hmc5883lRegB_Gain440,
	390.0:  hmc5883lRegB_Gain390,
	330.0:  hmc5883lRegB_Gain330,
	230.0:  hmc5883lRegB_Gain230,
}

// NewHMC5883LDriver creates a new driver with specified i2c interface
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//			i2c.WithBus(int):	bus to use with this driver
//			i2c.WithAddress(int):	address to use with this driver
//	   i2c.WithHMC5883LSamplesAveraged(int)
//	   i2c.WithHMC5883LDataOutputRate(int)
//	   i2c.WithHMC5883LMeasurementFlow(int)
//	   i2c.WithHMC5883LGain(int)
func NewHMC5883LDriver(c Connector, options ...func(Config)) *HMC5883LDriver {
	h := &HMC5883LDriver{
		Driver:          NewDriver(c, "HMC5883L", hmc5883lDefaultAddress),
		samplesAvg:      8,
		outputRate:      15000,
		applyBias:       0,
		measurementMode: hmc5883lRegM_Continuous,
		gain:            390,
	}
	h.afterStart = h.initialize

	for _, option := range options {
		option(h)
	}

	return h
}

// WithHMC5883LSamplesAveraged option sets the number of samples averaged per measurement.
// Valid settings are 1, 2, 4, 8.
func WithHMC5883LSamplesAveraged(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*HMC5883LDriver)
		if ok {
			if err := hmc5883lValidateSamplesAveraged(val); err != nil {
				panic(err)
			}
			d.samplesAvg = uint8(val)
		} else if hmc5883lDebug {
			log.Printf("Trying to set samples averaged for non-HMC5883LDriver %v", c)
		}
	}
}

// WithHMC5883LDataOutputRate option sets the data output rate in mHz.
// Valid settings are 750, 1500, 3000, 7500, 15000, 30000, 75000.
func WithHMC5883LDataOutputRate(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*HMC5883LDriver)
		if ok {
			if err := hmc5883lValidateOutputRate(val); err != nil {
				panic(err)
			}
			d.outputRate = uint32(val)
		} else if hmc5883lDebug {
			log.Printf("Trying to set data output rate for non-HMC5883LDriver %v", c)
		}
	}
}

// WithHMC5883LApplyBias option sets to apply a measurement bias.
// Valid settings are -1 (negative bias), 0 (normal), 1 (positive bias).
func WithHMC5883LApplyBias(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*HMC5883LDriver)
		if ok {
			if err := hmc5883lValidateApplyBias(val); err != nil {
				panic(err)
			}
			d.applyBias = int8(val)
		} else if hmc5883lDebug {
			log.Printf("Trying to set measurement flow for non-HMC5883LDriver %v", c)
		}
	}
}

// WithHMC5883LGain option sets the gain.
// Valid settings are 1370, 1090, 820, 660, 440, 390, 330 230 in 1/Gauss.
func WithHMC5883LGain(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*HMC5883LDriver)
		if ok {
			if err := hmc5883lValidateGain(val); err != nil {
				panic(err)
			}
			d.gain = float64(val)
		} else if hmc5883lDebug {
			log.Printf("Trying to set gain for non-HMC5883LDriver %v", c)
		}
	}
}

// Read reads the values X, Y, Z in Gauss
func (h *HMC5883LDriver) Read() (x float64, y float64, z float64, err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	xr, yr, zr, err := h.readRawData()
	if err != nil {
		return
	}
	return float64(xr) / h.gain, float64(yr) / h.gain, float64(zr) / h.gain, nil
}

// Heading returns the current heading in radians
func (h *HMC5883LDriver) Heading() (heading float64, err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var x, y int16
	x, y, _, err = h.readRawData()
	if err != nil {
		return
	}
	heading = math.Atan2(float64(y), float64(x))
	if heading > 2*math.Pi {
		heading -= 2 * math.Pi
	}
	if heading < 0 {
		heading += 2 * math.Pi
	}
	return
}

// readRawData reads the raw values from the X, Y, and Z registers
func (h *HMC5883LDriver) readRawData() (x int16, y int16, z int16, err error) {
	// read the data, starting from the initial register
	data := make([]byte, 6)
	if err = h.connection.ReadBlockData(hmc5883lAxisX, data); err != nil {
		return
	}

	unsignedX := (uint16(data[0]) << 8) | uint16(data[1])
	unsignedZ := (uint16(data[2]) << 8) | uint16(data[3])
	unsignedY := (uint16(data[4]) << 8) | uint16(data[5])

	return twosComplement16Bit(unsignedX), twosComplement16Bit(unsignedY), twosComplement16Bit(unsignedZ), nil
}

func (h *HMC5883LDriver) initialize() (err error) {
	regA := hmc5883lMeasurementFlowBits[h.applyBias]
	regA |= hmc5883lOutputRateBits[h.outputRate] << 2
	regA |= hmc5883lSamplesAvgBits[h.samplesAvg] << 5
	if err := h.connection.WriteByteData(hmc5883lRegA, uint8(regA)); err != nil {
		return err
	}
	regB := hmc5883lGainBits[h.gain] << 5
	if err := h.connection.WriteByteData(hmc5883lRegB, uint8(regB)); err != nil {
		return err
	}
	if err := h.connection.WriteByteData(hmc5883lRegMode, uint8(h.measurementMode)); err != nil {
		return err
	}
	return
}

func hmc5883lValidateSamplesAveraged(samplesAvg int) (err error) {
	if _, ok := hmc5883lSamplesAvgBits[uint8(samplesAvg)]; ok {
		return
	}

	keys := []int{}
	for k := range hmc5883lSamplesAvgBits {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	err = fmt.Errorf("Samples averaged must be one of: %d", keys)
	return
}

func hmc5883lValidateOutputRate(outputRate int) (err error) {
	if _, ok := hmc5883lOutputRateBits[uint32(outputRate)]; ok {
		return
	}

	keys := []int{}
	for k := range hmc5883lOutputRateBits {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	err = fmt.Errorf("Data output rate must be one of: %d", keys)
	return
}

func hmc5883lValidateApplyBias(applyBias int) (err error) {
	if _, ok := hmc5883lMeasurementFlowBits[int8(applyBias)]; ok {
		return
	}

	keys := []int{}
	for k := range hmc5883lMeasurementFlowBits {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	err = fmt.Errorf("Apply measurement bias must be one of: %d", keys)
	return
}

func hmc5883lValidateGain(gain int) (err error) {
	if _, ok := hmc5883lGainBits[float64(gain)]; ok {
		return
	}

	keys := []int{}
	for k := range hmc5883lGainBits {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	err = fmt.Errorf("Gain must be one of: %d", keys)
	return
}
