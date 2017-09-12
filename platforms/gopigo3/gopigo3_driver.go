// Based on https://github.com/DexterInd/GoPiGo3/blob/master/Software/Python/gopigo3.py
// Uou will need to run the following script if using a stock raspbian image before this library will work:
// https://www.dexterindustries.com/GoPiGo/get-started-with-the-gopigo3-raspberry-pi-robot/3-program-your-raspberry-pi-robot/python-programming-language/
package gopigo3

import (
	"fmt"
	"math"

	"gobot.io/x/gobot"
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
	WHEEL_BASE_WIDTH             = 117                                                       // distance (mm) from left wheel to right wheel. This works with the initial GPG3 prototype. Will need to be adjusted.
	WHEEL_DIAMETER               = 66.5                                                      // wheel diameter (mm)
	WHEEL_BASE_CIRCUMFERENCE     = WHEEL_BASE_WIDTH * math.Pi                                // circumference of the circle the wheels will trace while turning (mm)
	WHEEL_CIRCUMFERENCE          = WHEEL_DIAMETER * math.Pi                                  // circumference of the wheels (mm)
	MOTOR_GEAR_RATIO             = 120                                                       // motor gear ratio
	ENCODER_TICKS_PER_ROTATION   = 6                                                         // encoder ticks per motor rotation (number of magnet positions)
	MOTOR_TICKS_PER_DEGREE       = ((MOTOR_GEAR_RATIO * ENCODER_TICKS_PER_ROTATION) / 360.0) // encoder ticks per output shaft rotation degree
	GROVE_I2C_LENGTH_LIMIT       = 16
	MOTOR_FLOAT                  = -128
	GROVE_INPUT_DIGITAL          = 0
	GROVE_OUTPUT_DIGITAL         = 1
	GROVE_INPUT_DIGITAL_PULLUP   = 2
	GROVE_INPUT_DIGITAL_PULLDOWN = 3
	GROVE_INPUT_ANALOG           = 4
	GROVE_OUTPUT_PWM             = 5
	GROVE_INPUT_ANALOG_PULLUP    = 6
	GROVE_INPUT_ANALOG_PULLDOWN  = 7
	GROVE_LOW                    = 0
	GROVE_HIGH                   = 1
)

// Servo contains the address for the 2 servo ports.
type Servo byte

