// Package gopigo3 is based on https://github.com/DexterInd/GoPiGo3/blob/master/Software/Python/gopigo3.py
// You will need to run the following commands if using a stock raspbian image before this library will work:
// sudo curl -kL dexterindustries.com/update_gopigo3 | bash
// sudo reboot
package gopigo3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/spi"
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
	// distance (mm) from left wheel to right wheel. This works with the initial GPG3 prototype. Will need to be adjusted.
	WHEEL_BASE_WIDTH = 117
	// wheel diameter (mm)
	WHEEL_DIAMETER = 66.5
	// circumference of the circle the wheels will trace while turning (mm)
	WHEEL_BASE_CIRCUMFERENCE = WHEEL_BASE_WIDTH * math.Pi
	// circumference of the wheels (mm)
	WHEEL_CIRCUMFERENCE = WHEEL_DIAMETER * math.Pi
	// motor gear ratio
	MOTOR_GEAR_RATIO = 120
	// encoder ticks per motor rotation (number of magnet positions)
	ENCODER_TICKS_PER_ROTATION = 6
	// encoder ticks per output shaft rotation degree
	MOTOR_TICKS_PER_DEGREE = ((MOTOR_GEAR_RATIO * ENCODER_TICKS_PER_ROTATION) / 360.0)
	GROVE_I2C_LENGTH_LIMIT = 16
	MOTOR_FLOAT            = -128
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
	AD11     string = "AD_1_1"
	AD12     string = "AD_1_2"
	AD21     string = "AD_2_1"
	AD22     string = "AD_2_2"
	AD_1_1_G Grove  = 0x01 // default pin for most grove devices, A0/D0
	AD_1_2_G Grove  = 0x02
	AD_2_1_G Grove  = 0x04 // default pin for most grove devices, A0/D0
	AD_2_2_G Grove  = 0x08
	AD_1_G   Grove  = AD_1_1_G + AD_1_2_G
	AD_2_G   Grove  = AD_2_1_G + AD_2_2_G
)

// GroveType represents the type of a grove device.
type GroveType int

const (
	CUSTOM        GroveType = 1
	IR_DI_REMOTE  GroveType = 2
	IR_EV3_REMOTE GroveType = 3
	US            GroveType = 4
	I2C           GroveType = 5
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
	spi.Config
}

