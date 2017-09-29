// Based on https://github.com/DexterInd/GoPiGo3/blob/master/Software/Python/gopigo3.py
// You will need to run the following script if using a stock raspbian image before this library will work:
// https://www.dexterindustries.com/GoPiGo/get-started-with-the-gopigo3-raspberry-pi-robot/3-program-your-raspberry-pi-robot/python-programming-language/
package gopigo3

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
)

// spi address for gopigo3
const goPiGo3Address = 0x08

// register addresses for gopigo3
const (
	NONE byte = iota
	GET_MANUFACTURER
	GET_NAME
	GET_HARDWARE_VERSION
	GET_FIRMWARE_VERSION
	GET_ID
	SET_LED
	GET_VOLTAGE_5V
	GET_VOLTAGE_VCC
	SET_SERVO
	SET_MOTOR_PWM
	SET_MOTOR_POSITION
	SET_MOTOR_POSITION_KP
	SET_MOTOR_POSITION_KD
	SET_MOTOR_DPS
	SET_MOTOR_LIMITS
	OFFSET_MOTOR_ENCODER
	GET_MOTOR_ENCODER_LEFT
	GET_MOTOR_ENCODER_RIGHT
	GET_MOTOR_STATUS_LEFT
	GET_MOTOR_STATUS_RIGHT
	SET_GROVE_TYPE
	SET_GROVE_MODE
	SET_GROVE_STATE
	SET_GROVE_PWM_DUTY
	SET_GROVE_PWM_FREQUENCY
	GET_GROVE_VALUE_1
	GET_GROVE_VALUE_2
	GET_GROVE_STATE_1_1
	GET_GROVE_STATE_1_2
	GET_GROVE_STATE_2_1
	GET_GROVE_STATE_2_2
	GET_GROVE_VOLTAGE_1_1
	GET_GROVE_VOLTAGE_1_2
	GET_GROVE_VOLTAGE_2_1
	GET_GROVE_VOLTAGE_2_2
	GET_GROVE_ANALOG_1_1
	GET_GROVE_ANALOG_1_2
	GET_GROVE_ANALOG_2_1
	GET_GROVE_ANALOG_2_2
	START_GROVE_I2C_1
	START_GROVE_I2C_2
)

const (
	WHEEL_BASE_WIDTH           = 117                                                       // distance (mm) from left wheel to right wheel. This works with the initial GPG3 prototype. Will need to be adjusted.
	WHEEL_DIAMETER             = 66.5                                                      // wheel diameter (mm)
	WHEEL_BASE_CIRCUMFERENCE   = WHEEL_BASE_WIDTH * math.Pi                                // circumference of the circle the wheels will trace while turning (mm)
	WHEEL_CIRCUMFERENCE        = WHEEL_DIAMETER * math.Pi                                  // circumference of the wheels (mm)
	MOTOR_GEAR_RATIO           = 120                                                       // motor gear ratio
	ENCODER_TICKS_PER_ROTATION = 6                                                         // encoder ticks per motor rotation (number of magnet positions)
	MOTOR_TICKS_PER_DEGREE     = ((MOTOR_GEAR_RATIO * ENCODER_TICKS_PER_ROTATION) / 360.0) // encoder ticks per output shaft rotation degree
	GROVE_I2C_LENGTH_LIMIT     = 16
	MOTOR_FLOAT                = -128
)

// GroveMode sets the mode of AD pins on the gopigo3.
type GroveMode byte

const (
	GROVE_INPUT_DIGITAL          GroveMode = 0
	GROVE_OUTPUT_DIGITAL         GroveMode = 1
	GROVE_INPUT_DIGITAL_PULLUP   GroveMode = 2
	GROVE_INPUT_DIGITAL_PULLDOWN GroveMode = 3
	GROVE_INPUT_ANALOG           GroveMode = 4
	GROVE_OUTPUT_PWM             GroveMode = 5
	GROVE_INPUT_ANALOG_PULLUP    GroveMode = 6
	GROVE_INPUT_ANALOG_PULLDOWN  GroveMode = 7
)

// Servo contains the address for the 2 servo ports.
type Servo byte

const (
	SERVO_1 Servo = 0x01
	SERVO_2 Servo = 0x02
)

// Motor contains the address for the left and right motors.
type Motor byte

const (
	MOTOR_LEFT  Motor = 0x01
	MOTOR_RIGHT Motor = 0x02
)

// Led contains the addresses for all leds on the board.
type Led byte

const (
	LED_EYE_RIGHT     Led = 0x01
	LED_EYE_LEFT      Led = 0x02
	LED_BLINKER_LEFT  Led = 0x04
	LED_BLINKER_RIGHT Led = 0x08
	LED_WIFI          Led = 0x80
)

