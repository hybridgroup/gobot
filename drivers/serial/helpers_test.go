package serial

type serialTestAdaptor struct {
	isConnected bool
	name        string
}

func newSerialTestAdaptor() *serialTestAdaptor {
	return &serialTestAdaptor{}
}

func (t *serialTestAdaptor) IsConnected() bool {
	return t.isConnected
}

func (t *serialTestAdaptor) SerialRead(b []byte) (int, error) {
	return len(b), nil
}

func (t *serialTestAdaptor) SerialWrite(b []byte) (int, error) {
	return len(b), nil
}

// gobot.Adaptor interfaces
func (t *serialTestAdaptor) Connect() error   { return nil }
func (t *serialTestAdaptor) Finalize() error  { return nil }
func (t *serialTestAdaptor) Name() string     { return t.name }
func (t *serialTestAdaptor) SetName(n string) { t.name = n }