// NewDriver creates a new Gobot Driver for the GoPiGo3 board.
//
// Params:
//
//	a *Adaptor - the Adaptor to use with this Driver
//
// Optional params:
//
//	 spi.WithBusNumber(int):  bus to use with this driver
//		spi.WithChipNumber(int): chip to use with this driver
//	 spi.WithMode(int):    	 mode to use with this driver
//	 spi.WithBitCount(int):   number of bits to use with this driver
//	 spi.WithSpeed(int64):    speed in Hz to use with this driver
func NewDriver(a spi.Connector, options ...func(spi.Config)) *Driver {
	spiConfig := spi.NewConfig()
	// use /dev/spidev0.1
	spiConfig.SetBusNumber(0)
	spiConfig.SetChipNumber(1)
	d := &Driver{
		name:      gobot.DefaultName("GoPiGo3"),
		connector: a,
		Config:    spiConfig,
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Name returns the name of the device.
func (d *Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *Driver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.name)
	return nil
}

// Halt stops the driver.
func (d *Driver) Halt() error {
	err := d.resetAll()
	time.Sleep(10 * time.Millisecond)
	return err
}

// Start initializes the GoPiGo3
func (d *Driver) Start() error {
	bus := d.GetBusNumberOrDefault(d.connector.SpiDefaultBusNumber())
	chip := d.GetChipNumberOrDefault(d.connector.SpiDefaultChipNumber())
	mode := d.GetModeOrDefault(d.connector.SpiDefaultMode())
	bits := d.GetBitCountOrDefault(d.connector.SpiDefaultBitCount())
	maxSpeed := d.GetSpeedOrDefault(d.connector.SpiDefaultMaxSpeed())

	var err error
	d.connection, err = d.connector.GetSpiConnection(bus, chip, mode, bits, maxSpeed)
	return err
}

// GetManufacturerName returns the manufacturer from the firmware.
func (d *Driver) GetManufacturerName() (string, error) {
	// read 24 bytes to get manufacturer name
	response, err := d.readBytes(goPiGo3Address, GET_MANUFACTURER, 24)
	if err != nil {
		return "", err
	}
	if err := d.responseValid(response); err != nil {
		return "", err
	}
	mf := response[4:23]
	mf = bytes.Trim(mf, "\x00")
	return string(mf), nil
}

// GetBoardName returns the board name from the firmware.
func (d *Driver) GetBoardName() (string, error) {
	// read 24 bytes to get board name
	response, err := d.readBytes(goPiGo3Address, GET_NAME, 24)
	if err != nil {
		return "", err
	}
	if err := d.responseValid(response); err != nil {
		return "", err
	}
	mf := response[4:23]
	mf = bytes.Trim(mf, "\x00")
	return string(mf), nil
}

// GetHardwareVersion returns the hardware version from the firmware.
func (d *Driver) GetHardwareVersion() (string, error) {
	response, err := d.readUint32(goPiGo3Address, GET_HARDWARE_VERSION)
	if err != nil {
		return "", err
	}
	major := response / 1000000
	minor := response / 1000 % 1000
	patch := response % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

// GetFirmwareVersion returns the current firmware version.
func (d *Driver) GetFirmwareVersion() (string, error) {
	response, err := d.readUint32(goPiGo3Address, GET_FIRMWARE_VERSION)
	if err != nil {
		return "", err
	}
	major := response / 1000000
	minor := response / 1000 % 1000
	patch := response % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

// GetSerialNumber returns the 128-bit hardware serial number of the board.
func (d *Driver) GetSerialNumber() (string, error) {
	// read 20 bytes to get the serial number
	response, err := d.readBytes(goPiGo3Address, GET_ID, 20)
	if err != nil {
		return "", err
	}
	if err := d.responseValid(response); err != nil {
		return "", err
	}
	return fmt.Sprintf("%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X", response[4], response[5],
		response[6], response[7], response[8], response[9], response[10], response[11], response[12], response[13],
		response[14], response[15], response[16], response[17], response[18], response[19]), nil
}

// Get5vVoltage returns the current voltage on the 5v line.
func (d *Driver) Get5vVoltage() (float32, error) {
	val, err := d.readUint16(goPiGo3Address, GET_VOLTAGE_5V)
	return (float32(val) / 1000.0), err
}

// GetBatteryVoltage gets the battery voltage from the main battery pack (7v-12v).
func (d *Driver) GetBatteryVoltage() (float32, error) {
	val, err := d.readUint16(goPiGo3Address, GET_VOLTAGE_VCC)
	return (float32(val) / 1000.0), err
}

// SetLED sets rgb values from 0 to 255.
func (d *Driver) SetLED(led Led, red, green, blue uint8) error {
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_LED,
		byte(led),
		red,
		green,
		blue,
	})
}

// SetServo sets a servo's position in microseconds (0-16666).
func (d *Driver) SetServo(srvo Servo, us uint16) error {
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_SERVO,
		byte(srvo),
		byte((us >> 8) & 0xFF),
		byte(us & 0xFF),
	})
}

// ServoWrite writes an angle (0-180) to the given servo (servo 1 or servo 2).
// Must implement the ServoWriter interface of gpio package.
func (d *Driver) ServoWrite(port string, angle byte) error {
	srvo := SERVO_1 // default for unknown ports
	if port == "2" || port == "SERVO_2" {
		srvo = SERVO_2
	}

	pulseWidthRange := 1850
	if angle > 180 {
		angle = 180
	}
	pulseWidth := ((1500 - (pulseWidthRange / 2)) + ((pulseWidthRange / 180) * int(angle)))
	return d.SetServo(srvo, uint16(pulseWidth)) //nolint:gosec // TODO: fix later
}

