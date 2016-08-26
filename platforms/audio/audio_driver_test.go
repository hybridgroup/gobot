// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)

package audio

import (
	"os/exec"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Driver = (*AudioDriver)(nil)

func TestAudioDriver(t *testing.T) {
	d := NewAudioDriver(NewAudioAdaptor("conn"), "dev", "../../examples/laser.mp3")

	gobottest.Assert(t, d.Name(), "dev")
	gobottest.Assert(t, d.Filename(), "../../examples/laser.mp3")

	gobottest.Assert(t, d.Connection().Name(), "conn")

	gobottest.Assert(t, len(d.Start()), 0)

	gobottest.Assert(t, len(d.Halt()), 0)
}

func TestAudioDriverSoundWithNoFilename(t *testing.T) {
	d := NewAudioDriver(NewAudioAdaptor("conn"), "dev", "")

	errors := d.Sound("")
	gobottest.Assert(t, errors[0].Error(), "Requires filename for audio file.")
}

func TestAudioDriverSoundWithDefaultFilename(t *testing.T) {
	execCommand = gobottest.ExecCommand
	defer func() { execCommand = exec.Command }()

	d := NewAudioDriver(NewAudioAdaptor("conn"), "dev", "../../examples/laser.mp3")

	errors := d.Play()
	gobottest.Assert(t, len(errors), 0)
}
