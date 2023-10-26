// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)

package audio

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func TestAudioDriver(t *testing.T) {
	d := NewDriver(NewAdaptor(), "../../examples/laser.mp3")

	assert.Equal(t, "../../examples/laser.mp3", d.Filename())

	assert.NotNil(t, d.Connection())

	assert.NoError(t, d.Start())

	assert.NoError(t, d.Halt())
}

func TestAudioDriverName(t *testing.T) {
	d := NewDriver(NewAdaptor(), "")
	assert.True(t, strings.HasPrefix(d.Name(), "Audio"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestAudioDriverSoundWithNoFilename(t *testing.T) {
	d := NewDriver(NewAdaptor(), "")

	errors := d.Sound("")
	assert.Equal(t, "Requires filename for audio file.", errors[0].Error())
}

func TestAudioDriverSoundWithDefaultFilename(t *testing.T) {
	execCommand = myExecCommand
	defer func() { execCommand = exec.Command }()

	d := NewDriver(NewAdaptor(), "../../examples/laser.mp3")

	errors := d.Play()
	assert.Equal(t, 0, len(errors))
}
