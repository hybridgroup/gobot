package i2c

// Device documentation at
// https://ae-bst.resource.bosch.com/media/_tech/media/datasheets/BST-BME280_DS001-11.pdf

import (
	"github.com/hybridgroup/gobot"

	"bytes"
	"encoding/binary"
	"time"
	"errors"
)

const bme280Address = 0x76

var _ gobot.Driver = (*BME280Driver)(nil)

const BME280_REGISTER_PRESSURE_MSB = 0xF7
const BME280_REGISTER_PRESSURE_LSB = 0xF8
const BME280_REGISTER_PRESSURE_XLSB = 0xF9

const BME280_REGISTER_TEMP_MSB = 0xFA
const BME280_REGISTER_TEMP_LSB = 0xFB
const BME280_REGISTER_TEMP_XLSB = 0xFC

const BME280_REGISTER_HUM_MSB = 0xFD
const BME280_REGISTER_HUM_LSB = 0xFE

const BME280_PRESSURE_TEMPERATURE_CALIB_DATA_LENGTH = 26
const BME280_HUMIDITY_CALIB_DATA_LENGTH = 7

const BME280_DATA_FRAME_LENGTH = 12

const BME280_REGISTER_DIG_T1_LSB = 0x88
const BME280_REGISTER_DIG_T1_MSB = 0x89

const BME280_REGISTER_DIG_T2_LSB = 0x8A
const BME280_REGISTER_DIG_T2_MSB = 0x8B

const BME280_REGISTER_DIG_T3_LSB = 0x8C
const BME280_REGISTER_DIG_T3_MSB = 0x8D

const BME280_REGISTER_DIG_P1_LSB = 0x8E
const BME280_REGISTER_DIG_P1_MSB = 0x8F

const BME280_REGISTER_DIG_P2_LSB = 0x90
const BME280_REGISTER_DIG_P2_MSB = 0x91

const BME280_REGISTER_DIG_P3_LSB = 0x92
const BME280_REGISTER_DIG_P3_MSB = 0x93

const BME280_REGISTER_DIG_P4_LSB = 0x94
const BME280_REGISTER_DIG_P4_MSB = 0x95

const BME280_REGISTER_DIG_P5_LSB = 0x96
const BME280_REGISTER_DIG_P5_MSB = 0x97

const BME280_REGISTER_DIG_P6_LSB = 0x98
const BME280_REGISTER_DIG_P6_MSB = 0x99

const BME280_REGISTER_DIG_P7_LSB = 0x9A
const BME280_REGISTER_DIG_P7_MSB = 0x9B

const BME280_REGISTER_DIG_P8_LSB = 0x9C
const BME280_REGISTER_DIG_P8_MSB = 0x9D

const BME280_REGISTER_DIG_P9_LSB = 0x9E
const BME280_REGISTER_DIG_P9_MSB = 0x9F

const BME280_REGISTER_DIG_H1 = 0xA1

const BME280_REGISTER_DIG_H2_LSB = 0xE1
const BME280_REGISTER_DIG_H2_MSB = 0xE2

const BME280_REGISTER_DIG_H3 = 0xE3

const BME280_REGISTER_DIG_H4_LSB = 0xE4
const BME280_REGISTER_DIG_H4_MSB = 0xE5

const BME280_REGISTER_DIG_H5_MSB = 0xE6

const BME280_REGISTER_DIG_H6 = 0xE7

const BME280_REGISTER_STARTCONVERSION = 0x12
const BME280_REGISTER_CTL_MEAS = 0xF4
const BME280_REGISTER_CTL_HUMIDITY = 0xF2
const BME280_REGISTER_CTL_CONFIG = 0xF5
const BME280_REGISTER_STATUS = 0xF3
const BME280_REGISTER_RESET = 0xE0

