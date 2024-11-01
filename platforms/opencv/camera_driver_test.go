//go:build gocv
// +build gocv

package opencv

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*CameraDriver)(nil)

func initTestCameraDriver() *CameraDriver {
	d := NewCameraDriver("")
	d.start = func(c *CameraDriver) error {
		d.camera = &testCapture{}
		return nil
	}
	return d
}

func TestCameraDriver(t *testing.T) {
	d := initTestCameraDriver()
	assert.Equal(t, "Camera", d.Name())
	assert.Equal(t, (gobot.Connection)(nil), d.Connection())
}

func TestCameraDriverName(t *testing.T) {
	d := initTestCameraDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Camera"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestCameraDriverStart(t *testing.T) {
	sem := make(chan bool)
	d := initTestCameraDriver()
	require.NoError(t, d.Start())
	d.On(d.Event("frame"), func(data interface{}) {
		sem <- true
	})
	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "Event \"frame\" was not published")
	}

	d = NewCameraDriver("")
	require.NoError(t, d.Start())

	d = NewCameraDriver(true)
	assert.NotNil(t, d.Start())
}

func TestCameraDriverHalt(t *testing.T) {
	d := initTestCameraDriver()
	require.NoError(t, d.Halt())
}
