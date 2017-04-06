// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)

package audio

import (
	"os/exec"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func TestAudioDriver(t *testing.T) {
	d := NewDriver(NewAdaptor(), "../../examples/laser.mp3")

	gobottest.Assert(t, d.Filename(), "../../examples/laser.mp3")

	gobottest.Refute(t, d.Connection(), nil)

	gobottest.Assert(t, d.Start(), nil)

	gobottest.Assert(t, d.Halt(), nil)
}

func TestAudioDriverName(t *testing.T) {
	d := NewDriver(NewAdaptor(), "")
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Audio"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestAudioDriverSoundWithNoFilename(t *testing.T) {
	d := NewDriver(NewAdaptor(), "")

	errors := d.Sound("")
	gobottest.Assert(t, errors[0].Error(), "Requires filename for audio file.")
}

func TestAudioDriverSoundWithDefaultFilename(t *testing.T) {
	execCommand = gobottest.ExecCommand
	defer func() { execCommand = exec.Command }()

	d := NewDriver(NewAdaptor(), "../../examples/laser.mp3")

	errors := d.Play()
	gobottest.Assert(t, len(errors), 0)
}
