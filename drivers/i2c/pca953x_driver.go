package i2c

import (
	"fmt"
)

// pca953xDefaultAddress is set to variant PCA9533/2
const (
	pca953xDebug          = false
	pca953xDefaultAddress = 0x63
)

type pca953xRegister uint8

// PCA953xGPIOMode is used to set the mode while write GPIO
type PCA953xGPIOMode uint8

const (
	pca953xRegInp  pca953xRegister = 0x00 // input register
	pca953xRegPsc0 pca953xRegister = 0x01 // r,   frequency prescaler 0
	pca953xRegPwm0 pca953xRegister = 0x02 // r/w, PWM register 0
	pca953xRegPsc1 pca953xRegister = 0x03 // r/w, frequency prescaler 1
	pca953xRegPwm1 pca953xRegister = 0x04 // r/w, PWM register 1
	pca953xRegLs0  pca953xRegister = 0x05 // r/w, LED selector 0
	pca953xRegLs1  pca953xRegister = 0x06 // r/w, LED selector 1 (only in PCA9531, PCA9532)

	pca953xAiMask = 0x10 // autoincrement bit

	// PCA953xModeHighImpedance set the GPIO to high (LED off)
	PCA953xModeHighImpedance PCA953xGPIOMode = 0x00
	// PCA953xModeLowImpedance set the GPIO to low (LED on)
	PCA953xModeLowImpedance PCA953xGPIOMode = 0x01
	// PCA953xModePwm0 set the GPIO to PWM (PWM0 & PSC0)
	PCA953xModePwm0 PCA953xGPIOMode = 0x02
	// PCA953xModePwm1 set the GPIO to PWM (PWM1 & PSC1)
	PCA953xModePwm1 PCA953xGPIOMode = 0x03
)

var (
	errToSmallPeriod    = fmt.Errorf("Given Period to small, must be at least 1/152s (~6.58ms) or 152Hz")
	errToBigPeriod      = fmt.Errorf("Given Period to high, must be max. 256/152s (~1.68s) or 152/256Hz (~0.6Hz)")
	errToSmallDutyCycle = fmt.Errorf("Given Duty Cycle to small, must be at least 0%%")
	errToBigDutyCycle   = fmt.Errorf("Given Duty Cycle to high, must be max. 100%%")
)

// PCA953xDriver is a Gobot Driver for LED Dimmer PCA9530 (2-bit), PCA9533 (4-bit), PCA9531 (8-bit), PCA9532 (16-bit)
// Although this is designed for LED's it can be used as a GPIO (read, write, pwm).
// The names of the public functions reflect this.
//
// please refer to data sheet: https://www.nxp.com/docs/en/data-sheet/PCA9533.pdf
//
// Address range:
// * PCA9530   0x60-0x61 (96-97 dec)
// * PCA9531   0x60-0x67 (96-103 dec)
// * PCA9532   0x60-0x67 (96-103 dec)
// * PCA9533/1 0x62      (98 dec)
// * PCA9533/2 0x63      (99 dec)
//
// each new command must start by setting the register and the AI flag
// 0 0 0 AI | 0 R2 R1 R0
// AI=1 means auto incrementing for R0-R2, which enable reading/writing all registers sequentially
// when AI=1 and reading, then R!=0
// this means: do not start with reading input register, writing input register is recognized but has no effect
// => when AI=1 in general start with R>0
type PCA953xDriver struct {
	*Driver
}

// NewPCA953xDriver creates a new driver with specified i2c interface
// Params:
//
//	conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewPCA953xDriver(c Connector, options ...func(Config)) *PCA953xDriver {
	d := &PCA953xDriver{
		Driver: NewDriver(c, "PCA953x", pca953xDefaultAddress),
	}

	for _, option := range options {
		option(d)
	}

	// TODO: API commands
	return d
}

// SetLED sets the mode (LED off, on, PWM0, PWM1) for the LED output (index 0-7)
func (d *PCA953xDriver) SetLED(idx uint8, mode PCA953xGPIOMode) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.writeLED(idx, mode)
}

// WriteGPIO writes a value to a gpio output (index 0-7)
func (d *PCA953xDriver) WriteGPIO(idx uint8, val uint8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	mode := PCA953xModeLowImpedance
	if val > 0 {
		mode = PCA953xModeHighImpedance
	}

	return d.writeLED(idx, mode)
}

// ReadGPIO reads a gpio input (index 0-7) to a value
func (d *PCA953xDriver) ReadGPIO(idx uint8) (uint8, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// read input register
	val, err := d.readRegister(pca953xRegInp)
	// create return bit
	if err != nil {
		return val, err
	}
	val = 1 << idx & val
	if val > 1 {
		val = 1
	}
	return val, nil
}

