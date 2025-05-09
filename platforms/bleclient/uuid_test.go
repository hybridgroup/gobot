package bleclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_convertUUID(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    string
		wantErr string
	}{
		"32_bit": {
			input: "12345678-4321-1234-4321-123456789abc",
			want:  "12345678-4321-1234-4321-123456789abc",
		},
		"16_bit": {
			input: "12f4",
			want:  "000012f4-0000-1000-8000-00805f9b34fb",
		},
		"32_bit_without_dashes": {
			input: "0123456789abcdef012345678abcdefc",
			want:  "01234567-89ab-cdef-0123-45678abcdefc",
		},
		"error_bad_chacters_16bit": {
			input:   "123g",
			wantErr: "'123g' is not a valid 16-bit Bluetooth UUID",
		},
		"error_bad_chacters_32bit": {
			input:   "12345678-4321-1234-4321-123456789abg",
			wantErr: "'12345678-4321-1234-4321-123456789abg' is not a valid 128-bit Bluetooth UUID",
		},
		"error_too_long": {
			input:   "12345678-4321-1234-4321-123456789abcd",
			wantErr: "'12345678-4321-1234-4321-123456789abcd' is not a valid 128-bit Bluetooth UUID",
		},
		"error_invalid": {
			input:   "12345",
			wantErr: "'12345' is not a valid 128-bit Bluetooth UUID",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// act
			got, err := convertUUID(tc.input)
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
