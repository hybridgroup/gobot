#ifndef LITTLEWIRE_H
#define LITTLEWIRE_H

/*
  Cross platform computer interface library for Little Wire project

  http://littlewire.github.io

  Copyright (C) <2013> ihsan Kehribar <ihsan@kehribar.me>
  Copyright (C) <2013> Omer Kilic <omerkilic@gmail.com>

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

#ifdef _WIN32
   #include <lusb0_usb.h>   // this is libusb, see http://libusb.sourceforge.net/
#else
   #include <usb.h>       // this is libusb, see http://libusb.sourceforge.net/
#endif
#include "opendevice.h"			// common code moved to separate module
#include "littleWire_util.h"
#include <stdio.h>

#define	VENDOR_ID 0x1781
#define	PRODUCT_ID 0x0c9f
#define USB_TIMEOUT 5000
#define RX_BUFFER_SIZE 64

#define INPUT 1
#define OUTPUT 0

#define ENABLE 1
#define DISABLE 0

#define HIGH 1
#define LOW	 0

#define AUTO_CS 1
#define MANUAL_CS 0

// Voltage ref definition
#define VREF_VCC 0
#define VREF_1100mV 1
#define VREF_2560mV 2 

// I2C definition
#define END_WITH_STOP 1
#define NO_STOP 0
#define READ 1
#define WRITE 0

// General Purpose Pins
#define PIN1 1
#define PIN2 2
#define PIN3 5
#define PIN4 0

// ADC Channels
#define ADC_PIN3 0
#define ADC_PIN2 1
#define ADC_TEMP_SENS 2

// PWM Pins
#define PWM1 PIN4
#define PWM2 PIN1

// Aliases
#define ADC0 ADC_PIN3
#define ADC1 ADC_PIN2
#define ADC2 ADC_TEMP_SENS
#define PWMA PWM1
#define PWMB PWM2

// 'AVR ISP' Pins
#define SCK_PIN PIN2
#define MISO_PIN PIN1
#define MOSI_PIN PIN4
#define RESET_PIN PIN3

extern unsigned char rxBuffer[RX_BUFFER_SIZE]; /* This has to be unsigned for the data's sake */
extern unsigned char ROM_NO[8];
extern int lwStatus;

/*! \addtogroup General
*  @brief General library functions
*  @{
*/

typedef usb_dev_handle littleWire;

typedef struct lwCollection
{
  struct usb_device* lw_device;
  int serialNumber;  
}lwCollection;

extern lwCollection lwResults[16];

extern int lw_totalDevices;

/**
  * Tries to cache all the littleWire devices and stores them in lwResults array. \n
  * Don't actually connects to any of the device(s).
  *
  * @param (none)
  * @return Total number of littleWire devices found in the USB system.
  */
int littlewire_search();

/**
  * Tries to connect to the spesific littleWire device by array id. 
  *
  * @param desiredID array index of the lwResults array.
  * @return littleWire pointer for healthy connection, NULL for a failed trial.
  */
littleWire* littlewire_connect_byID(int desiredID);

/**
  * Tries to connect to the spesific littleWire with a given serial number. \n
  * If multiple devices have the same serial number, it connects to the last one it finds 
  *
  * @param mySerial Serial number of the desired littlewire device.
  * @return littleWire pointer for healthy connection, NULL for a failed trial.
  */
littleWire* littlewire_connect_bySerialNum(int mySerial);

/**
  * Tries to connect to the first littleWire device that libusb can find. 
  *
  * @param (none)
  * @return littleWire pointer for healthy connection, NULL for a failed trial.
  */
littleWire* littleWire_connect();

/**
  * Reads the firmware version of the Little Wire \n
  * Format: 0xXY => X: Primary version Y: Minor version
  *
  * @param (none)
  * @return Firmware version
  */ 
unsigned char readFirmwareVersion(littleWire* lwHandle);

/**
  * Changes the USB serial number of the Little Wire 
  * 
  * @param serialNumber Serial number integer value (100-99)
  * @return (none)
  */
void changeSerialNumber(littleWire* lwHandle,int serialNumber);

