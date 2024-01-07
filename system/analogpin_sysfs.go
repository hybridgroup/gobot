package system

import "fmt"

type analogPinSysFs struct {
	sysfsPath string
	r, w      bool
	sfa       *sysfsFileAccess
}

func newAnalogPinSysfs(sfa *sysfsFileAccess, path string, r, w bool) *analogPinSysFs {
	p := &analogPinSysFs{
		sysfsPath: path,
		sfa:       sfa,
		r:         r,
		w:         w,
	}
	return p
}

// Read reads a value from sysf path
func (p *analogPinSysFs) Read() (int, error) {
	if !p.r {
		return 0, fmt.Errorf("the pin '%s' is not allowed to read", p.sysfsPath)
	}

	return p.sfa.readInteger(p.sysfsPath)
}

// Write writes a value to sysf path
func (p *analogPinSysFs) Write(val int) error {
	if !p.w {
		return fmt.Errorf("the pin '%s' is not allowed to write (val: %v)", p.sysfsPath, val)
	}

	return p.sfa.writeInteger(p.sysfsPath, val)
}
