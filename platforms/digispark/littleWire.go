package digispark

//#cgo pkg-config: libusb
//#include "littleWire.h"
//#include "littleWire_servo.h"
//typedef usb_dev_handle littleWire;
import "C"

import (
	"errors"
	"fmt"
)

type lw interface {
	digitalWrite(uint8, uint8) error
	pinMode(uint8, uint8) error
	pwmInit() error
	pwmStop() error
	pwmUpdateCompare(uint8, uint8) error
	pwmUpdatePrescaler(uint) error
	servoInit() error
	servoUpdateLocation(uint8, uint8) error
	i2cInit() error
	i2cStart(address7bit uint8, direction uint8) error
	i2cWrite(sendBuffer []byte, length int, endWithStop uint8) error
	i2cRead(readBuffer []byte, length int, endWithStop uint8) error
	i2cUpdateDelay(duration uint) error
	error() error
}

type littleWire struct {
	lwHandle *C.littleWire
}

func littleWireConnect() *littleWire {
	return &littleWire{
		lwHandle: C.littleWire_connect(),
	}
}

func (l *littleWire) digitalWrite(pin uint8, state uint8) error {
	C.digitalWrite(l.lwHandle, C.uchar(pin), C.uchar(state))
	return l.error()
}

func (l *littleWire) pinMode(pin uint8, mode uint8) error {
	C.pinMode(l.lwHandle, C.uchar(pin), C.uchar(mode))
	return l.error()
}

func (l *littleWire) pwmInit() error {
	C.pwm_init(l.lwHandle)
	return l.error()
}

func (l *littleWire) pwmStop() error {
	C.pwm_stop(l.lwHandle)
	return l.error()
}

func (l *littleWire) pwmUpdateCompare(channelA uint8, channelB uint8) error {
	C.pwm_updateCompare(l.lwHandle, C.uchar(channelA), C.uchar(channelB))
	return l.error()
}

func (l *littleWire) pwmUpdatePrescaler(value uint) error {
	C.pwm_updatePrescaler(l.lwHandle, C.uint(value))
	return l.error()
}

func (l *littleWire) servoInit() error {
	C.servo_init(l.lwHandle)
	return l.error()
}

func (l *littleWire) servoUpdateLocation(locationA uint8, locationB uint8) error {
	C.servo_updateLocation(l.lwHandle, C.uchar(locationA), C.uchar(locationB))
	return l.error()
}

func (l *littleWire) i2cInit() error {
	C.i2c_init(l.lwHandle)
	return l.error()
}

// i2cStart starts the i2c communication; set direction to 1 for reading, 0 for writing
func (l *littleWire) i2cStart(address7bit uint8, direction uint8) error {
	if C.i2c_start(l.lwHandle, C.uchar(address7bit), C.uchar(direction)) == 1 {
		return nil
	}
	if err := l.error(); err != nil {
		return err
	}
	return fmt.Errorf("Littlewire i2cStart failed for %d in direction %d", address7bit, direction)
}

// i2cWrite sends byte(s) over i2c with a given length <= 4
func (l *littleWire) i2cWrite(sendBuffer []byte, length int, endWithStop uint8) error {
	C.i2c_write(l.lwHandle, (*C.uchar)(&sendBuffer[0]), C.uchar(length), C.uchar(endWithStop))
	return l.error()
}

// i2cRead reads byte(s) over i2c with a given length <= 8
func (l *littleWire) i2cRead(readBuffer []byte, length int, endWithStop uint8) error {
	C.i2c_read(l.lwHandle, (*C.uchar)(&readBuffer[0]), C.uchar(length), C.uchar(endWithStop))
	return l.error()
}

// i2cUpdateDelay updates i2c signal delay amount. Tune if neccessary to fit your requirements
func (l *littleWire) i2cUpdateDelay(duration uint) error {
	C.i2c_updateDelay(l.lwHandle, C.uint(duration))
	return l.error()
}

func (l *littleWire) error() error {
	str := C.GoString(C.littleWire_errorName())
	if str != "" {
		return errors.New(str)
	}
	return nil
}
