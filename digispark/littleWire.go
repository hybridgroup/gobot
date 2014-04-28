package digispark

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
	return &LittleWire{
		lwHandle: C.littleWire_connect(),
	}
}

//unsigned char readFirmwareVersion(littleWire* lwHandle);
func (l *LittleWire) ReadFirmwareVersion() string {
	version := uint8(C.readFirmwareVersion(l.lwHandle))
	return fmt.Sprintf("%v.%v", version&0xF0>>4, version&0x0F)
}

//void changeSerialNumber(littleWire* lwHandle,int serialNumber);
func (l *LittleWire) ChangeSerialNumber(serialNumber int) {
	C.changeSerialNumber(l.lwHandle, C.int(serialNumber))
}

//int customMessage(littleWire* lwHandle,unsigned char* receiveBuffer,unsigned char command,unsigned char d1,unsigned char d2, unsigned char d3, unsigned char d4);
//func (littleWire *LittleWire) CustomMessage(receiveBuffer *[]uint8, command uint8, d1 uint8, d2 uint8, d3 uint8, d4 uint8) int {
//  return int(C.customMessage(littleWire.lwHandle, receiveBuffer, command, d1, d2, d3, d4))
//}

//int littleWire_error ();
func (l *LittleWire) LittleWireError() int {
	return int(C.littleWire_error())
}

//char *littleWire_errorName ();
//func LittleWireErrorName() string{
//  return string(C.littleWire_errorName())
//}

//void digitalWrite(littleWire* lwHandle, unsigned char pin, unsigned char state);
func (l *LittleWire) DigitalWrite(pin uint8, state uint8) {
	C.digitalWrite(l.lwHandle, C.uchar(pin), C.uchar(state))
}

//void pinMode(littleWire* lwHandle, unsigned char pin, unsigned char mode);
func (l *LittleWire) PinMode(pin uint8, mode uint8) {
	C.pinMode(l.lwHandle, C.uchar(pin), C.uchar(mode))
}

//unsigned char digitalRead(littleWire* lwHandle, unsigned char pin);
func (l *LittleWire) DigitalRead(pin uint8) uint8 {
	return uint8(C.digitalRead(l.lwHandle, C.uchar(pin)))
}

//void internalPullup(littleWire* lwHandle, unsigned char pin, unsigned char state);
func (l *LittleWire) InternalPullup(pin uint8, state uint8) {
	C.internalPullup(l.lwHandle, C.uchar(pin), C.uchar(state))
}

//void analog_init(littleWire* lwHandle, unsigned char voltageRef);
func (l *LittleWire) AnalogInit(voltageRef uint8) {
	C.analog_init(l.lwHandle, C.uchar(voltageRef))
}

//unsigned int analogRead(littleWire* lwHandle, unsigned char channel);
func (l *LittleWire) AnalogRead(channel uint8) uint {
	return uint(C.analogRead(l.lwHandle, C.uchar(channel)))
}

//void pwm_init(littleWire* lwHandle);
func (l *LittleWire) PwmInit() {
	C.pwm_init(l.lwHandle)
}

//void pwm_stop(littleWire* lwHandle);
func (l *LittleWire) PwmStop() {
	C.pwm_stop(l.lwHandle)
}

//void pwm_updateCompare(littleWire* lwHandle, unsigned char channelA, unsigned char channelB);
func (l *LittleWire) PwmUpdateCompare(channelA uint8, channelB uint8) {
	C.pwm_updateCompare(l.lwHandle, C.uchar(channelA), C.uchar(channelB))
}

//void pwm_updatePrescaler(littleWire* lwHandle, unsigned int value);
func (l *LittleWire) PwmUpdatePrescaler(value uint) {
	C.pwm_updatePrescaler(l.lwHandle, C.uint(value))
}

//void spi_init(littleWire* lwHandle);
func (l *LittleWire) SpiInit() {
	C.spi_init(l.lwHandle)
}

//void spi_sendMessage(littleWire* lwHandle, unsigned char * sendBuffer, unsigned char * inputBuffer, unsigned char length ,unsigned char mode);
//func (littleWire *LittleWire) SpiSendMessage(sendBuffer *[]uint8, inputBuffer *[]uint8, length uint8, mode uint8) {
//  C.spi_sendMessage(littleWire.lwHandle, sendBuffer, inputBuffer, length, mode)
//}

