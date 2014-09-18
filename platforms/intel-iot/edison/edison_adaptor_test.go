package edison

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestEdisonAdaptor() *EdisonAdaptor {
	return NewEdisonAdaptor("myAdaptor")
}

func TestEdisonAdaptorConnect(t *testing.T) {
	t.Skip()
	a := initTestEdisonAdaptor()
	gobot.Assert(t, a.Connect(), true)
}

func TestEdisonAdaptorFinalize(t *testing.T) {
	t.Skip()
	a := initTestEdisonAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
