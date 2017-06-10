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

/******************************************************************************
* See the littleWire.h for the function descriptions/comments
/*****************************************************************************/
#include "littleWire.h"

unsigned char	crc8;
int		LastDiscrepancy;
int 	LastFamilyDiscrepancy;
int 	LastDeviceFlag;


unsigned char rxBuffer[RX_BUFFER_SIZE]; /* This has to be unsigned for the data's sake */
unsigned char ROM_NO[8];
int lwStatus;
lwCollection lwResults[16];
int lw_totalDevices;


/******************************************************************************
/ Taken from: http://www.maxim-ic.com/appnotes.cfm/appnote_number/187
/*****************************************************************************/
static unsigned char dscrc_table[] = {
        0, 94,188,226, 97, 63,221,131,194,156,126, 32,163,253, 31, 65,
      157,195, 33,127,252,162, 64, 30, 95,  1,227,189, 62, 96,130,220,
       35,125,159,193, 66, 28,254,160,225,191, 93,  3,128,222, 60, 98,
      190,224,  2, 92,223,129, 99, 61,124, 34,192,158, 29, 67,161,255,
       70, 24,250,164, 39,121,155,197,132,218, 56,102,229,187, 89,  7,
      219,133,103, 57,186,228,  6, 88, 25, 71,165,251,120, 38,196,154,
      101, 59,217,135,  4, 90,184,230,167,249, 27, 69,198,152,122, 36,
      248,166, 68, 26,153,199, 37,123, 58,100,134,216, 91,  5,231,185,
      140,210, 48,110,237,179, 81, 15, 78, 16,242,172, 47,113,147,205,
       17, 79,173,243,112, 46,204,146,211,141,111, 49,178,236, 14, 80,
      175,241, 19, 77,206,144,114, 44,109, 51,209,143, 12, 82,176,238,
       50,108,142,208, 83, 13,239,177,240,174, 76, 18,145,207, 45,115,
      202,148,118, 40,171,245, 23, 73,  8, 86,180,234,105, 55,213,139,
       87,  9,235,181, 54,104,138,212,149,203, 41,119,244,170, 72, 22,
      233,183, 85, 11,136,214, 52,106, 43,117,151,201, 74, 20,246,168,
      116, 42,200,150, 21, 75,169,247,182,232, 10, 84,215,137,107, 53};
/*****************************************************************************/

int littlewire_search()
{
  struct usb_bus *bus;
  struct usb_device *dev;

  usb_init();
  usb_find_busses();
  usb_find_devices();

  lw_totalDevices = 0;

  for (bus = usb_busses; bus; bus = bus->next)
  {
    for (dev = bus->devices; dev; dev = dev->next)
    {
      usb_dev_handle *udev;
      char description[256];
      char string[256];
      int ret, i;

      if((dev->descriptor.idVendor == VENDOR_ID) && (dev->descriptor.idProduct == PRODUCT_ID))
      {
        udev = usb_open(dev);
        if (udev)
        {
          if (dev->descriptor.iSerialNumber)
          {
            ret = usb_get_string_simple(udev, dev->descriptor.iSerialNumber, string, sizeof(string));
            if (ret > 0)
            {
              lwResults[lw_totalDevices].serialNumber = atoi(string);
              lwResults[lw_totalDevices].lw_device = dev;
            }
          }
          usb_close(udev);
          lw_totalDevices++;
        }
      }
    }
  }

  return lw_totalDevices;
}

littleWire* littlewire_connect_byID(int desiredID)
{
  littleWire  *tempHandle = NULL;

  if(desiredID > (lw_totalDevices-1))
  {
    return tempHandle;
  }

  tempHandle = usb_open(lwResults[desiredID].lw_device);

  return tempHandle;
}

littleWire* littlewire_connect_bySerialNum(int mySerial)
{
  littleWire  *tempHandle = NULL;
  int temp_id = 0xDEAF;
  int i;

  for(i=0;i<lw_totalDevices;i++)
  {
    if(lwResults[i].serialNumber == mySerial)
    {
      temp_id = i;
    }
  }

  tempHandle = littlewire_connect_byID(temp_id);
  return tempHandle;
}

littleWire* littleWire_connect()
{
	littleWire  *tempHandle = NULL;

	usb_init();
	usbOpenDevice(&tempHandle, VENDOR_ID, "*", PRODUCT_ID, "*", "*", NULL, NULL );

	return tempHandle;
}

unsigned char readFirmwareVersion(littleWire* lwHandle)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 34, 0, 0, rxBuffer, 8, USB_TIMEOUT);

	return rxBuffer[0];
}

