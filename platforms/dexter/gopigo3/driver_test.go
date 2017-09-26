package gopigo3

import (
	"encoding/hex"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/gobottest"
	xspi "golang.org/x/exp/io/spi"
)

var _ gobot.Driver = (*Driver)(nil)
var negativeEncoder = false

func initTestDriver() *Driver {
	d := NewDriver(&TestConnector{})
	d.Start()
	return d
}

func TestDriverStart(t *testing.T) {
	d := initTestDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestDriverHalt(t *testing.T) {
	d := initTestDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestDriverManufacturerName(t *testing.T) {
	expectedName := "Dexter Industries"
	d := initTestDriver()
	name, err := d.GetManufacturerName()
	if err != nil {
		t.Error(err)
	}
	if name != expectedName {
		t.Errorf("Expected name: %x, got: %x", expectedName, name)
	}
}

func TestDriverBoardName(t *testing.T) {
	expectedBoardName := "GoPiGo3"
	d := initTestDriver()
	name, err := d.GetBoardName()
	if err != nil {
		t.Error(err)
	}
	if name != expectedBoardName {
		t.Errorf("Expected name: %s, got: %s", expectedBoardName, name)
	}
}

func TestDriverHardwareVersion(t *testing.T) {
	expectedHwVer := "3.1.3"
	d := initTestDriver()
	ver, err := d.GetHardwareVersion()
	if err != nil {
		t.Error(err)
	}
	if ver != expectedHwVer {
		t.Errorf("Expected hw ver: %s, got: %s", expectedHwVer, ver)
	}
}

func TestDriverFirmareVersion(t *testing.T) {
	expectedFwVer := "0.3.4"
	d := initTestDriver()
	ver, err := d.GetFirmwareVersion()
	if err != nil {
		t.Error(err)
	}
	if ver != expectedFwVer {
		t.Errorf("Expected fw ver: %s, got: %s", expectedFwVer, ver)
	}
}

func TestGetSerialNumber(t *testing.T) {
	expectedSerialNumber := "E0180A54514E343732202020FF112137"
	d := initTestDriver()
	serial, err := d.GetSerialNumber()
	if err != nil {
		t.Error(err)
	}
	if serial != expectedSerialNumber {
		t.Errorf("Expected serial number: %s, got: %s", expectedSerialNumber, serial)
	}
}

func TestGetFiveVolts(t *testing.T) {
	expectedVoltage := float32(5.047000)
	d := initTestDriver()
	voltage, err := d.Get5vVoltage()
	if err != nil {
		t.Error(err)
	}
	if voltage != expectedVoltage {
		t.Errorf("Expected 5v voltage: %f, got: %f", expectedVoltage, voltage)
	}
}

func TestGetBatVolts(t *testing.T) {
	expectedBatVoltage := float32(15.411000)
	d := initTestDriver()
	voltage, err := d.GetBatteryVoltage()
	if err != nil {
		t.Error(err)
	}
	if voltage != expectedBatVoltage {
		t.Errorf("Expected battery voltage: %f, got: %f", expectedBatVoltage, voltage)
	}
}

func TestGetMotorStatus(t *testing.T) {
	expectedPower := uint16(65408)
	d := initTestDriver()
	_, power, _, _, err := d.GetMotorStatus(MOTOR_LEFT)
	if err != nil {
		t.Error(err)
	}
	if power != expectedPower {
		t.Errorf("Expected power: %d, got: %d", expectedPower, power)
	}
}

func TestGetEncoderStatusPos(t *testing.T) {
	negativeEncoder = false
	expectedEncoderValue := int64(127)
	d := initTestDriver()
	encoderValue, err := d.GetMotorEncoder(MOTOR_LEFT)
	if err != nil {
		t.Error(err)
	}
	if encoderValue != expectedEncoderValue {
		t.Errorf("Expected encoder value: %d, got: %d", expectedEncoderValue, encoderValue)
	}
}

func TestGetEncoderStatusNeg(t *testing.T) {
	negativeEncoder = true
	expectedEncoderValue := int64(-128)
	d := initTestDriver()
	encoderValue, err := d.GetMotorEncoder(MOTOR_LEFT)
	if err != nil {
		t.Error(err)
	}
	if encoderValue != expectedEncoderValue {
		t.Errorf("Expected encoder value: %d, got: %d", expectedEncoderValue, encoderValue)
	}
}

func TestAnalogRead(t *testing.T) {
	expectedVal := 160
	d := initTestDriver()
	val, err := d.AnalogRead("AD_1_1")
	if err != nil {
		t.Error(err)
	}
	if val != expectedVal {
		t.Errorf("Expected value: %d, got: %d", expectedVal, val)
	}
}

func TestDigitalRead(t *testing.T) {
	expectedVal := 1
	d := initTestDriver()
	val, err := d.DigitalRead("AD_1_2")
	if err != nil {
		t.Error(err)
	}
	if val != expectedVal {
		t.Errorf("Expected value: %d, got: %d", expectedVal, val)
	}
}

func TestServoWrite(t *testing.T) {
	d := initTestDriver()
	err := d.ServoWrite("SERVO_1", 0x7F)
	if err != nil {
		t.Error(err)
	}
}

func TestSetMotorPosition(t *testing.T) {
	d := initTestDriver()
	err := d.SetMotorPosition(MOTOR_LEFT, 12.0)
	if err != nil {
		t.Error(err)
	}
}

func TestSetMotorPower(t *testing.T) {
	d := initTestDriver()
	err := d.SetMotorPower(MOTOR_LEFT, 127)
	if err != nil {
		t.Error(err)
	}
}

func TestSetMotorDps(t *testing.T) {
	d := initTestDriver()
	err := d.SetMotorDps(MOTOR_LEFT, 12.0)
	if err != nil {
		t.Error(err)
	}
}

func TestOffsetMotorEncoder(t *testing.T) {
	d := initTestDriver()
	err := d.OffsetMotorEncoder(MOTOR_LEFT, 12.0)
	if err != nil {
		t.Error(err)
	}
}

func TestSetPWMDuty(t *testing.T) {
	d := initTestDriver()
	err := d.SetPWMDuty(AD_1_1, 80)
	if err != nil {
		t.Error(err)
	}
}

func TestSetPWMfreq(t *testing.T) {
	d := initTestDriver()
	err := d.SetPWMFreq(AD_1_2, 100)
	if err != nil {
		t.Error(err)
	}
}

func TestPwmWrite(t *testing.T) {
	d := initTestDriver()
	err := d.PwmWrite("AD_2_2", 80)
	if err != nil {
		t.Error(err)
	}
}

func TestDigitalWrite(t *testing.T) {
	d := initTestDriver()
	err := d.DigitalWrite("AD_2_1", 1)
	if err != nil {
		t.Error(err)
	}
}

type TestConnector struct{}

func (ctr *TestConnector) GetSpiConnection(busNum, mode int, maxSpeed int64) (device spi.Connection, err error) {
	return spi.NewConnection(&TestSpiDevice{}), nil
}

func (ctr *TestConnector) GetSpiDefaultBus() int {
	return 0
}

func (ctr *TestConnector) GetSpiDefaultMode() int {
	return 0
}

func (ctr *TestConnector) GetSpiDefaultMaxSpeed() int64 {
	return 0
}

type TestSpiDevice struct {
	bus spi.SPIDevice
}

func (c *TestSpiDevice) Close() error {
	return nil
}

func (c *TestSpiDevice) SetBitOrder(o xspi.Order) error {
	return nil
}

func (c *TestSpiDevice) SetBitsPerWord(bits int) error {
	return nil
}

func (c *TestSpiDevice) SetCSChange(leaveEnabled bool) error {
	return nil
}

func (c *TestSpiDevice) SetDelay(t time.Duration) error {
	return nil
}

func (c *TestSpiDevice) SetMaxSpeed(speed int) error {
	return nil
}

func (c *TestSpiDevice) SetMode(mode xspi.Mode) error {
	return nil
}

func (c *TestSpiDevice) Tx(w, r []byte) error {
	manName, _ := hex.DecodeString("ff0000a544657874657220496e6475737472696573000000")
	boardName, _ := hex.DecodeString("ff0000a5476f5069476f3300000000000000000000000000")
	hwVersion, _ := hex.DecodeString("ff0000a5002dcaab")
	fwVersion, _ := hex.DecodeString("ff0000a500000bbc")
	serialNum, _ := hex.DecodeString("ff0000a5e0180a54514e343732202020ff112137")
	fiveVoltVoltage, _ := hex.DecodeString("ff0000a513b7")
	batteryVoltage, _ := hex.DecodeString("ff0000a53c33")
	negMotorEncoder, _ := hex.DecodeString("ff0000a5ffffff00")
	motorEncoder, _ := hex.DecodeString("ff0000a5000000ff")
	motorStatus, _ := hex.DecodeString("ff0000a50080000000000000")
	analogValue, _ := hex.DecodeString("ff0000a50000a0")
	buttonPush, _ := hex.DecodeString("ff0000a50001")
	switch w[1] {
	case 1:
		copy(r, manName)
		return nil
	case 2:
		copy(r, boardName)
		return nil
	case 3:
		copy(r, hwVersion)
		return nil
	case 4:
		copy(r, fwVersion)
		return nil
	case 5:
		copy(r, serialNum)
		return nil
	case 7:
		copy(r, fiveVoltVoltage)
		return nil
	case 8:
		copy(r, batteryVoltage)
		return nil
	case 17:
		if negativeEncoder {
			copy(r, negMotorEncoder)
			return nil
		}
		copy(r, motorEncoder)
		return nil
	case 19:
		copy(r, motorStatus)
		return nil
	case 29:
		copy(r, buttonPush)
		return nil
	case 36:
		copy(r, analogValue)
		return nil
	}
	return nil
}
