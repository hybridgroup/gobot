package onewire

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDS18B20Driver(t *testing.T) {
	// arrange & act
	a := newOneWireTestAdaptor()
	d := NewDS18B20Driver(a, 2345)
	// assert
	assert.IsType(t, &DS18B20Driver{}, d)
	assert.NotNil(t, d.driver)
	assert.NotNil(t, d.ds18b20Cfg)
	assert.Equal(t, uint8(12), d.ds18b20Cfg.resolution)
	assert.Equal(t, uint16(750), d.ds18b20Cfg.conversionTime)
	assert.InDelta(t, float32(1), d.ds18b20Cfg.scaleUnit(1000), 0.0)
}

func TestDS18B20Start(t *testing.T) {
	tests := map[string]struct {
		cfgResolution uint8
		cfgConvTime   uint16
		simulateErr   bool
		wantCommands  []string
		wantErr       string
	}{
		"start_ok": {
			cfgResolution: 12,
			cfgConvTime:   750,
		},
		"start_change_resolution": {
			cfgResolution: 9,
			cfgConvTime:   750,
			wantCommands:  []string{"resolution"},
		},
		"start_change_conversiontime": {
			cfgResolution: 12,
			cfgConvTime:   250,
			wantCommands:  []string{"conv_time"},
		},
		"start_change_all": {
			cfgResolution: 8,
			cfgConvTime:   150,
			wantCommands:  []string{"resolution", "conv_time"},
		},
		"error_start": {
			simulateErr: true,
			wantErr:     "GetOneWireConnection error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newOneWireTestAdaptor()
			d := NewDS18B20Driver(a, 987654321)
			d.ds18b20Cfg.resolution = tc.cfgResolution
			d.ds18b20Cfg.conversionTime = tc.cfgConvTime
			a.retErr = tc.simulateErr
			// act
			err := d.Start()
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantCommands, a.sendCommands)
		})
	}
}

func TestDS18B20Halt(t *testing.T) {
	tests := map[string]struct {
		cfgResolution uint8
		cfgConvTime   uint16
		simulateErr   bool
		wantCommands  []string
		wantErr       string
	}{
		"start_ok": {
			cfgResolution: 12,
			cfgConvTime:   750,
		},
		"start_change_resolution": {
			cfgResolution: 9,
			cfgConvTime:   750,
			wantCommands:  []string{"resolution"},
		},
		"start_change_conversiontime": {
			cfgResolution: 12,
			cfgConvTime:   250,
			wantCommands:  []string{"conv_time"},
		},
		"start_change_all": {
			cfgResolution: 8,
			cfgConvTime:   150,
			wantCommands:  []string{"resolution", "conv_time"},
		},
		"error_halt": {
			cfgResolution: 8, // to force writing
			simulateErr:   true,
			wantCommands:  []string{"resolution"},
			wantErr:       "WriteInteger error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newOneWireTestAdaptor()
			d := NewDS18B20Driver(a, 987654321)
			require.NoError(t, d.Start())
			d.ds18b20Cfg.resolution = tc.cfgResolution
			d.ds18b20Cfg.conversionTime = tc.cfgConvTime
			a.retErr = tc.simulateErr
			// act
			err := d.Halt()
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantCommands, a.sendCommands)
		})
	}
}

func TestDS18B20Temperature(t *testing.T) {
	const readValue = 24500
	tests := map[string]struct {
		simulateErr bool
		wantVal     float32
		wantErr     string
	}{
		"read_ok": {
			wantVal: 24.5,
		},
		"error_read": {
			simulateErr: true,
			wantErr:     "ReadInteger error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newOneWireTestAdaptor()
			a.lastValue = readValue
			d := NewDS18B20Driver(a, 987654321)
			require.NoError(t, d.Start())
			a.retErr = tc.simulateErr
			// act
			got, err := d.Temperature()
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, []string{"temperature"}, a.sendCommands)
			assert.InDelta(t, tc.wantVal, got, 0.0)
		})
	}
}