const (
	SERVO_1 = 0x01
	SERVO_2 = 0x02
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

// Grove contains the addresses for the optional grove devices.
type Grove byte

const (
	GROVE_1_1 = 0x01
	GROVE_1_2 = 0x02
	GROVE_2_1 = 0x04
	GROVE_2_2 = 0x08
	GROVE_1   = GROVE_1_1 + GROVE_1_2
	GROVE_2   = GROVE_2_1 + GROVE_2_2
)

// GroveType represents the type of a grove device.
type GroveType int

const (
	CUSTOM GroveType = iota
	IR_DI_REMOTE
	IR_EV3_REMOTE
	US
	I2C
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

var _ gobot.Driver = (*GoPiGo3Driver)(nil)

// GoPiGo3Driver is a Gobot Driver for the GoPiGo3 board.
type GoPiGo3Driver struct {
	name       string
	connection *Adaptor
}

// NewGoPiGo3Driver creates a new Gobot Driver for the GoPiGo3 board.
//
// Params:
//		a *Adaptor - the Adaptor to use with this Driver
//
func NewGoPiGo3Driver(a *Adaptor) *GoPiGo3Driver {
	g := &GoPiGo3Driver{
		name:       gobot.DefaultName("GoPiGo3"),
		connection: a,
	}
	return g
}

// Name returns the name of the device.
func (g *GoPiGo3Driver) Name() string { return g.name }

// SetName sets the name of the device.
func (g *GoPiGo3Driver) SetName(n string) { g.name = n }

// Connection returns the Connection of the device.
func (g *GoPiGo3Driver) Connection() gobot.Connection { return gobot.Connection(g.connection) }

// Halt stops the driver.
func (g *GoPiGo3Driver) Halt() (err error) { return }

// Start initializes the GoPiGo3
func (g *GoPiGo3Driver) Start() (err error) {
	return nil
}

// GetManufacturerName returns the manufacturer from the firmware.
func (g *GoPiGo3Driver) GetManufacturerName() (mName string, err error) {
	// read 24 bytes to get manufacturer name
	response, err := g.connection.ReadBytes(goPiGo3Address, GET_MANUFACTURER, 24)
	if err != nil {
		return mName, err
	}
	if response[3] == 0xA5 {
		mf := response[4:23]
		return fmt.Sprintf("%s", string(mf)), nil
	}
	return mName, nil
}

// GetBoardName returns the board name from the firmware.
func (g *GoPiGo3Driver) GetBoardName() (bName string, err error) {
	// read 24 bytes to get board name
	response, err := g.connection.ReadBytes(goPiGo3Address, GET_NAME, 24)
	if err != nil {
		return bName, err
	}
	if response[3] == 0xA5 {
		mf := response[4:23]
		return fmt.Sprintf("%s", string(mf)), nil
	}
	return bName, nil
}

// GetHardwareVersion returns the hardware version from the firmware.
func (g *GoPiGo3Driver) GetHardwareVersion() (hVer string, err error) {
	response, err := g.connection.ReadUint32(goPiGo3Address, GET_HARDWARE_VERSION)
	if err != nil {
		return hVer, err
	}
	major := response / 1000000
	minor := response / 1000 % 1000
	patch := response % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

// GetFirmwareVersion returns the current firmware version.
func (g *GoPiGo3Driver) GetFirmwareVersion() (fVer string, err error) {
	response, err := g.connection.ReadUint32(goPiGo3Address, GET_FIRMWARE_VERSION)
	if err != nil {
		return fVer, err
	}
	major := response / 1000000
	minor := response / 1000 % 1000
	patch := response % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

// GetSerialNumber returns the 128-bit hardware serial number of the board.
func (g *GoPiGo3Driver) GetSerialNumber() (sNum string, err error) {
	// read 20 bytes to get the serial number
	response, err := g.connection.ReadBytes(goPiGo3Address, GET_ID, 20)
	if err != nil {
		return sNum, err
	}
	if response[3] == 0xA5 {
		return fmt.Sprintf("%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X", response[4], response[5], response[6], response[7], response[8], response[9], response[10], response[11], response[12], response[13], response[14], response[15], response[16], response[17], response[18], response[19]), nil

	}
	return sNum, nil
}

// Get5vVoltage returns the current voltage on the 5v line.
func (g *GoPiGo3Driver) Get5vVoltage() (voltage float32, err error) {
	val, err := g.connection.ReadUint16(goPiGo3Address, GET_VOLTAGE_5V)
	return (float32(val) / 1000.0), err
}

// GetBatteryVoltage gets the battery voltage from the main battery pack (7v-12v).
func (g *GoPiGo3Driver) GetBatteryVoltage() (voltage float32, err error) {
	val, err := g.connection.ReadUint16(goPiGo3Address, GET_VOLTAGE_VCC)
	return (float32(val) / 1000.0), err
}

// SetLED sets rgb values from 0 to 255.
func (g *GoPiGo3Driver) SetLED(led Led, red, green, blue uint8) error {
	return g.connection.WriteBytes([]byte{
		goPiGo3Address,
		SET_LED,
		byte(led),
		byte(red),
		byte(green),
		byte(blue),
	})
}

// SetServo sets a servo's position in microseconds (0-16666).
func (g *GoPiGo3Driver) SetServo(servo Servo, us uint16) error {
	return g.connection.WriteBytes([]byte{
		goPiGo3Address,
		SET_SERVO,
		byte(servo),
		byte(us),
	})
}

// SetMotorPower from -128 to 127.
func (g *GoPiGo3Driver) SetMotorPower(motor Motor, power int8) error {
	return g.connection.WriteBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_PWM,
		byte(motor),
		byte(power),
	})
}

// SetMotorPosition sets the motor's position in degrees.
func (g *GoPiGo3Driver) SetMotorPosition(motor Motor, position float64) error {
	positionRaw := math.Float64bits(position * MOTOR_TICKS_PER_DEGREE)
	return g.connection.WriteBytes([]byte{
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
func (g *GoPiGo3Driver) SetMotorDps(motor Motor, dps float64) error {
	dpsUint := math.Float64bits(dps * MOTOR_TICKS_PER_DEGREE)
	return g.connection.WriteBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_DPS,
		byte(motor),
		byte((dpsUint >> 8) & 0xFF),
		byte(dpsUint & 0xFF),
	})
}

// SetMotorLimits sets the speed limits for a motor.
func (g *GoPiGo3Driver) SetMotorLimits(motor Motor, power int8, dps float64) error {
	dpsUint := math.Float64bits(dps * MOTOR_TICKS_PER_DEGREE)
	return g.connection.WriteBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_LIMITS,
		byte(motor),
		byte(power),
		byte((dpsUint >> 8) & 0xFF),
		byte(dpsUint & 0xFF),
	})
}

// GetMotorStatus returns the status for the given motor.
func (g *GoPiGo3Driver) GetMotorStatus(motor Motor) (flags uint8, power uint16, encoder, dps float64, err error) {
	message := GET_MOTOR_STATUS_RIGHT
	if motor == MOTOR_LEFT {
		message = GET_MOTOR_STATUS_LEFT
	}
	response, err := g.connection.ReadBytes(goPiGo3Address, message, 12)
	if err != nil {
		return flags, power, encoder, dps, err
	}
	if response[3] == 0xA5 {
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
	return flags, power, encoder, dps, nil
}

// GetMotorEncoder reads a motor's encoder in degrees.
func (g *GoPiGo3Driver) GetMotorEncoder(motor Motor) (encoder float64, err error) {
	message := GET_MOTOR_ENCODER_RIGHT
	if motor == MOTOR_LEFT {
		message = GET_MOTOR_ENCODER_LEFT
	}
	response, err := g.connection.ReadUint32(goPiGo3Address, message)
	if err != nil {
		return encoder, err
	}
	if response&0x80000000 == 0x80000000 {
		encoder = float64(encoder - 0x100000000)
	}
	return encoder, nil
}

// OffsetMotorEncoder offsets a motor's encoder for calibration purposes.
func (g *GoPiGo3Driver) OffsetMotorEncoder(motor Motor, offset float64) error {
	offsetUint := math.Float64bits(offset * MOTOR_TICKS_PER_DEGREE)
	return g.connection.WriteBytes([]byte{
		goPiGo3Address,
		OFFSET_MOTOR_ENCODER,
		byte(motor),
		byte((offsetUint >> 24) & 0xFF),
		byte((offsetUint >> 16) & 0xFF),
		byte((offsetUint >> 8) & 0xFF),
		byte(offsetUint & 0xFF),
	})
}

//TODO: add grove functions
