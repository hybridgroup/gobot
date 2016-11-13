// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)
package audio

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
)

type Adaptor struct {
	name string
}

func NewAdaptor() *Adaptor {
	return &Adaptor{name: "Audio"}
}

func (a *Adaptor) Name() string { return a.name }

func (a *Adaptor) SetName(n string) { a.name = n }

func (a *Adaptor) Connect() error { return nil }

func (a *Adaptor) Finalize() error { return nil }

func (a *Adaptor) Sound(fileName string) []error {
	var errorsList []error

	if fileName == "" {
		log.Println("Requires filename for audio file.")
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
	commandName, err := CommandName(fileName)
	if err != nil {
		log.Println(err)
		errorsList = append(errorsList, err)
		return errorsList
	}

	err = RunCommand(commandName, fileName)
	if err != nil {
		log.Println(err)
		errorsList = append(errorsList, err)
		return errorsList
	}

	// Need to return to fulfill function sig, even though returning an empty
	return nil
}

func CommandName(fileName string) (commandName string, err error) {
	fileType := path.Ext(fileName)
	if fileType == ".mp3" {
		return "mpg123", nil
	} else if fileType == ".wav" {
		return "aplay", nil
	} else {
		return "", errors.New("Unknown filetype for audio file.")
	}
}

var execCommand = exec.Command

func RunCommand(audioCommand string, filename string) error {
	cmd := execCommand(audioCommand, filename)
	err := cmd.Start()
	return err
}