type BME280Driver struct {
	name       string
	connection I2c
	interval   time.Duration
	gobot.Eventer

	calT1       uint16
	calT2       int16
	calT3       int16

	calP1       uint16
	calP2       int16
	calP3       int16
	calP4       int16
	calP5       int16
	calP6       int16
	calP7       int16
	calP8       int16
	calP9       int16

	calH1       uint8
	calH2       int16
	calH3       uint8
	calH4       int16
	calH5       int16
	calH6       int8

	calTFine    int32

	Pressure    float32
	Humidity    float32
	Temperature float32
}

// NewBME280Driver creates a new driver with specified name and i2c interface
func NewBME280Driver(a I2c, name string, v ...time.Duration) *BME280Driver {
	m := &BME280Driver{
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
	}

	if len(v) > 0 {
		m.interval = v[0]
	}
	m.AddEvent(Error)
	return m
}

func (h *BME280Driver) Name() string                 { return h.name }
func (h *BME280Driver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start writes initialization bytes and reads from adaptor
// using specified interval to accelerometer andtemperature data
func (h *BME280Driver) Start() (errs []error) {
	var tempMSB		uint8
	var tempLSB		uint8
	var tempXSB		uint8
	var temperature 	uint32
	var tcvar1		float32
	var tcvar2		float32
	var temperatureComp 	float32

	var presMSB		uint8
	var presLSB		uint8
	var presXSB		uint8
	var pressure		uint32

	var pcvar1		float32
	var pcvar2		float32
	var pressureComp	float32

	var humidity16 		uint16
	var humidity 		int32
	var hcvar1		float32


	if err := h.initialization(); err != nil {
		return []error{err}
	}

	go func() {
		for {
			/*
			if err := h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_STARTCONVERSION, 0}); err != nil {
				gobot.Publish(h.Event(Error), err)
				continue

			}
			<-time.After(5 * time.Millisecond)
			*/

			if err := h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_STATUS}); err != nil {
				gobot.Publish(h.Event(Error), err)
				continue
			}

			ret, err := h.connection.I2cRead(bme280Address, BME280_DATA_FRAME_LENGTH)
			if err != nil {
				gobot.Publish(h.Event(Error), err)
				continue
			}
			if len(ret) == BME280_DATA_FRAME_LENGTH {
				buf := bytes.NewBuffer(ret)

				statusByte := uint8(0)
				binary.Read(buf, binary.BigEndian, &statusByte)    // F3
				binary.Read(buf, binary.BigEndian, &presMSB)    // F4
				binary.Read(buf, binary.BigEndian, &presMSB)    // F5
				binary.Read(buf, binary.BigEndian, &presMSB)    // F6

				
				if statusByte != 0x04 {
					continue
				}

				binary.Read(buf, binary.BigEndian, &presMSB)    // F7
				binary.Read(buf, binary.BigEndian, &presLSB)    // F8
				binary.Read(buf, binary.BigEndian, &presXSB)    // F9

				binary.Read(buf, binary.BigEndian, &tempMSB)    // FA
				binary.Read(buf, binary.BigEndian, &tempLSB)    // FB
				binary.Read(buf, binary.BigEndian, &tempXSB)    // FC

				binary.Read(buf, binary.BigEndian, &humidity16) // FD-FE
				humidity = int32(humidity16)
				

				presArray := []uint8{0, presMSB, presLSB, presXSB}

				pressure = binary.BigEndian.Uint32(presArray)

				pressure = pressure >> 4

				temperature = 0
				temperature = (uint32(tempMSB) << 12) | (uint32(tempLSB) << 4) | (uint32(tempXSB) >> 4)

				// Calculations from Bosch https://ae-bst.resource.bosch.com/media/_tech/media/datasheets/BST-BME280_DS001-11.pdf
				// temperature compensation
				tcvar1 = ((float32(temperature) / 16384.0) - (float32(h.calT1) / 1024.0)) * float32(h.calT2)
				tcvar2 = (((float32(temperature) / 131072.0) - (float32(h.calT1)/8192.0)) * ((float32(temperature) / 131072.0) - float32(h.calT1) / 8192.0)) * float32(h.calT3)
				temperatureComp = (tcvar1 + tcvar2) / 5120.0
				h.calTFine = int32(tcvar1 + tcvar2)
				h.Temperature = temperatureComp

				// pressure compensation
				pcvar1 = (float32(h.calTFine)/2.0) - 64000.0
				pcvar2 = pcvar1 * pcvar1 * (float32(h.calP6)) / 32768.0
				pcvar2 = pcvar2 + pcvar1 * (float32(h.calP5)) * 2.0
				pcvar2 = (pcvar2 / 4.0) + (float32(h.calP4) * 65536.0)
				pcvar1 = (float32(h.calP3) * pcvar1 * pcvar1 / 524288.0 + float32(h.calP2) * pcvar1) / 524288.0
				pcvar1 = (1.0 + pcvar1 / 32768.0) * (float32(h.calP1))
				if (pcvar1 == 0.0) { // avoid divide by zero
					return
				}
				pressureComp = 1048576.0 - float32(pressure)
				pressureComp = (pressureComp - (pcvar2 / 4096.0)) * (6250.0 / pcvar1)
				pcvar1 = float32(h.calP9) * pressureComp * pressureComp /  2147483648.0
				pcvar2 = pressureComp * float32(h.calP8) / 32768.0
				pressureComp = pressureComp + (pcvar1 + pcvar2 + float32(h.calP7)) / 16.0
				h.Pressure = pressureComp
				
				// Humidity compensation

/*
// Returns humidity in %rH as as double. Output value of "46.332" represents 46.332 %rH
double bme280_compensate_H_double(BME280_S32_t adc_H);
{
double var_H;
var_H = (((double)t_fine) - 76800.0);
var_H = (adc_H - (((double)dig_H4) * 64.0 + ((double)dig_H5) / 16384.0 * var_H)) *
(((double)dig_H2) / 65536.0 * 
(1.0 + ((double)dig_H6) / 67108864.0 * var_H *
(1.0 + ((double)dig_H3) / 67108864.0 * var_H)));

var_H = var_H * (1.0 - ((double)dig_H1) * var_H / 524288.0);
if (var_H > 100.0)
var_H = 100.0;
else if (var_H < 0.0)
var_H = 0.0;
return var_H;
}
*/

				hcvar1 = float32(h.calTFine) - 76800.0
				hcvar1 = (float32(humidity) - (float32(h.calH4) * 64.0 + 
						float32(h.calH5) / 16384.0 * float32(humidity))) *
					(float32(h.calH2)) / 65536.0 *
					(1.0 + float32(h.calH6) / 67108864.0 * float32(humidity) *
						(1.0 + float32(h.calH3) / 67108864.0 * float32(humidity)))

				hcvar1 = hcvar1 * (1.0 - float32(h.calH1) * hcvar1 / 524288.0)
				if (hcvar1 > 100.0) {
					hcvar1 = 100.0
				} else if (hcvar1 < 0.0) {
					hcvar1 = 0.0
				}
				
				h.Humidity = hcvar1
				
			}
			<-time.After(h.interval)
		}
	}()
	return
}

