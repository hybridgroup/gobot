package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var b *BlinkMDriver

func init() {
	b = NewBlinkMDriver(TestAdaptor{}, "bot")
}

func TestBlinkMStart(t *testing.T) {
	gobot.Expect(t, b.Start(), true)
}
