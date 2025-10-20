package gpio

import (
    "errors"
    "gobot.io/x/gobot"
    "math"
)

var (
    ErrBidirectionalOperation = errors.New("operation possible in bidirectional mode only")
    ErrA1NotSet = errors.New("pin to driver's terminal A1 must be set")
    ErrA2NotSet = errors.New("pin to driver's terminal A2 must be set")
    ErrUnidirectionalCCW = errors.New("pins to driver's terminals A1 and A2 must be set to apply a CCW throttle in unidirectional mode")
)

// SabertoothDIPSwitches represents the DIP switch positions on the sabertooth driver
type SabertoothDIPSwitches uint8

// NewSabertoothDIPSwitches creates a new SabertoothDIPSwitches given the
// position (true: ON, false: OFF) of the DIP switches 3, 4, 5 and 6.
func NewSabertoothDIPSwitches(powerSupply, mixed, linear, bidirectional bool) SabertoothDIPSwitches {
    // Transforms booleans into uint8
    i := func(b bool) SabertoothDIPSwitches {
        if b { return 1 } else { return 0 }
    }
    
    // The first two DIP switches sets the driver in analog mode
    return 0b0000_0011 |
        i(powerSupply) << 2 |
        i(mixed) << 3 |
        i(linear) << 4 |
        i(bidirectional) << 5
}

// DIPSwitch1 returns true if the DIP switch 1 is ON
func (d SabertoothDIPSwitches) DIPSwitch1() bool { return (0b0000_0001 & d) > 0 }
// DIPSwitch2 returns true if the DIP switch 2 is ON
func (d SabertoothDIPSwitches) DIPSwitch2() bool { return (0b0000_0010 & d) > 0 }
// DIPSwitch3 returns true if the DIP switch 3 is ON
func (d SabertoothDIPSwitches) DIPSwitch3() bool { return (0b0000_0100 & d) > 0 }
// DIPSwitch4 returns true if the DIP switch 4 is ON
func (d SabertoothDIPSwitches) DIPSwitch4() bool { return (0b0000_1000 & d) > 0 }
// DIPSwitch5 returns true if the DIP switch 5 is ON
func (d SabertoothDIPSwitches) DIPSwitch5() bool { return (0b0001_0000 & d) > 0 }
// DIPSwitch6 returns true if the DIP switch 6 is ON
func (d SabertoothDIPSwitches) DIPSwitch6() bool { return (0b0010_0000 & d) > 0 }

// PowerSupply returns true if the DIP switch 3 is ON
func (d SabertoothDIPSwitches) PowerSupply() bool { return d.DIPSwitch3() }

// Battery returns true if the DIP switch 3 is OFF
func (d SabertoothDIPSwitches) Battery() bool { return !d.DIPSwitch3() }

// Mixed returns true if the DIP switch 4 is ON
func (d SabertoothDIPSwitches) Mixed() bool { return d.DIPSwitch4() }

// Independent returns true if the DIP switch 4 is OFF
func (d SabertoothDIPSwitches) Independent() bool { return !d.DIPSwitch4() }

// Linear returns true if the DIP switch 5 is ON
func (d SabertoothDIPSwitches) Linear() bool { return d.DIPSwitch5() }

// Exponential returns true if the DIP switch 5 is OFF
func (d SabertoothDIPSwitches) Exponential() bool { return !d.DIPSwitch5() }

// Bidirectional returns true if the DIP switch 6 is ON
func (d SabertoothDIPSwitches) Bidirectional() bool { return d.DIPSwitch6() }

// Unidirectional returns true if the DIP switch 6 is OFF
func (d SabertoothDIPSwitches) Unidirectional() bool { return !d.DIPSwitch6() }

type sabertoothAnalogDriver struct {
    name         string
    switches     SabertoothDIPSwitches
    s1           string
    s2           string
    a1           string
    a2           string
    throttle     float64
    differential float64
    
    // Used for differential change
    m1Ratio, m2Ratio float64
    
    connection PwmWriter
    gobot.Commander
}

