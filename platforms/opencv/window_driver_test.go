//go:build gocv
// +build gocv

package opencv

import (
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gocv.io/x/gocv"
)

var _ gobot.Driver = (*WindowDriver)(nil)

func initTestWindowDriver() *WindowDriver {
	d := NewWindowDriver()
	return d
}

func TestWindowDriver(t *testing.T) {
	d := initTestWindowDriver()
	assert.Equal(t, "Window", d.Name())
	assert.Equal(t, (gobot.Connection)(nil), d.Connection())
}

func TestWindowDriverName(t *testing.T) {
	d := initTestWindowDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Window"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestWindowDriverStart(t *testing.T) {
	d := initTestWindowDriver()
	require.NoError(t, d.Start())
}

func TestWindowDriverHalt(t *testing.T) {
	d := initTestWindowDriver()
	require.NoError(t, d.Halt())
}

func TestWindowDriverShowImage(t *testing.T) {
	d := initTestWindowDriver()
	_, currentfile, _, _ := runtime.Caller(0)
	image := gocv.IMRead(path.Join(path.Dir(currentfile), "lena-256x256.jpg"), gocv.IMReadColor)
	d.Start()
	d.ShowImage(image)
}