// Halt returns true if devices is halted successfully
func (h *BME280Driver) Halt() (err []error) { return }

// reset returns true if devices is reset successfully
func (h *BME280Driver) Reset() (err error) { 

	// Send the RESET instruction
	if err = h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_RESET, 0xB6}); err != nil {
		return errors.New("bme280 reset failed")
	}
	<-time.After(h.interval)

	// put the device to sleep
	if err = h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_CTL_MEAS, 0x00}); err != nil {
		return errors.New("bme280 sleep failed")
	}
	<-time.After(h.interval)


	// Configure the device's filtering, sampling, and SPI interface
	// Sample rate should be 000 - 0.5ms per sample because why not
	// Filter should be 16 (100)
	// SPI should be off

	// Therefore the control byte is
	// 0 1 2 3 4 5 6 7
	// 0 0 1 0 0 0 0 0
	// 0x20
	// of course we are wrong
	// 0 1 2 3 4 5 6 7
	// 0 0 0 0 0 1 0 0
	// 0x04
	// i give up, set 0x00

	if err = h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_CTL_CONFIG, 0x04}); err != nil {
		return errors.New("bme280 operation reconfig failed")
	}
	<-time.After(h.interval)

	// const BME280_REGISTER_CTL_HUMIDITY = 0xF2
	// following operation cribbed from adafruit
	// https://github.com/adafruit/Adafruit_BME280_Library

	if err = h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_CTL_HUMIDITY, 0x05}); err != nil {
		return errors.New("bme280 measurement reconfig failed")
	}
	<-time.After(h.interval)

	// Configure the device to active with oversampling x16 on all sensors
	if err = h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_CTL_MEAS, 0xFF}); err != nil {
		return errors.New("bme280 measurement reconfig failed")
	}
	<-time.After(h.interval)

	return 
}


