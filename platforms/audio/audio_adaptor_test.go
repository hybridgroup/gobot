// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)
package audio

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func TestAudioAdaptor(t *testing.T) {
	a := NewAdaptor()

	require.NoError(t, a.Connect())
	require.NoError(t, a.Finalize())
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
	require.ErrorContains(t, err, "unknown filetype for audio file")
}

func TestAudioAdaptorSoundWithNoFilename(t *testing.T) {
	a := NewAdaptor()

	errors := a.Sound("")
	require.ErrorContains(t, errors[0], "requires filename for audio file")
}

func TestAudioAdaptorSoundWithNonexistingFilename(t *testing.T) {
	a := NewAdaptor()

	errors := a.Sound("doesnotexist.mp3")
	require.ErrorContains(t, errors[0], "stat doesnotexist.mp3: no such file or directory")
}

func TestAudioAdaptorSoundWithValidMP3Filename(t *testing.T) {
	execCommand = myExecCommand

	a := NewAdaptor()
	defer func() { execCommand = exec.Command }()

	errors := a.Sound("../../examples/laser.mp3")

	assert.Empty(t, errors)
}

func myExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...) //nolint:gosec // ok for test
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}
