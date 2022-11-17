package system

import (
	"gobot.io/x/gobot"
)

type mockDigitalPinHandler struct {
	fs *MockFilesystem
}

type digitalPinMock struct{}

func (h *mockDigitalPinHandler) isSupported() bool { return true }

func (h *mockDigitalPinHandler) createPin(chip string, pin int,
	o ...func(gobot.DigitalPinOptioner) bool) gobot.DigitalPinner {
	dpm := &digitalPinMock{}
	return dpm
}

func (h *mockDigitalPinHandler) setFs(fs filesystem) {
	// do nothing
	return
}

func (d *digitalPinMock) ApplyOptions(options ...func(gobot.DigitalPinOptioner) bool) error {
	return nil
}

func (d *digitalPinMock) DirectionBehavior() string {
	return ""
}

// Write writes the given value to the character device
func (d *digitalPinMock) Write(b int) error {
	return nil
}

// Read reads the given value from character device
func (d *digitalPinMock) Read() (n int, err error) {
	return 0, err
}

// Export sets the pin as exported with the configured direction
func (d *digitalPinMock) Export() error {
	return nil
}

// Unexport release the pin
func (d *digitalPinMock) Unexport() error {
	return nil
}
