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
	"time"

	"gobot.io/x/gobot"
)

const (

	// TH02Address is the default address of device
	TH02Address = 0x40

	//TH02ConfigReg is the configuration register
	TH02ConfigReg = 0x03
)

//Accuracy constants for the TH02 devices
const (
	TH02HighAccuracy = 0 //High Accuracy
	TH02LowAccuracy  = 1 //Lower Accuracy
)

// TH02Driver is a Driver for a TH02 humidity and temperature sensor
type TH02Driver struct {
	Units      string
	name       string
	connector  Connector
	connection Connection
	Config
	addr     byte
	accuracy byte
	heating  bool

	delay time.Duration
}

// NewTH02Driver creates a new driver with specified i2c interface.
// Defaults to:
//	- Using high accuracy (lower speed) measurements cycles.
//  - Emitting values in "C". If you want F, set Units to "F"
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewTH02Driver(a Connector, options ...func(Config)) *TH02Driver {
	s := &TH02Driver{
		Units:     "C",
		name:      gobot.DefaultName("TH02"),
		connector: a,
		addr:      TH02Address,
		Config:    NewConfig(),
		heating:   false,
	}

	s.SetAccuracy(1)

	for _, option := range options {
		option(s)
	}

	return s
}

// Name returns the name for this Driver
func (s *TH02Driver) Name() string { return s.name }

// SetName sets the name for this Driver
func (s *TH02Driver) SetName(n string) { s.name = n }

// Connection returns the connection for this Driver
func (s *TH02Driver) Connection() gobot.Connection { return s.connector.(gobot.Connection) }

// Start initializes the TH02
func (s *TH02Driver) Start() (err error) {
	bus := s.GetBusOrDefault(s.connector.GetDefaultBus())
	address := s.GetAddressOrDefault(int(s.addr))

	s.connection, err = s.connector.GetConnection(address, bus)
	return err
}

// Halt returns true if devices is halted successfully
func (s *TH02Driver) Halt() (err error) { return }

// SetAddress sets the address of the device
func (s *TH02Driver) SetAddress(address int) { s.addr = byte(address) }

// Accuracy returns the accuracy of the sampling
func (s *TH02Driver) Accuracy() byte { return s.accuracy }

// SetAccuracy sets the accuracy of the sampling.  It will only be used on the next
// measurment request.  Invalid value will use the default of High
func (s *TH02Driver) SetAccuracy(a byte) {
	if a == TH02LowAccuracy {
		s.accuracy = a
	} else {
		s.accuracy = TH02HighAccuracy
	}
}

// SerialNumber returns the serial number of the chip
func (s *TH02Driver) SerialNumber() (sn uint32, err error) {
	ret, err := s.readRegister(0x11)
	return uint32(ret) >> 4, err
}

// Heater returns true if the heater is enabled
func (s *TH02Driver) Heater() (status bool, err error) {
	st, err := s.readRegister(0x11)
	return (0x02 & st) == 0x02, err
}

func (s *TH02Driver) applysettings(base byte) byte {
	if s.accuracy == TH02LowAccuracy {
		base = base & 0xd5
	} else {
		base = base | 0x20
	}
	if s.heating {
		base = base & 0xfd
	} else {
		base = base | 0x02
	}
	base = base | 0x01 //set the "sample" bit
	return base
}

// Sample returns the temperature in celsius and relative humidity for one sample
func (s *TH02Driver) Sample() (temperature float32, relhumidity float32, _ error) {

	if err := s.writeRegister(TH02ConfigReg, s.applysettings(0x10)); err != nil {
		return 0, 0, err
	}

	rawrh, err := s.readData()
	if err != nil {
		return 0, 0, err
	}
	relhumidity = float32(rawrh>>4)/16.0 - 24.0

	if err := s.writeRegister(TH02ConfigReg, s.applysettings(0x00)); err != nil {
		return 0, relhumidity, err
	}
	rawt, err := s.readData()
	if err != nil {
		return 0, relhumidity, err
	}
	temperature = float32(rawt>>2)/32.0 - 50.0

	switch s.Units {
	case "F":
		temperature = 9.0/5.0*temperature + 32.0
	}

	return temperature, relhumidity, nil

}

//writeRegister writes the value to the register.
func (s *TH02Driver) writeRegister(reg, value byte) error {
	_, err := s.connection.Write([]byte{reg, value})
	return err
}

//readRegister returns the value of a single regusterm and a non-nil error on problem
func (s *TH02Driver) readRegister(reg byte) (byte, error) {
	if _, err := s.connection.Write([]byte{reg}); err != nil {
		return 0, err
	}
	rcvd := make([]byte, 1)
	_, err := s.connection.Read(rcvd)
	return rcvd[0], err
}

/*waitForReady blocks for up to the passed duration (which defaults to 50mS if nil)
until the ~RDY bit is cleared, meanign a sample has been fully sampled and is ready for reading.

This is greedy.
*/
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

		//yes, i am eating the error.
		if reg, _ := s.readRegister(0x00); reg == 0 {
			return nil
		}
	}
}

/*readData fetches the data from the data 'registers'*/
func (s *TH02Driver) readData() (uint16, error) {
	if err := s.waitForReady(nil); err != nil {
		return 0, err
	}

	if n, err := s.connection.Write([]byte{0x01}); err != nil || n != 1 {
		return 0, fmt.Errorf("n=%d not 1, or err = %v", n, err)
	}
	rcvd := make([]byte, 3)
	n, err := s.connection.Read(rcvd)
	if err != nil || n != 3 {
		return 0, fmt.Errorf("n=%d not 3, or err = %v", n, err)
	}
	return uint16(rcvd[1])<<8 + uint16(rcvd[2]), nil

}
