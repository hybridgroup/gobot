#ifndef LITTLEWIRE_SERVO_H
#define LITTLEWIRE_SERVO_H
/*
	Higher level servo driving library for Little Wire.

	Copyright (C) <2013> ihsan Kehribar <ihsan@kehribar.me>
	
	Permission is hereby granted, free of charge, to any person obtaining a copy of
	this software and associated documentation files (the "Software"), to deal in
	the Software without restriction, including without limitation the rights to
	use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
	of the Software, and to permit persons to whom the Software is furnished to do
	so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

#include "littleWire.h"

/*! \addtogroup Servo
  *  @brief Servo functions. Higher level access to PWM module.
  *  @{
  */

/**
  * Initialize the PWM module on the Little-Wire with the Servo special settings.
  *
  * @param lwHandle littleWire device pointer
  * @return (none)
  */
void servo_init(littleWire* lwHandle);

/**
  * Update servo locations 
  *
  * @param lwHandle littleWire device pointer
  * @param locationChannelA Location of servo connected to channel A ( in degrees )
  * @param locationChannelB Location of servo connected to channel B ( in degrees )
  * @return (none)
  */
void servo_updateLocation(littleWire* lwHandle,unsigned char locationChannelA,unsigned char locationChannelB);

/*! @} */

#endif