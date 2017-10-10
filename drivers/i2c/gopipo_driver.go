package i2c

import (
	"time"

	"github.com/hybridgroup/gobot"
)

// debug=0
//
// us_cmd				=[117]		#Read the distance from the ultrasonic sensor
// en_com_timeout_cmd	=[80]		#Enable communication timeout
// dis_com_timeout_cmd	=[81]		#Disable communication timeout
// timeout_status_cmd	=[82]		#Read the timeout status
// enc_read_cmd		=[53]		#Read encoder values
// trim_test_cmd		=[30]		#Test the trim values
// trim_write_cmd		=[31]		#Write the trim values
// trim_read_cmd		=[32]
//
// ir_read_cmd			=[21]
// ir_recv_pin_cmd		=[22]
// cpu_speed_cmd		=[25]

// v16_thresh=790
const (
	GOPIGO_ADDRESS = 0x08 // i2c address for the GoPiGo

	GOPIGO_STOP             = 0x78 // stop the gopigo
	GOPIGO_FIRMWARE_VERSION = 0x10 // get the version of the GoPiGo

	GOPIGO_MOTOR1           = 0x6F // control motor 1
	GOPIGO_MOTOR2           = 0x70 // control motor 2
	GOPIGO_READ_MOTOR_SPEED = 0x72 // read motor speed back

	GOPIGO_SET_LEFT_SPEED  = 0x46 // set speed for the left motor
	GOPIGO_SET_RIGHT_SPEED = 0x47 // set speed for the left motor
	GOPIGO_INCREASE_SPEED  = 0x74 // increase the speed by 10
	GOPIGO_DECREASE_SPEED  = 0x67 // iecrease the speed by 10

	GOPIGO_FORWARD        = 0x77 // move forward with PID control
	GOPIGO_MOTOR_FORWARD  = 0x69 // move forward without PID control
	GOPIGO_BACKWARD       = 0x73 // move backward with PID control
	GOPIGO_MOTOR_BACKWARD = 0x6B // move backward without PID control

	GOPIGO_TURN_LEFT    = 0x73 // turn left by turning off one motor
	GOPIGO_TURN_RIGHT   = 0x64 // turn right by turning off one motor
	GOPIGO_ROTATE_LEFT  = 0x62 // rotate left by running both motors is opposite direction
	GOPIGO_ROTATE_RIGHT = 0x6E // rotate Right by running both motors is opposite direction

	GOPIGO_ENABLE_ENCODER  = 0x33 // enable the encoders
	GOPIGO_DISABLE_ENCODER = 0x34 // disable the encoders
	GOPIGO_ENCODER_TARGET  = 0x30 // set the encoder targeting
	GOPIGO_ENCODER_STATUS  = 0x35 // read encoder status

	GOPIGO_ENABLE_SERVO  = 0x3d // enable the servo's
	GOPIGO_DISABLE_SERVO = 0x3c // disable the servo's
	GOPIGO_SERVO         = 0x65 // rotate the servo
)

type GoPiGoDriver struct {
	name       string
	connection I2c
}

func NewGoPiGoDriver(a I2c) *GoPiGoDriver {

	return &GoPiGoDriver{
		name:       "GoPiGo",
		connection: a,
	}
}

func (g *GoPiGoDriver) SetName(name string) {
	g.name = name
}

// Name identifies this driver objectName identifies this driver object
func (g *GoPiGoDriver) Name() string { return g.name }

func (g *GoPiGoDriver) Start() error {
	return g.connection.I2cStart(GOPIGO_ADDRESS)
}

// Halt terminates the driver
func (g *GoPiGoDriver) Halt() error {
	err := g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_STOP, 0, 0, 0})
	if err != nil {
		return err
	}
	return g.connection.Finalize()
}

func (g *GoPiGoDriver) Connection() gobot.Connection {
	return g.connection
}

func (g *GoPiGoDriver) FirmwareVersion() (int, error) {
	err := g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_FIRMWARE_VERSION, 0, 0, 0})
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	d, err := g.connection.I2cRead(GOPIGO_ADDRESS, 2)
	if err != nil {
		return 0, err
	}
	return int(d[0]), nil
}

// Motor1  controls motor 1
func (g *GoPiGoDriver) Motor1(direction byte, speed byte) error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_MOTOR1, direction, speed, 0})
}

// Motor2 controls motor 2
func (g *GoPiGoDriver) Motor2(direction byte, speed byte) error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_MOTOR2, direction, speed, 0})
}

