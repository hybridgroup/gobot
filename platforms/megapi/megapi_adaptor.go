package megapi

import (
	"github.com/hybridgroup/gobot"
	"github.com/tarm/serial"
	"io"
	"time"
)

var _ gobot.Adaptor = (*MegaPiAdaptor)(nil)

// MegaPiAdaptor is the Gobot adaptor for the MakeBlock MegaPi board
type MegaPiAdaptor struct {
	name              string
	connection        io.ReadWriteCloser
	serialConfig      *serial.Config
	writeBytesChannel chan []byte
	finalizeChannel   chan struct{}
}

// NewMegaPiAdaptor returns a new MegaPiAdaptor with specified name and specified serial port used to talk to the MegaPi with a baud rate of 115200
func NewMegaPiAdaptor(name string, device string) *MegaPiAdaptor {
	c := &serial.Config{Name: device, Baud: 115200}
	return &MegaPiAdaptor{
		name:              name,
		connection:        nil,
		serialConfig:      c,
		writeBytesChannel: make(chan []byte),
		finalizeChannel:   make(chan struct{}),
	}
}

// Name returns the name of this adaptor
func (megaPi *MegaPiAdaptor) Name() string {
	return megaPi.name
}

// Connect starts a connection to the board
func (megaPi *MegaPiAdaptor) Connect() (errs []error) {
	if megaPi.connection == nil {
		sp, err := serial.OpenPort(megaPi.serialConfig)
		if err != nil {
			return []error{err}
		}

		// sleeping is required to give the board a chance to reset
		time.Sleep(2 * time.Second)
		megaPi.connection = sp
	}

	// kick off thread to send bytes to the board
	go func() {
		for {
			select {
			case bytes := <-megaPi.writeBytesChannel:
				megaPi.connection.Write(bytes)
				time.Sleep(10 * time.Millisecond)
			case <-megaPi.finalizeChannel:
				megaPi.finalizeChannel <- struct{}{}
				return
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
	return
}

// Finalize terminates the connection to the board
func (megaPi *MegaPiAdaptor) Finalize() (errs []error) {
	megaPi.finalizeChannel <- struct{}{}
	<-megaPi.finalizeChannel
	if err := megaPi.connection.Close(); err != nil {
		return []error{err}
	}
	return
}
