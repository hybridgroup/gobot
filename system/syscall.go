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
// Note: It would be possible to transfer the address as an unsafe.Pointer to e.g. a byte, uint16 or integer variable.
// The unpack process here would be as follows:
// * convert the payload back to the pointer: addrPtr := (*byte)(payload)
// * call with the content converted to uintptr: r1, r2, errNo = unix.Syscall(trap, f.Fd(), signal, uintptr(*addrPtr))
// This has the main disadvantage, that if someone change the type of the address at caller side, the compiler will not
// detect this problem and this unpack procedure would cause unpredictable results.
// So the decision was taken to give the address here as a separate parameter, although it is not used in every call.
// Note also, that the size of the address variable at Kernel side is u16, therefore uint16 is used here.
//
//nolint:nonamedreturns // useful here
func (sys *nativeSyscall) syscall(
	trap uintptr,
	f File,
	signal uintptr,
	payload unsafe.Pointer,
	address uint16,
) (r1, r2 uintptr, err SyscallErrno) {
	var errNo unix.Errno
	if signal == I2C_TARGET {
		// this is the setup for the address, it just needs to be converted to an uintptr,
		// the given payload is not used in this case, see the comment on the function
		r1, r2, errNo = unix.Syscall(trap, f.Fd(), signal, uintptr(address))
	} else {
		r1, r2, errNo = unix.Syscall(trap, f.Fd(), signal, uintptr(payload))
	}

	return r1, r2, SyscallErrno(errNo)
}

// Error implements the error interface. It wraps the "unix.Errno.Error()".
func (e SyscallErrno) Error() string {
	return unix.Errno(e).Error()
}
