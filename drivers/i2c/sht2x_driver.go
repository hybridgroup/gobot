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

// SHT2xDriver is a driver for the SHT2x based devices.
//
// This module was tested with Sensirion SHT21 Breakout.

import (
	"errors"
	"time"

	"github.com/sigurn/crc8"
	"gobot.io/x/gobot"
)

const (
	// SHT2xDefaultAddress is the default I2C address for SHT2x
	SHT2xDefaultAddress = 0x40

	// SHT2xAccuracyLow is the faster, but lower accuracy sample setting
	//  0/1 = 8bit RH, 12bit Temp
	SHT2xAccuracyLow = byte(0x01)

	// SHT2xAccuracyMedium is the medium accuracy and speed sample setting
	//  1/0 = 10bit RH, 13bit Temp
	SHT2xAccuracyMedium = byte(0x80)

	// SHT2xAccuracyHigh is the high accuracy and slowest sample setting
	//  0/0 = 12bit RH, 14bit Temp
	//  Power on default is 0/0
	SHT2xAccuracyHigh = byte(0x00)

	// SHT2xTriggerTempMeasureHold is the command for measureing temperature in hold master mode
	SHT2xTriggerTempMeasureHold = 0xe3

	// SHT2xTriggerHumdMeasureHold is the command for measureing humidity in hold master mode
	SHT2xTriggerHumdMeasureHold = 0xe5

	// SHT2xTriggerTempMeasureNohold is the command for measureing humidity in no hold master mode
	SHT2xTriggerTempMeasureNohold = 0xf3

	// SHT2xTriggerHumdMeasureNohold is the command for measureing humidity in no hold master mode
	SHT2xTriggerHumdMeasureNohold = 0xf5

	// SHT2xWriteUserReg is the command for writing user register
	SHT2xWriteUserReg = 0xe6

	// SHT2xReadUserReg is the command for reading user register
	SHT2xReadUserReg = 0xe7

	// SHT2xReadUserReg is the command for reading user register
	SHT2xSoftReset = 0xfe
)

// SHT2xDriver is a Driver for a SHT2x humidity and temperature sensor
type SHT2xDriver struct {
	Units string

	name       string
	connector  Connector
	connection Connection
	Config
	sht2xAddress int
	accuracy     byte
	delay        time.Duration
	crcTable     *crc8.Table
}

// NewSHT2xDriver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewSHT2xDriver(a Connector, options ...func(Config)) *SHT2xDriver {
	// From the document "CRC Checksum Calculation -- For Safe Communication with SHT2x Sensors":
	crc8Params := crc8.Params{0x31, 0x00, false, false, 0x00, 0x00, "CRC-8/SENSIRION-SHT2x"}
	s := &SHT2xDriver{
		Units:     "C",
		name:      gobot.DefaultName("SHT2x"),
		connector: a,
		Config:    NewConfig(),
		crcTable:  crc8.MakeTable(crc8Params),
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// Name returns the name for this Driver
func (d *SHT2xDriver) Name() string { return d.name }

// SetName sets the name for this Driver
func (d *SHT2xDriver) SetName(n string) { d.name = n }

// Connection returns the connection for this Driver
func (d *SHT2xDriver) Connection() gobot.Connection { return d.connector.(gobot.Connection) }

// Start initializes the SHT2x
func (d *SHT2xDriver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(SHT2xDefaultAddress)

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return
	}

	if err = d.Reset(); err != nil {
		return
	}

	d.sendAccuracy()

	return
}

// Halt returns true if devices is halted successfully
func (d *SHT2xDriver) Halt() (err error) { return }

func (d *SHT2xDriver) Accuracy() byte { return d.accuracy }

// SetAccuracy sets the accuracy of the sampling
func (d *SHT2xDriver) SetAccuracy(acc byte) (err error) {
	d.accuracy = acc

	if d.connection != nil {
		err = d.sendAccuracy()
	}

	return
}

// Reset does a software reset of the device
func (d *SHT2xDriver) Reset() (err error) {
	if err = d.connection.WriteByte(SHT2xSoftReset); err != nil {
		return
	}

	time.Sleep(15 * time.Millisecond) // 15ms delay (from the datasheet 5.5)

	return
}

// Temperature returns the current temperature, in celsius degrees.
func (d *SHT2xDriver) Temperature() (temp float32, err error) {
	var rawT uint16
	if rawT, err = d.readSensor(SHT2xTriggerTempMeasureNohold); err != nil {
		return
	}

	// From the datasheet 6.2:
	// T[C] = -46.85 + 175.72 * St / 2^16
	temp = -46.85 + 175.72/65536.0*float32(rawT)

	return
}

// Humidity returns the current humidity in percentage of relative humidity
func (d *SHT2xDriver) Humidity() (humidity float32, err error) {
	var rawH uint16
	if rawH, err = d.readSensor(SHT2xTriggerHumdMeasureNohold); err != nil {
		return
	}

	// From the datasheet 6.1:
	// RH = -6 + 125 * Srh / 2^16
	humidity = -6.0 + 125.0/65536.0*float32(rawH)

	return
}

// sendCommandDelayGetResponse is a helper function to reduce duplicated code
func (d *SHT2xDriver) readSensor(cmd byte) (read uint16, err error) {
	if err = d.connection.WriteByte(cmd); err != nil {
		return
	}

	//Hang out while measurement is taken. 85ms max, page 9 of datasheet.
	time.Sleep(85 * time.Millisecond)

	//Comes back in three bytes, data(MSB) / data(LSB) / Checksum
	buf := make([]byte, 3)
	counter := 0
	for {
		var got int
		got, err = d.connection.Read(buf)
		counter++
		if counter > 50 {
			return
		}
		if err == nil {
			if got != 3 {
				err = ErrNotEnoughBytes
				return
			}
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

	//Store the result
	crc := crc8.Checksum(buf[0:2], d.crcTable)
	if buf[2] != crc {
		err = errors.New("Invalid crc")
		return
	}
	read = uint16(buf[0])<<8 | uint16(buf[1])
	read &= 0xfffc // clear two low bits (status bits)

	return
}

func (d *SHT2xDriver) sendAccuracy() (err error) {
	if err = d.connection.WriteByte(SHT2xReadUserReg); err != nil {
		return
	}
	userRegister, err := d.connection.ReadByte()
	if err != nil {
		return
	}

	userRegister &= 0x7e //Turn off the resolution bits
	acc := d.accuracy
	acc &= 0x81         //Turn off all other bits but resolution bits
	userRegister |= acc //Mask in the requested resolution bits

	//Request a write to user register
	_, err = d.connection.Write([]byte{SHT2xWriteUserReg, userRegister})
	if err != nil {
		return
	}

	userRegister, err = d.connection.ReadByte()
	if err != nil {
		return
	}

	return
}
