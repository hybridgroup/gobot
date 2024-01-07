package gopigo3

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/spi"
)

var (
	_               gobot.Driver = (*Driver)(nil)
	negativeEncoder              = false
)

func initTestDriver() *Driver {
	d := NewDriver(&TestConnector{})
	_ = d.Start()
	return d
}

func TestDriverStart(t *testing.T) {
	d := initTestDriver()
	require.NoError(t, d.Start())
}

func TestDriverHalt(t *testing.T) {
	d := initTestDriver()
	require.NoError(t, d.Halt())
}

func TestDriverManufacturerName(t *testing.T) {
	wantName := "Dexter Industries"
	d := initTestDriver()
	name, err := d.GetManufacturerName()
	require.NoError(t, err)
	assert.Equal(t, wantName, name)
}

func TestDriverBoardName(t *testing.T) {
	wantBoardName := "GoPiGo3"
	d := initTestDriver()
	name, err := d.GetBoardName()
	require.NoError(t, err)
	assert.Equal(t, wantBoardName, name)
}

func TestDriverHardwareVersion(t *testing.T) {
	wantHwVer := "3.1.3"
	d := initTestDriver()
	ver, err := d.GetHardwareVersion()
	require.NoError(t, err)
	assert.Equal(t, wantHwVer, ver)
}

func TestDriverFirmareVersion(t *testing.T) {
	wantFwVer := "0.3.4"
	d := initTestDriver()
	ver, err := d.GetFirmwareVersion()
	require.NoError(t, err)
	assert.Equal(t, wantFwVer, ver)
}

func TestGetSerialNumber(t *testing.T) {
	wantSerialNumber := "E0180A54514E343732202020FF112137"
	d := initTestDriver()
	serial, err := d.GetSerialNumber()
	require.NoError(t, err)
	assert.Equal(t, wantSerialNumber, serial)
}

func TestGetFiveVolts(t *testing.T) {
	wantVoltage := float32(5.047000)
	d := initTestDriver()
	voltage, err := d.Get5vVoltage()
	require.NoError(t, err)
	assert.InDelta(t, wantVoltage, voltage, 0.0)
}

func TestGetBatVolts(t *testing.T) {
	wantBatVoltage := float32(15.411000)
	d := initTestDriver()
	voltage, err := d.GetBatteryVoltage()
	require.NoError(t, err)
	assert.InDelta(t, wantBatVoltage, voltage, 0.0)
}

func TestGetMotorStatus(t *testing.T) {
	wantPower := uint16(65408)
	d := initTestDriver()
	f1, power, f2, f3, err := d.GetMotorStatus(MOTOR_LEFT)
	require.NoError(t, err)
	assert.Equal(t, wantPower, power)
	assert.Equal(t, uint8(0), f1)
	assert.Equal(t, 0, f2)
	assert.Equal(t, 0, f3)
}

func TestGetEncoderStatusPos(t *testing.T) {
	negativeEncoder = false
	wantEncoderValue := int64(127)
	d := initTestDriver()
	encoderValue, err := d.GetMotorEncoder(MOTOR_LEFT)
	require.NoError(t, err)
	assert.Equal(t, wantEncoderValue, encoderValue)
}

func TestGetEncoderStatusNeg(t *testing.T) {
	negativeEncoder = true
	wantEncoderValue := int64(-128)
	d := initTestDriver()
	encoderValue, err := d.GetMotorEncoder(MOTOR_LEFT)
	require.NoError(t, err)
	assert.Equal(t, wantEncoderValue, encoderValue)
}

func TestAnalogRead(t *testing.T) {
	wantVal := 160
	d := initTestDriver()
	val, err := d.AnalogRead("AD_1_1")
	require.NoError(t, err)
	assert.Equal(t, wantVal, val)
}

func TestDigitalRead(t *testing.T) {
	wantVal := 1
	d := initTestDriver()
	val, err := d.DigitalRead("AD_1_2")
	require.NoError(t, err)
	assert.Equal(t, wantVal, val)
}

func TestServoWrite(t *testing.T) {
	d := initTestDriver()
	err := d.ServoWrite("SERVO_1", 0x7F)
	require.NoError(t, err)
}

func TestSetMotorPosition(t *testing.T) {
	d := initTestDriver()
	err := d.SetMotorPosition(MOTOR_LEFT, 12.0)
	require.NoError(t, err)
}

func TestSetMotorPower(t *testing.T) {
	d := initTestDriver()
	err := d.SetMotorPower(MOTOR_LEFT, 127)
	require.NoError(t, err)
}

func TestSetMotorDps(t *testing.T) {
	d := initTestDriver()
	err := d.SetMotorDps(MOTOR_LEFT, 12.0)
	require.NoError(t, err)
}

func TestOffsetMotorEncoder(t *testing.T) {
	d := initTestDriver()
	err := d.OffsetMotorEncoder(MOTOR_LEFT, 12.0)
	require.NoError(t, err)
}

func TestSetPWMDuty(t *testing.T) {
	d := initTestDriver()
	err := d.SetPWMDuty(AD_1_1_G, 80)
	require.NoError(t, err)
}

func TestSetPWMfreq(t *testing.T) {
	d := initTestDriver()
	err := d.SetPWMFreq(AD_1_2_G, 100)
	require.NoError(t, err)
}

func TestPwmWrite(t *testing.T) {
	d := initTestDriver()
	err := d.PwmWrite("AD_2_2", 80)
	require.NoError(t, err)
}

func TestDigitalWrite(t *testing.T) {
	d := initTestDriver()
	err := d.DigitalWrite("AD_2_1", 1)
	require.NoError(t, err)
}

type TestConnector struct{}

func (ctr *TestConnector) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (spi.Connection, error) {
	return TestSpiDevice{}, nil
}

func (ctr *TestConnector) SpiDefaultBusNumber() int {
	return 0
}

func (ctr *TestConnector) SpiDefaultChipNumber() int {
	return 0
}

func (ctr *TestConnector) SpiDefaultMode() int {
	return 0
}

func (ctr *TestConnector) SpiDefaultBitCount() int {
	return 8
}

func (ctr *TestConnector) SpiDefaultMaxSpeed() int64 {
	return 0
}

type TestSpiDevice struct{}

func (c TestSpiDevice) Close() error {
	return nil
}

func (c TestSpiDevice) ReadByteData(byte) (byte, error)   { return 0, nil }
func (c TestSpiDevice) ReadBlockData(byte, []byte) error  { return nil }
func (c TestSpiDevice) WriteByte(byte) error              { return nil }
func (c TestSpiDevice) WriteByteData(byte, byte) error    { return nil }
func (c TestSpiDevice) WriteBlockData(byte, []byte) error { return nil }
func (c TestSpiDevice) WriteBytes([]byte) error           { return nil }

func (c TestSpiDevice) ReadCommandData(w, r []byte) error {
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