void changeSerialNumber(littleWire* lwHandle,int serialNumber)
{
	char serBuf[4];

	if(serialNumber > 999)
	{
		serialNumber = 999;
	}
	else if(serialNumber < 100)
	{
		serialNumber = 100;
	}

	sprintf(serBuf,"%d",serialNumber);

	lwStatus=usb_control_msg(lwHandle, 0xC0, 55, (serBuf[1]<<8)|serBuf[0],serBuf[2], rxBuffer, 8, USB_TIMEOUT);
}

void digitalWrite(littleWire* lwHandle, unsigned char pin, unsigned char state)
{
	if(state){
		lwStatus=usb_control_msg(lwHandle, 0xC0, 18, pin, 0, rxBuffer, 8, USB_TIMEOUT);
	} else{
		lwStatus=usb_control_msg(lwHandle, 0xC0, 19, pin, 0, rxBuffer, 8, USB_TIMEOUT);
	}
}

void pinMode(littleWire* lwHandle, unsigned char pin, unsigned char mode)
{
	if(mode){
		lwStatus=usb_control_msg(lwHandle, 0xC0, 13, pin, 0, rxBuffer, 8, USB_TIMEOUT);
	} else {
		lwStatus=usb_control_msg(lwHandle, 0xC0, 14, pin, 0, rxBuffer, 8, USB_TIMEOUT);
	}
}

unsigned char digitalRead(littleWire* lwHandle, unsigned char pin)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 20, pin, 0, rxBuffer, 8, USB_TIMEOUT);

	return rxBuffer[0];
}

void internalPullup(littleWire* lwHandle, unsigned char pin, unsigned char state)
{
	if(state){
		lwStatus=usb_control_msg(lwHandle, 0xC0, 18, pin, 0, rxBuffer, 8, USB_TIMEOUT);
	} else{
		lwStatus=usb_control_msg(lwHandle, 0xC0, 19, pin, 0, rxBuffer, 8, USB_TIMEOUT);
	}
}

void analog_init(littleWire* lwHandle, unsigned char voltageRef)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 35, (voltageRef<<8) | 0x07, 0, rxBuffer, 8, USB_TIMEOUT);
}

unsigned int analogRead(littleWire* lwHandle, unsigned char channel)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 15, channel, 0, rxBuffer, 8, USB_TIMEOUT);

	return ((rxBuffer[1] *256) + (rxBuffer[0]));
}

void pwm_init(littleWire* lwHandle)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 16, 0, 0, rxBuffer, 8, USB_TIMEOUT);
}

void pwm_stop(littleWire* lwHandle)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 32, 0, 0, rxBuffer, 8, USB_TIMEOUT);
}

void pwm_updateCompare(littleWire* lwHandle, unsigned char channelA, unsigned char channelB)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 17, channelA, channelB, rxBuffer, 8, USB_TIMEOUT);
}

void pwm_updatePrescaler(littleWire* lwHandle, unsigned int value)
{
	switch(value)
	{
		case 1024:
			lwStatus=usb_control_msg(lwHandle, 0xC0, 22, 4, 0, rxBuffer, 8, USB_TIMEOUT);
		break;
		case 256:
			lwStatus=usb_control_msg(lwHandle, 0xC0, 22, 3, 0, rxBuffer, 8, USB_TIMEOUT);
		break;
		case 64:
			lwStatus=usb_control_msg(lwHandle, 0xC0, 22, 2, 0, rxBuffer, 8, USB_TIMEOUT);
		break;
		case 8:
			lwStatus=usb_control_msg(lwHandle, 0xC0, 22, 1, 0, rxBuffer, 8, USB_TIMEOUT);
		break;
		case 1:
			lwStatus=usb_control_msg(lwHandle, 0xC0, 22, 0, 0, rxBuffer, 8, USB_TIMEOUT);
		break;
	}
}

void spi_init(littleWire* lwHandle)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 23, 0, 0, rxBuffer, 8, USB_TIMEOUT);
}

void spi_sendMessage(littleWire* lwHandle, unsigned char * sendBuffer, unsigned char * inputBuffer, unsigned char length ,unsigned char mode)
{
	int i=0;
	if(length>4)
		length=4;
	lwStatus=usb_control_msg(lwHandle, 0xC0, (0xF0 + length + (mode<<3) ), (sendBuffer[1]<<8) + sendBuffer[0] , (sendBuffer[3]<<8) + sendBuffer[2], rxBuffer, 8, USB_TIMEOUT);
	lwStatus=usb_control_msg(lwHandle, 0xC0, 40, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	for(i=0;i<length;i++)
		inputBuffer[i]=rxBuffer[i];
}

unsigned char debugSpi(littleWire* lwHandle, unsigned char message)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 33, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	lwStatus=usb_control_msg(lwHandle, 0xC0, 40, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	return rxBuffer[0];
}

void spi_updateDelay(littleWire* lwHandle, unsigned int duration)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 31, duration, 0, rxBuffer, 8, USB_TIMEOUT);
}

