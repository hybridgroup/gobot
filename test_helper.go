package gobot

import "fmt"

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
	name       string
	pin        string
	connection Connection
	Commander
}

func (t *testDriver) Start() (errs []error)  { return }
func (t *testDriver) Halt() (errs []error)   { return }
func (t *testDriver) Name() string           { return t.name }
func (t *testDriver) Pin() string            { return t.pin }
func (t *testDriver) String() string         { return "testDriver" }
func (t *testDriver) Connection() Connection { return t.connection }
func (t *testDriver) ToJSON() *JSONDevice    { return &JSONDevice{} }

func NewTestDriver(name string, adaptor *testAdaptor) *testDriver {
	t := &testDriver{
		name:       name,
		connection: adaptor,
		pin:        "1",
		Commander:  NewCommander(),
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
	name string
	port string
}

func (t *testAdaptor) Finalize() (errs []error) { return }
func (t *testAdaptor) Connect() (errs []error)  { return }
func (t *testAdaptor) Name() string             { return t.name }
func (t *testAdaptor) Port() string             { return t.port }
func (t *testAdaptor) String() string           { return "testAdaptor" }
func (t *testAdaptor) ToJSON() *JSONConnection  { return &JSONConnection{} }

func NewTestAdaptor(name string) *testAdaptor {
	return &testAdaptor{
		name: name,
		port: "/dev/null",
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
