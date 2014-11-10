package sysfs

import (
	"syscall"
)

type SystemCaller interface {
	Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
}

type NativeSyscall struct{}
type MockSyscall struct{}

var sys SystemCaller = &NativeSyscall{}

func SetSyscall(s SystemCaller) {
	sys = s
}

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	return sys.Syscall(trap, a1, a2, a3)
}

func (sys *NativeSyscall) Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	return syscall.Syscall(trap, a1, a2, a3)
}

func (sys *MockSyscall) Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	return 0, 0, 0
}