func (h *BME280Driver) initialization() (err error) {

	if err = h.connection.I2cStart(bme280Address); err != nil {
		return
	}

	/*
	// Configure the device to active with oversampling x16 on all sensors
	if err = h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_CTL_MEAS, 0xFF}); err != nil {
		return
	}
	*/

	h.Reset()


	// Set the first byte to read configuration data from
	if err = h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_DIG_T1_LSB}); err != nil {
		return
	}

	// slurp 26 bytes of temp and pres calibration from bme280
	tp_ret, tp_err := h.connection.I2cRead(bme280Address, BME280_PRESSURE_TEMPERATURE_CALIB_DATA_LENGTH)
	if tp_err != nil {
		return
	}
	tp_buf := bytes.NewBuffer(tp_ret)

	binary.Read(tp_buf, binary.LittleEndian, &h.calT1) // 88-89

	binary.Read(tp_buf, binary.LittleEndian, &h.calT2) // 8A-8B

	binary.Read(tp_buf, binary.LittleEndian, &h.calT3) // 8C-8D

	binary.Read(tp_buf, binary.LittleEndian, &h.calP1) // 8E ...
	binary.Read(tp_buf, binary.LittleEndian, &h.calP2)
	binary.Read(tp_buf, binary.LittleEndian, &h.calP3)
	binary.Read(tp_buf, binary.LittleEndian, &h.calP4)
	binary.Read(tp_buf, binary.LittleEndian, &h.calP5)
	binary.Read(tp_buf, binary.LittleEndian, &h.calP6)
	binary.Read(tp_buf, binary.LittleEndian, &h.calP7)
	binary.Read(tp_buf, binary.LittleEndian, &h.calP8)
	binary.Read(tp_buf, binary.LittleEndian, &h.calP9) // ... 9F

	binary.Read(tp_buf, binary.LittleEndian, &h.calH1) // A0
	binary.Read(tp_buf, binary.LittleEndian, &h.calH1) // A1


	// Set the first byte to read humidity configuration data from
	if err = h.connection.I2cWrite(bme280Address, []byte{BME280_REGISTER_DIG_H2_LSB}); err != nil {
		return
	}

	// slurp 7 bytes of humidity calibration from bme280
	hu_ret, hu_err := h.connection.I2cRead(bme280Address, BME280_HUMIDITY_CALIB_DATA_LENGTH)
	if hu_err != nil {
		return
	}
	hu_buf := bytes.NewBuffer(hu_ret)

	// H4 and H5 laid out strangely on the bme280
	var addrE4	byte
	var addrE5	byte
	var addrE6	byte

	binary.Read(hu_buf, binary.LittleEndian, &h.calH2) // E1 ...
	binary.Read(hu_buf, binary.BigEndian, &h.calH3) // E3
	binary.Read(hu_buf, binary.BigEndian, &addrE4) // E4
	binary.Read(hu_buf, binary.BigEndian, &addrE5) // E5
	binary.Read(hu_buf, binary.BigEndian, &addrE6) // E6
	binary.Read(hu_buf, binary.BigEndian, &h.calH6) // ... E7

	h.calH4 = 0 + (int16(addrE4) << 4) | (int16(addrE5 & 0x0F))
	h.calH5 = 0 + (int16(addrE6) << 4) | (int16(addrE5) >> 4)
	return
}
