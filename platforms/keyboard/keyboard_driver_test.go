package keyboard

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestKeyboardDriver() *Driver {
	d := NewDriver()
	d.connect = func(k *Driver) error {
		k.stdin = &os.File{}
		return nil
	}
	d.listen = func(k *Driver) {}
	return d
}

func TestKeyboardDriver(t *testing.T) {
	d := initTestKeyboardDriver()
	assert.Equal(t, (gobot.Connection)(nil), d.Connection())
}

func TestKeyboardDriverName(t *testing.T) {
	d := initTestKeyboardDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Keyboard"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestKeyboardDriverStart(t *testing.T) {
	d := initTestKeyboardDriver()
	require.NoError(t, d.Start())
}

func TestKeyboardDriverHalt(t *testing.T) {
	d := initTestKeyboardDriver()
	require.NoError(t, d.Halt())
}
