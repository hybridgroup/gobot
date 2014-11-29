package gobot

type loopbackAdaptor struct {
	name string
	port string
}

func (t *loopbackAdaptor) Finalize() (errs []error) { return }
func (t *loopbackAdaptor) Connect() (errs []error)  { return }
func (t *loopbackAdaptor) Name() string             { return t.name }
func (t *loopbackAdaptor) Port() string             { return t.port }
func (t *loopbackAdaptor) String() string           { return "loopbackAdaptor" }
func (t *loopbackAdaptor) ToJSON() *JSONConnection  { return &JSONConnection{} }

func NewLoopbackAdaptor(name string) *loopbackAdaptor {
	return &loopbackAdaptor{
		name: name,
		port: "",
	}
}
