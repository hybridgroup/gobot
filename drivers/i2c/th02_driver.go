/*
 * Copyright (c) 2018 Nicholas Potts <nick@the-potts.com>
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
// This module was tested with a Grove Temperature and Humidty Sensor (High Accuracy)
// https://www.seeedstudio.com/Grove-Temperature-Humidity-Sensor-High-Accuracy-Min-p-1921.html

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
)

// TH02Address is the default address of device
const TH02Address = 0x40

//TH02ConfigReg is the configuration register
const TH02ConfigReg = 0x03

// TH02HighAccuracyTemp is the CONFIG write value to start reading temperature
const TH02HighAccuracyTemp = 0x11

//TH02HighAccuracyRH is the CONFIG write value to start reading high accuracy RH
const TH02HighAccuracyRH = 0x01

// TH02Driver is a Driver for a TH02 humidity and temperature sensor
type TH02Driver struct {
	Units      string
	name       string
	connector  Connector
	connection Connection
	Config
	addr     byte
	accuracy byte
	delay    time.Duration
}

// NewTH02Driver creates a new driver with specified i2c interface
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
	}
	//	s.SetAccuracy(TH02AccuracyHigh)

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

// SerialNumber returns the serial number of the chip
func (s *TH02Driver) SerialNumber() (sn uint32, err error) {
	ret, err := s.readRegister(0x11)
	return uint32(ret) >> 4, err
}

// Sample returns the temperature in celsius and relative humidity for one sample
func (s *TH02Driver) Sample() (temp float32, rh float32, err error) {
	if err := s.writeRegister(TH02ConfigReg, TH02HighAccuracyRH); err != nil {
		return 0, 0, err
	}

	rrh, err := s.readData()
	if err != nil {
		return 0, 0, err
	}
	rrh = rrh >> 4
	rh = float32(rrh)/16.0 - 24.0

	if err := s.writeRegister(TH02ConfigReg, TH02HighAccuracyTemp); err != nil {
		return 0, 0, err
	}

	rt, err := s.readData()
	if err != nil {
		return 0, rh, err
	}
	rt = rt / 4
	temp = float32(rt)/32.0 - 50.0

	switch s.Units {
	case "F":
		temp = 9.0/5.0 + 32.0
	}

	return temp, rh, nil

}

// getStatusRegister returns the device status register
func (s *TH02Driver) getStatusRegister() (status byte, err error) {
	return s.readRegister(TH02ConfigReg)
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

func (s *TH02Driver) waitForReady() error {
	start := time.Now()
	for {
		if time.Since(start) > 100*time.Millisecond {
			return fmt.Errorf("timeout on \\RDY")
		}
		reg, _ := s.readRegister(0x00)
		if reg == 0 {
			return nil
		}
	}
}

func (s *TH02Driver) readData() (uint16, error) {
	if err := s.waitForReady(); err != nil {
		return 1, err
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