// SabertoothAnalogDriver represents a sabertooth driver in analog mode
type SabertoothAnalogDriver interface {
    gobot.Driver
    gobot.Commander
    
    // PinS1 returns the pin connected to the S1 terminal of the motor driver.
    PinS1() string
    
    // PinS2 returns the pin connected to the S2 terminal of the motor driver.
    PinS2() string
    
    // PinA1 returns the pin connected to the A1 terminal of the motor driver.
    PinA1() string
    
    // PinA2 returns the pin connected to the A2 terminal of the motor driver.
    PinA2() string
    
    // Throttle sets the throttle of both motors.
    // The throttle value ranges from -1 to 1.
    //
    // A positive value turns the motors clockwise and a negative value turns the motors
    // counterclockwise or the opposite depending on the motors connection.
    // A value of 0 stops the motors.
    Throttle(throttle float64) error
    
    // GetThrottle returns the current throttle value.
    GetThrottle() float64
    
    // Differential sets the throttle difference between the two motors.
    // The differential value ranges from -1 to 1.
    //
    // A negative value applies a lower power to the motor 1 and a positive value applies
    // a lower power to the motor 2. A value of 0 applies the same power for both motors.
    Differential(differential float64) error
    
    // GetDifferential returns the current differential value.
    GetDifferential() float64
    
    // SetRampRate sets the speed ramping rate of both motors.
    //
    // NOTE: this operation is possible in bidirectional mode only.
    SetRampRate(rate float64) error
    
    // SetMaxPower sets the maximum output for both motors.
    //
    // NOTE: this operation is possible in bidirectional mode only.
    SetMaxPower(ratio float64) error
}

// NewSabertoothAnalogDriver creates a new SabertoothAnalogDriver given a PwmWriter,
// the DIP switch settings and the pins connected to the S1, S2, A1 and A2 terminals
// of the driver.
//
// Adds the following API commands:
//  throttle(value float64) - see SabertoothAnalogDriver.Throttle
//  differential(value float64) - see SabertoothAnalogDriver.Differential
func NewSabertoothAnalogDriver(w PwmWriter, d SabertoothDIPSwitches, s1, s2, a1, a2 string) SabertoothAnalogDriver {
    s := &sabertoothAnalogDriver{
        name:       gobot.DefaultName("SabertoothAnalogDriver"),
        switches:   d,
        s1:         s1,
        s2:         s2,
        a1:         a1,
        a2:         a2,
        m1Ratio:    1,
        m2Ratio:    1,
        connection: w,
        Commander:  gobot.NewCommander(),
    }
    
    s.AddCommand("throttle", func(params map[string]interface{}) interface{} {
        return s.Throttle(params["value"].(float64))
    })
    s.AddCommand("differential", func(params map[string]interface{}) interface{} {
        return s.Differential(params["value"].(float64))
    })
    
    return s
}

func (s *sabertoothAnalogDriver) Start() error { return nil }

func (s *sabertoothAnalogDriver) Halt() error { return nil }

func (s *sabertoothAnalogDriver) Name() string { return s.name }

func (s *sabertoothAnalogDriver) SetName(n string) { s.name = n }

func (s *sabertoothAnalogDriver) Connection() gobot.Connection {
    return s.connection.(gobot.Connection)
}

func (s *sabertoothAnalogDriver) PinS1() string { return s.s1 }

func (s *sabertoothAnalogDriver) PinS2() string { return s.s2 }

func (s *sabertoothAnalogDriver) PinA1() string { return s.a1 }

func (s *sabertoothAnalogDriver) PinA2() string { return s.a2 }

func (s *sabertoothAnalogDriver) GetThrottle() float64 { return s.throttle }

func (s *sabertoothAnalogDriver) GetDifferential() float64 { return s.differential }

func (s *sabertoothAnalogDriver) write(pin string, value float64) error {
    if pin != "" {
        value = math.Min(math.Max(value, 0), 1)
        value = gobot.Rescale(value, 0, 1, 0, math.MaxUint8)
        return s.connection.PwmWrite(pin, byte(math.Round(value)))
    }
    return nil
} // END write