// WritePeriod set the content of the frequency prescaler of the given index (0,1) with the given value in seconds
func (d *PCA953xDriver) WritePeriod(idx uint8, valSec float32) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// period is valid in range ~6.58ms..1.68s
	val, err := pca953xCalcPsc(valSec)
	if err != nil && pca953xDebug {
		fmt.Println(err, "value limited!")
	}
	var regPsc pca953xRegister = pca953xRegPsc0
	if idx > 0 {
		regPsc = pca953xRegPsc1
	}
	return d.writeRegister(regPsc, val)
}

// ReadPeriod reads the frequency prescaler in seconds of the given index (0,1)
func (d *PCA953xDriver) ReadPeriod(idx uint8) (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	regPsc := pca953xRegPsc0
	if idx > 0 {
		regPsc = pca953xRegPsc1
	}
	psc, err := d.readRegister(regPsc)
	if err != nil {
		return -1, err
	}
	return pca953xCalcPeriod(psc), nil
}

// WriteFrequency set the content of the frequency prescaler of the given index (0,1) with the given value in Hz
func (d *PCA953xDriver) WriteFrequency(idx uint8, valHz float32) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// frequency is valid in range ~0.6..152Hz
	val, err := pca953xCalcPsc(1 / valHz)
	if err != nil && pca953xDebug {
		fmt.Println(err, "value limited!")
	}
	regPsc := pca953xRegPsc0
	if idx > 0 {
		regPsc = pca953xRegPsc1
	}
	return d.writeRegister(regPsc, val)
}

// ReadFrequency read the frequency prescaler in Hz of the given index (0,1)
func (d *PCA953xDriver) ReadFrequency(idx uint8) (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	regPsc := pca953xRegPsc0
	if idx > 0 {
		regPsc = pca953xRegPsc1
	}
	psc, err := d.readRegister(regPsc)
	if err != nil {
		return -1, err
	}
	// valHz = 1/valSec
	return 1 / pca953xCalcPeriod(psc), nil
}

// WriteDutyCyclePercent set the PWM duty cycle of the given index (0,1) with the given value in percent
func (d *PCA953xDriver) WriteDutyCyclePercent(idx uint8, valPercent float32) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	val, err := pca953xCalcPwm(valPercent)
	if err != nil && pca953xDebug {
		fmt.Println(err, "value limited!")
	}
	regPwm := pca953xRegPwm0
	if idx > 0 {
		regPwm = pca953xRegPwm1
	}
	return d.writeRegister(regPwm, val)
}

// ReadDutyCyclePercent get the PWM duty cycle in percent of the given index (0,1)
func (d *PCA953xDriver) ReadDutyCyclePercent(idx uint8) (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	regPwm := pca953xRegPwm0
	if idx > 0 {
		regPwm = pca953xRegPwm1
	}
	pwm, err := d.readRegister(regPwm)
	if err != nil {
		return -1, err
	}
	// PWM=0..255
	return pca953xCalcDutyCyclePercent(pwm), nil
}

func (d *PCA953xDriver) writeLED(idx uint8, mode PCA953xGPIOMode) error {
	// prepare
	regLs := pca953xRegLs0
	if idx > 3 {
		regLs = pca953xRegLs1
		idx = idx - 4
	}
	regLsShift := idx * 2
	// read old value
	regLsVal, err := d.readRegister(regLs)
	if err != nil {
		return err
	}
	// reset 2 bits at LED position
	regLsVal &= ^uint8(0x03 << regLsShift)
	// set 2 bits according to mode at LED position
	regLsVal |= uint8(mode) << regLsShift
	// write new value
	return d.writeRegister(regLs, regLsVal)
}

func (d *PCA953xDriver) writeRegister(regAddress pca953xRegister, val uint8) error {
	// ensure AI bit is not set
	regAddress = regAddress &^ pca953xAiMask
	// write content of requested register
	return d.connection.WriteByteData(uint8(regAddress), val)
}

func (d *PCA953xDriver) readRegister(regAddress pca953xRegister) (uint8, error) {
	// ensure AI bit is not set
	regAddress = regAddress &^ pca953xAiMask
	return d.connection.ReadByteData(uint8(regAddress))
}

func pca953xCalcPsc(valSec float32) (uint8, error) {
	// valSec = (PSC+1)/152; (PSC=0..255)
	psc := 152*valSec - 1
	if psc < 0 {
		return 0, errToSmallPeriod
	}
	if psc > 255 {
		return 255, errToBigPeriod
	}
	// add 0.5 for better rounding experience
	return uint8(psc + 0.5), nil
}

func pca953xCalcPeriod(psc uint8) float32 {
	return (float32(psc) + 1) / 152
}

func pca953xCalcPwm(valPercent float32) (uint8, error) {
	// valPercent = PWM/256*(256/255*100); (PWM=0..255)
	pwm := 255 * valPercent / 100
	if pwm < 0 {
		return 0, errToSmallDutyCycle
	}
	if pwm > 255 {
		return 255, errToBigDutyCycle
	}
	// add 0.5 for better rounding experience
	return uint8(pwm + 0.5), nil
}

func pca953xCalcDutyCyclePercent(pwm uint8) float32 {
	return 100 * float32(pwm) / 255
}
