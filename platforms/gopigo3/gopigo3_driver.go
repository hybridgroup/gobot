package gopigo3

import (
	"errors"
	"fmt"

	"gobot.io/x/gobot"
)

type Motor byte
type Led byte

const (
	MOTOR_LEFT  Motor = 0x01
	MOTOR_RIGHT Motor = 0x02

	LED_EYE_LEFT      Led = 0x02
	LED_EYE_RIGHT     Led = 0x01
	LED_BLINKER_LEFT  Led = 0x04
	LED_BLINKER_RIGHT Led = 0x08
	LED_LEFT_EYE      Led = LED_EYE_LEFT
	LED_RIGHT_EYE     Led = LED_EYE_RIGHT
	LED_LEFT_BLINKER  Led = LED_BLINKER_LEFT
	LED_RIGHT_BLINKER Led = LED_BLINKER_RIGHT
	LED_WIFI          Led = 0x80
)

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

// GoPiGo3Driver represents a GoPiGo3
type GoPiGo3Driver struct {
	name       string
	connection gobot.Connection
}

// NewGoPiGo3Driver returns a new GoPiGo3Driver given a GoPiGo3Adaptor
func NewGoPiGo3Driver(connection *Adaptor) *GoPiGo3Driver {
	return &GoPiGo3Driver{
		name:       "GoPiGo",
		connection: connection,
	}
}

// Name identifies this driver objectName identifies this driver object
func (g *GoPiGo3Driver) Name() string {
	return g.name
}

func (g *GoPiGo3Driver) SetName(name string) {
	g.name = name
}

func (g *GoPiGo3Driver) Start() error {
	return nil
}

// Halt terminates the driver
func (g *GoPiGo3Driver) Halt() error {
	return nil
}

// Connection returns the Driver's Connection
func (g *GoPiGo3Driver) Connection() gobot.Connection {
	return g.connection
}

func (g *GoPiGo3Driver) adaptor() *Adaptor {
	return g.Connection().(*Adaptor)
}

// SetLed sets a specific led to an rgb color
func (g *GoPiGo3Driver) SetLed(led Led, red, green, blue uint8) error {
	if led < 0 || led > 255 {
		return errors.New(fmt.Sprintf("Led value can only be (0..255). Was: ", led))
	}

	return g.adaptor().connection.Write([]byte{
		GOPIGO_ADDRESS,
		SET_LED,
		byte(led),
		byte(red),
		byte(green),
		byte(blue),
	})
}

// SetMotorPower sets a motor's power in percent
func (g *GoPiGo3Driver) SetMotorPower(motor Motor, power int) error {
	if power > 100 {
		power = 100
	}

	if power < -100 {
		power = -100
	}

	return g.adaptor().connection.Write([]byte{
		GOPIGO_ADDRESS,
		SET_MOTOR_PWM,
		byte(motor),
		byte(power),
	})
}

func (g *GoPiGo3Driver) GetFirmwareVersion() string {
	response := g.adaptor().connection.Read32(GOPIGO_ADDRESS, GET_FIRMWARE_VERSION)
	major := response / 1000000
	minor := response / 1000 % 1000
	patch := response % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
