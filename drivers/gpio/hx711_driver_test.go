package gpio

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/hx711"

	gobot "gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

func initTestHX711DriverWithStubbedAdaptor() (*HX711Driver, *sensorMock) {
	const (
		clockPinID = "3"
		dataPinID  = "4"
	)
	a := newGpioTestAdaptor()
	_ = a.addDigitalPin(clockPinID)
	_ = a.addDigitalPin(dataPinID)
	d := NewHX711Driver(a, clockPinID, dataPinID)
	s := sensorMock{}
	d.newSensor = func(gobot.DigitalPinner, gobot.DigitalPinner) sensorer { return &s }
	if err := d.Start(); err != nil {
		panic(err)
	}

	return d, &s
}

func TestNewHX711Driver(t *testing.T) {
	// arrange
	const (
		clockPinID = "3"
		dataPinID  = "4"
	)
	a := newGpioTestAdaptor()
	_ = a.addDigitalPin(clockPinID)
	_ = a.addDigitalPin(dataPinID)
	// act
	d := NewHX711Driver(a, clockPinID, dataPinID)
	// assert
	assert.IsType(t, &HX711Driver{}, d)
	// assert: gpio.driver attributes
	assert.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "HX711"))
	assert.Equal(t, a, d.connection)
	assert.Nil(t, d.sensor)
	require.NoError(t, d.afterStart()) // TODO: takes some time without options to reduce timings
	assert.NotNil(t, d.sensor)
	require.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.IsType(t, &hx711.Device{}, d.sensor)
	assert.NotNil(t, d.newSensor)
}

func TestNewHX711Driver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "count up"
		newGc  = hx711.A64
	)
	panicFunc := func() {
		NewHX711Driver(newGpioTestAdaptor(), "1", "2", WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewHX711Driver(newGpioTestAdaptor(), "1", "2", WithName(myName), WithHX711GainAndChannelCfg(newGc))
	// assert
	assert.Equal(t, newGc, d.hx711Cfg.gainAndChannelCfg)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestHX711_WithHX711ReadTimeout(t *testing.T) {
	// arrange
	const myReadTimeout = 11 * time.Second
	cfg := hx711Configuration{}
	// act
	WithHX711ReadTimeout(myReadTimeout).apply(&cfg)
	// assert
	assert.Equal(t, myReadTimeout, cfg.readTimeout)
}

func TestHX711_WithHX711ReadTriesMax(t *testing.T) {
	// arrange
	const myReadTriesMax = uint8(12)
	cfg := hx711Configuration{}
	// act
	WithHX711ReadTriesMax(myReadTriesMax).apply(&cfg)
	// assert
	assert.Equal(t, myReadTriesMax, cfg.readTriesMax)
}

func TestHX711_WithHX711TickSleep(t *testing.T) {
	// arrange
	const myTickSleep = 13 * time.Nanosecond
	cfg := hx711Configuration{}
	// act
	WithHX711TickSleep(myTickSleep).apply(&cfg)
	// assert
	assert.Equal(t, myTickSleep, cfg.tickSleep)
}

func TestHX711Start(t *testing.T) {
	const (
		clockPinID = "13"
		dataPinID  = "14"
	)

	cfg := hx711.DefaultConfig
	cfg.TickSleep = defaultTickSleep

	tests := map[string]struct {
		simulateIdErr     string
		simulateConfigErr error
		wantCfgCalled     []*hx711.DeviceConfig
		wantErr           string
	}{
		"start_ok": {
			wantCfgCalled: []*hx711.DeviceConfig{&cfg},
		},
		"error_no_clock_pin": {
			simulateIdErr: clockPinID,
			wantErr:       "error on get clock pin: pin '13' not found",
		},
		"error_no_data_pin": {
			simulateIdErr: dataPinID,
			wantErr:       "error on get data pin: pin '14' not found",
		},
		"error_configure": {
			simulateConfigErr: errors.New("cfg error"),
			wantCfgCalled:     []*hx711.DeviceConfig{&cfg},
			wantErr:           "cfg error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newGpioTestAdaptor()
			var cpin gobot.DigitalPinner
			if tc.simulateIdErr != clockPinID {
				cpin = a.addDigitalPin(clockPinID)
			}
			var dpin gobot.DigitalPinner
			if tc.simulateIdErr != dataPinID {
				dpin = a.addDigitalPin(dataPinID)
			}
			d := NewHX711Driver(a, clockPinID, dataPinID)
			s := sensorMock{cfgErr: tc.simulateConfigErr}
			d.newSensor = func(clockPin, dataPin gobot.DigitalPinner) sensorer {
				s.clockPin = clockPin
				s.dataPin = dataPin

				return &s
			}
			// act
			err := d.Start()
			// assert
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, cpin, s.clockPin)
				assert.Equal(t, dpin, s.dataPin)
			}
			assert.Equal(t, tc.wantCfgCalled, s.cfgCalled)
		})
	}
}

func TestHX711Zero(t *testing.T) {
	tests := map[string]struct {
		secondReading bool
		simulateErr   error
		wantCalled    []bool
		wantErr       string
	}{
		"zero_first_reading": {
			wantCalled: []bool{false},
		},
		"zero_second_reading": {
			secondReading: true,
			wantCalled:    []bool{true},
		},
		"error_zero": {
			simulateErr: errors.New("zeroing error"),
			wantCalled:  []bool{false},
			wantErr:     "zeroing error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, s := initTestHX711DriverWithStubbedAdaptor()
			s.zeroErr = tc.simulateErr
			// act
			err := d.Zero(tc.secondReading)
			// assert
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantCalled, s.zeroCalled)
		})
	}
}

