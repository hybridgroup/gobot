package system

import (
	"strconv"

	"gobot.io/x/gobot/v2"
)

// sysfsDitalPinHandler represents the sysfs implementation
type sysfsDigitalPinAccess struct {
	sfa *sysfsFileAccess
}

// gpiodDigitalPinAccess represents the character device implementation
type gpiodDigitalPinAccess struct {
	fs    filesystem
	chips []string
}

func (dpa *sysfsDigitalPinAccess) isType(accesserType digitalPinAccesserType) bool {
	return accesserType == digitalPinAccesserTypeSysfs
}

func (dpa *sysfsDigitalPinAccess) isSupported() bool {
	// currently this is supported by all Kernels
	return true
}

func (dpa *sysfsDigitalPinAccess) createPin(chip string, pin int,
	o ...func(gobot.DigitalPinOptioner) bool,
) gobot.DigitalPinner {
	return newDigitalPinSysfs(dpa.sfa, strconv.Itoa(pin), o...)
}

func (dpa *sysfsDigitalPinAccess) setFs(fs filesystem) {
	dpa.sfa = &sysfsFileAccess{fs: fs, readBufLen: 2}
}

func (dpa *gpiodDigitalPinAccess) isType(accesserType digitalPinAccesserType) bool {
	return accesserType == digitalPinAccesserTypeGpiod
}

func (dpa *gpiodDigitalPinAccess) isSupported() bool {
	chips, err := dpa.fs.find("/dev", "gpiochip")
	if err != nil || len(chips) == 0 {
		return false
	}
	dpa.chips = chips
	return true
}

func (dpa *gpiodDigitalPinAccess) createPin(chip string, pin int,
	o ...func(gobot.DigitalPinOptioner) bool,
) gobot.DigitalPinner {
	return newDigitalPinGpiod(chip, pin, o...)
}

func (dpa *gpiodDigitalPinAccess) setFs(fs filesystem) {
	dpa.fs = fs
}
