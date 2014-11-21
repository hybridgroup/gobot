package gobot

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func logFailure(t *testing.T, message string) {
	_, file, line, _ := runtime.Caller(2)
	s := strings.Split(file, "/")
	t.Errorf("%v:%v: %v", s[len(s)-1], line, message)
}

// Assert ensures a and b are equal and of the same type, or else it
// causes a t.Error
func Assert(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		logFailure(t, fmt.Sprintf("%v - \"%v\", should equal,  %v - \"%v\"",
			a, reflect.TypeOf(a), b, reflect.TypeOf(b)))
	}
}

// Refute ensures a and b are not equal, causes a t.Error if they are equal
func Refute(t *testing.T, a interface{}, b interface{}) {
	if reflect.DeepEqual(a, b) {
		logFailure(t, fmt.Sprintf("%v - \"%v\", should not equal,  %v - \"%v\"",
			a, reflect.TypeOf(a), b, reflect.TypeOf(b)))
	}
}

type testStruct struct {
	i int
	f float64
}

func NewTestStruct() *testStruct {
	return &testStruct{
		i: 10,
		f: 0.2,
	}
}

func (t *testStruct) Hello(name string, message string) string {
	return fmt.Sprintf("Hello %v! %v", name, message)
}

type NullReadWriteCloser struct{}

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), nil
}

func (NullReadWriteCloser) Close() error {
	return nil
}

type testDriver struct {
	Driver
}

func (t *testDriver) Start() (errs []error) { return }
func (t *testDriver) Halt() (errs []error)  { return }

func NewTestDriver(name string, adaptor *testAdaptor) *testDriver {
	t := &testDriver{
		Driver: *NewDriver(
			name,
			"TestDriver",
			adaptor,
			"1",
			100*time.Millisecond,
		),
	}

	t.AddCommand("TestDriverCommand", func(params map[string]interface{}) interface{} {
		name := params["name"].(string)
		return fmt.Sprintf("hello %v", name)
	})

	t.AddCommand("DriverCommand", func(params map[string]interface{}) interface{} {
		name := params["name"].(string)
		return fmt.Sprintf("hello %v", name)
	})

	return t
}

type testAdaptor struct {
	Adaptor
}

func (t *testAdaptor) Finalize() (errs []error) { return }
func (t *testAdaptor) Connect() (errs []error)  { return }

func NewTestAdaptor(name string) *testAdaptor {
	return &testAdaptor{
		Adaptor: *NewAdaptor(
			name,
			"TestAdaptor",
			"/dev/null",
		),
	}
}

func NewTestRobot(name string) *Robot {
	adaptor1 := NewTestAdaptor("Connection1")
	adaptor2 := NewTestAdaptor("Connection2")
	adaptor3 := NewTestAdaptor("")
	driver1 := NewTestDriver("Device1", adaptor1)
	driver2 := NewTestDriver("Device2", adaptor2)
	driver3 := NewTestDriver("", adaptor3)
	work := func() {}
	r := NewRobot(name,
		[]Connection{adaptor1, adaptor2, adaptor3},
		[]Device{driver1, driver2, driver3},
		work,
	)
	r.AddCommand("robotTestFunction", func(params map[string]interface{}) interface{} {
		message := params["message"].(string)
		robot := params["robot"].(string)
		return fmt.Sprintf("hey %v, %v", robot, message)
	})
	return r
}

type loopbackAdaptor struct {
	Adaptor
}

func (t *loopbackAdaptor) Finalize() (errs []error) { return }
func (t *loopbackAdaptor) Connect() (errs []error)  { return }

func NewLoopbackAdaptor(name string) *loopbackAdaptor {
	return &loopbackAdaptor{
		Adaptor: *NewAdaptor(
			name,
			"Loopback",
		),
	}
}

type pingDriver struct {
	Driver
}

func (t *pingDriver) Start() (errs []error) { return }
func (t *pingDriver) Halt() (errs []error)  { return }

func NewPingDriver(adaptor *loopbackAdaptor, name string) *pingDriver {
	t := &pingDriver{
		Driver: *NewDriver(
			name,
			"Ping",
			adaptor,
		),
	}

	t.AddEvent("ping")

	t.AddCommand("ping", func(params map[string]interface{}) interface{} {
		return t.Ping()
	})

	return t
}

func (t *pingDriver) Ping() string {
	Publish(t.Event("ping"), "ping")
	return "pong"
}
