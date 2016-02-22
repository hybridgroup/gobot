package gobot

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func TestConnectionEach(t *testing.T) {
	r := newTestRobot("Robot1")

	i := 0
	r.Connections().Each(func(conn Connection) {
		i++
	})
	gobottest.Assert(t, r.Connections().Len(), i)
}

func initTestGobot() *Gobot {
	log.SetOutput(&NullReadWriteCloser{})
	g := NewGobot()
	g.trap = func(c chan os.Signal) {
		c <- os.Interrupt
	}
	g.AddRobot(newTestRobot("Robot1"))
	g.AddRobot(newTestRobot("Robot2"))
	g.AddRobot(newTestRobot(""))
	return g
}

func TestVersion(t *testing.T) {
	gobottest.Assert(t, version, Version())
}

func TestNullReadWriteCloser(t *testing.T) {
	n := &NullReadWriteCloser{}
	i, _ := n.Write([]byte{1, 2, 3})
	gobottest.Assert(t, i, 3)
	i, _ = n.Read(make([]byte, 10))
	gobottest.Assert(t, i, 10)
	gobottest.Assert(t, n.Close(), nil)
}

func TestGobotRobot(t *testing.T) {
	g := initTestGobot()
	gobottest.Assert(t, g.Robot("Robot1").Name, "Robot1")
	gobottest.Assert(t, g.Robot("Robot4"), (*Robot)(nil))
	gobottest.Assert(t, g.Robot("Robot4").Device("Device1"), (Device)(nil))
	gobottest.Assert(t, g.Robot("Robot4").Connection("Connection1"), (Connection)(nil))
	gobottest.Assert(t, g.Robot("Robot1").Device("Device4"), (Device)(nil))
	gobottest.Assert(t, g.Robot("Robot1").Device("Device1").Name(), "Device1")
	gobottest.Assert(t, g.Robot("Robot1").Devices().Len(), 3)
	gobottest.Assert(t, g.Robot("Robot1").Connection("Connection4"), (Connection)(nil))
	gobottest.Assert(t, g.Robot("Robot1").Connections().Len(), 3)
}

func TestGobotToJSON(t *testing.T) {
	g := initTestGobot()
	g.AddCommand("test_function", func(params map[string]interface{}) interface{} {
		return nil
	})
	json := NewJSONGobot(g)
	gobottest.Assert(t, len(json.Robots), g.Robots().Len())
	gobottest.Assert(t, len(json.Commands), len(g.Commands()))
}

func TestGobotStart(t *testing.T) {
	g := initTestGobot()
	gobottest.Assert(t, len(g.Start()), 0)
	gobottest.Assert(t, len(g.Stop()), 0)
}

func TestGobotStartErrors(t *testing.T) {
	log.SetOutput(&NullReadWriteCloser{})
	g := NewGobot()

	adaptor1 := newTestAdaptor("Connection1", "/dev/null")
	driver1 := newTestDriver(adaptor1, "Device1", "0")
	r := NewRobot("Robot1",
		[]Connection{adaptor1},
		[]Device{driver1},
	)

	g.AddRobot(r)

	testDriverStart = func() (errs []error) {
		return []error{
			errors.New("driver start error 1"),
		}
	}

	gobottest.Assert(t, len(g.Start()), 1)
	gobottest.Assert(t, len(g.Stop()), 0)

	testDriverStart = func() (errs []error) { return }
	testAdaptorConnect = func() (errs []error) {
		return []error{
			errors.New("adaptor start error 1"),
		}
	}

	gobottest.Assert(t, len(g.Start()), 1)
	gobottest.Assert(t, len(g.Stop()), 0)

	testDriverStart = func() (errs []error) { return }
	testAdaptorConnect = func() (errs []error) { return }

	g.trap = func(c chan os.Signal) {
		c <- os.Interrupt
	}

	testDriverHalt = func() (errs []error) {
		return []error{
			errors.New("driver halt error 1"),
		}
	}

	testAdaptorFinalize = func() (errs []error) {
		return []error{
			errors.New("adaptor finalize error 1"),
		}
	}

	gobottest.Assert(t, len(g.Start()), 0)
	gobottest.Assert(t, len(g.Stop()), 2)
}