// Grove contains the addresses for pins/ports of the AD1/AD2 grove connector.
type Grove byte

const (
	AD_1_1 Grove = 0x01 // default pin for most grove devices, A0/D0
	AD_1_2 Grove = 0x02
	AD_2_1 Grove = 0x04 // default pin for most grove devices, A0/D0
	AD_2_2 Grove = 0x08
	AD_1   Grove = AD_1_1 + AD_1_2
	AD_2   Grove = AD_2_1 + AD_2_2
)

// GroveType represents the type of a grove device.
type GroveType int

const (
	CUSTOM        GroveType = 1
	IR_DI_REMOTE            = 2
	IR_EV3_REMOTE           = 3
	US                      = 4
	I2C                     = 5
)

// GroveState contains the state of a grove device.
type GroveState int

const (
	VALID_DATA GroveState = iota
	NOT_CONFIGURED
	CONFIGURING
	NO_DATA
	I2C_ERROR
)

// Driver is a Gobot Driver for the GoPiGo3 board.
type Driver struct {
	name       string
	connector  spi.Connector
	connection spi.Connection
}

// NewDriver creates a new Gobot Driver for the GoPiGo3 board.
//
// Params:
//		a *Adaptor - the Adaptor to use with this Driver
//
func NewDriver(a spi.Connector) *Driver {
	g := &Driver{
		name:      gobot.DefaultName("GoPiGo3"),
		connector: a,
	}
	return g
}

// Name returns the name of the device.
func (g *Driver) Name() string { return g.name }

// SetName sets the name of the device.
func (g *Driver) SetName(n string) { g.name = n }

// Connection returns the Connection of the device.
func (g *Driver) Connection() gobot.Connection { return g.connection.(gobot.Connection) }

// Halt stops the driver.
func (g *Driver) Halt() (err error) {
	g.resetAll()
	time.Sleep(10 * time.Millisecond)
	return
}

// Start initializes the GoPiGo3
func (g *Driver) Start() (err error) {
	bus := g.connector.GetSpiDefaultBus()
	mode := g.connector.GetSpiDefaultMode()
	maxSpeed := g.connector.GetSpiDefaultMaxSpeed()
	g.connection, err = g.connector.GetSpiConnection(bus, mode, maxSpeed)
	if err != nil {
		return err
	}
	return nil
}

// GetManufacturerName returns the manufacturer from the firmware.
func (g *Driver) GetManufacturerName() (mName string, err error) {
	// read 24 bytes to get manufacturer name
	response, err := g.readBytes(goPiGo3Address, GET_MANUFACTURER, 24)
	if err != nil {
		return mName, err
	}
	if err := g.responseValid(response); err != nil {
		return mName, err
	}
	mf := response[4:23]
	mf = bytes.Trim(mf, "\x00")
	return fmt.Sprintf("%s", string(mf)), nil
}

// GetBoardName returns the board name from the firmware.
func (g *Driver) GetBoardName() (bName string, err error) {
	// read 24 bytes to get board name
	response, err := g.readBytes(goPiGo3Address, GET_NAME, 24)
	if err != nil {
		return bName, err
	}
	if err := g.responseValid(response); err != nil {
		return bName, err
	}
	mf := response[4:23]
	mf = bytes.Trim(mf, "\x00")
	return fmt.Sprintf("%s", string(mf)), nil
}

// GetHardwareVersion returns the hardware version from the firmware.
func (g *Driver) GetHardwareVersion() (hVer string, err error) {
	response, err := g.readUint32(goPiGo3Address, GET_HARDWARE_VERSION)
	if err != nil {
		return hVer, err
	}
	major := response / 1000000
	minor := response / 1000 % 1000
	patch := response % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

// GetFirmwareVersion returns the current firmware version.
func (g *Driver) GetFirmwareVersion() (fVer string, err error) {
	response, err := g.readUint32(goPiGo3Address, GET_FIRMWARE_VERSION)
	if err != nil {
		return fVer, err
	}
	major := response / 1000000
	minor := response / 1000 % 1000
	patch := response % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

// GetSerialNumber returns the 128-bit hardware serial number of the board.
func (g *Driver) GetSerialNumber() (sNum string, err error) {
	// read 20 bytes to get the serial number
	response, err := g.readBytes(goPiGo3Address, GET_ID, 20)
	if err != nil {
		return sNum, err
	}
	if err := g.responseValid(response); err != nil {
		return sNum, err
	}
	return fmt.Sprintf("%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X", response[4], response[5], response[6], response[7], response[8], response[9], response[10], response[11], response[12], response[13], response[14], response[15], response[16], response[17], response[18], response[19]), nil
}

// Get5vVoltage returns the current voltage on the 5v line.
func (g *Driver) Get5vVoltage() (voltage float32, err error) {
	val, err := g.readUint16(goPiGo3Address, GET_VOLTAGE_5V)
	return (float32(val) / 1000.0), err
}

// GetBatteryVoltage gets the battery voltage from the main battery pack (7v-12v).
func (g *Driver) GetBatteryVoltage() (voltage float32, err error) {
	val, err := g.readUint16(goPiGo3Address, GET_VOLTAGE_VCC)
	return (float32(val) / 1000.0), err
}

// SetLED sets rgb values from 0 to 255.
func (g *Driver) SetLED(led Led, red, green, blue uint8) error {
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_LED,
		byte(led),
		byte(red),
		byte(green),
		byte(blue),
	})
}

