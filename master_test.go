package gobot

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/gobottest"
)

func initTestMaster() *Master {
	log.SetOutput(&NullReadWriteCloser{})
	g := NewMaster()
	g.trap = func(c chan os.Signal) {
		c <- os.Interrupt
	}
	g.AddRobot(newTestRobot("Robot1"))
	g.AddRobot(newTestRobot("Robot2"))
	g.AddRobot(newTestRobot(""))
	return g
}

func initTestMaster1Robot() *Master {
	log.SetOutput(&NullReadWriteCloser{})
	g := NewMaster()
	g.trap = func(c chan os.Signal) {
		c <- os.Interrupt
	}
	g.AddRobot(newTestRobot("Robot99"))

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

func TestMasterRobot(t *testing.T) {
	g := initTestMaster()
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

func TestMasterToJSON(t *testing.T) {
	g := initTestMaster()
	g.AddCommand("test_function", func(params map[string]interface{}) interface{} {
		return nil
	})
	json := NewJSONMaster(g)
	gobottest.Assert(t, len(json.Robots), g.Robots().Len())
	gobottest.Assert(t, len(json.Commands), len(g.Commands()))
}

func TestMasterStart(t *testing.T) {
	g := initTestMaster()
	gobottest.Assert(t, g.Start(), nil)
	gobottest.Assert(t, g.Stop(), nil)
	gobottest.Assert(t, g.Running(), false)
}

func TestMasterStartAutoRun(t *testing.T) {
	g := NewMaster()
	g.AddRobot(newTestRobot("Robot99"))
	go g.Start()
	time.Sleep(10 * time.Millisecond)
	gobottest.Assert(t, g.Running(), true)

	// stop it
	gobottest.Assert(t, g.Stop(), nil)
	gobottest.Assert(t, g.Running(), false)
}

func TestMasterStartDriverErrors(t *testing.T) {
	g := initTestMaster1Robot()
	e := errors.New("driver start error 1")
	testDriverStart = func() (err error) {
		return e
	}

	var expected error
	expected = multierror.Append(expected, e)
	expected = multierror.Append(expected, e)
	expected = multierror.Append(expected, e)

	gobottest.Assert(t, g.Start(), expected)
	gobottest.Assert(t, g.Stop(), nil)

	testDriverStart = func() (err error) { return }
}

func TestMasterHaltFromRobotDriverErrors(t *testing.T) {
	g := initTestMaster1Robot()
	e := errors.New("driver halt error 1")
	testDriverHalt = func() (err error) {
		return e
	}

	var expected error
	expected = multierror.Append(expected, e)
	expected = multierror.Append(expected, e)
	expected = multierror.Append(expected, e)

	gobottest.Assert(t, g.Start(), nil)
	gobottest.Assert(t, g.Stop(), expected)

	testDriverHalt = func() (err error) { return }
}

func TestMasterStartRobotAdaptorErrors(t *testing.T) {
	g := initTestMaster1Robot()
	e := errors.New("adaptor start error 1")

	testAdaptorConnect = func() (err error) {
		return e
	}

	var expected error
	expected = multierror.Append(expected, e)
	expected = multierror.Append(expected, e)
	expected = multierror.Append(expected, e)

	gobottest.Assert(t, g.Start(), expected)
	gobottest.Assert(t, g.Stop(), nil)

	testAdaptorConnect = func() (err error) { return }
}

func TestMasterFinalizeErrors(t *testing.T) {
	g := initTestMaster1Robot()
	e := errors.New("adaptor finalize error 2")

	testAdaptorFinalize = func() (err error) {
		return e
	}

	var expected error
	expected = multierror.Append(expected, e)
	expected = multierror.Append(expected, e)
	expected = multierror.Append(expected, e)

	gobottest.Assert(t, g.Start(), nil)
	gobottest.Assert(t, g.Stop(), expected)

	testAdaptorFinalize = func() (err error) {
		return nil
	}
}
