// Package audio is based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)
package audio

import (
	"time"

	"gobot.io/x/gobot"
)

// Driver is gobot software device for audio playback
type Driver struct {
	name       string
	connection gobot.Connection
	interval   time.Duration
	halt       chan bool
	gobot.Eventer
	gobot.Commander
	filename string
}

// NewDriver returns a new audio Driver. It accepts:
//
// *Adaptor: The audio adaptor to use for the driver
//  string: The filename of the audio to start playing
//
func NewDriver(a *Adaptor, filename string) *Driver {
	return &Driver{
		name:       gobot.DefaultName("Audio"),
		connection: a,
		interval:   500 * time.Millisecond,
		filename:   filename,
		halt:       make(chan bool, 0),
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
	}
}

// Name returns the Driver Name
func (d *Driver) Name() string { return d.name }

// SetName sets the Driver Name
func (d *Driver) SetName(n string) { d.name = n }

// Filename returns the file name for the driver to playback
func (d *Driver) Filename() string { return d.filename }

// Connection returns the Driver Connection
func (d *Driver) Connection() gobot.Connection {
	return d.connection
}

// Sound plays back a sound file. It accepts:
//
//  string: The filename of the audio to start playing
func (d *Driver) Sound(fileName string) []error {
	return d.Connection().(*Adaptor).Sound(fileName)
}

// Play plays back the current sound file.
func (d *Driver) Play() []error {
	return d.Sound(d.Filename())
}

func (d *Driver) adaptor() *Adaptor {
	return d.Connection().(*Adaptor)
}

// Start starts the Driver
func (d *Driver) Start() (err error) {
	return
}

// Halt halts the Driver
func (d *Driver) Halt() (err error) {
	return
}
