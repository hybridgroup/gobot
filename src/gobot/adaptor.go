package gobot

type Adaptor struct {
  Name string
  Connected bool 
  Params map[string]string
}

func (Adaptor) NewAdaptor(a Adaptor) Adaptor {
  return a
}

func (a *Adaptor) Finalize() bool{
  if a.IsConnected() {
    a.Disconnect()
  }
  return true
}

func (a *Adaptor) Connect() bool {
  a.Connected = true
  return true
}

func (a *Adaptor) Disconnect() bool {
  a.Connected = false
  return true
}

func (a *Adaptor) Reconnect() bool {
  if !a.IsConnected(){
    return a.Connect()
  }
  return true
}

func (a *Adaptor) IsConnected() bool {
  return a.Connected
}
