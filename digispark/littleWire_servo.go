package gobotDigispark

//#include "littleWire_servo.h"
import "C"

//void servo_init(littleWire* lwHandle);
func (littleWire *LittleWire) ServoInit() {
	C.servo_init(littleWire.lwHandle)
}

//void servo_updateLocation(littleWire* lwHandle,unsigned char locationChannelA,unsigned char locationChannelB);
func (littleWire *LittleWire) ServoUpdateLocation(locationChannelA uint8, locationChannelB uint8) {
	C.servo_updateLocation(littleWire.lwHandle, C.uchar(locationChannelA), C.uchar(locationChannelB))
}