void i2c_init(littleWire* lwHandle)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 44, 0, 0, rxBuffer, 8, USB_TIMEOUT);
}

unsigned char i2c_start(littleWire* lwHandle, unsigned char address7bit, unsigned char direction)
{
	unsigned char temp;

	temp = (address7bit << 1) | direction;

	lwStatus=usb_control_msg(lwHandle, 0xC0, 45, temp, 0, rxBuffer, 8, USB_TIMEOUT);
	lwStatus=usb_control_msg(lwHandle, 0xC0, 40, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	return !rxBuffer[0];
}

void i2c_write(littleWire* lwHandle, unsigned char* sendBuffer, unsigned char length, unsigned char endWithStop)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, (0xE0 + length + (endWithStop<<3) ), (sendBuffer[1]<<8) + sendBuffer[0] , (sendBuffer[3]<<8) + sendBuffer[2], rxBuffer, 8, USB_TIMEOUT);
}

void i2c_read(littleWire* lwHandle, unsigned char* readBuffer, unsigned char length, unsigned char endWithStop)
{
	int i=0;

	if(endWithStop)
		lwStatus=usb_control_msg(lwHandle, 0xC0, 46, (length<<8) + 1, 1, rxBuffer, 8, USB_TIMEOUT);
	else
		lwStatus=usb_control_msg(lwHandle, 0xC0, 46, (length<<8) + 0, 0, rxBuffer, 8, USB_TIMEOUT);

	delay(3);

  lwStatus=usb_control_msg(lwHandle, 0xC0, 40, 0, 0, rxBuffer, 8, USB_TIMEOUT);

	for(i=0;i<length;i++)
		readBuffer[i]=rxBuffer[i];
}

void i2c_updateDelay(littleWire* lwHandle, unsigned int duration)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 49, duration, 0, rxBuffer, 8, USB_TIMEOUT);
}

void onewire_sendBit(littleWire* lwHandle, unsigned char bitValue)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 51, bitValue, 0, rxBuffer, 8, USB_TIMEOUT);
}

void onewire_writeByte(littleWire* lwHandle, unsigned char messageToSend)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 42, messageToSend, 0, rxBuffer, 8, USB_TIMEOUT);
	delay(3);
}

unsigned char onewire_readByte(littleWire* lwHandle)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 43, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	delay(3);
	lwStatus=usb_control_msg(lwHandle, 0xC0, 40, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	return rxBuffer[0];
}

unsigned char onewire_readBit(littleWire* lwHandle)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 50, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	lwStatus=usb_control_msg(lwHandle, 0xC0, 40, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	return rxBuffer[0];
}

unsigned char onewire_resetPulse(littleWire* lwHandle)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 41, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	delay(3);
	lwStatus=usb_control_msg(lwHandle, 0xC0, 40, 0, 0, rxBuffer, 8, USB_TIMEOUT);
	return rxBuffer[0];
}

void softPWM_state(littleWire* lwHandle,unsigned char state)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 47, state, 0, rxBuffer, 8, USB_TIMEOUT);
}

void softPWM_write(littleWire* lwHandle,unsigned char ch1,unsigned char ch2,unsigned char ch3)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 48, (ch2<<8) | ch1, ch3, rxBuffer, 8, USB_TIMEOUT);
}

void ws2812_write(littleWire* lwHandle, unsigned char pin, unsigned char r,unsigned char g,unsigned char b)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 54, (g<<8) | pin | 0x30, (b<<8) | r, rxBuffer, 8, USB_TIMEOUT);
}

void ws2812_flush(littleWire* lwHandle, unsigned char pin)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 54, pin | 0x10, 0, rxBuffer, 8, USB_TIMEOUT);
}

void ws2812_preload(littleWire* lwHandle, unsigned char r,unsigned char g,unsigned char b)
{
	lwStatus=usb_control_msg(lwHandle, 0xC0, 54, (g<<8) | 0x20, (b<<8) | r, rxBuffer, 8, USB_TIMEOUT);
}

int customMessage(littleWire* lwHandle,unsigned char* receiveBuffer,unsigned char command,unsigned char d1,unsigned char d2, unsigned char d3, unsigned char d4)
{
	int i;
	int rc;
	rc = lwStatus=usb_control_msg(lwHandle, 0xC0, command, (d2<<8)|d1, (d4<<8)|d3, rxBuffer, 8, USB_TIMEOUT);
	for(i=0;i<8;i++)
		receiveBuffer[i]=rxBuffer[i];
	return rc;
}

