package spark

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var s *SparkCoreAdaptor

func init() {
	s = NewSparkCoreAdaptor("bot", "", "")
}

func TestFinalize(t *testing.T) {
	gobot.Expect(t, s.Finalize(), true)
}
func TestConnect(t *testing.T) {
	gobot.Expect(t, s.Connect(), true)
}
