package digispark

//#cgo LDFLAGS: -lusb
//#include "littleWire.h"
//#include "littleWire_servo.h"
//typedef usb_dev_handle littleWire;
import "C"

type lw interface {
	digitalWrite(uint8, uint8)
	pinMode(uint8, uint8)
	pwmInit()
	pwmStop()
	pwmUpdateCompare(uint8, uint8)
	pwmUpdatePrescaler(uint)
	servoInit()
	servoUpdateLocation(uint8, uint8)
}

type littleWire struct {
	lwHandle *C.littleWire
}

func littleWireConnect() *littleWire {
	return &littleWire{
		lwHandle: C.littleWire_connect(),
	}
}

func (l *littleWire) digitalWrite(pin uint8, state uint8) {
	C.digitalWrite(l.lwHandle, C.uchar(pin), C.uchar(state))
}

func (l *littleWire) pinMode(pin uint8, mode uint8) {
	C.pinMode(l.lwHandle, C.uchar(pin), C.uchar(mode))
}

func (l *littleWire) pwmInit() {
	C.pwm_init(l.lwHandle)
}

func (l *littleWire) pwmStop() {
	C.pwm_stop(l.lwHandle)
}

func (l *littleWire) pwmUpdateCompare(channelA uint8, channelB uint8) {
	C.pwm_updateCompare(l.lwHandle, C.uchar(channelA), C.uchar(channelB))
}

func (l *littleWire) pwmUpdatePrescaler(value uint) {
	C.pwm_updatePrescaler(l.lwHandle, C.uint(value))
}

func (l *littleWire) servoInit() {
	C.servo_init(l.lwHandle)
}

func (l *littleWire) servoUpdateLocation(locationA uint8, locationB uint8) {
	C.servo_updateLocation(l.lwHandle, C.uchar(locationA), C.uchar(locationB))
}
