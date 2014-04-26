package gobotGPIO

type TestAdaptor struct{}

func (t TestAdaptor) AnalogWrite(string, byte)  {}
func (t TestAdaptor) DigitalWrite(string, byte) {}
func (t TestAdaptor) ServoWrite(string, byte)   {}
func (t TestAdaptor) PwmWrite(string, byte)     {}
func (t TestAdaptor) InitServo()                {}
func (t TestAdaptor) AnalogRead(string) int {
	return 99
}
func (t TestAdaptor) DigitalRead(string) int {
	return 1
}