/******************************************************************************
* Do the crc8 calculation
* Taken from: http://www.maxim-ic.com/appnotes.cfm/appnote_number/187
/*****************************************************************************/
unsigned char docrc8(unsigned char value)
{
   // See Maxim Application Note 27

   crc8 = dscrc_table[crc8 ^ value];

   return crc8;
}

int onewire_nextAddress(littleWire* lwHandle)
{
   int id_bit_number;
   int last_zero, rom_byte_number, search_result;
   int id_bit, cmp_id_bit;
   unsigned char rom_byte_mask, search_direction;

   // initialize for search
   id_bit_number = 1;
   last_zero = 0;
   rom_byte_number = 0;
   rom_byte_mask = 1;
   search_result = 0;
   crc8 = 0;

   // if the last call was not the last one
   if (!LastDeviceFlag)
   {
      // 1-Wire reset
      if (!onewire_resetPulse(lwHandle))
      {
         // reset the search
         LastDiscrepancy = 0;
         LastDeviceFlag = 0;
         LastFamilyDiscrepancy = 0;
         return 0;
      }

      // issue the search command
      onewire_writeByte(lwHandle,0xF0);

      // loop to do the search
      do
      {
         // read a bit and its complement
         id_bit = onewire_readBit(lwHandle);
         cmp_id_bit = onewire_readBit(lwHandle);

         // check for no devices on 1-wire
         if ((id_bit == 1) && (cmp_id_bit == 1))
            break;
         else
         {
            // all devices coupled have 0 or 1
            if (id_bit != cmp_id_bit)
               search_direction = id_bit;  // bit write value for search
            else
            {
               // if this discrepancy if before the Last Discrepancy
               // on a previous next then pick the same as last time
               if (id_bit_number < LastDiscrepancy)
                  search_direction = ((ROM_NO[rom_byte_number] & rom_byte_mask) > 0);
               else
                  // if equal to last pick 1, if not then pick 0
                  search_direction = (id_bit_number == LastDiscrepancy);

               // if 0 was picked then record its position in LastZero
               if (search_direction == 0)
               {
                  last_zero = id_bit_number;

                  // check for Last discrepancy in family
                  if (last_zero < 9)
                     LastFamilyDiscrepancy = last_zero;
               }
            }

            // set or clear the bit in the ROM byte rom_byte_number
            // with mask rom_byte_mask
            if (search_direction == 1)
              ROM_NO[rom_byte_number] |= rom_byte_mask;
            else
              ROM_NO[rom_byte_number] &= ~rom_byte_mask;

            // serial number search direction write bit
            onewire_sendBit(lwHandle,search_direction);

            // increment the byte counter id_bit_number
            // and shift the mask rom_byte_mask
            id_bit_number++;
            rom_byte_mask <<= 1;

            // if the mask is 0 then go to new SerialNum byte rom_byte_number and reset mask
            if (rom_byte_mask == 0)
            {
                docrc8(ROM_NO[rom_byte_number]);  // accumulate the CRC
                rom_byte_number++;
                rom_byte_mask = 1;
            }
         }
      }
      while(rom_byte_number < 8);  // loop until through all ROM bytes 0-7

      // if the search was successful then
      if (!((id_bit_number < 65) || (crc8 != 0)))
      {
         // search successful so set LastDiscrepancy,LastDeviceFlag,search_result
         LastDiscrepancy = last_zero;

         // check for last device
         if (LastDiscrepancy == 0)
            LastDeviceFlag = 1;

         search_result = 1;
      }
   }

   // if no device found then reset counters so next 'search' will be like a first
   if (!search_result || !ROM_NO[0])
   {
      LastDiscrepancy = 0;
      LastDeviceFlag = 0;
      LastFamilyDiscrepancy = 0;
      search_result = 0;
   }

   return search_result;
}

int onewire_firstAddress(littleWire* lwHandle)
{
	littleWire* temp = lwHandle;

   // reset the search state
   LastDiscrepancy = 0;
   LastDeviceFlag = 0;
   LastFamilyDiscrepancy = 0;

   return onewire_nextAddress(temp);
}

int littleWire_error () {
        if (lwStatus<0) return lwStatus;
        else return 0;
}

char *littleWire_errorName () {
        if (lwStatus<0) switch (lwStatus) {
                case -1: return "I/O Error"; break;
                case -2: return "Invalid paramenter"; break;
                case -3: return "Access error"; break;
                case -4: return "No device"; break;
                case -5: return "Not found"; break;
                case -6: return "Busy"; break;
                case -7: return "Timeout"; break;
                case -8: return "Overflow"; break;
                case -9: return "Pipe"; break;
                case -10: return "Interrupted"; break;
                case -11: return "No memory"; break;
                case -12: return "Not supported"; break;
                case -99: return "Other"; break;
                default: return "unknown";
        }
        else return 0;
}
