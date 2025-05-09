package adaptors

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/spi"
	"gobot.io/x/gobot/v2/system"
)

type spiBusNumberValidator func(busNumber int) error

// spiBusConfiguration contains all changeable attributes of the adaptor.
type spiBusConfiguration struct {
	debug                 bool
	spiGpioPinnerProvider gobot.DigitalPinnerProvider
	systemOptions         []system.AccesserOptionApplier
}

// SpiBusAdaptor is a adaptor for SPI bus, normally used for composition in platforms.
type SpiBusAdaptor struct {
	sys               *system.Accesser
	validateBusNumber spiBusNumberValidator
	defaultBusNumber  int
	defaultChipNumber int
	defaultMode       int
	defaultBitCount   int
	defaultMaxSpeed   int64 // Hz
	spiBusCfg         *spiBusConfiguration
	mutex             sync.Mutex
	connections       map[string]spi.Connection
}

// NewSpiBusAdaptor provides the access to SPI buses of the board. The validator is used to check the
// bus number (given by user) to the abilities of the board.
func NewSpiBusAdaptor(
	sys *system.Accesser,
	v spiBusNumberValidator,
	busNum, chipNum, mode, bits int,
	maxSpeed int64,
	spiGpioPinnerProvider gobot.DigitalPinnerProvider,
	opts ...SpiBusOptionApplier,
) *SpiBusAdaptor {
	a := SpiBusAdaptor{
		sys:               sys,
		validateBusNumber: v,
		defaultBusNumber:  busNum,
		defaultChipNumber: chipNum,
		defaultMode:       mode,
		defaultBitCount:   bits,
		defaultMaxSpeed:   maxSpeed,
		spiBusCfg:         &spiBusConfiguration{spiGpioPinnerProvider: spiGpioPinnerProvider},
	}

	for _, o := range opts {
		o.apply(a.spiBusCfg)
	}

	sys.AddSPISupport(a.spiBusCfg.systemOptions...)

	return &a
}

// WithSpiDebug can be used to switch on debugging for SPI implementation.
func WithSpiDebug() spiBusDebugOption {
	return spiBusDebugOption(true)
}

// WithSpiGpioAccess can be used to switch the default SPI implementation to GPIO usage.
func WithSpiGpioAccess(sclkPin, ncsPin, sdoPin, sdiPin string) spiBusDigitalPinsForSystemSpiOption {
	o := spiBusDigitalPinsForSystemSpiOption{
		sclkPin: sclkPin,
		ncsPin:  ncsPin,
		sdoPin:  sdoPin,
		sdiPin:  sdiPin,
	}

	return o
}

// Connect prepares the connection to SPI buses.
func (a *SpiBusAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.spiBusCfg.debug {
		fmt.Println("connect the SPI bus adaptor")
	}

	a.connections = make(map[string]spi.Connection)
	return nil
}

// Finalize closes all SPI connections.
func (a *SpiBusAdaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.spiBusCfg.debug {
		fmt.Println("finalize the SPI bus adaptor")
	}

	var err error
	for _, con := range a.connections {
		if con != nil {
			if e := con.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	a.connections = nil
	return err
}

// GetSpiConnection returns an spi connection to a device on a specified bus.
// Valid bus numbers range between 0 and 65536, valid chip numbers are 0 ... 255.
func (a *SpiBusAdaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (spi.Connection, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.spiBusCfg.debug {
		fmt.Println("get SPI connection")
	}

	if a.connections == nil {
		return nil, fmt.Errorf("not connected")
	}

	id := fmt.Sprintf("%d_%d", busNum, chipNum)

	con := a.connections[id]
	if con == nil {
		if err := a.validateBusNumber(busNum); err != nil {
			return nil, err
		}
		var err error
		bus, err := a.sys.NewSpiDevice(busNum, chipNum, mode, bits, maxSpeed)
		if err != nil {
			return nil, err
		}
		con = spi.NewConnection(bus)
		a.connections[id] = con
	}

	return con, nil
}

// SpiDefaultBusNumber returns the default bus number for this platform.
func (a *SpiBusAdaptor) SpiDefaultBusNumber() int {
	return a.defaultBusNumber
}

// SpiDefaultChipNumber returns the default chip number for this platform.
func (a *SpiBusAdaptor) SpiDefaultChipNumber() int {
	return a.defaultChipNumber
}

// SpiDefaultMode returns the default SPI mode for this platform.
func (a *SpiBusAdaptor) SpiDefaultMode() int {
	return a.defaultMode
}

// SpiDefaultBitCount returns the default number of bits used for this platform.
func (a *SpiBusAdaptor) SpiDefaultBitCount() int {
	return a.defaultBitCount
}

// SpiDefaultMaxSpeed returns the default maximal speed for this platform.
func (a *SpiBusAdaptor) SpiDefaultMaxSpeed() int64 {
	return a.defaultMaxSpeed
}
