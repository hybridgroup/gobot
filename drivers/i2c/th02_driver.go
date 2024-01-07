/*
 * Copyright (c) 2018 Nick Potts <nick@the-potts.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package i2c

// TH02Driver is a driver for the TH02-D based devices.
//
// This module was tested with a Grove "Temperature&Humidity Sensor (High-Accuracy & Mini ) v1.0"
// from https://www.seeedstudio.com/Grove-Temperature-Humidity-Sensor-High-Accuracy-Min-p-1921.htm
// Datasheet is at http://www.hoperf.com/upload/sensor/TH02_V1.1.pdf

import (
	"fmt"
	"log"
	"time"
)

const (
	th02Debug          = false
	th02DefaultAddress = 0x40
)

const (
	th02Reg_Status  = 0x00
	th02Reg_DataMSB = 0x01
	th02Reg_DataLSB = 0x02
	th02Reg_Config  = 0x03
	th02Reg_ID      = 0x11

	// th02Status_ReadyBit = 0x01 // D0 is /RDY

	th02Config_StartBit = 0x01 // D0 is START
	th02Config_HeatBit  = 0x02 // D1 is HEAT
	th02Config_TempBit  = 0x10 // D4 is TEMP (if not set read humidity)
	th02Config_FastBit  = 0x20 // D5 is FAST (if set use 18 ms, but lower accuracy T: 13 bit, H: 11 bit)
)

// Accuracy constants for the TH02 devices (deprecated, use WithFastMode() instead)
const (
	TH02HighAccuracy = 0 // High Accuracy (T: 14 bit, H: 12 bit), normal (35 ms)
	TH02LowAccuracy  = 1 // Lower Accuracy (T: 13 bit, H: 11 bit), fast (18 ms)
)

// TH02Driver is a Driver for a TH02 humidity and temperature sensor
type TH02Driver struct {
	*Driver
	Units    string
	heating  bool
	fastMode bool
}

// NewTH02Driver creates a new driver with specified i2c interface.
// Defaults to:
//   - Using high accuracy (lower speed) measurements cycles.
//   - Emitting values in "C". If you want F, set Units to "F"
//
// Params:
//
//	conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewTH02Driver(a Connector, options ...func(Config)) *TH02Driver {
	s := &TH02Driver{
		Driver:   NewDriver(a, "TH02", th02DefaultAddress, options...),
		Units:    "C",
		heating:  false,
		fastMode: false,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// WithTH02FastMode option sets the fast mode (leads to lower accuracy).
// Valid settings are <=0 (off), >0 (on).
func WithTH02FastMode(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*TH02Driver)
		if ok {
			d.fastMode = (val > 0)
		} else if th02Debug {
			log.Printf("Trying to set fast mode for non-TH02Driver %v", c)
		}
	}
}

// Accuracy returns the accuracy of the sampling (deprecated, use FastMode() instead)
func (s *TH02Driver) Accuracy() byte {
	if s.fastMode {
		return TH02LowAccuracy
	}
	return TH02HighAccuracy
}

// SetAccuracy sets the accuracy of the sampling. (deprecated, use WithFastMode() instead)
// It will only be used on the next measurement request.  Invalid value will use the default of High
func (s *TH02Driver) SetAccuracy(a byte) {
	s.fastMode = (a == TH02LowAccuracy)
}

// SerialNumber returns the serial number of the chip
func (s *TH02Driver) SerialNumber() (uint8, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ret, err := s.connection.ReadByteData(th02Reg_ID)
	return ret >> 4, err
}

// FastMode returns true if the fast mode is enabled in the device
func (s *TH02Driver) FastMode() (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	cfg, err := s.connection.ReadByteData(th02Reg_Config)
	return (th02Config_FastBit & cfg) == th02Config_FastBit, err
}

// SetHeater sets the heater of the device to the given state.
func (s *TH02Driver) SetHeater(state bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.heating = state
	return s.connection.WriteByteData(th02Reg_Config, s.createConfig(false, false))
}

// Heater returns true if the heater is enabled in the device
func (s *TH02Driver) Heater() (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	cfg, err := s.connection.ReadByteData(th02Reg_Config)
	return (th02Config_HeatBit & cfg) == th02Config_HeatBit, err
}

// Sample returns the temperature in celsius and relative humidity for one sample
//
//nolint:nonamedreturns // is sufficient here
func (s *TH02Driver) Sample() (temperature float32, relhumidity float32, _ error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// read humidity
	if err := s.connection.WriteByteData(th02Reg_Config, s.createConfig(true, false)); err != nil {
		return 0, 0, err
	}

	rawrh, err := s.waitAndReadData()
	if err != nil {
		return 0, 0, err
	}
	relhumidity = float32(rawrh>>4)/16.0 - 24.0

	// read temperature
	if err := s.connection.WriteByteData(th02Reg_Config, s.createConfig(true, true)); err != nil {
		return 0, relhumidity, err
	}
	rawt, err := s.waitAndReadData()
	if err != nil {
		return 0, relhumidity, err
	}
	temperature = float32(rawt>>2)/32.0 - 50.0

	if s.Units == "F" {
		temperature = 9.0/5.0*temperature + 32.0
	}

	return temperature, relhumidity, nil
}

func (s *TH02Driver) createConfig(measurement bool, readTemp bool) byte {
	cfg := byte(0x00)
	if measurement {
		cfg = cfg | th02Config_StartBit
		if readTemp {
			cfg = cfg | th02Config_TempBit
		}
		if s.fastMode {
			cfg = cfg | th02Config_FastBit
		}
	}
	if s.heating {
		cfg = cfg | th02Config_HeatBit
	}
	return cfg
}

func (s *TH02Driver) waitAndReadData() (uint16, error) {
	if err := s.waitForReady(nil); err != nil {
		return 0, err
	}

	rcvd := make([]byte, 2)
	err := s.connection.ReadBlockData(th02Reg_DataMSB, rcvd)
	if err != nil {
		return 0, err
	}
	return uint16(rcvd[0])<<8 + uint16(rcvd[1]), nil
}

// waitForReady blocks for up to the passed duration (which defaults to 50mS if nil)
// until the ~RDY bit is cleared, meaning a sample has been fully sampled and is ready for reading.
func (s *TH02Driver) waitForReady(dur *time.Duration) error {
	wait := 100 * time.Millisecond
	if dur != nil {
		wait = *dur
	}
	start := time.Now()
	for {
		if time.Since(start) > wait {
			return fmt.Errorf("timeout on \\RDY")
		}

		if reg, err := s.connection.ReadByteData(th02Reg_Status); (reg == 0) && (err == nil) {
			return nil
		}
		time.Sleep(wait / 10)
	}
}
