// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)

package audio

import (
	"os/exec"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func TestAudioDriver(t *testing.T) {
	d := NewAudioDriver(NewAudioAdaptor("conn"), "dev", nil)

	gobottest.Assert(t, d.Name(), "dev")
	gobottest.Assert(t, d.Connection().Name(), "conn")

	gobottest.Assert(t, len(d.Start()), 0)

	gobottest.Assert(t, len(d.Halt()), 0)

	_, err := exec.LookPath("mpg123")
	numErrsForTest := 0
	if err != nil {
		numErrsForTest = 1
	}
	gobottest.Assert(t, len(d.Sound("../resources/foo.mp3")), numErrsForTest)
}
