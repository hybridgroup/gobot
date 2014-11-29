package gobot

type pingDriver struct {
	name       string
	pin        string
	connection Connection
	Eventer
	Commander
}

func (t *pingDriver) Start() (errs []error)  { return }
func (t *pingDriver) Halt() (errs []error)   { return }
func (t *pingDriver) Name() string           { return t.name }
func (t *pingDriver) Pin() string            { return t.pin }
func (t *pingDriver) String() string         { return "pingDriver" }
func (t *pingDriver) Connection() Connection { return t.connection }
func (t *pingDriver) ToJSON() *JSONDevice    { return &JSONDevice{} }

func NewPingDriver(adaptor *loopbackAdaptor, name string) *pingDriver {
	t := &pingDriver{
		name:       name,
		connection: adaptor,
		pin:        "",
		Eventer:    NewEventer(),
		Commander:  NewCommander(),
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