func (g *GoPiGoDriver) SetLeftSpeed(speed byte) error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_SET_LEFT_SPEED, speed, 0, 0})
}

func (g *GoPiGoDriver) SetRightSpeed(speed byte) error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_SET_RIGHT_SPEED, speed, 0, 0})
}

// ReadMotorSpeed returns the current speed for motor 1 and 2
func (g *GoPiGoDriver) ReadMotorSpeed() (byte, byte, error) {
	err := g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_READ_MOTOR_SPEED, 0, 0, 0})
	if err != nil {
		return 0, 0, err
	}
	time.Sleep(100 * time.Millisecond)
	d, err := g.connection.I2cRead(GOPIGO_ADDRESS, 2)
	if err != nil {
		return 0, 0, err
	}
	return d[0], d[1], nil
}

func (g *GoPiGoDriver) IncreaseSpeed() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_INCREASE_SPEED, 0, 0, 0})
}

func (g *GoPiGoDriver) DecraesSpeed() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_DECREASE_SPEED, 0, 0, 0})
}

// Forward drives the robot forward with PID control
func (g *GoPiGoDriver) Forward() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_FORWARD, 0, 0, 0})
}

// MotorForward drives forward without PID control
func (g *GoPiGoDriver) MotorForward() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_MOTOR_FORWARD, 0, 0, 0})
}

// Backward drives the robot forward with PID control
func (g *GoPiGoDriver) Backward() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_BACKWARD, 0, 0, 0})
}

// MotorBackward moves the robot backword with out PDI control
func (g *GoPiGoDriver) MotorBackward() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_MOTOR_BACKWARD, 0, 0, 0})
}

// TurnLeft by turning off one motor
func (g *GoPiGoDriver) TurnLeft() error {
	// return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIPO_TURN_LEFT, 0, 0, 0})
	return nil
}

// TurnLeft by turning off one motor
func (g *GoPiGoDriver) TurnRight() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_TURN_RIGHT, 0, 0, 0})
}

// RotateLeft by running both motors is opposite direction
func (g *GoPiGoDriver) RotateLeft() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_ROTATE_LEFT, 0, 0, 0})
}

// RotateRight by running both motors is opposite direction
func (g *GoPiGoDriver) RotateRight() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_ROTATE_RIGHT, 0, 0, 0})
}

// // EncoderTurnLeft
// func EncoderTurnLeft(degrees int) error {
//
// }
//
// func EncoderTurnRight(degrees int) error {
//
// }
//
// func EncoderForward(distance int) error {
//
// }
//
// func EncoderBackward(distance int) error {
//
// }

// EnableEncoder turns the encoders on
func (g *GoPiGoDriver) EnableEncoder() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_ENABLE_ENCODER, 0, 0, 0})
}

// EnableEncoder turns the encoders off
func (g *GoPiGoDriver) DisableEncoder() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_DISABLE_ENCODER, 0, 0, 0})
}

// Set encoder targeting on for motor 1 and/or motor 2 with the number of encoder pulses
func (g *GoPiGoDriver) EncoderTarget(motor1 bool, motor2 bool, pulses int) error {
	motorSelect := byte(0)
	if motor1 {
		motorSelect += 2
	}
	if motor2 {
		motorSelect += 1
	}
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_ENCODER_TARGET, motorSelect, byte(pulses / 256), byte(pulses % 256)})
}

// TODO: figure this one out
// func (g *GoPiGoDriver) EncoderStatus()

// Volt reads the voltage in V
func (g *GoPiGoDriver) Volt() (float64, error) {
	// err := g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_VOLT, 0, 0, 0})
	// if err != nil {
	// 	return err
	// }
	// time.Sleep(100 * time.Millisecond)
	// d, err := g.connection.I2cRead(GOPIGO_ADDRESS, 2)
	// if err != nil {
	// 	return err
	// }
	//
	// return (d[0]*265 + d[1]) / 1024 / 0.4
	return 0, nil
}

// Enable turns the servo on
func (g *GoPiGoDriver) EnableServo() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_ENABLE_SERVO, 0, 0, 0})
}

// DisableServo turns the servo off
func (g *GoPiGoDriver) DisableServo() error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_DISABLE_SERVO, 0, 0, 0})
}

// Servo sets the servo to angle with degrees
func (g *GoPiGoDriver) Servo(degrees byte) error {
	return g.connection.I2cWrite(GOPIGO_ADDRESS, []byte{GOPIGO_SERVO, degrees, 0, 0})
}
