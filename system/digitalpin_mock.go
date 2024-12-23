package system

import (
	"errors"
	"strconv"

	"gobot.io/x/gobot/v2"
)

type simulateErrors struct {
	applyOption bool
	export      bool
	write       bool
	read        bool
	unexport    bool
}

type mockDigitalPinAccess struct {
	underlyingDigitalPinAccess digitalPinAccesser
	values                     map[string]int
	simulateErrors             map[string]simulateErrors // key is the pin-key
	pins                       map[string]*digitalPinMock
}

type digitalPinMock struct {
	chip           string
	pin            int
	appliedOptions int
	exported       int
	written        []int
	value          int
	simulateErrors simulateErrors
}

func newMockDigitalPinAccess(underlyingDigitalPinAccess digitalPinAccesser) *mockDigitalPinAccess {
	dpa := mockDigitalPinAccess{
		underlyingDigitalPinAccess: underlyingDigitalPinAccess,
		values:                     make(map[string]int),
		simulateErrors:             make(map[string]simulateErrors),
		pins:                       make(map[string]*digitalPinMock),
	}
	return &dpa
}

func (dpa *mockDigitalPinAccess) isType(accesserType digitalPinAccesserType) bool {
	return dpa.underlyingDigitalPinAccess.isType(accesserType)
}

func (dpa *mockDigitalPinAccess) isSupported() bool { return true }

func (dpa *mockDigitalPinAccess) createPin(chip string, pin int,
	o ...func(gobot.DigitalPinOptioner) bool,
) gobot.DigitalPinner {
	dpm := &digitalPinMock{chip: chip, pin: pin}

	key := getDigitalPinMockKey(chip, strconv.Itoa(pin))
	if v, ok := dpa.values[key]; ok {
		dpm.value = v
	}

	if v, ok := dpa.simulateErrors[key]; ok {
		dpm.simulateErrors = v
	}

	dpa.pins[key] = dpm
	return dpm
}

func (dpa *mockDigitalPinAccess) setFs(fs filesystem) {
	panic("setFs() for mockDigitalPinAccess not supported")
}

func (dpa *mockDigitalPinAccess) AppliedOptions(chip, pin string) int {
	return dpa.pins[getDigitalPinMockKey(chip, pin)].appliedOptions
}

func (dpa *mockDigitalPinAccess) Written(chip, pin string) []int {
	return dpa.pins[getDigitalPinMockKey(chip, pin)].written
}

func (dpa *mockDigitalPinAccess) Exported(chip, pin string) int {
	return dpa.pins[getDigitalPinMockKey(chip, pin)].exported
}

func (dpa *mockDigitalPinAccess) UseValue(chip, pin string, value int) {
	key := getDigitalPinMockKey(chip, pin)
	if pin, ok := dpa.pins[key]; ok {
		pin.value = value
	}

	// for creation and re-creation
	dpa.values[key] = value
}

func (dpa *mockDigitalPinAccess) UseUnexportError(chip, pin string) {
	key := getDigitalPinMockKey(chip, pin)
	if pin, ok := dpa.pins[key]; ok {
		pin.simulateErrors.unexport = true
	}

	// for creation and re-creation
	simErrs, ok := dpa.simulateErrors[key]
	if !ok {
		simErrs = simulateErrors{}
	}

	simErrs.unexport = true
	dpa.simulateErrors[getDigitalPinMockKey(chip, pin)] = simErrs
}

func (dp *digitalPinMock) ApplyOptions(options ...func(gobot.DigitalPinOptioner) bool) error {
	dp.appliedOptions = dp.appliedOptions + len(options)

	if dp.simulateErrors.applyOption {
		return errors.New("applyOption error")
	}

	return nil
}

func (dp *digitalPinMock) DirectionBehavior() string {
	panic("DirectionBehavior() for digitalPinMock needs do be implemented now")
}

// Write writes the given value to the character device
func (dp *digitalPinMock) Write(b int) error {
	dp.written = append(dp.written, b)

	if dp.simulateErrors.write {
		return errors.New("write error")
	}

	return nil
}

// Read reads the given value from character device
func (dp *digitalPinMock) Read() (int, error) {
	if dp.simulateErrors.read {
		return -1, errors.New("read error")
	}

	return dp.value, nil
}

// Export sets the pin as exported with the configured direction
func (dp *digitalPinMock) Export() error {
	dp.exported++
	if dp.simulateErrors.export {
		return errors.New("export error")
	}

	return nil
}

// Unexport release the pin
func (dp *digitalPinMock) Unexport() error {
	dp.exported--
	if dp.simulateErrors.unexport {
		return errors.New("unexport error")
	}

	return nil
}

func getDigitalPinMockKey(chip, pin string) string {
	return chip + "_" + pin
}