func (s *sabertoothAnalogDriver) SetRampRate(rate float64) error {
    if !s.switches.Bidirectional() { return ErrBidirectionalOperation }
    if s.a1 == "" { return ErrA1NotSet }
    rate = math.Min(math.Max(rate, 0), 1)
    return s.write(s.a1, rate)
} // END SetRampRate

func (s *sabertoothAnalogDriver) SetMaxPower(ratio float64) error {
    if !s.switches.Bidirectional() { return ErrBidirectionalOperation }
    if s.a2 == "" { return ErrA2NotSet }
    ratio = math.Min(math.Max(ratio, 0), 1)
    return s.write(s.a2, ratio)
} // END SetMaxPower

func (s *sabertoothAnalogDriver) Throttle(throttle float64) error {
    throttle = math.Min(math.Max(throttle, -1), 1)
    s.throttle = throttle
    
    if s.switches.Bidirectional() {
        if s.switches.Mixed() {
            // Bidirectional + Mixed
            return s.write(s.s1, (throttle / 2) + 0.5)
        } else {
            // Bidirectional + Independent
            
            // Applies ratio and transforms from a range of -1 to 1 into a range of 0 to 1
            s1Throttle := ((throttle * s.m1Ratio) / 2) + 0.5
            s2Throttle := ((throttle * s.m2Ratio) / 2) + 0.5
            
            if err := s.write(s.s1, s1Throttle); err != nil { return err }
            return s.write(s.s2, s2Throttle)
        }
    } else {
        if throttle < 0 && (s.a1 == "" || s.a2 == "") {
            return ErrUnidirectionalCCW
        }
        
        if s.switches.Mixed() {
            // Unidirectional + Mixed
            a := 1.0
            if throttle < 0 { a = 0.0 }
            
            if err := s.write(s.a1, a); err != nil { return err }
            if err := s.write(s.a2, a); err != nil { return err }
            return s.write(s.s1, math.Abs(throttle))
        } else {
            // Unidirectional + Independent
            m1Throttle := throttle * s.m1Ratio
            m2Throttle := throttle * s.m2Ratio
            
            applyThrottle := func(t float64, sPin, aPin string) error {
                a := 1.0
                if t < 0 { a = 0 }
                
                if err := s.write(aPin, a); err != nil { return err }
                return s.write(sPin, math.Abs(t))
            }
            
            if err := applyThrottle(m1Throttle, s.s1, s.a1); err != nil { return err }
            return applyThrottle(m2Throttle, s.s2, s.a2)
        }
    }
} // END Throttle

func (s *sabertoothAnalogDriver) Differential(angle float64) error {
    angle = math.Min(math.Max(angle, -1), 1)
    s.differential = angle
    
    if s.switches.Mixed() {
        return s.write(s.s2, (angle / 2) + 0.5)
    } else if s.switches.Unidirectional() && (s.a1 == "" || s.a2 == "") {
        // CCW (reversed) throttle not possible (minimum ratio == 0)
        if angle < 0 { // Affects M1
            s.m1Ratio, s.m2Ratio = angle + 1, 1
        } else if angle > 0 { // Affects M2
            s.m1Ratio, s.m2Ratio = 1, angle
        } else { // angle == 0 (centered // no differential change)
            s.m1Ratio, s.m2Ratio = 1, 1
        }
    } else {
        if angle < 0 { // Affects M1
            angle = gobot.Rescale(angle, -1, 0, -1, 1)
            s.m1Ratio, s.m2Ratio = angle, 1
        } else if angle > 0 { // Affects M2
            angle = gobot.Rescale(angle, 0, 1, -1, 1)
            s.m1Ratio, s.m2Ratio = 1, angle
        } else { // angle == 0 (centered // no differential change)
            s.m1Ratio, s.m2Ratio = 1, 1
        }
    }
    return s.Throttle(s.throttle) // Update motors throttle
} // END Differential