//unsigned char debugSpi(littleWire* lwHandle, unsigned char message);
func (l *LittleWire) DebugSpi(message uint8) uint8 {
	return uint8(C.debugSpi(l.lwHandle, C.uchar(message)))
}

//void spi_updateDelay(littleWire* lwHandle, unsigned int duration);
func (l *LittleWire) SpiUpdateDelay(duration uint) {
	C.spi_updateDelay(l.lwHandle, C.uint(duration))
}

//void i2c_init(littleWire* lwHandle);
func (l *LittleWire) I2cInit() {
	C.i2c_init(l.lwHandle)
}

//unsigned char i2c_start(littleWire* lwHandle, unsigned char address7bit, unsigned char direction);
func (l *LittleWire) I2cStart(address7bit uint8, direction uint8) uint8 {
	return uint8(C.i2c_start(l.lwHandle, C.uchar(address7bit), C.uchar(direction)))
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
func (l *LittleWire) I2cUpdateDelay(duration uint) {
	C.i2c_updateDelay(l.lwHandle, C.uint(duration))
}

//void onewire_sendBit(littleWire* lwHandle, unsigned char bitValue);
func (l *LittleWire) OneWireSendBit(bitValue uint8) {
	C.onewire_sendBit(l.lwHandle, C.uchar(bitValue))
}

//void onewire_writeByte(littleWire* lwHandle, unsigned char messageToSend);
func (l *LittleWire) OneWireWriteByte(messageToSend uint8) {
	C.onewire_writeByte(l.lwHandle, C.uchar(messageToSend))
}

//unsigned char onewire_readByte(littleWire* lwHandle);
func (l *LittleWire) OneWireReadByte() uint8 {
	return uint8(C.onewire_readByte(l.lwHandle))
}

//unsigned char onewire_readBit(littleWire* lwHandle);
func (l *LittleWire) OneWireReadBit() uint8 {
	return uint8(C.onewire_readBit(l.lwHandle))
}

//unsigned char onewire_resetPulse(littleWire* lwHandle);
func (l *LittleWire) OneWireResetPulse() uint8 {
	return uint8(C.onewire_resetPulse(l.lwHandle))
}

//int onewire_firstAddress(littleWire* lwHandle);
func (l *LittleWire) OneWireFirstAddress() int {
	return int(C.onewire_firstAddress(l.lwHandle))
}

//int onewire_nextAddress(littleWire* lwHandle);
func (l *LittleWire) OneWireNextAddress() int {
	return int(C.onewire_nextAddress(l.lwHandle))
}

//void softPWM_state(littleWire* lwHandle,unsigned char state);
func (l *LittleWire) SoftPWMState(state uint8) {
	C.softPWM_state(l.lwHandle, C.uchar(state))
}

//void softPWM_write(littleWire* lwHandle,unsigned char ch1,unsigned char ch2,unsigned char ch3);
func (l *LittleWire) SoftPWMWrite(ch1 uint8, ch2 uint8, ch3 uint8) {
	C.softPWM_write(l.lwHandle, C.uchar(ch1), C.uchar(ch2), C.uchar(ch3))
}

//void ws2812_write(littleWire* lwHandle, unsigned char pin,unsigned char r,unsigned char g,unsigned char b);
func (l *LittleWire) Ws2812Write(pin uint8, r uint8, g uint8, b uint8) {
	C.ws2812_write(l.lwHandle, C.uchar(pin), C.uchar(r), C.uchar(g), C.uchar(b))
}

//void ws2812_flush(littleWire* lwHandle, unsigned char pin);
func (l *LittleWire) Ws2812Flush(pin uint8) {
	C.ws2812_flush(l.lwHandle, C.uchar(pin))
}

//void ws2812_preload(littleWire* lwHandle, unsigned char r,unsigned char g,unsigned char b);
func (l *LittleWire) Ws2812Preload(r uint8, g uint8, b uint8) {
	C.ws2812_preload(l.lwHandle, C.uchar(r), C.uchar(g), C.uchar(b))
}
