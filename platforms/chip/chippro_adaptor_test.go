package chip

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/system"
)

func initConnectedTestProAdaptorWithMockedFilesystem() (*Adaptor, *system.MockFilesystem) {
	a := NewProAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestNewProAdaptor(t *testing.T) {
	a := NewProAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "CHIP Pro"))
	assert.True(t, a.sys.IsSysfsDigitalPinAccess())
}

func TestProDigitalIO(t *testing.T) {
	a, fs := initConnectedTestProAdaptorWithMockedFilesystem()

	require.NoError(t, a.DigitalWrite("CSID7", 1))
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio139/value"].Contents)

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, err := a.DigitalRead("TWI2-SDA")
	assert.Equal(t, 1, i)
	require.NoError(t, err)

	require.ErrorContains(t, a.DigitalWrite("XIO-P0", 1), "'XIO-P0' is not a valid id for a digital pin")
	require.NoError(t, a.Finalize())
}