func TestHX711Calibrate(t *testing.T) {
	tests := map[string]struct {
		setValue                int32
		secondReading           bool
		simulateErr             error
		wantSecondReadingCalled []bool
		wantErr                 string
	}{
		"calibrate_first_reading": {
			setValue:                1234,
			wantSecondReadingCalled: []bool{false},
		},
		"calibrate_second_reading": {
			setValue:                -2345,
			secondReading:           true,
			wantSecondReadingCalled: []bool{true},
		},
		"error_calibrate": {
			simulateErr:             errors.New("calibrate error"),
			wantSecondReadingCalled: []bool{false},
			wantErr:                 "calibrate error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, s := initTestHX711DriverWithStubbedAdaptor()
			s.calibrateErr = tc.simulateErr
			// act
			err := d.Calibrate(tc.setValue, tc.secondReading)
			// assert
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.InDelta(t, tc.setValue, s.calibrateSetValueCalled, 0)
			assert.Equal(t, tc.wantSecondReadingCalled, s.calibrateCalled)
		})
	}
}

func TestHX711SetOffsetAndCalibrationFactor(t *testing.T) {
	tests := map[string]struct {
		setOffsVal              int32
		setFacVal               float32
		secondReading           bool
		simulateErr             error
		wantSecondReadingCalled []bool
		wantErr                 string
	}{
		"setcal_first_reading": {
			setOffsVal:              7654,
			setFacVal:               123.4,
			wantSecondReadingCalled: []bool{false},
		},
		"setcal_second_reading": {
			setOffsVal:              654,
			setFacVal:               -23.45,
			secondReading:           true,
			wantSecondReadingCalled: []bool{true},
		},
		"error_calibrate": {
			setOffsVal:              54,
			setFacVal:               3,
			simulateErr:             errors.New("calibrate error"),
			wantSecondReadingCalled: []bool{false},
			wantErr:                 "calibrate error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, s := initTestHX711DriverWithStubbedAdaptor()
			s.setCalValsErr = tc.simulateErr
			// act
			err := d.SetOffsetAndCalibrationFactor(tc.setOffsVal, tc.setFacVal, tc.secondReading)
			gotOffs, gotFac := d.OffsetAndCalibrationFactor(tc.secondReading)
			// assert
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.setOffsVal, s.setCalValsOffs)
			assert.InDelta(t, tc.setFacVal, s.setCalValsFac, 0)
			assert.Equal(t, tc.wantSecondReadingCalled, s.setCalValsCalled)
			assert.Equal(t, tc.setOffsVal, gotOffs)
			assert.InDelta(t, tc.setFacVal, gotFac, 0)
			assert.Equal(t, tc.wantSecondReadingCalled, s.getCalValsCalled)
		})
	}
}

func TestHX711Measure(t *testing.T) {
	tests := map[string]struct {
		secondReading    bool
		simulateErr      error
		wantUpdateCalled []drivers.Measurement
		wantValuesCalled int
		wantErr          string
	}{
		"measure_ok": {
			wantUpdateCalled: []drivers.Measurement{0},
		},
		"error_measure": {
			simulateErr:      errors.New("measure error"),
			wantUpdateCalled: []drivers.Measurement{0},
			wantErr:          "measure error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const (
				v1 = 12345
				v2 = -5432
			)
			d, s := initTestHX711DriverWithStubbedAdaptor()
			s.updateErr = tc.simulateErr
			s.retVal1 = v1
			s.retVal2 = v2
			// act
			got1, got2, err := d.Measure()
			// assert
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.InDelta(t, v1, got1, 0.0001)
				assert.InDelta(t, v2, got2, 0.0001)
			}
			assert.Equal(t, tc.wantUpdateCalled, s.updateCalled)
		})
	}
}

type sensorMock struct {
	clockPin gobot.DigitalPinner
	dataPin  gobot.DigitalPinner

	cfgCalled               []*hx711.DeviceConfig
	cfgErr                  error
	zeroCalled              []bool
	zeroErr                 error
	calibrateCalled         []bool
	calibrateSetValueCalled int32
	calibrateErr            error
	setCalValsCalled        []bool
	setCalValsOffs          int32
	setCalValsFac           float32
	setCalValsErr           error
	getCalValsCalled        []bool
	updateCalled            []drivers.Measurement
	updateErr               error
	valuesCalled            int
	retVal1                 int64
	retVal2                 int64
}

func (s *sensorMock) Configure(cfg *hx711.DeviceConfig) error {
	s.cfgCalled = append(s.cfgCalled, cfg)
	return s.cfgErr
}

func (s *sensorMock) Zero(secondReading bool) error {
	s.zeroCalled = append(s.zeroCalled, secondReading)
	return s.zeroErr
}

func (s *sensorMock) Calibrate(setValue int32, secondReading bool) error {
	s.calibrateSetValueCalled = setValue
	s.calibrateCalled = append(s.calibrateCalled, secondReading)
	return s.calibrateErr
}

func (s *sensorMock) Update(m drivers.Measurement) error {
	s.updateCalled = append(s.updateCalled, m)
	return s.updateErr
}

func (s *sensorMock) Values() (int64, int64) {
	s.valuesCalled++
	return s.retVal1, s.retVal2
}

func (s *sensorMock) OffsetAndCalibrationFactor(secondReading bool) (int32, float32) {
	s.getCalValsCalled = append(s.getCalValsCalled, secondReading)
	return s.setCalValsOffs, s.setCalValsFac
}

func (s *sensorMock) SetOffsetAndCalibrationFactor(offset int32, calibrationFactor float32, secondReading bool) error {
	s.setCalValsOffs = offset
	s.setCalValsFac = calibrationFactor
	s.setCalValsCalled = append(s.setCalValsCalled, secondReading)
	return s.setCalValsErr
}