// SetMotorPower sets a motor's power from -128 to 127.
func (d *Driver) SetMotorPower(motor Motor, power int8) error {
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_PWM,
		byte(motor),
		byte(power),
	})
}

// SetMotorPosition sets the motor's position in degrees.
func (d *Driver) SetMotorPosition(motor Motor, position int) error {
	positionRaw := position * MOTOR_TICKS_PER_DEGREE
	return d.writeBytes([]byte{
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
func (d *Driver) SetMotorDps(motor Motor, dps int) error {
	mdps := dps * MOTOR_TICKS_PER_DEGREE
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_DPS,
		byte(motor),
		byte((mdps >> 8) & 0xFF),
		byte(mdps & 0xFF),
	})
}

// SetMotorLimits sets the speed limits for a motor.
func (d *Driver) SetMotorLimits(motor Motor, power int8, dps int) error {
	dpsUint := dps * MOTOR_TICKS_PER_DEGREE
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_MOTOR_LIMITS,
		byte(motor),
		byte(power),
		byte((dpsUint >> 8) & 0xFF),
		byte(dpsUint & 0xFF),
	})
}

// GetMotorStatus returns the status for the given motor.
//
//nolint:nonamedreturns // sufficient here
func (d *Driver) GetMotorStatus(motor Motor) (flags uint8, power uint16, encoder, dps int, err error) {
	message := GET_MOTOR_STATUS_RIGHT
	if motor == MOTOR_LEFT {
		message = GET_MOTOR_STATUS_LEFT
	}
	response, err := d.readBytes(goPiGo3Address, message, 12)
	if err != nil {
		return flags, power, encoder, dps, err
	}
	if err := d.responseValid(response); err != nil {
		return flags, power, encoder, dps, err
	}
	// get flags
	flags = response[4]
	// get power
	power = uint16(response[5])
	if power&0x80 == 0x80 {
		power = power - 0x100
	}
	// get encoder
	enc := make([]byte, 4)
	enc[3] = response[6]
	enc[2] = response[7]
	enc[1] = response[8]
	enc[0] = response[9]
	e := binary.LittleEndian.Uint32(enc)
	encoder = int(e)
	if e&0x80000000 == 0x80000000 {
		encoder = int(uint64(e) - 0x100000000) //nolint:gosec // TODO: fix later
	}
	// get dps
	dpsRaw := make([]byte, 4)
	dpsRaw[1] = response[10]
	dpsRaw[0] = response[11]
	ds := binary.LittleEndian.Uint32(dpsRaw)
	dps = int(ds)
	if ds&0x8000 == 0x8000 {
		dps = int(ds - 0x10000)
	}
	return flags, power, encoder / MOTOR_TICKS_PER_DEGREE, dps / MOTOR_TICKS_PER_DEGREE, nil
}

// GetMotorEncoder reads a motor's encoder in degrees.
func (d *Driver) GetMotorEncoder(motor Motor) (int64, error) {
	message := GET_MOTOR_ENCODER_RIGHT
	if motor == MOTOR_LEFT {
		message = GET_MOTOR_ENCODER_LEFT
	}
	response, err := d.readUint32(goPiGo3Address, message)
	if err != nil {
		return 0, err
	}
	encoder := int64(response)
	if response&0x80000000 != 0 {
		encoder = encoder - 0x100000000
	}
	encoder = encoder / MOTOR_TICKS_PER_DEGREE
	return encoder, nil
}

// OffsetMotorEncoder offsets a motor's encoder for calibration purposes.
func (d *Driver) OffsetMotorEncoder(motor Motor, offset float64) error {
	offsetUint := math.Float64bits(offset * MOTOR_TICKS_PER_DEGREE)
	return d.writeBytes([]byte{
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
func (d *Driver) SetGroveType(port Grove, gType GroveType) error {
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_TYPE,
		byte(port),
		byte(gType),
	})
}

// SetGroveMode sets the mode a given pin/port of the grove connector.
func (d *Driver) SetGroveMode(port Grove, mode GroveMode) error {
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_MODE,
		byte(port),
		byte(mode),
	})
}

