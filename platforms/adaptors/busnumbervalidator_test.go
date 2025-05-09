package adaptors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBusNumberValidator(t *testing.T) {
	// arrange
	validNums := []int{5, 8}
	// act
	bnv := NewBusNumberValidator(validNums)
	// assert
	assert.IsType(t, &BusNumberValidator{}, bnv)
	assert.Equal(t, validNums, bnv.validNumbers)
}

func TestBusNumberValidatorValidate(t *testing.T) {
	tests := map[string]struct {
		validNumbers []int
		busNr        int
		wantErr      error
	}{
		"number_negative_error": {
			validNumbers: []int{0, 1, 2, 3, 4},
			busNr:        -1,
			wantErr:      fmt.Errorf("Bus number -1 out of range [0 1 2 3 4]"),
		},
		"number_0_ok": {
			validNumbers: []int{0, 1, 2, 3, 4},
			busNr:        0,
		},
		"number_1_ok": {
			validNumbers: []int{0, 1, 2, 3, 4},
			busNr:        1,
		},
		"number_2_ok": {
			validNumbers: []int{0, 1, 2, 3, 4},
			busNr:        2,
		},
		"number_3_ok": {
			validNumbers: []int{0, 1, 2, 3, 4},
			busNr:        3,
		},
		"number_4_ok": {
			validNumbers: []int{0, 1, 2, 3, 4},
			busNr:        4,
		},
		"number_5_error": {
			validNumbers: []int{0, 1, 2, 3, 4},
			busNr:        5,
			wantErr:      fmt.Errorf("Bus number 5 out of range [0 1 2 3 4]"),
		},
		"number_negative_error_0_2": {
			validNumbers: []int{0, 2},
			busNr:        -1,
			wantErr:      fmt.Errorf("Bus number -1 out of range [0 2]"),
		},
		"number_0_ok_0_2": {
			validNumbers: []int{0, 2},
			busNr:        0,
		},
		"number_1_error_0_2": {
			validNumbers: []int{0, 2},
			busNr:        1,
			wantErr:      fmt.Errorf("Bus number 1 out of range [0 2]"),
		},
		"number_2_ok_0_2": {
			validNumbers: []int{0, 2},
			busNr:        2,
		},
		"number_3_error_0_2": {
			validNumbers: []int{0, 2},
			busNr:        3,
			wantErr:      fmt.Errorf("Bus number 3 out of range [0 2]"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			bnv := NewBusNumberValidator(tc.validNumbers)
			// act
			err := bnv.Validate(tc.busNr)
			// assert
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
