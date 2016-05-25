// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)
package audio

import (
	"os"
	"os/exec"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func TestAudioAdaptor(t *testing.T) {
	a := NewAudioAdaptor("tester")

	gobottest.Assert(t, a.Name(), "tester")

	gobottest.Assert(t, len(a.Connect()), 0)

	gobottest.Assert(t, len(a.Finalize()), 0)
}

func TestAudioAdaptorCommandsWav(t *testing.T) {
	cmd, _ := CommandName("whatever.wav")
	gobottest.Assert(t, cmd, "aplay")
}

func TestAudioAdaptorCommandsMp3(t *testing.T) {
	cmd, _ := CommandName("whatever.mp3")
	gobottest.Assert(t, cmd, "mpg123")
}

func TestAudioAdaptorCommandsUnknown(t *testing.T) {
	cmd, err := CommandName("whatever.unk")
	gobottest.Refute(t, cmd, "mpg123")
	gobottest.Assert(t, err.Error(), "Unknown filetype for audio file.")
}

func TestAudioAdaptorSoundWithNoFilename(t *testing.T) {
	a := NewAudioAdaptor("tester")

	errors := a.Sound("")
	gobottest.Assert(t, errors[0].Error(), "Requires filename for audio file.")
}

func TestAudioAdaptorSoundWithNonexistingFilename(t *testing.T) {
	a := NewAudioAdaptor("tester")

	errors := a.Sound("doesnotexist.mp3")
	gobottest.Assert(t, errors[0].Error(), "stat doesnotexist.mp3: no such file or directory")
}

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestAudioAdaptorSoundWithValidMP3Filename(t *testing.T) {
	execCommand = fakeExecCommand

	a := NewAudioAdaptor("tester")
	defer func() { execCommand = exec.Command }()

	errors := a.Sound("../../examples/laser.mp3")

	gobottest.Assert(t, len(errors), 0)
}
