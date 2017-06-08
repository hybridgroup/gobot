package ollie

import (
	"sync"

	"gobot.io/x/gobot/platforms/ble"
)

var _ ble.BLEConnector = (*bleTestClientAdaptor)(nil)

type bleTestClientAdaptor struct {
	name            string
	address         string
	mtx             sync.Mutex
	withoutReponses bool

	testReadCharacteristic  func(string) ([]byte, error)
	testWriteCharacteristic func(string, []byte) error
}

func (t *bleTestClientAdaptor) Connect() (err error)     { return }
func (t *bleTestClientAdaptor) Reconnect() (err error)   { return }
func (t *bleTestClientAdaptor) Disconnect() (err error)  { return }
func (t *bleTestClientAdaptor) Finalize() (err error)    { return }
func (t *bleTestClientAdaptor) Name() string             { return t.name }
func (t *bleTestClientAdaptor) SetName(n string)         { t.name = n }
func (t *bleTestClientAdaptor) Address() string          { return t.address }
func (t *bleTestClientAdaptor) WithoutReponses(use bool) { t.withoutReponses = use }

func (t *bleTestClientAdaptor) ReadCharacteristic(cUUID string) (data []byte, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testReadCharacteristic(cUUID)
}

func (t *bleTestClientAdaptor) WriteCharacteristic(cUUID string, data []byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testWriteCharacteristic(cUUID, data)
}

func (t *bleTestClientAdaptor) Subscribe(cUUID string, f func([]byte, error)) (err error) {
	// TODO: implement this...
	return
}

func (t *bleTestClientAdaptor) TestReadCharacteristic(f func(cUUID string) (data []byte, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testReadCharacteristic = f
}

func (t *bleTestClientAdaptor) TestWriteCharacteristic(f func(cUUID string, data []byte) (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testWriteCharacteristic = f
}

func NewBleTestAdaptor() *bleTestClientAdaptor {
	return &bleTestClientAdaptor{
		address: "01:02:03:04:05:06",
		testReadCharacteristic: func(cUUID string) (data []byte, e error) {
			return
		},
		testWriteCharacteristic: func(cUUID string, data []byte) (e error) {
			return
		},
	}
}
