// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)
package audio

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Adaptor = (*AudioAdaptor)(nil)

type AudioAdaptor struct {
	name string
}

func NewAudioAdaptor(name string) *AudioAdaptor {
	return &AudioAdaptor{
		name: name,
	}
}

func (a *AudioAdaptor) Name() string { return a.name }

func (a *AudioAdaptor) Connect() []error { return nil }

func (a *AudioAdaptor) Finalize() []error { return nil }

func (a *AudioAdaptor) Sound(fileName string) []error {

	var errorsList []error

	if fileName == "" {
		log.Println("Require filename for audio file.")
		errorsList = append(errorsList, errors.New("Requires filename for audio file."))
		return errorsList
	}

	_, err := os.Stat(fileName)
	if err != nil {
		log.Println(err)
		errorsList = append(errorsList, err)
		return errorsList
	}

	// command to play audio file based on file type
	fileType := path.Ext(fileName)
	var commandName string
	if fileType == ".mp3" {
		commandName = "mpg123"
	} else if fileType == ".wav" {
		commandName = "aplay"
	} else {
		log.Println("Unknown filetype for audio file.")
		errorsList = append(errorsList, errors.New("Unknown filetype for audio file."))
		return errorsList
	}

	cmd := exec.Command(commandName, fileName)
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		errorsList = append(errorsList, err)
		return errorsList
	}

	// Need to return to fulfill function sig, even though returning an empty
	return nil
}
