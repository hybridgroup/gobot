// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)
package audio

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func TestAudioAdaptor(t *testing.T) {
	a := NewAdaptor()

	assert.NoError(t, a.Connect())
	assert.NoError(t, a.Finalize())
}

func TestAudioAdaptorName(t *testing.T) {
	a := NewAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Audio"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestAudioAdaptorCommandsWav(t *testing.T) {
	cmd, _ := CommandName("whatever.wav")
	assert.Equal(t, "aplay", cmd)
}

func TestAudioAdaptorCommandsMp3(t *testing.T) {
	cmd, _ := CommandName("whatever.mp3")
	assert.Equal(t, "mpg123", cmd)
}

func TestAudioAdaptorCommandsUnknown(t *testing.T) {
	cmd, err := CommandName("whatever.unk")
	assert.NotEqual(t, "mpg123", cmd)
	assert.ErrorContains(t, err, "Unknown filetype for audio file.")
}

func TestAudioAdaptorSoundWithNoFilename(t *testing.T) {
	a := NewAdaptor()

	errors := a.Sound("")
	assert.Equal(t, "Requires filename for audio file.", errors[0].Error())
}

func TestAudioAdaptorSoundWithNonexistingFilename(t *testing.T) {
	a := NewAdaptor()

	errors := a.Sound("doesnotexist.mp3")
	assert.Equal(t, "stat doesnotexist.mp3: no such file or directory", errors[0].Error())
}

func TestAudioAdaptorSoundWithValidMP3Filename(t *testing.T) {
	execCommand = myExecCommand

	a := NewAdaptor()
	defer func() { execCommand = exec.Command }()

	errors := a.Sound("../../examples/laser.mp3")

	assert.Equal(t, 0, len(errors))
}

func myExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...) //nolint:gosec // ok for test
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}
