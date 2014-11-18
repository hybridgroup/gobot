package firmata

import (
  "fmt"
  "strconv"
  "time"
  "net"

  "github.com/hybridgroup/gobot"
)

type FirmataEthernetAdaptor struct {
  gobot.Adaptor
  board      *board
  i2cAddress byte
  connect    func(*FirmataEthernetAdaptor)
}

// NewFirmataEthernetAdaptor returns a new firmata adaptor with specified name
// Generates a connect function that opens serial communication in specified port
func NewFirmataEthernetAdaptor(name, sockAddress string) *FirmataEthernetAdaptor {
  return &FirmataEthernetAdaptor{
    Adaptor: *gobot.NewAdaptor(
      name,
      "FirmataEthernetAdaptor",
      sockAddress,
    ),
    connect: func(f *FirmataEthernetAdaptor) {
      fmt.Printf("****** %s *******\n", f.Port())
      conn, err := net.Dial("tcp", f.Port())
      if err != nil {
        panic(err)
      }
      fmt.Printf("connected, creating board %s\n", conn)
      f.board = newBoard(conn)
      fmt.Printf("board crated\n")
    },
  }
}

// Connect returns true if connection to board is succesfull
func (f *FirmataEthernetAdaptor) Connect() bool {
  f.connect(f)

  gobot.Once(f.board.events["report_version"], func(data interface{}) {
    f.board.queryFirmware()
  })

  f.board.connect()
  f.SetConnected(true)
  return true
}

// close finishes connection to serial port
// Prints error message on error
func (f *FirmataEthernetAdaptor) Disconnect() bool {
  err := f.board.serial.Close()
  if err != nil {
    fmt.Println(err)
  }
  return true
}

// Finalize disconnects firmata adaptor
func (f *FirmataEthernetAdaptor) Finalize() bool { return f.Disconnect() }

// InitServo (not yet implemented)
func (f *FirmataEthernetAdaptor) InitServo() {}

// ServoWrite sets angle form 0 to 360 to specified servo pin
func (f *FirmataEthernetAdaptor) ServoWrite(pin string, angle byte) {
  p, _ := strconv.Atoi(pin)

  f.board.setPinMode(byte(p), servo)
  f.board.analogWrite(byte(p), angle)
}

// PwmWrite writes analog value to specified pin
func (f *FirmataEthernetAdaptor) PwmWrite(pin string, level byte) {
  p, _ := strconv.Atoi(pin)

  f.board.setPinMode(byte(p), pwm)
  f.board.analogWrite(byte(p), level)
}

// DigitalWrite writes digital values to specified pin
func (f *FirmataEthernetAdaptor) DigitalWrite(pin string, level byte) {
  p, _ := strconv.Atoi(pin)

  f.board.setPinMode(byte(p), output)
  f.board.digitalWrite(byte(p), level)
}

// DigitalRead retrieves digital value from specified pin
// Returns -1 if response from board is timed out
func (f *FirmataEthernetAdaptor) DigitalRead(pin string) int {
  ret := make(chan int)

  p, _ := strconv.Atoi(pin)
  f.board.setPinMode(byte(p), input)
  f.board.togglePinReporting(byte(p), high, reportDigital)
  f.board.readAndProcess()

  gobot.Once(f.board.events[fmt.Sprintf("digital_read_%v", pin)], func(data interface{}) {
    ret <- int(data.([]byte)[0])
  })

  select {
  case data := <-ret:
    return data
  case <-time.After(10 * time.Millisecond):
  }
  return -1
}

// AnalogRead retrieves value from analog pin.
// NOTE pins are numbered A0-A5, which translate to digital pins 14-19
func (f *FirmataEthernetAdaptor) AnalogRead(pin string) int {
  ret := make(chan int)

  p, _ := strconv.Atoi(pin)
  p = f.digitalPin(p)
  f.board.setPinMode(byte(p), analog)
  f.board.togglePinReporting(byte(p), high, reportAnalog)
  f.board.readAndProcess()

  gobot.Once(f.board.events[fmt.Sprintf("analog_read_%v", pin)], func(data interface{}) {
    b := data.([]byte)
    ret <- int(uint(b[0])<<24 | uint(b[1])<<16 | uint(b[2])<<8 | uint(b[3]))
  })

  select {
  case data := <-ret:
    return data
  case <-time.After(10 * time.Millisecond):
  }
  return -1
}

// AnalogWrite writes value to ananlog pin
func (f *FirmataEthernetAdaptor) AnalogWrite(pin string, level byte) {
  f.PwmWrite(pin, level)
}

// digitalPin converts pin number to digital mapping
func (f *FirmataEthernetAdaptor) digitalPin(pin int) int {
  return pin + 14
}

// I2cStart initializes board with i2c configuration
func (f *FirmataEthernetAdaptor) I2cStart(address byte) {
  f.i2cAddress = address
  f.board.i2cConfig([]byte{0})
}

// I2cRead reads from I2c specified size
// Returns empty byte array if response is timed out
func (f *FirmataEthernetAdaptor) I2cRead(size uint) []byte {
  ret := make(chan []byte)
  f.board.i2cReadRequest(f.i2cAddress, size)

  f.board.readAndProcess()

  gobot.Once(f.board.events["i2c_reply"], func(data interface{}) {
    ret <- data.(map[string][]byte)["data"]
  })

  select {
  case data := <-ret:
    return data
  case <-time.After(10 * time.Millisecond):
  }
  return []byte{}
}

// I2cWrite retrieves i2c data
func (f *FirmataEthernetAdaptor) I2cWrite(data []byte) {
  f.board.i2cWriteRequest(f.i2cAddress, data)
}
