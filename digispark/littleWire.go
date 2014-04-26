package gobotDigispark

//#cgo LDFLAGS: -lusb
//#include "littleWire.h"
//typedef usb_dev_handle littleWire;
import "C"
import "fmt"

type LittleWire struct {
	lwHandle *C.littleWire
}

//int littlewire_search();
//func LittleWireSearch() interface{} {
//  return int(C.littlewire_search())
//}

//littleWire* littlewire_connect_byID(int desiredID);
//func littleWireConnectByID(desiredID int) *littleWire {
//  littleWire = new(LittleWire)
//  C.littlewire_connect_byID(desiredID)
//}

//littleWire* littlewire_connect_bySerialNum(int mySerial);
//func littleWireConnectBySerialNum(mySerial int) *littleWire {
//  return C.littlewire_connect_bySerialNum(mySerial)
//}

//littleWire* littleWire_connect();
func LittleWireConnect() *LittleWire {
	littleWire := new(LittleWire)
	littleWire.lwHandle = C.littleWire_connect()
	return littleWire
}

//unsigned char readFirmwareVersion(littleWire* lwHandle);
func (littleWire *LittleWire) ReadFirmwareVersion() string {
	version := uint8(C.readFirmwareVersion(littleWire.lwHandle))
	return fmt.Sprintf("%v.%v", version&0xF0>>4, version&0x0F)
}

//void changeSerialNumber(littleWire* lwHandle,int serialNumber);
func (littleWire *LittleWire) ChangeSerialNumber(serialNumber int) {
	C.changeSerialNumber(littleWire.lwHandle, C.int(serialNumber))
}

//int customMessage(littleWire* lwHandle,unsigned char* receiveBuffer,unsigned char command,unsigned char d1,unsigned char d2, unsigned char d3, unsigned char d4);
//func (littleWire *LittleWire) CustomMessage(receiveBuffer *[]uint8, command uint8, d1 uint8, d2 uint8, d3 uint8, d4 uint8) int {
//  return int(C.customMessage(littleWire.lwHandle, receiveBuffer, command, d1, d2, d3, d4))
//}

//int littleWire_error ();
func (littleWire *LittleWire) LittleWireError() int {
	return int(C.littleWire_error())
}

//char *littleWire_errorName ();
//func LittleWireErrorName() string{
//  return string(C.littleWire_errorName())
//}

//void digitalWrite(littleWire* lwHandle, unsigned char pin, unsigned char state);
func (littleWire *LittleWire) DigitalWrite(pin uint8, state uint8) {
	C.digitalWrite(littleWire.lwHandle, C.uchar(pin), C.uchar(state))
}

//void pinMode(littleWire* lwHandle, unsigned char pin, unsigned char mode);
func (littleWire *LittleWire) PinMode(pin uint8, mode uint8) {
	C.pinMode(littleWire.lwHandle, C.uchar(pin), C.uchar(mode))
}

//unsigned char digitalRead(littleWire* lwHandle, unsigned char pin);
func (littleWire *LittleWire) DigitalRead(pin uint8) uint8 {
	return uint8(C.digitalRead(littleWire.lwHandle, C.uchar(pin)))
}

//void internalPullup(littleWire* lwHandle, unsigned char pin, unsigned char state);
func (littleWire *LittleWire) InternalPullup(pin uint8, state uint8) {
	C.internalPullup(littleWire.lwHandle, C.uchar(pin), C.uchar(state))
}

//void analog_init(littleWire* lwHandle, unsigned char voltageRef);
func (littleWire *LittleWire) AnalogInit(voltageRef uint8) {
	C.analog_init(littleWire.lwHandle, C.uchar(voltageRef))
}

//unsigned int analogRead(littleWire* lwHandle, unsigned char channel);
func (littleWire *LittleWire) AnalogRead(channel uint8) uint {
	return uint(C.analogRead(littleWire.lwHandle, C.uchar(channel)))
}

//void pwm_init(littleWire* lwHandle);
func (littleWire *LittleWire) PwmInit() {
	C.pwm_init(littleWire.lwHandle)
}

//void pwm_stop(littleWire* lwHandle);
func (littleWire *LittleWire) PwmStop() {
	C.pwm_stop(littleWire.lwHandle)
}

//void pwm_updateCompare(littleWire* lwHandle, unsigned char channelA, unsigned char channelB);
func (littleWire *LittleWire) PwmUpdateCompare(channelA uint8, channelB uint8) {
	C.pwm_updateCompare(littleWire.lwHandle, C.uchar(channelA), C.uchar(channelB))
}

//void pwm_updatePrescaler(littleWire* lwHandle, unsigned int value);
func (littleWire *LittleWire) PwmUpdatePrescaler(value uint) {
	C.pwm_updatePrescaler(littleWire.lwHandle, C.uint(value))
}

//void spi_init(littleWire* lwHandle);
func (littleWire *LittleWire) SpiInit() {
	C.spi_init(littleWire.lwHandle)
}

//void spi_sendMessage(littleWire* lwHandle, unsigned char * sendBuffer, unsigned char * inputBuffer, unsigned char length ,unsigned char mode);
//func (littleWire *LittleWire) SpiSendMessage(sendBuffer *[]uint8, inputBuffer *[]uint8, length uint8, mode uint8) {
//  C.spi_sendMessage(littleWire.lwHandle, sendBuffer, inputBuffer, length, mode)
//}

