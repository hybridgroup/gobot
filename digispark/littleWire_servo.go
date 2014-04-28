package digispark

//#include "littleWire_servo.h"
import "C"

//void servo_init(littleWire* lwHandle);
func (l *LittleWire) ServoInit() {
	C.servo_init(l.lwHandle)
}

//void servo_updateLocation(littleWire* lwHandle,unsigned char locationChannelA,unsigned char locationChannelB);
func (l *LittleWire) ServoUpdateLocation(locationChannelA uint8, locationChannelB uint8) {
	C.servo_updateLocation(l.lwHandle, C.uchar(locationChannelA), C.uchar(locationChannelB))
}
