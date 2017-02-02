package gpio

import (
	"time"

	"gobot.io/x/gobot"
)

const (
	Whole   = 4
	Half    = 2
	Quarter = 1
	Eighth  = 0.500
)

const (
	Rest = 0
	C0   = 16.35
	Db0  = 17.32
	D0   = 18.35
	Eb0  = 19.45
	E0   = 20.60
	F0   = 21.83
	Gb0  = 23.12
	G0   = 24.50
	Ab0  = 25.96
	A0   = 27.50
	Bb0  = 29.14
	B0   = 30.87
	C1   = 32.70
	Db1  = 34.65
	D1   = 36.71
	Eb1  = 38.89
	E1   = 41.20
	F1   = 43.65
	Gb1  = 46.25
	G1   = 49.00
	Ab1  = 51.91
	A1   = 55.00
	Bb1  = 58.27
	B1   = 61.74
	C2   = 65.41
	Db2  = 69.30
	D2   = 73.42
	Eb2  = 77.78
	E2   = 82.41
	F2   = 87.31
	Gb2  = 92.50
	G2   = 98.00
	Ab2  = 103.83
	A2   = 110.00
	Bb2  = 116.54
	B2   = 123.47
	C3   = 130.81
	Db3  = 138.59
	D3   = 146.83
	Eb3  = 155.56
	E3   = 164.81
	F3   = 174.61
	Gb3  = 185.00
	G3   = 196.00
	Ab3  = 207.65
	A3   = 220.00
	Bb3  = 233.08
	B3   = 246.94
	C4   = 261.63
	Db4  = 277.18
	D4   = 293.66
	Eb4  = 311.13
	E4   = 329.63
	F4   = 349.23
	Gb4  = 369.99
	G4   = 392.00
	Ab4  = 415.30
	A4   = 440.00
	Bb4  = 466.16
	B4   = 493.88
	C5   = 523.25
	Db5  = 554.37
	D5   = 587.33
	Eb5  = 622.25
	E5   = 659.25
	F5   = 698.46
	Gb5  = 739.99
	G5   = 783.99
	Ab5  = 830.61
	A5   = 880.00
	Bb5  = 932.33
	B5   = 987.77
	C6   = 1046.50
	Db6  = 1108.73
	D6   = 1174.66
	Eb6  = 1244.51
	E6   = 1318.51
	F6   = 1396.91
	Gb6  = 1479.98
	G6   = 1567.98
	Ab6  = 1661.22
	A6   = 1760.00
	Bb6  = 1864.66
	B6   = 1975.53
	C7   = 2093.00
	Db7  = 2217.46
	D7   = 2349.32
	Eb7  = 2489.02
	E7   = 2637.02
	F7   = 2793.83
	Gb7  = 2959.96
	G7   = 3135.96
	Ab7  = 3322.44
	A7   = 3520.00
	Bb7  = 3729.31
	B7   = 3951.07
	C8   = 4186.01
	Db8  = 4434.92
	D8   = 4698.63
	Eb8  = 4978.03
	E8   = 5274.04
	F8   = 5587.65
	Gb8  = 5919.91
	G8   = 6271.93
	Ab8  = 6644.88
	A8   = 7040.00
	Bb8  = 7458.62
	B8   = 7902.13
)

var _ gobot.Driver = (*BuzzerDriver)(nil)

// BuzzerDriver represents a digital buzzer
type BuzzerDriver struct {
	pin        string
	name       string
	connection DigitalWriter
	high       bool
	BPM        float64
}

// NewBuzzerDriver return a new BuzzerDriver given a DigitalWriter and pin.
func NewBuzzerDriver(a DigitalWriter, pin string) *BuzzerDriver {
	l := &BuzzerDriver{
		name:       gobot.DefaultName("Buzzer"),
		pin:        pin,
		connection: a,
		high:       false,
		BPM:        96.0,
	}

	return l
}

// Start implements the Driver interface
func (l *BuzzerDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (l *BuzzerDriver) Halt() (err error) { return }

// Name returns the BuzzerDrivers name
func (l *BuzzerDriver) Name() string { return l.name }

// SetName sets the BuzzerDrivers name
func (l *BuzzerDriver) SetName(n string) { l.name = n }

// Pin returns the BuzzerDrivers name
func (l *BuzzerDriver) Pin() string { return l.pin }

// Connection returns the BuzzerDrivers Connection
func (l *BuzzerDriver) Connection() gobot.Connection {
	return l.connection.(gobot.Connection)
}

// State return true if the buzzer is On and false if the led is Off
func (l *BuzzerDriver) State() bool {
	return l.high
}

// On sets the buzzer to a high state.
func (l *BuzzerDriver) On() (err error) {
	if err = l.connection.DigitalWrite(l.Pin(), 1); err != nil {
		return
	}
	l.high = true
	return
}

// Off sets the buzzer to a low state.
func (l *BuzzerDriver) Off() (err error) {
	if err = l.connection.DigitalWrite(l.Pin(), 0); err != nil {
		return
	}
	l.high = false
	return
}

// Toggle sets the buzzer to the opposite of it's current state
func (l *BuzzerDriver) Toggle() (err error) {
	if l.State() {
		err = l.Off()
	} else {
		err = l.On()
	}
	return
}

func (l *BuzzerDriver) Tone(hz, duration float64) (err error) {
	// calculation based off https://www.arduino.cc/en/Tutorial/Melody
	tone := (1.0 / (2.0 * hz)) * 1000000.0

	tempo := ((60 / l.BPM) * (duration * 1000))

	for i := 0.0; i < tempo*1000; i += tone * 2.0 {
		if err = l.On(); err != nil {
			return
		}
		time.Sleep(time.Duration(tone) * time.Microsecond)

		if err = l.Off(); err != nil {
			return
		}
		time.Sleep(time.Duration(tone) * time.Microsecond)
	}

	return
}
