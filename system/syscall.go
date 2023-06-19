package system

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// SyscallErrno wraps the "unix.Errno"
type SyscallErrno unix.Errno

// wrapping for used constants of unix package
const (
	Syscall_SYS_IOCTL = unix.SYS_IOCTL
	Syscall_EINVAL    = unix.EINVAL
	Syscall_EBUSY     = unix.EBUSY
	Syscall_EFAULT    = unix.EFAULT
)

// nativeSyscall represents the native Syscall
type nativeSyscall struct{}

// Syscall calls the native unix.Syscall, implements the SystemCaller interface
func (sys *nativeSyscall) syscall(trap uintptr, f File, signal uintptr, payload unsafe.Pointer) (r1, r2 uintptr, err SyscallErrno) {
	r1, r2, errNo := unix.Syscall(trap, f.Fd(), signal, uintptr(payload))
	return r1, r2, SyscallErrno(errNo)
}

// Error implements the error interface. It wraps the "unix.Errno.Error()".
func (e SyscallErrno) Error() string {
	return unix.Errno(e).Error()
}