// SetPWMDuty sets the pwm duty cycle for the given pin/port.
func (d *Driver) SetPWMDuty(port Grove, duty uint16) error {
	if duty > 100 {
		duty = 100
	}
	duty = duty * 10
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_PWM_DUTY,
		byte(port),
		byte((duty >> 8) & 0xFF),
		byte(duty & 0xFF),
	})
}

// SetPWMFreq setst the pwm frequency for the given pin/port.
func (d *Driver) SetPWMFreq(port Grove, freq uint16) error {
	if freq < 3 {
		freq = 3
	}
	if freq > 48000 {
		freq = 48000
	}
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_PWM_FREQUENCY,
		byte(port),
		byte((freq >> 8) & 0xFF),
		byte(freq & 0xFF),
	})
}

// PwmWrite implents the pwm interface for the gopigo3.
func (d *Driver) PwmWrite(pin string, val byte) error {
	var (
		grovePin, grovePort Grove
		err                 error
	)
	if grovePin, grovePort, _, _, err = getGroveAddresses(pin); err != nil {
		return err
	}
	if err := d.SetGroveType(grovePort, CUSTOM); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err = d.SetGroveMode(grovePin, GROVE_OUTPUT_PWM); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err = d.SetPWMFreq(grovePin, 24000); err != nil {
		return err
	}
	val64 := math.Float64frombits(uint64(val))
	dutyCycle := uint16(math.Float64bits((100.0 / 255.0) * val64)) //nolint:gosec // TODO: fix later
	return d.SetPWMDuty(grovePin, dutyCycle)
}

