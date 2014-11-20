package mavlink

import (
	"io"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

var _ gobot.AdaptorInterface = (*MavlinkAdaptor)(nil)

type MavlinkAdaptor struct {
	gobot.Adaptor
	sp      io.ReadWriteCloser
	connect func(*MavlinkAdaptor) (err error)
}

// NewMavLinkAdaptor creates a new mavlink adaptor with specified name and port
func NewMavlinkAdaptor(name string, port string) *MavlinkAdaptor {
	return &MavlinkAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"mavlink.MavlinkAdaptor",
			port,
		),
		connect: func(m *MavlinkAdaptor) (err error) {
			s, err := serial.OpenPort(&serial.Config{Name: m.Port(), Baud: 57600})
			if err != nil {
				return err
			}
			m.sp = s
			return
		},
	}
}

// Connect returns true if connection to device is successful
func (m *MavlinkAdaptor) Connect() error {
	return m.connect(m)
}

// Finalize returns true if connection to devices is closed successfully
func (m *MavlinkAdaptor) Finalize() error {
	return m.sp.Close()
}
