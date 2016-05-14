// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)
package audio

import (
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func TestAudioAdaptor(t *testing.T) {
	a := NewAudioAdaptor("tester")

	gobottest.Assert(t, a.Name(), "tester")

	gobottest.Assert(t, len(a.Connect()), 0)

	gobottest.Assert(t, len(a.Finalize()), 0)
}
