package gobot

import (
  "fmt"
  "math/rand"
  "time"
)

type Connection struct {
  ConnectionId string
  Name string
  Adaptor string
  Port string
  Params map[string]string
  Robot Robot
}

type ConnectionType struct {
  ConnectionId string
  Name string
  Adaptor Adaptor
  Port Port
  Robot Robot
  Params map[string]string
}

func NewConnection(c Connection) *ConnectionType {
  ct := new(ConnectionType)
  if c.ConnectionId == "" {
    rand.Seed( time.Now().UTC().UnixNano())
    i := rand.Int()
    ct.ConnectionId = fmt.Sprintf("%v", i)
  } else {
    ct.ConnectionId = c.ConnectionId
  }
  ct.Name = c.Name
  //ct.Port = Port.New(c.Port)
  ct.Robot = c.Robot
  ct.Params = c.Params
  ct.Adaptor = Adaptor{ Port: ct.Port, Robot: ct.Robot, Params: ct.Params, }

  return ct
}

func (ct *ConnectionType) Connect() {
  fmt.Println("Connecting to "+ ct.Name + " on port " + ct.Port.ToString() + "...")
  ct.Adaptor.Connect()
}

func (ct *ConnectionType) Disconnect() {
  fmt.Println("Diconnecting from "+ ct.Name + " on port " + ct.Port.ToString() + "...")
  ct.Adaptor.Disconnect()
}

// @return [Boolean] Connection status
func (ct *ConnectionType) IsConnected() bool {
  return ct.Adaptor.IsConnected()
}

// @return [String] Adaptor class name
func (ct *ConnectionType) AdaptorName() string {
  return ct.Adaptor.Name
}

//    # Redirects missing methods to adaptor,
//    # attemps reconnection if adaptor not connected
//    def method_missing(method_name, *arguments, &block)
//      unless adaptor.connected?
//        Logger.warn "Cannot call unconnected adaptor '#{name}', attempting to reconnect..."
//        adaptor.reconnect
//        return nil
//      end
//      adaptor.send(method_name, *arguments, &block)
//    rescue Exception => e
//      Logger.error e.message
//      Logger.error e.backtrace.inspect
//      return nil
//    end