/**
  * Sends a custom message to the device. \n
  * Useful when developing new features in the firmware.
  *
  * @param receiveBuffer Returned data buffer
  * @param command Firmware command
  * @param d1 data[0] for the command
  * @param d2 data[1] for the command
  * @param d3 data[2] for the command
  * @param d4 data[3] for the command
  * @return status
  */
int customMessage(littleWire* lwHandle,unsigned char* receiveBuffer,unsigned char command,unsigned char d1,unsigned char d2, unsigned char d3, unsigned char d4);

/**
  * Returns the numeric value of the status of the last communication attempt
  *
  * @param (none)
  * @return Numeric value of the status of the last communication attempt
  */
int littleWire_error ();

/**
  * Returns the string version of the last communication attempt status if there was an error
  *
  * @param (none)
  * @return String version of the last communication attempt status if there was an error
  */
char *littleWire_errorName ();

/*! @} */

/*! \addtogroup GPIO
*  @brief GPIO library functions with Arduino-like syntax
*  @{
*/

/**
  * Set pin value
  *
  * @param lwHandle littleWire device pointer
  * @param pin Pin name (\b PIN1 , \b PIN2 , \b PIN3 or \b PIN4 )
  * @param state Pin state (\b HIGH or \b LOW)
  * @return (none)
  */
void digitalWrite(littleWire* lwHandle, unsigned char pin, unsigned char state);

/**
  * Set pin as input/output
  *
  * @param lwHandle littleWire device pointer
  * @param pin Pin name (\b PIN1 , \b PIN2 , \b PIN3 or \b PIN4 )
  * @param mode Mode of pin (\b INPUT or \b OUTPUT)
  * @return (none)
  */
void pinMode(littleWire* lwHandle, unsigned char pin, unsigned char mode);

/**
  * Read pin value
  *
  * @param lwHandle littleWire device pointer
  * @param pin Pin name (\b PIN1 , \b PIN2 , \b PIN3 or \b PIN4 )
  * @return Pin state (\b HIGH or \b LOW)
  */
unsigned char digitalRead(littleWire* lwHandle, unsigned char pin);

/**
  * Sets the state of the internal pullup resistor. 
  * \n Call this function after you assign the pin as an input.
  *
  * @param lwHandle littleWire device pointer
  * @param pin Pin name (\b PIN1 , \b PIN2 , \b PIN3 or \b PIN4 )
  * @rparam state (\b ENABLE or \b DISABLE )
  * @return (none)
  */
void internalPullup(littleWire* lwHandle, unsigned char pin, unsigned char state);

/*! @} */

/*! \addtogroup ADC
*  @brief Analog to digital converter functions.
*  @{
*/

/**
  * Initialize the analog module. VREF_VCC is the standard voltage coming from the USB plug
  * \n Others are the Attiny's internal voltage references.
  *
  * @param lwHandle littleWire device pointer
  * @param voltageRef (\b VREF_VCC , \b VREF_110mV or \b VREF_2560mV )
  * @return (none)
  */
void analog_init(littleWire* lwHandle, unsigned char voltageRef);

/**
  * Read analog voltage. Analog voltage reading from ADC_PIN3 isn't advised (it is a bit noisy) but supported. Use it at your own risk.
  * \n For more about internal temperature sensor, look at the Attiny85 datasheet.
  *
  * @param lwHandle littleWire device pointer
  * @param channel Source of ADC reading (\b ADC_PIN2 , \b ADC_PIN3 or \b ADC_TEMP_SENS )
  * @return 10 bit ADC result
  */
unsigned int analogRead(littleWire* lwHandle, unsigned char channel);

/*! @} */

/*! \addtogroup PWM
*  @brief Pulse width modulation functions.
*  @{
*/

/**
  * Initialize the PWM module on the Little-Wire 
  *
  * @param lwHandle littleWire device pointer
  * @return (none)
  */
void pwm_init(littleWire* lwHandle);

/**
  * Stop the PWM module on the Little-Wire 
  *
  * @param lwHandle littleWire device pointer
  * @return (none)
  */
void pwm_stop(littleWire* lwHandle);

/**
  * Update the compare values of the PWM output pins. Resolution is 8 bit.
  *
  * @param lwHandle littleWire device pointer
  * @param channelA Compare value of \b PWMA pin
  * @param channelB Compare value of \b PWMB pin
  * @return (none)
  */
