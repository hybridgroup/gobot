package digispark

//#include "littleWire_servo.h"
import "C"

func (l *LittleWire) ServoInit() {
	C.servo_init(l.lwHandle)
}

func (l *LittleWire) ServoUpdateLocation(locationChannelA uint8, locationChannelB uint8) {
	C.servo_updateLocation(l.lwHandle, C.uchar(locationChannelA), C.uchar(locationChannelB))
}
