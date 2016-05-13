// Based on aplay audio adaptor written by @colemanserious (https://github.com/colemanserious)

package audio

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"time"
)

var _ gobot.Driver = (*AudioDriver)(nil)

type AudioDriver struct {
	name       string
	connection gobot.Connection
	interval   time.Duration
	halt       chan bool
	gobot.Eventer
	gobot.Commander
	queue chan string
}

func NewAudioDriver(a *AudioAdaptor, name string, queue chan string) *AudioDriver {
	d := &AudioDriver{
		name:       name,
		connection: a,
		interval:   500 * time.Millisecond,
		queue:      queue,
		halt:       make(chan bool, 0),
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
	}
	return d
}

func (d *AudioDriver) Name() string { return d.name }

func (d *AudioDriver) Connection() gobot.Connection {
	return d.connection
}

func (d *AudioDriver) Sound(fileName string) []error {
	return d.Connection().(*AudioAdaptor).Sound(fileName)
}

func (d *AudioDriver) adaptor() *AudioAdaptor {
	return d.Connection().(*AudioAdaptor)
}

func (d *AudioDriver) Start() (err []error) {
	go d.serve(d.queue)
	return
}

func (d *AudioDriver) Halt() (err []error) {
	return
}

// Use semaphore to control how many sounds might be playing at a time
var sem = make(chan int, 1)

// See example at https://golang.org/doc/effective_go.html#concurrency
// Purpose: receive messages on channel, but throttle execution of playing
func (d *AudioDriver) serve(queue chan string) {
	for req := range queue {
		sem <- 1
		go func(req string) {
			fmt.Printf("Playing sound %v\n", req)
			d.Sound(req)
			<-sem
		}(req)
	}
}
