package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestBlinkMDriver() *BlinkMDriver {
	return NewBlinkMDriver(TestAdaptor{}, "bot")
}

func TestBlinkMDriverStart(t *testing.T) {
	d := initTestBlinkMDriver()
	gobot.Expect(t, d.Start(), true)
}
