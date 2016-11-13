package keyboard

import (
	"log"
	"os"

	"github.com/hybridgroup/gobot"
)

const (
	// Keyboard event
	Key = "key"
)

type Driver struct {
	name    string
	connect func(*Driver) (err error)
	listen  func(*Driver)
	stdin   *os.File
	gobot.Eventer
}

func NewDriver() *Driver {
	k := &Driver{
		name: "Keyboard",
		connect: func(k *Driver) (err error) {
			if err := configure(); err != nil {
				return err
			}

			k.stdin = os.Stdin
			return
		},
		listen: func(k *Driver) {
			ctrlc := bytes{3}

			for {
				var keybuf bytes
				k.stdin.Read(keybuf[0:3])

				if keybuf == ctrlc {
					proc, err := os.FindProcess(os.Getpid())
					if err != nil {
						log.Fatal(err)
					}

					proc.Signal(os.Interrupt)
					break
				}

				k.Publish(Key, Parse(keybuf))

			}
		},
		Eventer: gobot.NewEventer(),
	}

	k.AddEvent(Key)

	return k
}

func (k *Driver) Name() string                 { return k.name }
func (k *Driver) SetName(n string)             { k.name = n }
func (k *Driver) Connection() gobot.Connection { return nil }

// Start initializes keyboard by grabbing key events as they come in and
// publishing a key event
func (k *Driver) Start() (err error) {
	if err = k.connect(k); err != nil {
		return err
	}

	go k.listen(k)

	return
}

// Halt stops camera driver
func (k *Driver) Halt() (err error) {
	if originalState != "" {
		return restore()
	}
	return
}
