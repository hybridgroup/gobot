package i2c

import (
	"fmt"
	"math"
	"time"
)

type Channel int
type ServoMode int
type ContinuousDirection int

const (
	Ch0 Channel = iota
	Ch1
	Ch2
	Ch3
	Ch4
	Ch5
	Ch6
	Ch7
	Ch8
	Ch9
	Ch10
	Ch11
	Ch12
	Ch13
	Ch14
	Ch15
)
const (
	RotationContinuousServo ServoMode = iota
	RotationPositionalServo
	LinearServo
)

func (sm ServoMode) String() string {
	switch sm {
	case RotationContinuousServo:
		return "RotationContinuous"
	case RotationPositionalServo:
		return "RotationPositional"
	case LinearServo:
		return "Linear"
	}
	return ""
}

const (
	defaultClockFreq   = 25e6   // 25MHz
	defaultUpdateRate  = 200    // 200Hz
	defaultRampingRate = 100000 // 100kHz

	defaultMinPulse = 1000 // 1000µs (1ms)
	defaultMaxPulse = 2000 // 2000µs (2ms) - depending on the servo it could be as high as 3ms

	defaultAddress = 0x40
	defaultBus     = 1

	minPWMUpdateRate = 24   // Hz
	maxPWMUpdateRate = 1526 // Hz

	maxOSCClock = 50e6 // 50MHz max (from external clock input)
)

const (
	Stop ContinuousDirection = iota
	CCW
	CW
)

func (cd ContinuousDirection) String() string {
	switch cd {
	case Stop:
		return "Stop"
	case CCW:
		return "CCW"
	case CW:
		return "CW"
	}
	return ""
}

var channelAddress = map[Channel]struct {
	Start uint8
	Stop  uint8
}{
	Ch0:  {0x06, 0x08},
	Ch1:  {0x0A, 0x0C},
	Ch2:  {0x0E, 0x10},
	Ch3:  {0x12, 0x14},
	Ch4:  {0x16, 0x18},
	Ch5:  {0x1A, 0x1C},
	Ch6:  {0x1E, 0x20},
	Ch7:  {0x22, 0x24},
	Ch8:  {0x26, 0x28},
	Ch9:  {0x2A, 0x2C},
	Ch10: {0x2E, 0x30},
	Ch11: {0x32, 0x34},
	Ch12: {0x36, 0x38},
	Ch13: {0x3A, 0x3C},
	Ch14: {0x3E, 0x40},
	Ch15: {0x42, 0x44},
}

type RampingFreq int

const (
	LogRamping RampingFreq = iota
	LinearRamping
	ExponentialRamping
)

func (r RampingFreq) String() string {
	switch r {
	case LogRamping:
		return "Log"
	case LinearRamping:
		return "Linear"
	case ExponentialRamping:
		return "Exponential"
	}
	return ""
}

type ChannelConfigure interface {
	CenterOffset(units int)
	InstantPower(power bool)
	RampingRate(hz uint)
	RampingUpdateFreq(mode RampingFreq)
	MinimalPulse(microseconds uint)
	MaximalPulse(microseconds uint)
	Mode(mode ServoMode)
}

type channelConfig struct {
	centerOffset      int  //due to variability the calculated center slightly incorrect
	instantPower      bool // if true no ramping will occur (increases power fluctuations)
	rampingRate       uint // in Hz
	rampingUpdateFreq RampingFreq
	currentPos        uint16
	minPulse          uint //µs
	maxPulse          uint //µs
	mode              ServoMode
	calCenterPos      uint16
	calMinPos         uint16
	calMaxPos         uint16
	calPosFactor      float64
}

func (cc *channelConfig) CenterOffset(units int) {
	cc.centerOffset = units
}

func (cc *channelConfig) InstantPower(power bool) {
	cc.instantPower = power
}

func (cc *channelConfig) RampingRate(hz uint) {
	cc.rampingRate = hz
}

func (cc *channelConfig) RampingUpdateFreq(mode RampingFreq) {
	cc.rampingUpdateFreq = mode
}

func (cc *channelConfig) MinimalPulse(microseconds uint) {
	cc.minPulse = microseconds
}

func (cc *channelConfig) MaximalPulse(microseconds uint) {
	cc.maxPulse = microseconds
}

func (cc *channelConfig) Mode(mode ServoMode) {
	cc.mode = mode
}

type SparkFunDev15316 struct {
	Config
	updateRate uint // Hz in the range [24,1526]
	oscClock   uint // Hz
	channels   map[Channel]*channelConfig
	Connector  Connector
	connection Connection
}

