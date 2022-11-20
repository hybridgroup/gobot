package gobot

// DigitalPinner is the interface for system gpio interactions
type DigitalPinner interface {
	// Export exports the pin for use by the operating system
	Export() error
	// Unexport unexports the pin and releases the pin from the operating system
	Unexport() error
	// Direction sets the direction for the pin
	Direction(string) error
	// Read reads the current value of the pin
	Read() (int, error)
	// Write writes to the pin
	Write(int) error
}

// DigitalPinnerProvider is the interface that an Adaptor should implement to allow
// clients to obtain access to any DigitalPin's available on that board.
type DigitalPinnerProvider interface {
	DigitalPin(string, string) (DigitalPinner, error)
}

// PWMPinner is the interface for system PWM interactions
type PWMPinner interface {
	// Export exports the pin for use by the operating system
	Export() error
	// Unexport unexports the pin and releases the pin from the operating system
	Unexport() error
	// Enable enables/disables the PWM pin
	// TODO: rename to "SetEnable(bool)" according to golang style and allow "Enable()" to be the getter function
	Enable(bool) (err error)
	// Polarity returns the polarity either normal or inverted
	Polarity() (polarity string, err error)
	// SetPolarity writes value to pwm polarity path
	SetPolarity(value string) (err error)
	// InvertPolarity sets the polarity to inverted if called with true
	InvertPolarity(invert bool) (err error)
	// Period returns the current PWM period for pin
	Period() (period uint32, err error)
	// SetPeriod sets the current PWM period for pin
	SetPeriod(period uint32) (err error)
	// DutyCycle returns the duty cycle for the pin
	DutyCycle() (duty uint32, err error)
	// SetDutyCycle writes the duty cycle to the pin
	SetDutyCycle(duty uint32) (err error)
}

// PWMPinnerProvider is the interface that an Adaptor should implement to allow
// clients to obtain access to any PWMPin's available on that board.
type PWMPinnerProvider interface {
	PWMPin(string) (PWMPinner, error)
}

// Adaptor is the interface that describes an adaptor in gobot
type Adaptor interface {
	// Name returns the label for the Adaptor
	Name() string
	// SetName sets the label for the Adaptor
	SetName(n string)
	// Connect initiates the Adaptor
	Connect() error
	// Finalize terminates the Adaptor
	Finalize() error
}

// Porter is the interface that describes an adaptor's port
type Porter interface {
	Port() string
}