// AnalogRead returns the analog value of the given pin.
func (d *Driver) AnalogRead(pin string) (int, error) {
	grovePin, grovePort, analogCmd, _, err := getGroveAddresses(pin)
	if err != nil {
		return 0, err
	}
	if err := d.SetGroveType(grovePort, CUSTOM); err != nil {
		return 0, err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.SetGroveMode(grovePin, GROVE_INPUT_ANALOG); err != nil {
		return 0, err
	}
	time.Sleep(10 * time.Millisecond)
	response, err := d.readBytes(goPiGo3Address, analogCmd, 7)
	if err != nil {
		return 0, err
	}
	if err := d.responseValid(response); err != nil {
		return 0, err
	}
	if err := d.valueValid(response); err != nil {
		return 0, err
	}
	highBytes := uint16(response[5])
	lowBytes := uint16(response[6])
	return int((highBytes<<8)&0xFF00) | int(lowBytes&0xFF), nil
}

// DigitalWrite writes a 0/1 value to the given pin.
func (d *Driver) DigitalWrite(pin string, val byte) error {
	var (
		grovePin, grovePort Grove
		err                 error
	)
	grovePin, grovePort, _, _, err = getGroveAddresses(pin)
	if err != nil {
		return err
	}
	if err := d.SetGroveType(grovePort, CUSTOM); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := d.SetGroveMode(grovePin, GROVE_OUTPUT_DIGITAL); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	return d.writeBytes([]byte{
		goPiGo3Address,
		SET_GROVE_STATE,
		byte(grovePin),
		val,
	})
}

// DigitalRead reads the 0/1 value from the given pin.
func (d *Driver) DigitalRead(pin string) (int, error) {
	grovePin, grovePort, _, stateCmd, err := getGroveAddresses(pin)
	if err != nil {
		return 0, err
	}
	err = d.SetGroveType(grovePort, CUSTOM)
	if err != nil {
		return 0, err
	}
	time.Sleep(10 * time.Millisecond)
	err = d.SetGroveMode(grovePin, GROVE_INPUT_DIGITAL)
	if err != nil {
		return 0, err
	}
	time.Sleep(10 * time.Millisecond)
	response, err := d.readBytes(goPiGo3Address, stateCmd, 6)
	if err != nil {
		return 0, err
	}
	if err := d.responseValid(response); err != nil {
		return 0, err
	}
	if err := d.valueValid(response); err != nil {
		return 0, err
	}
	return int(response[5]), nil
}

//nolint:nonamedreturns // sufficient here
func getGroveAddresses(pin string) (gPin, gPort Grove, analog, state byte, err error) {
	switch pin {
	case "AD_1_1":
		gPin = AD_1_1_G
		gPort = AD_1_G
		analog = GET_GROVE_ANALOG_1_1
		state = GET_GROVE_STATE_1_1
	case "AD_1_2":
		gPin = AD_1_2_G
		gPort = AD_1_G
		analog = GET_GROVE_ANALOG_1_2
		state = GET_GROVE_STATE_1_2
	case "AD_2_1":
		gPin = AD_2_1_G
		gPort = AD_2_G
		analog = GET_GROVE_ANALOG_2_1
		state = GET_GROVE_STATE_2_1
	case "AD_2_2":
		gPin = AD_2_2_G
		gPort = AD_2_G
		analog = GET_GROVE_ANALOG_2_2
		state = GET_GROVE_STATE_2_2
	default:
		err = fmt.Errorf("Invalid grove pin name")
	}
	return gPin, gPort, analog, state, err
}

func (d *Driver) responseValid(response []byte) error {
	if response[3] != 0xA5 {
		return fmt.Errorf("No SPI response, response not valid")
	}
	return nil
}

func (d *Driver) valueValid(value []byte) error {
	if value[4] != byte(VALID_DATA) {
		return fmt.Errorf("Invalid value")
	}
	return nil
}

func (d *Driver) readBytes(address byte, msg byte, numBytes int) ([]byte, error) {
	w := make([]byte, numBytes)
	w[0] = address
	w[1] = msg
	r := make([]byte, len(w))
	if err := d.connection.ReadCommandData(w, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (d *Driver) readUint16(address, msg byte) (uint16, error) {
	r, err := d.readBytes(address, msg, 8)
	if err != nil {
		return 0, err
	}
	if err := d.responseValid(r); err != nil {
		return 0, err
	}
	return uint16(r[4])<<8 | uint16(r[5]), nil
}

func (d *Driver) readUint32(address, msg byte) (uint32, error) {
	r, err := d.readBytes(address, msg, 8)
	if err != nil {
		return 0, err
	}
	if err := d.responseValid(r); err != nil {
		return 0, err
	}
	return uint32(r[4])<<24 | uint32(r[5])<<16 | uint32(r[6])<<8 | uint32(r[7]), nil
}

func (d *Driver) writeBytes(w []byte) error {
	return d.connection.WriteBytes(w)
}

func (d *Driver) resetAll() error {
	err := d.SetGroveType(AD_1_G+AD_2_G, CUSTOM)
	time.Sleep(10 * time.Millisecond)
	if e := d.SetGroveMode(AD_1_G+AD_2_G, GROVE_INPUT_DIGITAL); e != nil {
		err = multierror.Append(err, e)
	}
	time.Sleep(10 * time.Millisecond)
	if e := d.SetMotorPower(MOTOR_LEFT+MOTOR_RIGHT, 0.0); e != nil {
		err = multierror.Append(err, e)
	}
	if e := d.SetMotorLimits(MOTOR_LEFT+MOTOR_RIGHT, 0, 0); e != nil {
		err = multierror.Append(err, e)
	}
	if e := d.SetServo(SERVO_1+SERVO_2, 0); e != nil {
		err = multierror.Append(err, e)
	}
	if e := d.SetLED(LED_EYE_LEFT+LED_EYE_RIGHT+LED_BLINKER_LEFT+LED_BLINKER_RIGHT, 0, 0, 0); e != nil {
		err = multierror.Append(err, e)
	}

	return err
}