func NewSparkFunDev15316(conn Connector, options ...func(config Config)) *SparkFunDev15316 {
	driver := &SparkFunDev15316{
		Config:    NewConfig(),
		Connector: conn,
	}

	for _, option := range options {
		option(driver)
	}
	return driver
}

func (sf *SparkFunDev15316) Init() (err error) {
	return sf.InitWithRateAndClock(defaultUpdateRate, defaultClockFreq)
}

func (sf *SparkFunDev15316) InitWithRate(updateRate uint) (err error) {
	return sf.InitWithRateAndClock(updateRate, defaultClockFreq)
}

func (sf *SparkFunDev15316) InitWithRateAndClock(updateRate, oscClock uint) error {
	if sf == nil {
		return fmt.Errorf("can not operated on null object")
	}

	if updateRate < minPWMUpdateRate || updateRate > maxPWMUpdateRate {
		return fmt.Errorf("updateRate (%v) is out of range [%v,%v]", updateRate, minPWMUpdateRate, maxPWMUpdateRate)
	}
	sf.updateRate = updateRate

	if oscClock > maxOSCClock {
		return fmt.Errorf("the external clock %v exceeds allowable limit of %v", oscClock, maxOSCClock)
	}
	sf.oscClock = oscClock

	sf.channels = make(map[Channel]*channelConfig)

	var err error
	addr := sf.GetAddressOrDefault(defaultAddress)
	bus := sf.GetBusOrDefault(defaultBus)
	sf.connection, err = sf.Connector.GetConnection(addr, bus)
	if err != nil {
		return err
	}

	//now we enable the pHat and set the prescalar he enable the pHat again
	//first enable the chip
	err = sf.connection.WriteByteData(0, 0x20)
	if err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond) // pause for completion

	//now that we have a connection send over the prescalar
	prescalar := uint8(uint(math.Ceil(float64(sf.oscClock)/(4096*float64(sf.updateRate)))) - 1)

	//enable prescaling
	if err = sf.connection.WriteByteData(0, 0x10); err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond) // pause for completion

	//set the prescalar
	if err = sf.connection.WriteByteData(0xfe, prescalar); err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond) // pause for completion

	//finally enable the chip again
	if err = sf.connection.WriteByteData(0, 0x20); err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond) // pause for completion

	return nil
}

func (sf *SparkFunDev15316) Configure(channelNum Channel, options ...func(config ChannelConfigure)) error {
	if sf == nil {
		return fmt.Errorf("can not operated on null object")
	}

	if channelNum < Ch0 || channelNum > Ch15 {
		return fmt.Errorf("channel is out of range[%v,%v]", Ch0, Ch15)
	}

	if sf.channels == nil {
		return fmt.Errorf("driver must be initilized before configuration")
	}

	channel, has := sf.channels[channelNum]
	if !has {
		channel = &channelConfig{
			rampingRate: defaultRampingRate,
			minPulse:    defaultMinPulse,
			maxPulse:    defaultMaxPulse,
		}
		sf.channels[channelNum] = channel
	}
	for _, option := range options {
		option(channel)
	}

	//we need to configure the center(neutral) point
	var min = int(float64(channel.minPulse) / 1e6 * float64(sf.updateRate) * 4096)
	var max = int(float64(channel.maxPulse) / 1e6 * float64(sf.updateRate) * 4096)
	var center = int(float64(channel.minPulse+channel.maxPulse) / (1e6 * 2.0) * float64(sf.updateRate) * 4096)

	if channel.centerOffset+center < min || min+center < 0 {
		return fmt.Errorf("the centerOffset(%v) is too large for the expected center (%v) and min (%v)", channel.centerOffset, center, min)
	} else if channel.centerOffset+center > max {
		return fmt.Errorf("the centerOffset(%v) is too large for the expected center (%v) and max (%v)", channel.centerOffset, center, max)
	}

	channel.calMinPos = uint16(min + channel.centerOffset)
	channel.calMaxPos = uint16(max + channel.centerOffset)
	channel.calCenterPos = uint16(center + channel.centerOffset)
	channel.calPosFactor = float64(channel.calMaxPos-channel.calMinPos) / 100.0

	//by default we start pwm at the start of the wave form
	err := sf.connection.WriteByteData(channelAddress[channelNum].Start, 0)
	if err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond) // pause for completion

	channel.currentPos, err = sf.connection.ReadWordData(channelAddress[channelNum].Stop)
	if err != nil {
		return err
	}

	err = sf.connection.WriteWordData(channelAddress[channelNum].Stop, channel.currentPos)
	if err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond) // pause for completion

	return nil
}

