package spark

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestSparkCoreAdaptor() *SparkCoreAdaptor {
	return NewSparkCoreAdaptor("bot", "", "")
}

func TestSparkCoreAdaptorConnect(t *testing.T) {
	a := initTestSparkCoreAdaptor()
	gobot.Assert(t, a.Connect(), true)
}
func TestSparkCoreAdaptorFinalize(t *testing.T) {
	a := initTestSparkCoreAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
