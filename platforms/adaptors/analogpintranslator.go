package adaptors

import (
	"fmt"

	"gobot.io/x/gobot/v2/system"
)

type AnalogPinDefinition struct {
	Path       string
	W          bool   // writable
	ReadBufLen uint16 // readable if buffer > 0
}

type AnalogPinDefinitions map[string]AnalogPinDefinition

type AnalogPinTranslator struct {
	sys            *system.Accesser
	pinDefinitions AnalogPinDefinitions
}

// NewAnalogPinTranslator creates a new instance of a translator for analog pins, suitable for the most cases.
func NewAnalogPinTranslator(sys *system.Accesser, pinDefinitions AnalogPinDefinitions) *AnalogPinTranslator {
	return &AnalogPinTranslator{sys: sys, pinDefinitions: pinDefinitions}
}

// Translate returns the sysfs path for the given id.
func (pt *AnalogPinTranslator) Translate(id string) (string, bool, uint16, error) {
	pinInfo, ok := pt.pinDefinitions[id]
	if !ok {
		return "", false, 0, fmt.Errorf("'%s' is not a valid id for an analog pin", id)
	}

	path := pinInfo.Path
	info, err := pt.sys.Stat(path)
	if err != nil {
		return "", false, 0, fmt.Errorf("Error (%v) on access '%s'", err, path)
	}
	if info.IsDir() {
		return "", false, 0, fmt.Errorf("The item '%s' is a directory, which is not expected", path)
	}

	return path, pinInfo.W, pinInfo.ReadBufLen, nil
}
