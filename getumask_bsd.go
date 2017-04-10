// +build !cgo
// +build darwin dragonfly freebsd netbsd openbsd

package mode

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

var sysctlErr error

func osGetUmask() (modet, error) {
	if sysctlErr != nil {
		return 0, sysctlErr
	}

	// Note that this does not work if the sysctl
	// security.bsd.unprivileged_proc_debug is set to 0.
	b, err := unix.SysctlRaw("kern.proc.umask", os.Getpid())
	//log.Println("kern.proc.umask", os.Getpid(), b, err)
	if err != nil {
		sysctlErr = err
		return 0, err
	}
	if len(b) != 2 {
		err = syscall.EIO
		sysctlErr = err
		return 0, err
	}

	// It would be nicer if we could use SysctlUint16 but only
	// SysctlUint32/64 currently exists. SysctlRaw does an extra
	// syscall (to get the size of 2), does an extra allocation,
	// and leaves us to do the uint16 conversion ourselves.

	return *(*modet)(unsafe.Pointer(&b[0])), nil
}
