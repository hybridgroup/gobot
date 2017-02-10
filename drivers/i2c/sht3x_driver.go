/*
 * Copyright (c) 2016-2017 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
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

// SHT3xDriver is a driver for the SHT3x-D based devices.
//
// This module was tested with AdaFruit Sensiron SHT32-D Breakout.
// https://www.adafruit.com/products/2857

import (
	"errors"
	"time"

	"github.com/sigurn/crc8"
	"gobot.io/x/gobot"
)

// SHT3xAddressA is the default address of device
const SHT3xAddressA = 0x44

// SHT3xAddressB is the optional address of device
const SHT3xAddressB = 0x45

// SHT3xAccuracyLow is the faster, but lower accuracy sample setting
const SHT3xAccuracyLow = 0x16

// SHT3xAccuracyMedium is the medium accuracy and speed sample setting
const SHT3xAccuracyMedium = 0x0b

// SHT3xAccuracyHigh is the high accuracy and slowest sample setting
const SHT3xAccuracyHigh = 0x00

var (
	crc8Params         = crc8.Params{0x31, 0xff, false, false, 0x00, 0xf7, "CRC-8/SENSIRON"}
	ErrInvalidAccuracy = errors.New("Invalid accuracy")
	ErrInvalidCrc      = errors.New("Invalid crc")
	ErrInvalidTemp     = errors.New("Invalid temperature units")
)

// SHT3xDriver is a Driver for a SHT3x humidity and temperature sensor
type SHT3xDriver struct {
	Units string

	name       string
	connector  Connector
	connection Connection
	Config
	sht3xAddress int
	accuracy     byte
	delay        time.Duration
	crcTable     *crc8.Table
}

// NewSHT3xDriver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewSHT3xDriver(a Connector, options ...func(Config)) *SHT3xDriver {
	s := &SHT3xDriver{
		Units:        "C",
		name:         gobot.DefaultName("SHT3x"),
		connector:    a,
		Config:       NewConfig(),
		sht3xAddress: SHT3xAddressA,
		crcTable:     crc8.MakeTable(crc8Params),
	}
	s.SetAccuracy(SHT3xAccuracyHigh)

	for _, option := range options {
		option(s)
	}

	return s
}

// Name returns the name for this Driver
func (s *SHT3xDriver) Name() string { return s.name }

// SetName sets the name for this Driver
func (s *SHT3xDriver) SetName(n string) { s.name = n }

// Connection returns the connection for this Driver
func (s *SHT3xDriver) Connection() gobot.Connection { return s.connector.(gobot.Connection) }

// Start initializes the SHT3x
func (s *SHT3xDriver) Start() (err error) {
	bus := s.GetBusOrDefault(s.connector.GetDefaultBus())
	address := s.GetAddressOrDefault(s.sht3xAddress)

	s.connection, err = s.connector.GetConnection(address, bus)
	return
}

// Halt returns true if devices is halted successfully
func (s *SHT3xDriver) Halt() (err error) { return }

// SetAddress sets the address of the device
func (s *SHT3xDriver) SetAddress(address int) { s.sht3xAddress = address }

// Accuracy returns the accuracy of the sampling
func (s *SHT3xDriver) Accuracy() byte { return s.accuracy }

// SetAccuracy sets the accuracy of the sampling
func (s *SHT3xDriver) SetAccuracy(a byte) (err error) {
	switch a {
	case SHT3xAccuracyLow:
		s.delay = 5 * time.Millisecond // Actual max is 4, wait 1 ms longer
	case SHT3xAccuracyMedium:
		s.delay = 7 * time.Millisecond // Actual max is 6, wait 1 ms longer
	case SHT3xAccuracyHigh:
		s.delay = 16 * time.Millisecond // Actual max is 15, wait 1 ms longer
	default:
		err = ErrInvalidAccuracy
		return
	}

	s.accuracy = a

	return
}

// SerialNumber returns the serial number of the chip
func (s *SHT3xDriver) SerialNumber() (sn uint32, err error) {
	ret, err := s.sendCommandDelayGetResponse([]byte{0x37, 0x80}, nil, 2)
	if nil == err {
		sn = (uint32(ret[0]) << 16) | uint32(ret[1])
	}

	return
}

// Heater returns true if the heater is enabled
func (s *SHT3xDriver) Heater() (status bool, err error) {
	sr, err := s.getStatusRegister()
	if err == nil {
		if (1 << 13) == (sr & (1 << 13)) {
			status = true
		}
	}
	return
}

// SetHeater enables or disables the heater on the device
func (s *SHT3xDriver) SetHeater(enabled bool) (err error) {
	out := []byte{0x30, 0x66}
	if true == enabled {
		out[1] = 0x6d
	}
	_, err = s.connection.Write(out)
	return
}

// Sample returns the temperature in celsius and relative humidity for one sample
func (s *SHT3xDriver) Sample() (temp float32, rh float32, err error) {
	ret, err := s.sendCommandDelayGetResponse([]byte{0x24, s.accuracy}, &s.delay, 2)
	if nil != err {
		return
	}

	// From the datasheet:
	// RH = 100 * Srh / (2^16 - 1)
	rhSample := uint64(ret[1])
	rh = float32((uint64(1000000)*rhSample)/uint64(0xffff)) / 10000.0

	tempSample := uint64(ret[0])
	switch s.Units {
	case "C":
		// From the datasheet:
		// T[C] = -45 + 175 * (St / (2^16 - 1))
		temp = float32((uint64(1750000)*tempSample)/uint64(0xffff)-uint64(450000)) / 10000.0
	case "F":
		// From the datasheet:
		// T[F] = -49 + 315 * (St / (2^16 - 1))
		temp = float32((uint64(3150000)*tempSample)/uint64(0xffff)-uint64(490000)) / 10000.0
	default:
		err = ErrInvalidTemp
	}

	return
}

// getStatusRegister returns the device status register
func (s *SHT3xDriver) getStatusRegister() (status uint16, err error) {
	ret, err := s.sendCommandDelayGetResponse([]byte{0xf3, 0x2d}, nil, 1)
	if nil == err {
		status = ret[0]
	}
	return
}

// sendCommandDelayGetResponse is a helper function to reduce duplicated code
func (s *SHT3xDriver) sendCommandDelayGetResponse(send []byte, delay *time.Duration, expect int) (read []uint16, err error) {
	if _, err = s.connection.Write(send); err != nil {
		return
	}

	if nil != delay {
		time.Sleep(*delay)
	}

	buf := make([]byte, 3*expect)
	got, err := s.connection.Read(buf)
	if err != nil {
		return
	}
	if got != (3 * expect) {
		err = ErrNotEnoughBytes
		return
	}

	read = make([]uint16, expect)
	for i := 0; i < expect; i++ {
		crc := crc8.Checksum(buf[i*3:i*3+2], s.crcTable)
		if buf[i*3+2] != crc {
			err = ErrInvalidCrc
			return
		}
		read[i] = uint16(buf[i*3])<<8 | uint16(buf[i*3+1])
	}

	return
}