// SetServo sets a servo's position in microseconds (0-16666).
func (g *Driver) SetServo(servo Servo, us uint16) error {
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_SERVO,
		byte(servo),
		byte((us >> 8) & 0xFF),
		byte(us & 0xFF),
	})
}

func (g *Driver) ServoWrite(pin string, angle byte) error {
	pulseWidthRange := 1850
	if angle > 180 {
		angle = 180
	}
	if angle < 0 {
		angle = 0
	}
	pulseWidth := ((1500 - (pulseWidthRange / 2)) + ((pulseWidthRange / 180) * int(angle)))
	servoPin, err := strToServo(pin)
	if err != nil {
		return err
	}
	return g.SetServo(servoPin, uint16(pulseWidth))
}

// SetMotorPower sets a motor's power from -128 to 127.
func (g *Driver) SetMotorPower(motor Motor, power int8) error {
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_PWM,
		byte(motor),
		byte(power),
	})
}

// SetMotorPosition sets the motor's position in degrees.
func (g *Driver) SetMotorPosition(motor Motor, position float64) error {
	positionRaw := math.Float64bits(position * MOTOR_TICKS_PER_DEGREE)
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_POSITION,
		byte(motor),
		byte((positionRaw >> 24) & 0xFF),
		byte((positionRaw >> 16) & 0xFF),
		byte((positionRaw >> 8) & 0xFF),
		byte(positionRaw & 0xFF),
	})
}

// SetMotorDps sets the motor target speed in degrees per second.
func (g *Driver) SetMotorDps(motor Motor, dps float64) error {
	dpsUint := math.Float64bits(dps * MOTOR_TICKS_PER_DEGREE)
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_DPS,
		byte(motor),
		byte((dpsUint >> 8) & 0xFF),
		byte(dpsUint & 0xFF),
	})
}

// SetMotorLimits sets the speed limits for a motor.
func (g *Driver) SetMotorLimits(motor Motor, power int8, dps float64) error {
	dpsUint := math.Float64bits(dps * MOTOR_TICKS_PER_DEGREE)
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_LIMITS,
		byte(motor),
		byte(power),
		byte((dpsUint >> 8) & 0xFF),
		byte(dpsUint & 0xFF),
	})
}

// GetMotorStatus returns the status for the given motor.
func (g *Driver) GetMotorStatus(motor Motor) (flags uint8, power uint16, encoder, dps float64, err error) {
	message := GET_MOTOR_STATUS_RIGHT
	if motor == MOTOR_LEFT {
		message = GET_MOTOR_STATUS_LEFT
	}
	response, err := g.readBytes(goPiGo3Address, message, 12)
	if err != nil {
		return flags, power, encoder, dps, err
	}
	if err := g.responseValid(response); err != nil {
		return flags, power, encoder, dps, err
	}
	flags = uint8(response[4])
	power = uint16(response[5])
	if power&0x80 == 0x80 {
		power = power - 0x100
	}
	enc := uint64(response[6]<<24 | response[7]<<16 | response[8]<<8 | response[9])
	if enc&0x80000000 == 0x80000000 {
		encoder = float64(enc - 0x100000000)
	}
	d := uint64(response[10]<<8 | response[11])
	if d&0x8000 == 0x8000 {
		dps = float64(d - 0x10000)
	}
	return flags, power, encoder / MOTOR_TICKS_PER_DEGREE, dps / MOTOR_TICKS_PER_DEGREE, nil

}