void pwm_updateCompare(littleWire* lwHandle, unsigned char channelA, unsigned char channelB);

/**
  * Update the prescaler of the PWM module. Adjust this value according to your need for speed in PWM output. Default is 1024. Lower prescale means higher frequency PWM output.
  *
  * @param lwHandle littleWire device pointer
  * @param value Presecaler value (\b 1024, \b 256, \b 64, \b 8 or \b 1) 
  * @return (none)
  */
void pwm_updatePrescaler(littleWire* lwHandle, unsigned int value);

/*! @} */

/*! \addtogroup SPI
*  @brief Serial peripheral interface functions.
*  @{
*/

/**
  * Initialize the SPI module on the Little-Wire 
  *
  * @param lwHandle littleWire device pointer
  * @return (none)
  */
void spi_init(littleWire* lwHandle);

/**
  * Send SPI message(s). SPI Mode is 0. 
  *
  * @param lwHandle littleWire device pointer
  * @param sendBuffer Message array to send
  * @param inputBuffer Returned answer message
  * @param length Message length - maximum 4
  * @param mode \b AUTO_CS or \b MANUAL_CS
  * @return (none)
  */
void spi_sendMessage(littleWire* lwHandle, unsigned char * sendBuffer, unsigned char * inputBuffer, unsigned char length ,unsigned char mode);

/**
  * Send one byte SPI message over MOSI pin. Slightly slower than the actual one.
  * \n There isn't any chip select control involved. Useful for debug console app
  *
  * @param lwHandle littleWire device pointer
  * @param message Message to send
  * @return Received SPI message
  */
unsigned char debugSpi(littleWire* lwHandle, unsigned char message);

/**
  * Change the SPI message frequency by adjusting delay duration. By default, Little-Wire sends the SPI messages with maximum speed. 
  * \n If your hardware can't catch up with the speed, increase the duration value to lower the SPI speed.
  *
  * @param lwHandle littleWire device pointer
  * @param duration Amount of delay. 
  * @return (none)
  */
void spi_updateDelay(littleWire* lwHandle, unsigned int duration);

/*! @} */

/*! \addtogroup I2C
  *  @brief Inter IC communication functions.
  *  @{
  */

/**
  * Initialize the I2C module on the Little-Wire 
  *
  * @return (none)
  */
void i2c_init(littleWire* lwHandle);

/**
  * Start the i2c communication
  *
  * @param lwHandle littleWire device pointer
  * @param address 7 bit slave address.
  * @param direction ( \b READ or \b WRITE )
  * @return 1 if received ACK
  */
unsigned char i2c_start(littleWire* lwHandle, unsigned char address7bit, unsigned char direction);

/**
  * Send byte(s) over i2c bus
  *
  * @param lwHandle littleWire device pointer
  * @param sendBuffer Message array to send
  * @param length Message length -> Max = 4
  * @param endWithStop Should we send a STOP condition after this buffer? ( \b END_WITH_STOP or \b NO_STOP )
  * @return (none)
  */
void i2c_write(littleWire* lwHandle, unsigned char* sendBuffer, unsigned char length, unsigned char endWithStop);

/**
  * Read byte(s) over i2c bus
  * 
  * @param lwHandle littleWire device pointer
  * @param readBuffer Returned message array
  * @param length Message length -> Max = 8
  * @param endWithStop Should we send a STOP condition after this buffer? ( \b END_WITH_STOP or \b NO_STOP )
  * @return (none)
  */
void i2c_read(littleWire* lwHandle, unsigned char* readBuffer, unsigned char length, unsigned char endWithStop);

/**
  * Update i2c signal delay amount. Tune if neccessary to fit your requirements.
  * 
  * @param lwHandle littleWire device pointer
  * @param duration Delay amount
  * @return (none)
  */
void i2c_updateDelay(littleWire* lwHandle, unsigned int duration);

/*! @} */ 

/*! \addtogroup Onewire
 *  @brief Onewire functions.
 *  @{
 */

/**
  * Send a single bit over onewire bus.
  *
  * @param lwHandle littleWire device pointer  
  * @param bitValue \b 1 or \b 0
  * @return (none)
  */
