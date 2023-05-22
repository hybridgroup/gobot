package system

import (
	"syscall"
	"unsafe"
)

// nativeSyscall represents the native Syscall
type nativeSyscall struct{}

// Syscall calls the native syscall.Syscall, implements the SystemCaller interface
func (sys *nativeSyscall) syscall(trap uintptr, f File, signal uintptr, payload unsafe.Pointer) (r1, r2 uintptr, err syscall.Errno) {
	return syscall.Syscall(trap, f.Fd(), signal, uintptr(payload))
}