func (sf *SparkFunDev15316) SetRadialPos(channelNum Channel, percent float64) error {
	if channelNum < Ch0 || channelNum > Ch15 {
		return fmt.Errorf("channel is out of range[%v,%v]", Ch0, Ch15)
	}

	if sf.channels == nil {
		return fmt.Errorf("driver must be initilized before using")
	}

	channel, has := sf.channels[channelNum]
	if !has {
		return fmt.Errorf("channel must first be configured")
	}

	if channel.mode != RotationPositionalServo {
		return fmt.Errorf("can not set rotation position")
	}

	return ramp(sf.connection, channel, channelAddress[channelNum].Stop, channel.calMinPos+uint16(channel.calPosFactor*percent))
}

func (sf *SparkFunDev15316) SetContinuous(channelNum Channel, direction ContinuousDirection) error {
	if channelNum < Ch0 || channelNum > Ch15 {
		return fmt.Errorf("channel is out of range[%v,%v]", Ch0, Ch15)
	}

	if sf.channels == nil {
		return fmt.Errorf("driver must be initilized before using")
	}

	if direction < Stop || direction > CW {
		return fmt.Errorf("a valid direction expected but found %v", direction)
	}

	channel, has := sf.channels[channelNum]
	if !has {
		return fmt.Errorf("channel must first be configured")
	}

	if channel.mode != RotationContinuousServo {
		return fmt.Errorf("can not set rotation position")
	}

	var targetPosition uint16
	switch direction {
	case CCW:
		targetPosition = channel.calMaxPos
	case CW:
		targetPosition = channel.calMinPos
	default:
		targetPosition = channel.calCenterPos
	}

	return ramp(sf.connection, channel, channelAddress[channelNum].Stop, targetPosition)
}

func ramp(conn Connection, channel *channelConfig, address uint8, targetPosition uint16) (err error) {
	index := 0.0
	position := int(channel.currentPos)
	for position != int(targetPosition) && !channel.instantPower {
		err = conn.WriteWordData(address, uint16(position))
		if err != nil {
			return
		}
		time.Sleep(time.Duration(1.0 / float64(channel.rampingRate) * float64(time.Second)))
		if channel.currentPos < targetPosition {
			switch channel.rampingUpdateFreq {
			case LinearRamping:
				position += 1
			case LogRamping:
				exp := math.Ceil(math.Exp(index))
				if exp+float64(position) >= float64(targetPosition) {
					position = int(targetPosition)
				} else {
					position += int(exp)
				}
			case ExponentialRamping:
				exp := math.Ceil(math.Log(index + 2))
				if exp+float64(position) >= float64(targetPosition) {
					position = int(targetPosition)
				} else {
					position += int(exp)
				}
			}
		} else {
			switch channel.rampingUpdateFreq {
			case LinearRamping:
				position -= 1
			case LogRamping:
				exp := math.Ceil(math.Exp(index))
				if float64(position)-exp <= float64(targetPosition) {
					position = int(targetPosition)
				} else {
					position -= int(exp)
				}
			case ExponentialRamping:
				exp := math.Ceil(math.Log(index + 2))
				if float64(position)-exp <= float64(targetPosition) {
					position = int(targetPosition)
				} else {
					position -= int(exp)
				}
			}
		}
		index += 1.0
	}
	channel.currentPos = targetPosition
	err = conn.WriteWordData(address, targetPosition)
	time.Sleep(time.Duration(1.0 / float64(channel.rampingRate) * float64(time.Second)))
	return
}

func (sf *SparkFunDev15316) Halt() error {
	return nil
}

//ChannelConfigures
func CenterOffset(units int) func(config ChannelConfigure) {
	return func(config ChannelConfigure) {
		config.CenterOffset(units)
	}
}

func InstantPower(power bool) func(config ChannelConfigure) {
	return func(config ChannelConfigure) {
		config.InstantPower(power)
	}
}

func RampingRate(hertz uint) func(config ChannelConfigure) {
	return func(config ChannelConfigure) {
		config.RampingRate(hertz)
	}
}

func MinimalPulse(microseconds uint) func(config ChannelConfigure) {
	return func(config ChannelConfigure) {
		config.MinimalPulse(microseconds)
	}
}

func MaximalPulse(microseconds uint) func(config ChannelConfigure) {
	return func(config ChannelConfigure) {
		config.MaximalPulse(microseconds)
	}
}

func Mode(mode ServoMode) func(config ChannelConfigure) {
	return func(config ChannelConfigure) {
		config.Mode(mode)
	}
}

func RampingUpdateFreq(mode RampingFreq) func(config ChannelConfigure) {
	return func(config ChannelConfigure) {
		config.RampingUpdateFreq(mode)
	}
}
