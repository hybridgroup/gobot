package system

import (
	"strconv"

	"gobot.io/x/gobot/v2"
)

// sysfsDitalPinHandler represents the sysfs implementation
type sysfsDigitalPinAccess struct {
	sfa *sysfsFileAccess
}

// cdevDigitalPinAccess represents the character device implementation
type cdevDigitalPinAccess struct {
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

func (dpa *cdevDigitalPinAccess) isType(accesserType digitalPinAccesserType) bool {
	return accesserType == digitalPinAccesserTypeCdev
}

func (dpa *cdevDigitalPinAccess) isSupported() bool {
	chips, err := dpa.fs.find("/dev", "gpiochip")
	if err != nil || len(chips) == 0 {
		return false
	}
	dpa.chips = chips
	return true
}

func (dpa *cdevDigitalPinAccess) createPin(chip string, pin int,
	o ...func(gobot.DigitalPinOptioner) bool,
) gobot.DigitalPinner {
	return newDigitalPinCdev(chip, pin, o...)
}

func (dpa *cdevDigitalPinAccess) setFs(fs filesystem) {
	dpa.fs = fs
}