// GetMotorEncoder reads a motor's encoder in degrees.
func (g *Driver) GetMotorEncoder(motor Motor) (encoder int64, err error) {
	message := GET_MOTOR_ENCODER_RIGHT
	if motor == MOTOR_LEFT {
		message = GET_MOTOR_ENCODER_LEFT
	}
	response, err := g.readUint32(goPiGo3Address, message)
	if err != nil {
		return encoder, err
	}
	encoder = int64(response)
	if response&0x80000000 != 0 {
		encoder = encoder - 0x100000000
	}
	encoder = encoder / MOTOR_TICKS_PER_DEGREE
	return encoder, nil
}

// OffsetMotorEncoder offsets a motor's encoder for calibration purposes.
func (g *Driver) OffsetMotorEncoder(motor Motor, offset float64) error {
	offsetUint := math.Float64bits(offset * MOTOR_TICKS_PER_DEGREE)
	return g.writeBytes([]byte{
		goPiGo3Address,
		OFFSET_MOTOR_ENCODER,
		byte(motor),
		byte((offsetUint >> 24) & 0xFF),
		byte((offsetUint >> 16) & 0xFF),
		byte((offsetUint >> 8) & 0xFF),
		byte(offsetUint & 0xFF),
	})
}

// SetGroveType sets the given port to a grove device type.
func (g *Driver) SetGroveType(port Grove, gType GroveType) error {
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_TYPE,
		byte(port),
		byte(gType),
	})
}

// SetGroveMode sets the mode a given pin/port of the grove connector.
func (g *Driver) SetGroveMode(port Grove, mode GroveMode) error {
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_MODE,
		byte(port),
		byte(mode),
	})
}

// SetPWMDuty sets the pwm duty cycle for the given pin/port.
func (g *Driver) SetPWMDuty(port Grove, duty uint8) (err error) {
	if duty < 0 {
		duty = 0
	}
	if duty > 100 {
		duty = 100
	}
	duty = duty * 10
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_PWM_DUTY,
		byte(port),
		byte((duty >> 8) & 0xFF),
		byte(duty & 0xFF),
	})
}

// SetPWMFreq setst the pwm frequency for the given pin/port.
func (g *Driver) SetPWMFreq(port Grove, freq uint16) error {
	if freq < 3 {
		freq = 3
	}
	if freq > 48000 {
		freq = 48000
	}
	return g.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_PWM_FREQUENCY,
		byte(port),
		byte((freq >> 8) & 0xFF),
		byte(freq & 0xFF),
	})
}

// PwmWrite implents the pwm interface for the gopigo3.
func (g *Driver) PwmWrite(pin string, val byte) (err error) {
	var (
		grovePin, grovePort Grove
	)
	grovePin, grovePort, _, _, err = getGroveAddresses(pin)
	if err != nil {
		return err
	}
	err = g.SetGroveType(grovePort, CUSTOM)
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	err = g.SetGroveMode(grovePin, GROVE_OUTPUT_PWM)
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	err = g.SetPWMFreq(grovePin, 24000)
	if err != nil {
		return err
	}
	val64 := math.Float64frombits(uint64(val))
	dutyCycle := uint8(math.Float64bits((100.0 / 255.0) * val64))
	return g.SetPWMDuty(grovePin, dutyCycle)
}

// AnalogRead returns the analog value of the given pin.
func (g *Driver) AnalogRead(pin string) (value int, err error) {
	var (
		grovePin, grovePort Grove
		analogCmd           byte
	)
	grovePin, grovePort, analogCmd, _, err = getGroveAddresses(pin)
	if err != nil {
		return value, err
	}
	err = g.SetGroveType(grovePort, CUSTOM)
	if err != nil {
		return value, err
	}
	time.Sleep(10 * time.Millisecond)
	err = g.SetGroveMode(grovePin, GROVE_INPUT_ANALOG)
	if err != nil {
		return value, err
	}
	time.Sleep(10 * time.Millisecond)
	response, err := g.readBytes(goPiGo3Address, analogCmd, 7)
	if err != nil {
		return value, err
	}
	if err := g.responseValid(response); err != nil {
		return value, err
	}
	if err := g.valueValid(response); err != nil {
		return value, err
	}
	return int((uint64(response[5]<<8) & 0xFF00) | uint64(response[6]&0xFF)), nil
	return
}

// DigitalWrite writes a 0/1 value to the given pin.
func (g *Driver) DigitalWrite(pin string, val byte) (err error) {
	var (
		grovePin, grovePort Grove
	)
	grovePin, grovePort, _, _, err = getGroveAddresses(pin)
	if err != nil {
		return err
	}
	err = g.SetGroveType(grovePort, CUSTOM)
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	err = g.SetGroveMode(grovePin, GROVE_OUTPUT_DIGITAL)
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	err = g.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_STATE,
		byte(grovePin),
		byte(val),
	})
	return err
}

