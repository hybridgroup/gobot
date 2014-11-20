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
	PwmWrite(string, byte) (err error)
}
type Servo interface {
	InitServo() (err error)
	ServoWrite(string, byte) (err error)
}
type AnalogWriter interface {
	AnalogWrite(string, byte) (err error)
}
type AnalogReader interface {
	AnalogRead(string) (val int, err error)
}
type DigitalWriter interface {
	DigitalWrite(string, byte) (err error)
}
type DigitalReader interface {
	DigitalRead(string) (val int, err error)
}
