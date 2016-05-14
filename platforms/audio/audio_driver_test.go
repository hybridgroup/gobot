// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)

package audio

import (
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func TestAudioDriver(t *testing.T) {
	d := NewAudioDriver(NewAudioAdaptor("conn"), "dev", nil)

	gobottest.Assert(t, d.Name(), "dev")
	gobottest.Assert(t, d.Connection().Name(), "conn")

	gobottest.Assert(t, len(d.Start()), 0)

	gobottest.Assert(t, len(d.Halt()), 0)
}
