// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)

package audio

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func TestAudioDriver(t *testing.T) {
	d := NewDriver(NewAdaptor(), "../../examples/laser.mp3")

	assert.Equal(t, "../../examples/laser.mp3", d.Filename())

	assert.NotNil(t, d.Connection())

	require.NoError(t, d.Start())

	require.NoError(t, d.Halt())
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
	require.ErrorContains(t, errors[0], "requires filename for audio file")
}

func TestAudioDriverSoundWithDefaultFilename(t *testing.T) {
	execCommand = myExecCommand
	defer func() { execCommand = exec.Command }()

	d := NewDriver(NewAdaptor(), "../../examples/laser.mp3")

	errors := d.Play()
	assert.Empty(t, errors)
}
