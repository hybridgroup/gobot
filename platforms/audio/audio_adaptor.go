// Package audio is based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)
package audio

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"

	"gobot.io/x/gobot"
)

// Adaptor is gobot Adaptor connection to audio playback
type Adaptor struct {
	name string
}

// NewAdaptor returns a new audio Adaptor
//
func NewAdaptor() *Adaptor {
	return &Adaptor{name: gobot.DefaultName("Audio")}
}

// Name returns the Adaptor Name
func (a *Adaptor) Name() string { return a.name }

// SetName sets the Adaptor Name
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect establishes a connection to the Audio adaptor
func (a *Adaptor) Connect() error { return nil }

// Finalize terminates the connection to the Audio adaptor
func (a *Adaptor) Finalize() error { return nil }

// Sound plays a sound and accepts:
//
//  string: The filename of the audio to start playing
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

// CommandName defines the playback command for a sound and accepts:
//
//  string: The filename of the audio that needs playback
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

// RunCommand executes the playback command for a sound file and accepts:
//
//  string: The audio command to be use for playback
//  string: The filename of the audio that needs playback
func RunCommand(audioCommand string, filename string) error {
	cmd := execCommand(audioCommand, filename)
	err := cmd.Start()
	return err
}
