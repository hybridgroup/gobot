// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)
package audio

import (
	"os/exec"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func TestAudioAdaptor(t *testing.T) {
	a := NewAudioAdaptor("tester")

	gobottest.Assert(t, a.Name(), "tester")

	gobottest.Assert(t, len(a.Connect()), 0)

	_, err := exec.LookPath("mpg123")
	numErrsForTest := 0
	if err != nil {
		numErrsForTest = 1
	}
	gobottest.Assert(t, len(a.Sound("../resources/foo.wav")), numErrsForTest)

	gobottest.Assert(t, len(a.Connect()), 0)

	gobottest.Assert(t, len(a.Finalize()), 0)
}
