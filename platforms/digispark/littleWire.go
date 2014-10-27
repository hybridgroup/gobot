package digispark

//#cgo LDFLAGS: -lusb
//#include "littleWire.h"
//typedef usb_dev_handle littleWire;
import "C"
import "fmt"

type LittleWire struct {
	lwHandle *C.littleWire
}

func LittleWireConnect() *LittleWire {
	return &LittleWire{
		lwHandle: C.littleWire_connect(),
	}
}

func (l *LittleWire) ReadFirmwareVersion() string {
	version := uint8(C.readFirmwareVersion(l.lwHandle))
	return fmt.Sprintf("%v.%v", version&0xF0>>4, version&0x0F)
}

func (l *LittleWire) ChangeSerialNumber(serialNumber int) {
	C.changeSerialNumber(l.lwHandle, C.int(serialNumber))
}

func (l *LittleWire) LittleWireError() int {
	return int(C.littleWire_error())
}

func (l *LittleWire) DigitalWrite(pin uint8, state uint8) {
	C.digitalWrite(l.lwHandle, C.uchar(pin), C.uchar(state))
}

func (l *LittleWire) PinMode(pin uint8, mode uint8) {
	C.pinMode(l.lwHandle, C.uchar(pin), C.uchar(mode))
}

func (l *LittleWire) DigitalRead(pin uint8) uint8 {
	return uint8(C.digitalRead(l.lwHandle, C.uchar(pin)))
}

func (l *LittleWire) InternalPullup(pin uint8, state uint8) {
	C.internalPullup(l.lwHandle, C.uchar(pin), C.uchar(state))
}

func (l *LittleWire) AnalogInit(voltageRef uint8) {
	C.analog_init(l.lwHandle, C.uchar(voltageRef))
}

func (l *LittleWire) AnalogRead(channel uint8) uint {
	return uint(C.analogRead(l.lwHandle, C.uchar(channel)))
}

func (l *LittleWire) PwmInit() {
	C.pwm_init(l.lwHandle)
}

func (l *LittleWire) PwmStop() {
	C.pwm_stop(l.lwHandle)
}

func (l *LittleWire) PwmUpdateCompare(channelA uint8, channelB uint8) {
	C.pwm_updateCompare(l.lwHandle, C.uchar(channelA), C.uchar(channelB))
}

func (l *LittleWire) PwmUpdatePrescaler(value uint) {
	C.pwm_updatePrescaler(l.lwHandle, C.uint(value))
}

func (l *LittleWire) SpiInit() {
	C.spi_init(l.lwHandle)
}

func (l *LittleWire) DebugSpi(message uint8) uint8 {
	return uint8(C.debugSpi(l.lwHandle, C.uchar(message)))
}

func (l *LittleWire) SpiUpdateDelay(duration uint) {
	C.spi_updateDelay(l.lwHandle, C.uint(duration))
}

func (l *LittleWire) I2cInit() {
	C.i2c_init(l.lwHandle)
}

func (l *LittleWire) I2cStart(address7bit uint8, direction uint8) uint8 {
	return uint8(C.i2c_start(l.lwHandle, C.uchar(address7bit), C.uchar(direction)))
}

func (l *LittleWire) I2cUpdateDelay(duration uint) {
	C.i2c_updateDelay(l.lwHandle, C.uint(duration))
}

func (l *LittleWire) OneWireSendBit(bitValue uint8) {
	C.onewire_sendBit(l.lwHandle, C.uchar(bitValue))
}

func (l *LittleWire) OneWireWriteByte(messageToSend uint8) {
	C.onewire_writeByte(l.lwHandle, C.uchar(messageToSend))
}

func (l *LittleWire) OneWireReadByte() uint8 {
	return uint8(C.onewire_readByte(l.lwHandle))
}

func (l *LittleWire) OneWireReadBit() uint8 {
	return uint8(C.onewire_readBit(l.lwHandle))
}

func (l *LittleWire) OneWireResetPulse() uint8 {
	return uint8(C.onewire_resetPulse(l.lwHandle))
}

func (l *LittleWire) OneWireFirstAddress() int {
	return int(C.onewire_firstAddress(l.lwHandle))
}

func (l *LittleWire) OneWireNextAddress() int {
	return int(C.onewire_nextAddress(l.lwHandle))
}

func (l *LittleWire) SoftPWMState(state uint8) {
	C.softPWM_state(l.lwHandle, C.uchar(state))
}

func (l *LittleWire) SoftPWMWrite(ch1 uint8, ch2 uint8, ch3 uint8) {
	C.softPWM_write(l.lwHandle, C.uchar(ch1), C.uchar(ch2), C.uchar(ch3))
}

func (l *LittleWire) Ws2812Write(pin uint8, r uint8, g uint8, b uint8) {
	C.ws2812_write(l.lwHandle, C.uchar(pin), C.uchar(r), C.uchar(g), C.uchar(b))
}

func (l *LittleWire) Ws2812Flush(pin uint8) {
	C.ws2812_flush(l.lwHandle, C.uchar(pin))
}

func (l *LittleWire) Ws2812Preload(r uint8, g uint8, b uint8) {
	C.ws2812_preload(l.lwHandle, C.uchar(r), C.uchar(g), C.uchar(b))
}
