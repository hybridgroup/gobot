package sysfs

import (
	"syscall"
)

// SystemCaller represents a Syscall
type SystemCaller interface {
	Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
}

// NativeSyscall represents the native Syscall
type NativeSyscall struct{}

// MockSyscall represents the mock Syscall
type MockSyscall struct {
	Impl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
}

var sys SystemCaller = &NativeSyscall{}

// SetSyscall sets the Syscall implementation
func SetSyscall(s SystemCaller) {
	sys = s
}

// Syscall calls either the NativeSyscall or user defined Syscall
func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	return sys.Syscall(trap, a1, a2, a3)
}

// Syscall calls syscall.Syscall
func (sys *NativeSyscall) Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	return syscall.Syscall(trap, a1, a2, a3)
}

// Syscall implements the SystemCaller interface
func (sys *MockSyscall) Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	if sys.Impl != nil {
		return sys.Impl(trap, a1, a2, a3)
	}
	return 0, 0, 0
}
