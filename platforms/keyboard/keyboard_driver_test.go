package keyboard

import (
	"os"
	"testing"

	"github.com/hybridgroup/gobot"
)

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
	gobot.Assert(t, d.Name(), "keyboard")
	gobot.Assert(t, d.Connection(), (gobot.Connection)(nil))
}

func TestKeyboardDriverStart(t *testing.T) {
	d := initTestKeyboardDriver()
	gobot.Assert(t, len(d.Start()), 0)
}

func TestKeyboardDriverHalt(t *testing.T) {
	d := initTestKeyboardDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}