//unsigned char debugSpi(littleWire* lwHandle, unsigned char message);
func (littleWire *LittleWire) DebugSpi(message uint8) uint8 {
	return uint8(C.debugSpi(littleWire.lwHandle, C.uchar(message)))
}

//void spi_updateDelay(littleWire* lwHandle, unsigned int duration);
func (littleWire *LittleWire) SpiUpdateDelay(duration uint) {
	C.spi_updateDelay(littleWire.lwHandle, C.uint(duration))
}

//void i2c_init(littleWire* lwHandle);
func (littleWire *LittleWire) I2cInit() {
	C.i2c_init(littleWire.lwHandle)
}

//unsigned char i2c_start(littleWire* lwHandle, unsigned char address7bit, unsigned char direction);
func (littleWire *LittleWire) I2cStart(address7bit uint8, direction uint8) uint8 {
	return uint8(C.i2c_start(littleWire.lwHandle, C.uchar(address7bit), C.uchar(direction)))
}

//void i2c_write(littleWire* lwHandle, unsigned char* sendBuffer, unsigned char length, unsigned char endWithStop);
//func (littleWire *LittleWire) I2cWrite(sendBuffer *[]uint8, length uint8, endWithStop uint8) {
//  C.i2c_write(littleWire.lwHandle, sendBuffer, length, endWithStop)
//}

//void i2c_read(littleWire* lwHandle, unsigned char* readBuffer, unsigned char length, unsigned char endWithStop);
//func (littleWire *LittleWire) I2cRead(readBuffer *[]uint8, length uint8, endWithStop uint8) {
//  C.i2c_read(littleWire.lwHandle, readBuffer, length, endWithStop)
//}

//void i2c_updateDelay(littleWire* lwHandle, unsigned int duration);
func (littleWire *LittleWire) I2cUpdateDelay(duration uint) {
	C.i2c_updateDelay(littleWire.lwHandle, C.uint(duration))
}

//void onewire_sendBit(littleWire* lwHandle, unsigned char bitValue);
func (littleWire *LittleWire) OneWireSendBit(bitValue uint8) {
	C.onewire_sendBit(littleWire.lwHandle, C.uchar(bitValue))
}

//void onewire_writeByte(littleWire* lwHandle, unsigned char messageToSend);
func (littleWire *LittleWire) OneWireWriteByte(messageToSend uint8) {
	C.onewire_writeByte(littleWire.lwHandle, C.uchar(messageToSend))
}

//unsigned char onewire_readByte(littleWire* lwHandle);
func (littleWire *LittleWire) OneWireReadByte() uint8 {
	return uint8(C.onewire_readByte(littleWire.lwHandle))
}

//unsigned char onewire_readBit(littleWire* lwHandle);
func (littleWire *LittleWire) OneWireReadBit() uint8 {
	return uint8(C.onewire_readBit(littleWire.lwHandle))
}

//unsigned char onewire_resetPulse(littleWire* lwHandle);
func (littleWire *LittleWire) OneWireResetPulse() uint8 {
	return uint8(C.onewire_resetPulse(littleWire.lwHandle))
}

//int onewire_firstAddress(littleWire* lwHandle);
func (littleWire *LittleWire) OneWireFirstAddress() int {
	return int(C.onewire_firstAddress(littleWire.lwHandle))
}

//int onewire_nextAddress(littleWire* lwHandle);
func (littleWire *LittleWire) OneWireNextAddress() int {
	return int(C.onewire_nextAddress(littleWire.lwHandle))
}

//void softPWM_state(littleWire* lwHandle,unsigned char state);
func (littleWire *LittleWire) SoftPWMState(state uint8) {
	C.softPWM_state(littleWire.lwHandle, C.uchar(state))
}

//void softPWM_write(littleWire* lwHandle,unsigned char ch1,unsigned char ch2,unsigned char ch3);
func (littleWire *LittleWire) SoftPWMWrite(ch1 uint8, ch2 uint8, ch3 uint8) {
	C.softPWM_write(littleWire.lwHandle, C.uchar(ch1), C.uchar(ch2), C.uchar(ch3))
}

//void ws2812_write(littleWire* lwHandle, unsigned char pin,unsigned char r,unsigned char g,unsigned char b);
func (littleWire *LittleWire) Ws2812Write(pin uint8, r uint8, g uint8, b uint8) {
	C.ws2812_write(littleWire.lwHandle, C.uchar(pin), C.uchar(r), C.uchar(g), C.uchar(b))
}

//void ws2812_flush(littleWire* lwHandle, unsigned char pin);
func (littleWire *LittleWire) Ws2812Flush(pin uint8) {
	C.ws2812_flush(littleWire.lwHandle, C.uchar(pin))
}

//void ws2812_preload(littleWire* lwHandle, unsigned char r,unsigned char g,unsigned char b);
func (littleWire *LittleWire) Ws2812Preload(r uint8, g uint8, b uint8) {
	C.ws2812_preload(littleWire.lwHandle, C.uchar(r), C.uchar(g), C.uchar(b))
}