func TestDS18B20Resolution(t *testing.T) {
	tests := map[string]struct {
		readValue   int
		simulateErr bool
		wantVal     float32
		wantErr     string
	}{
		"read_ok": {
			readValue: 9,
			wantVal:   9,
		},
		"error_below_min": {
			readValue: 8,
			wantErr:   "the read value '8' is out of range (9, 10, 11, 12)",
		},
		"error_above_max": {
			readValue: 13,
			wantErr:   "the read value '13' is out of range (9, 10, 11, 12)",
		},
		"error_read": {
			simulateErr: true,
			wantErr:     "ReadInteger error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newOneWireTestAdaptor()
			a.lastValue = tc.readValue
			d := NewDS18B20Driver(a, 987654321)
			require.NoError(t, d.Start())
			a.retErr = tc.simulateErr
			// act
			got, err := d.Resolution()
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, []string{"resolution"}, a.sendCommands)
			assert.InDelta(t, tc.wantVal, got, 0.0)
		})
	}
}

func TestDS18B20IsExternalPowered(t *testing.T) {
	tests := map[string]struct {
		readValue   int
		simulateErr bool
		wantVal     bool
		wantErr     string
	}{
		"read_true": {
			readValue: 1,
			wantVal:   true,
		},
		"read_false": {
			readValue: 0,
			wantVal:   false,
		},
		"error_read": {
			simulateErr: true,
			wantErr:     "ReadInteger error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newOneWireTestAdaptor()
			a.lastValue = tc.readValue
			d := NewDS18B20Driver(a, 987654321)
			require.NoError(t, d.Start())
			a.retErr = tc.simulateErr
			// act
			got, err := d.IsExternalPowered()
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, []string{"ext_power"}, a.sendCommands)
			assert.Equal(t, tc.wantVal, got)
		})
	}
}

func TestDS18B20ConversionTime(t *testing.T) {
	tests := map[string]struct {
		readValue   int
		simulateErr bool
		wantVal     float32
		wantErr     string
	}{
		"read_ok": {
			readValue: 20,
			wantVal:   20,
		},
		"error_below_min": {
			readValue: -1,
			wantErr:   "the read value '-1' is out of range (uint16)",
		},
		"error_above_max": {
			readValue: 65536,
			wantErr:   "the read value '65536' is out of range (uint16)",
		},
		"error_read": {
			simulateErr: true,
			wantErr:     "ReadInteger error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newOneWireTestAdaptor()
			a.lastValue = tc.readValue
			d := NewDS18B20Driver(a, 987654321)
			require.NoError(t, d.Start())
			a.retErr = tc.simulateErr
			// act
			got, err := d.ConversionTime()
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, []string{"conv_time"}, a.sendCommands)
			assert.InDelta(t, tc.wantVal, got, 0.0)
		})
	}
}

func TestDS18B20WithName(t *testing.T) {
	// This is a general test, that parent options are applied by using the WithName() option.
	// All other configuration options can also be tested by With..(val).apply(cfg).
	// arrange
	const newName = "new name"
	a := newOneWireTestAdaptor()
	// act
	d := NewDS18B20Driver(a, 2345, WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestDS18B20WithResolution(t *testing.T) {
	// This is a general test, that options are applied by using the WithResolution() option.
	// All other configuration options can also be tested by With..(val).apply(cfg).
	// arrange
	const newValue = uint8(9)
	a := newOneWireTestAdaptor()
	// act
	d := NewDS18B20Driver(a, 2345, WithResolution(newValue))
	// assert
	assert.Equal(t, newValue, d.ds18b20Cfg.resolution)
}

func TestDS18B20WithConversionTime(t *testing.T) {
	// arrange
	const newValue = uint16(93)
	cfg := ds18b20Configuration{conversionTime: 15}
	// act
	WithConversionTime(newValue).apply(&cfg)
	// assert
	assert.Equal(t, newValue, cfg.conversionTime)
}

func TestDS18B20WithFahrenheit(t *testing.T) {
	// arrange
	cfg := ds18b20Configuration{}
	// act
	WithFahrenheit().apply(&cfg)
	// assert
	assert.InDelta(t, 33.8, cfg.scaleUnit(1000), 0.01)
}
