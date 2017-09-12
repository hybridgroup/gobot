package gopigo3

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	xspi "golang.org/x/exp/io/spi"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

type Adaptor struct {
	name               string
	port               string
	spiDefaultMode     xspi.Mode
	spiDefaultMaxSpeed int64
	connection         spi.SPIDevice
	writeBytesChannel  chan []byte
	finalizeChannel    chan struct{}
}

func NewAdaptor() *Adaptor {
	return &Adaptor{
		name:               "GoPiGo3",
		port:               "/dev/spidev0.1", //gopigo3 uses chip select CE1 on raspberry pi
		spiDefaultMode:     xspi.Mode0,
		spiDefaultMaxSpeed: 500000,
		connection:         nil,
		writeBytesChannel:  make(chan []byte),
		finalizeChannel:    make(chan struct{}),
	}
}

func (a *Adaptor) Name() string {
	return a.name
}

func (a *Adaptor) SetName(name string) {
	a.name = name
}

func (a *Adaptor) Connect() error {
	if a.connection == nil {
		devfs := &xspi.Devfs{
			Dev:      a.port,
			Mode:     a.spiDefaultMode,
			MaxSpeed: a.spiDefaultMaxSpeed,
		}
		con, err := xspi.Open(devfs)
		if err != nil {
			return err
		}
		a.connection = con
	}
	go func() {
		for {
			select {
			case bytes := <-a.writeBytesChannel:
				a.connection.Tx(bytes, nil)
				time.Sleep(10 * time.Millisecond)
			case <-a.finalizeChannel:
				a.finalizeChannel <- struct{}{}
				return
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
	return nil
}

func (a *Adaptor) Finalize() error {
	a.finalizeChannel <- struct{}{}
	<-a.finalizeChannel
	if err := a.connection.Close(); err != nil {
		return err
	}
	return nil
}

func (a *Adaptor) ReadBytes(address byte, msg byte, numBytes int) (val []byte, err error) {
	w := make([]byte, numBytes)
	w[0] = address
	w[1] = msg
	r := make([]byte, len(w))
	err = a.connection.Tx(w, r)
	if err != nil {
		return val, err
	}
	return r, nil
}

func (a *Adaptor) ReadUint8(address, msg byte) (val uint8, err error) {
	r, err := a.ReadBytes(address, msg, 8)
	if err != nil {
		return val, err
	}
	return uint8(r[4]) << 8, nil
}

func (a *Adaptor) ReadUint16(address, msg byte) (val uint16, err error) {
	r, err := a.ReadBytes(address, msg, 8)
	if err != nil {
		return val, err
	}
	return uint16(r[4])<<8 | uint16(r[5]), nil
}

func (a *Adaptor) ReadUint32(address, msg byte) (val uint32, err error) {
	r, err := a.ReadBytes(address, msg, 8)
	if err != nil {
		return val, err
	}
	return uint32(r[4])<<24 | uint32(r[5])<<16 | uint32(r[6])<<8 | uint32(r[7]), nil
}

func (a *Adaptor) WriteBytes(w []byte) (err error) {
	return a.connection.Tx(w, nil)
}
