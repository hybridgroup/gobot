package keyboard

import (
	"os"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Driver = (*KeyboardDriver)(nil)

func initTestKeyboardDriver() *KeyboardDriver {
	d := NewKeyboardDriver("keyboard")
	d.connect = func(k *KeyboardDriver) (err error) {
		k.stdin = &os.File{}
		return nil
	}
	d.listen = func(k *KeyboardDriver) {}
	return d
}

func TestKeyboardDriver(t *testing.T) {
	d := initTestKeyboardDriver()
	gobottest.Assert(t, d.Name(), "keyboard")
	gobottest.Assert(t, d.Connection(), (gobot.Connection)(nil))
}

func TestKeyboardDriverStart(t *testing.T) {
	d := initTestKeyboardDriver()
	gobottest.Assert(t, len(d.Start()), 0)
}

func TestKeyboardDriverHalt(t *testing.T) {
	d := initTestKeyboardDriver()
	gobottest.Assert(t, len(d.Halt()), 0)
}
