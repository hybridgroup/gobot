package adaptors

import (
	"fmt"

	"gobot.io/x/gobot/v2/system"
)

type CdevPin struct {
	Chip uint8
	Line uint8
}

type DigitalPinDefinition struct {
	Sysfs int
	Cdev  CdevPin
}

type DigitalPinDefinitions map[string]DigitalPinDefinition

type DigitalPinTranslator struct {
	sys            *system.Accesser
	pinDefinitions DigitalPinDefinitions
}

// NewDigitalPinTranslator creates a new instance of a translator for digital GPIO pins, suitable for the most cases.
func NewDigitalPinTranslator(sys *system.Accesser, pinDefinitions DigitalPinDefinitions) *DigitalPinTranslator {
	return &DigitalPinTranslator{sys: sys, pinDefinitions: pinDefinitions}
}

// Translate returns the chip and the line or for legacy sysfs usage the pin number from the given id.
func (pt *DigitalPinTranslator) Translate(id string) (string, int, error) {
	pindef, ok := pt.pinDefinitions[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
	}
	if pt.sys.HasDigitalPinSysfsAccess() {
		return "", pindef.Sysfs, nil
	}
	chip := fmt.Sprintf("gpiochip%d", pindef.Cdev.Chip)
	line := int(pindef.Cdev.Line)
	return chip, line, nil
}
