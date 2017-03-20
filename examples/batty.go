// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
)

func main() {
	gbot := gobot.NewMaster()

	api.NewAPI(gbot).Start()

	gbot.AddCommand("echo", func(params map[string]interface{}) interface{} {
		return params["a"]
	})

	loopback := NewLoopbackAdaptor("/dev/null")
	ping := NewPingDriver(loopback, "1")

	work := func() {
		gobot.Every(5*time.Second, func() {
			fmt.Println(ping.Ping())
		})
	}
	r := gobot.NewRobot("TestBot",
		[]gobot.Connection{loopback},
		[]gobot.Device{ping},
		work,
	)

	r.AddCommand("hello", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("Hello, %v!", params["greeting"])
	})

	gbot.AddRobot(r)
	gbot.Start()
}

var _ gobot.Adaptor = (*loopbackAdaptor)(nil)

type loopbackAdaptor struct {
	name string
	port string
}

func (t *loopbackAdaptor) Finalize() (err error) { return }
func (t *loopbackAdaptor) Connect() (err error)  { return }
func (t *loopbackAdaptor) Name() string          { return t.name }
func (t *loopbackAdaptor) SetName(n string)      { t.name = n }
func (t *loopbackAdaptor) Port() string          { return t.port }

func NewLoopbackAdaptor(port string) *loopbackAdaptor {
	return &loopbackAdaptor{
		name: "Loopback",
		port: port,
	}
}

var _ gobot.Driver = (*pingDriver)(nil)

type pingDriver struct {
	name       string
	pin        string
	connection gobot.Connection
	gobot.Eventer
	gobot.Commander
}

func (t *pingDriver) Start() (err error)           { return }
func (t *pingDriver) Halt() (err error)            { return }
func (t *pingDriver) Name() string                 { return t.name }
func (t *pingDriver) SetName(n string)             { t.name = n }
func (t *pingDriver) Pin() string                  { return t.pin }
func (t *pingDriver) Connection() gobot.Connection { return t.connection }

func NewPingDriver(adaptor *loopbackAdaptor, pin string) *pingDriver {
	t := &pingDriver{
		name:       "Ping",
		connection: adaptor,
		pin:        pin,
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
	}

	t.AddEvent("ping")

	t.AddCommand("ping", func(params map[string]interface{}) interface{} {
		return t.Ping()
	})

	return t
}

func (t *pingDriver) Ping() string {
	t.Publish(t.Event("ping"), "ping")
	return "pong"
}
