package gobot

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func Expect(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		_, file, line, _ := runtime.Caller(1)
		s := strings.Split(file, "/")
		t.Errorf("%v:%v Got %v - type %v, Expected %v - type %v", s[len(s)-1], line, a, reflect.TypeOf(a), b, reflect.TypeOf(b))
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

func (t *testDriver) Init() bool  { return true }
func (t *testDriver) Start() bool { return true }
func (t *testDriver) Halt() bool  { return true }

func NewTestDriver(name string, adaptor *testAdaptor) *testDriver {
	t := &testDriver{
		Driver: Driver{
			commands: make(map[string]func(map[string]interface{}) interface{}),
			name:     name,
			adaptor:  adaptor,
		},
	}

	t.Driver.AddCommand("TestDriverCommand", func(params map[string]interface{}) interface{} {
		name := params["name"].(string)
		return fmt.Sprintf("hello %v", name)
	})

	t.Driver.AddCommand("DriverCommand", func(params map[string]interface{}) interface{} {
		name := params["name"].(string)
		return fmt.Sprintf("hello %v", name)
	})

	return t
}

type testAdaptor struct {
	Adaptor
}

func (t *testAdaptor) Finalize() bool { return true }
func (t *testAdaptor) Connect() bool  { return true }

func NewTestAdaptor(name string) *testAdaptor {
	return &testAdaptor{
		Adaptor: Adaptor{
			Name: name,
			Params: map[string]interface{}{
				"param1": "1",
				"param2": 2,
			},
		},
	}
}

func NewTestRobot(name string) *Robot {
	adaptor1 := NewTestAdaptor("Connection 1")
	adaptor2 := NewTestAdaptor("Connection 2")
	adaptor3 := NewTestAdaptor("Connection 3")
	driver1 := NewTestDriver("Device 1", adaptor1)
	driver2 := NewTestDriver("Device 2", adaptor2)
	driver3 := NewTestDriver("Device 3", adaptor3)
	work := func() {}
	r := NewRobot(name, []Connection{adaptor1, adaptor2, adaptor3}, []Device{driver1, driver2, driver3}, work)
	r.AddCommand("robotTestFunction", func(params map[string]interface{}) interface{} {
		message := params["message"].(string)
		robot := params["robot"].(string)
		return fmt.Sprintf("hey %v, %v", robot, message)
	})
	return r
}