void onewire_sendBit(littleWire* lwHandle, unsigned char bitValue);

/**
  * Send a byte over onewire bus.
  *
  * @param lwHandle littleWire device pointer
  * @param messageToSend Message to send
  * @return (none)
  */
void onewire_writeByte(littleWire* lwHandle, unsigned char messageToSend);

/**
  * Read a byte over onewire bus.
  *
  * @param lwHandle littleWire device pointer
  * @return Read byte
  */
unsigned char onewire_readByte(littleWire* lwHandle);

/**
  * Read a single bit over onewire bus
  *
  * @param lwHandle littleWire device pointer
  * @return Read bit ( \b 1 or \b 0 )
  */
unsigned char onewire_readBit(littleWire* lwHandle);

/**
  * Send a reset pulse over onewire bus
  * 
  * @param lwHandle littleWire device pointer
  * @return Nonzero if any device presents on the bus
  */
unsigned char onewire_resetPulse(littleWire* lwHandle);

/**
  * Start searching for device address on the onewire bus.
  * \n Read the 8 byte address from \b ROM_NO array
  *
  * @param lwHandle littleWire device pointer  
  * @return Nonzero if any device found
  */
int onewire_firstAddress(littleWire* lwHandle);

/**
  * Try to find the next adress on the onewire bus.
  * \n Read the 8 byte address from \b ROM_NO array
  *
  * @param lwHandle littleWire device pointer  
  * @return Nonzero if any new device found
  */
int onewire_nextAddress(littleWire* lwHandle);

/*! @} */ 

/*! \addtogroup SOFT_PWM
  *  @brief Software PWM functions. Designed to be used with RGB LEDs.
  *  @{
  */

/**
  * Sets the state of the softPWM module
  * 
  * @param lwHandle littleWire device pointer
  * @param state State of the softPWM module ( \b ENABLE or \b DISABLE )
  * @return (none)
  */
void softPWM_state(littleWire* lwHandle,unsigned char state);

/**
  * Updates the values of softPWM modules
  *
  * @param lwHandle littleWire device pointer
  * @param ch1 Value of channel 1 - \b PIN4
  * @param ch2 Value of channel 2 - \b PIN1
  * @param ch3 Value of channel 3 - \b PIN2
  * @return (none)
  */
void softPWM_write(littleWire* lwHandle,unsigned char ch1,unsigned char ch2,unsigned char ch3);

/*! @} */

/*! @} */ 

/*! \addtogroup WS2812
  *  @brief WS2812 programmable RGB-LED support
  *  @{
  */

  /**
  * Writes to a WS2812 RGB-LED. This function writes the passed rgb values to a WS2812 led string
  * connected to the given pin. Use this if you want to control a single LED.
  *
  * If RGB values were preloaded with the preload call, the values passed in this call are added
  * to the buffer and the entire buffer is written to the LED string. This feature can be used to
  * reduce the number of USB transmissions for a string.
  *
  * @param lwHandle littleWire device pointer
  * @param r value of the red channel
  * @param g value of the green channel
  * @param b value of the blue channel
  * @param pin Pin name (\b PIN1 , \b PIN2 , \b PIN3 or \b PIN4 )
  * @return (none)
  */
void ws2812_write(littleWire* lwHandle, unsigned char pin,unsigned char r,unsigned char g,unsigned char b);

  /**
  * This function flushes the contents of the littlewire internal RGB buffer to the LED string.
  *
  * @param lwHandle littleWire device pointer
  * @param r value of the red channel
  * @param g value of the green channel
  * @param b value of the blue channel
  * @return (none)
  */
void ws2812_flush(littleWire* lwHandle, unsigned char pin);

  /**
  * Preloads a RGB value to the internal buffer. Up to 64 values can be preloaded. Further writes will be ignored
  *
  * @param lwHandle littleWire device pointer
  * @param r value of the red channel
  * @param g value of the green channel
  * @param b value of the blue channel
  * @return (none)
  */
void ws2812_preload(littleWire* lwHandle, unsigned char r,unsigned char g,unsigned char b);

  /*! @} */
  

/**
* @mainpage Introduction

\htmlinclude intro.html

*/

#endif
