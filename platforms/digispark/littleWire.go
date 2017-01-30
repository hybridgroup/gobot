package digispark

//#cgo pkg-config: libusb
//#include "littleWire.h"
//#include "littleWire_servo.h"
//typedef usb_dev_handle littleWire;
import "C"

import "errors"

type lw interface {
	digitalWrite(uint8, uint8) error
	pinMode(uint8, uint8) error
	pwmInit() error
	pwmStop() error
	pwmUpdateCompare(uint8, uint8) error
	pwmUpdatePrescaler(uint) error
	servoInit() error
	servoUpdateLocation(uint8, uint8) error
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

func (l *littleWire) error() error {
	str := C.GoString(C.littleWire_errorName())
	if str != "" {
		return errors.New(str)
	}
	return nil
}
