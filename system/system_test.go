package system

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccesser(t *testing.T) {
	// arrange & act
	a := NewAccesser()
	// assert
	assert.NotNil(t, a)
	assert.NotNil(t, a.accesserCfg)
	assert.Nil(t, a.sys)
	assert.Nil(t, a.fs)
	assert.Nil(t, a.digitalPinAccess)
	assert.Nil(t, a.spiAccess)
}

func TestAccesserAddAnalogSupport(t *testing.T) {
	// arrange
	a := NewAccesser()
	// act
	a.AddAnalogSupport()
	// assert
	assert.Nil(t, a.sys)
	assert.Nil(t, a.digitalPinAccess)
	assert.Nil(t, a.spiAccess)
	require.NotNil(t, a.fs)
	assert.IsType(t, &nativeFilesystem{}, a.fs)
}

func TestAccesserAddPWMSupport(t *testing.T) {
	// arrange
	a := NewAccesser()
	// act
	a.AddPWMSupport()
	// assert
	assert.Nil(t, a.sys)
	assert.Nil(t, a.digitalPinAccess)
	assert.Nil(t, a.spiAccess)
	require.NotNil(t, a.fs)
	assert.IsType(t, &nativeFilesystem{}, a.fs)
}

func TestAccesserAddDigitalPinSupport(t *testing.T) {
	// arrange
	a := NewAccesser()
	// act
	a.AddDigitalPinSupport()
	// assert
	assert.Nil(t, a.sys)
	assert.Nil(t, a.spiAccess)
	require.NotNil(t, a.fs)
	assert.IsType(t, &nativeFilesystem{}, a.fs)
	require.NotNil(t, a.digitalPinAccess)
	assert.IsType(t, &cdevDigitalPinAccess{}, a.digitalPinAccess)
}

func TestAccesserAddI2CSupport(t *testing.T) {
	// assert
	a := NewAccesser()
	// act
	a.AddI2CSupport()
	// assert
	assert.Nil(t, a.digitalPinAccess)
	assert.Nil(t, a.spiAccess)
	require.NotNil(t, a.fs)
	assert.IsType(t, &nativeFilesystem{}, a.fs)
	require.NotNil(t, a.sys)
	assert.IsType(t, &nativeSyscall{}, a.sys)
}

func TestAccesserAddSPISupport(t *testing.T) {
	// arrange
	a := NewAccesser()
	// act
	a.AddSPISupport() // this writes a message, but we need this test case to assert a.fs
	// assert
	assert.Nil(t, a.sys)
	assert.Nil(t, a.digitalPinAccess)
	require.NotNil(t, a.fs)
	assert.IsType(t, &nativeFilesystem{}, a.fs)
	// arrange for periphio needed
	a.UseMockFilesystem([]string{"/dev/spidev"})
	// act for apply periphio
	a.AddSPISupport()
	// assert
	require.NotNil(t, a.spiAccess)
	assert.IsType(t, &periphioSpiAccess{}, a.spiAccess)
}

func TestAccesserAddOneWireSupport(t *testing.T) {
	// arrange
	a := NewAccesser()
	// act
	a.AddOneWireSupport()
	// assert
	assert.Nil(t, a.sys)
	assert.Nil(t, a.digitalPinAccess)
	assert.Nil(t, a.spiAccess)
	require.NotNil(t, a.fs)
	assert.IsType(t, &nativeFilesystem{}, a.fs)
}

func TestAccesserNewDigitalPin(t *testing.T) {
	// arrange
	const (
		chip = "14"
		line = 13
	)
	a := NewAccesser()
	dpa := a.UseMockDigitalPinAccess()
	// act
	dp := a.NewDigitalPin(chip, line)
	// assert
	assert.Len(t, dpa.pins, 1)
	assert.IsType(t, &digitalPinMock{}, dp)
	myPin := dpa.pins["14_13"]
	assert.Same(t, myPin, dp)
	assert.Equal(t, chip, myPin.chip)
	assert.Equal(t, line, myPin.pin)
}

//nolint:forcetypeassert // ok for this test
func TestAccesserNewPWMPin(t *testing.T) {
	// arrange
	const (
		path         = "path"
		pin          = 12
		polNormIdent = "norm"
		polInvIdent  = "inv"
	)
	a := NewAccesser()
	// act
	pp := a.NewPWMPin(path, pin, polNormIdent, polInvIdent)
	// assert
	assert.IsType(t, &pwmPinSysFs{}, pp)
	assert.Equal(t, path, pp.(*pwmPinSysFs).path)
	assert.Equal(t, strconv.Itoa(pin), pp.(*pwmPinSysFs).pin)
	assert.Equal(t, polNormIdent, pp.(*pwmPinSysFs).polarityNormalIdentifier)
	assert.Equal(t, polInvIdent, pp.(*pwmPinSysFs).polarityInvertedIdentifier)
}

//nolint:forcetypeassert // ok for this test
func TestAccesserNewAnalogPin(t *testing.T) {
	// arrange
	const (
		path       = "path"
		w          = true
		readBufLen = 123
	)
	a := NewAccesser()
	// act
	ap := a.NewAnalogPin(path, w, readBufLen)
	// assert
	assert.IsType(t, &analogPinSysFs{}, ap)
	assert.Equal(t, path, ap.(*analogPinSysFs).sysfsPath)
	assert.True(t, ap.(*analogPinSysFs).r)
	assert.True(t, ap.(*analogPinSysFs).w)
	assert.IsType(t, &sysfsFileAccess{}, ap.(*analogPinSysFs).sfa)
}

func TestAccesserNewSpiDevice(t *testing.T) {
	// arrange
	const (
		busNum   = 15
		chipNum  = 14
		mode     = 13
		bits     = 12
		maxSpeed = int64(11)
	)
	a := NewAccesser()
	spi := a.UseMockSpi()
	// act
	con, err := a.NewSpiDevice(busNum, chipNum, mode, bits, maxSpeed)
	// assert
	require.NoError(t, err)
	assert.NotNil(t, con)
	assert.Equal(t, busNum, spi.busNum)
	assert.Equal(t, chipNum, spi.chipNum)
	assert.Equal(t, mode, spi.mode)
	assert.Equal(t, bits, spi.bits)
	assert.Equal(t, maxSpeed, spi.maxSpeed)
}

//nolint:forcetypeassert // ok for this test
func TestAccesserNewOneWireDevice(t *testing.T) {
	// arrange
	const (
		familyCode   = 255
		serialNumber = 12345678987654321
		wantID       = "ff-2bdc546291f4b1"
	)
	a := NewAccesser()
	// act
	con, err := a.NewOneWireDevice(familyCode, serialNumber)
	// assert
	require.NoError(t, err)
	assert.IsType(t, &onewireDeviceSysfs{}, con)
	assert.Equal(t, wantID, con.(*onewireDeviceSysfs).id)
	assert.Equal(t, "/sys/bus/w1/devices/"+wantID, con.(*onewireDeviceSysfs).sysfsPath)
	assert.IsType(t, &sysfsFileAccess{}, con.(*onewireDeviceSysfs).sfa)
}
