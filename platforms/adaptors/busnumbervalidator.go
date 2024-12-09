package adaptors

import "fmt"

type BusNumberValidator struct {
	validNumbers []int
}

// NewBusNumberValidator creates a new instance for a bus number validator, used for I2C and SPI.
func NewBusNumberValidator(validNumbers []int) *BusNumberValidator {
	return &BusNumberValidator{validNumbers: validNumbers}
}

func (bnv *BusNumberValidator) Validate(busNr int) error {
	for _, validNumber := range bnv.validNumbers {
		if validNumber == busNr {
			return nil
		}
	}

	return fmt.Errorf("Bus number %d out of range %v", busNr, bnv.validNumbers)
}
