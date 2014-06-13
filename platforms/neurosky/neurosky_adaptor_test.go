package neurosky

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var adaptor *NeuroskyAdaptor

func init() {
	adaptor = NewNeuroskyAdaptor("bot", "/dev/null")
}

func TestFinalize(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, adaptor.Finalize(), true)
}
func TestConnect(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, adaptor.Connect(), true)
}
