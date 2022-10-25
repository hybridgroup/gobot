package jetson

import (
	"errors"
	"strings"
	"testing"

	"runtime"
	"strconv"
	"sync"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ sysfs.DigitalPinnerProvider = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)
var _ spi.Connector = (*Adaptor)(nil)

func initTestAdaptor() *Adaptor {

	a := NewAdaptor()
	a.Connect()
	return a
}

func TestJetsonAdaptorName(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "JetsonNano"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestAdaptor(t *testing.T) {

	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "JetsonNano"), true)
	gobottest.Assert(t, a.i2cDefaultBus, 1)
	gobottest.Assert(t, a.revision, "")

}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()

	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/dev/i2c-1",
		"/dev/i2c-0",
		"/dev/spidev0.0",
		"/dev/spidev0.1",
	})

	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	a.DigitalWrite("3", 1)

	a.GetConnection(0xff, 0)
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAdaptorDigitalIO(t *testing.T) {
	a := initTestAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio216/value",
		"/sys/class/gpio/gpio216/direction",
		"/sys/class/gpio/gpio14/value",
		"/sys/class/gpio/gpio14/direction",
	})

	sysfs.SetFilesystem(fs)

	a.DigitalWrite("7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio216/value"].Contents, "1")

	a.DigitalWrite("13", 1)
	i, _ := a.DigitalRead("13")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("notexist", 1), errors.New("Not a valid pin"))

	fs.WithReadError = true
	_, err := a.DigitalRead("7")
	gobottest.Assert(t, err, errors.New("read error"))

	fs.WithWriteError = true
	_, err = a.DigitalRead("13")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorI2c(t *testing.T) {
	a := initTestAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	con, err := a.GetConnection(0xff, 1)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	_, err = a.GetConnection(0xff, 51)
	gobottest.Assert(t, err, errors.New("Bus number 51 out of range"))

	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestAdaptorSPI(t *testing.T) {
	a := initTestAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/spidev0.1",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	gobottest.Assert(t, a.GetSpiDefaultBus(), 0)
	gobottest.Assert(t, a.GetSpiDefaultChip(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMode(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMaxSpeed(), int64(10000000))

	_, err := a.GetSpiConnection(10, 0, 0, 8, 10000000)
	gobottest.Assert(t, err.Error(), "Bus number 10 out of range")

	// TODO: test tx/rx here...
}

func TestAdaptorDigitalPinConcurrency(t *testing.T) {

	oldProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(8)

	for retry := 0; retry < 20; retry++ {

		a := initTestAdaptor()
		var wg sync.WaitGroup

		for i := 0; i < 20; i++ {
			wg.Add(1)
			pinAsString := strconv.Itoa(i)
			go func(pin string) {
				defer wg.Done()
				a.DigitalPin(pin, sysfs.IN)
			}(pinAsString)
		}

		wg.Wait()
	}

	runtime.GOMAXPROCS(oldProcs)

}