// DigitalRead reads the 0/1 value from the given pin.
func (g *Driver) DigitalRead(pin string) (value int, err error) {
	var (
		grovePin, grovePort Grove
		stateCmd            byte
	)
	grovePin, grovePort, _, stateCmd, err = getGroveAddresses(pin)
	if err != nil {
		return value, err
	}
	err = g.SetGroveType(grovePort, CUSTOM)
	if err != nil {
		return value, err
	}
	time.Sleep(10 * time.Millisecond)
	err = g.SetGroveMode(grovePin, GROVE_INPUT_DIGITAL)
	if err != nil {
		return value, err
	}
	time.Sleep(10 * time.Millisecond)
	response, err := g.readBytes(goPiGo3Address, stateCmd, 6)
	if err != nil {
		return value, err
	}
	if err := g.responseValid(response); err != nil {
		return value, err
	}
	if err := g.valueValid(response); err != nil {
		return value, err
	}
	return int(response[5]), nil
}

func strToServo(pin string) (servo Servo, err error) {
	switch pin {
	case "SERVO_1":
		servo = SERVO_1
	case "SERVO_2":
		servo = SERVO_2
	default:
		err = fmt.Errorf("Invalid servo pin name")
	}
	return servo, err
}

func getGroveAddresses(pin string) (gPin, gPort Grove, analog, state byte, err error) {
	switch pin {
	case "AD_1_1":
		gPin = AD_1_1
		gPort = AD_1
		analog = GET_GROVE_ANALOG_1_1
		state = GET_GROVE_STATE_1_1
	case "AD_1_2":
		gPin = AD_1_2
		gPort = AD_1
		analog = GET_GROVE_ANALOG_1_2
		state = GET_GROVE_STATE_1_2
	case "AD_2_1":
		gPin = AD_2_1
		gPort = AD_2
		analog = GET_GROVE_ANALOG_2_1
		state = GET_GROVE_STATE_2_1
	case "AD_2_2":
		gPin = AD_2_2
		gPort = AD_2
		analog = GET_GROVE_ANALOG_2_2
		state = GET_GROVE_STATE_2_2
	default:
		err = fmt.Errorf("Invalid grove pin name")
	}
	return gPin, gPort, analog, state, err
}

func (g *Driver) responseValid(response []byte) error {
	if response[3] != 0xA5 {
		return fmt.Errorf("No SPI response, response not valid")
	}
	return nil
}

func (g *Driver) valueValid(value []byte) error {
	if value[4] != byte(VALID_DATA) {
		return fmt.Errorf("Invalid value")
	}
	return nil
}

func (g *Driver) readBytes(address byte, msg byte, numBytes int) (val []byte, err error) {
	w := make([]byte, numBytes)
	w[0] = address
	w[1] = msg
	r := make([]byte, len(w))
	err = g.connection.Tx(w, r)
	if err != nil {
		return val, err
	}
	return r, nil
}

func (g *Driver) readUint16(address, msg byte) (val uint16, err error) {
	r, err := g.readBytes(address, msg, 8)
	if err != nil {
		return val, err
	}
	if err := g.responseValid(r); err != nil {
		return val, err
	}
	return uint16(r[4])<<8 | uint16(r[5]), nil
}

func (g *Driver) readUint32(address, msg byte) (val uint32, err error) {
	r, err := g.readBytes(address, msg, 8)
	if err != nil {
		return val, err
	}
	if err := g.responseValid(r); err != nil {
		return val, err
	}
	return uint32(r[4])<<24 | uint32(r[5])<<16 | uint32(r[6])<<8 | uint32(r[7]), nil
}

func (g *Driver) writeBytes(w []byte) (err error) {
	return g.connection.Tx(w, nil)
}

func (g *Driver) resetAll() {
	g.SetGroveType(AD_1+AD_2, CUSTOM)
	time.Sleep(10 * time.Millisecond)
	g.SetGroveMode(AD_1+AD_2, GROVE_INPUT_DIGITAL)
	time.Sleep(10 * time.Millisecond)
	g.SetMotorPower(MOTOR_LEFT+MOTOR_RIGHT, 0.0)
	g.SetMotorLimits(MOTOR_LEFT+MOTOR_RIGHT, 0, 0)
	g.SetServo(SERVO_1+SERVO_2, 0)
	g.SetLED(LED_EYE_LEFT+LED_EYE_RIGHT+LED_BLINKER_LEFT+LED_BLINKER_RIGHT, 0, 0, 0)
}
