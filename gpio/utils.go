package gpio

type PwmDigitalWriter interface {
	DigitalWriter
	Pwm
}
type DirectPin interface {
	DigitalWriter
	DigitalReader
	Pwm
	Servo
	AnalogWriter
	AnalogReader
}
type Pwm interface {
	PwmWrite(string, byte)
}
type Servo interface {
	InitServo()
	ServoWrite(string, byte)
}
type AnalogWriter interface {
	AnalogWrite(string, byte)
}
type AnalogReader interface {
	AnalogRead(string) int
}
type DigitalWriter interface {
	DigitalWrite(string, byte)
}
type DigitalReader interface {
	DigitalRead(string) int
}

// convert to PWM value from analog reading
func ToPwm(i int) byte {
	return byte((255 / 1023.0) * float64(i))
}
