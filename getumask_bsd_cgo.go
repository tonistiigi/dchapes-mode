// Only tested on FreeBSD.
//go:build cgo && (darwin || dragonfly || freebsd || netbsd || openbsd)
// +build cgo
// +build darwin dragonfly freebsd netbsd openbsd

package mode

/*
// From FreBSD /usr/src/lib/libc/gen/setmode.c
#include <signal.h>
#include <sys/stat.h>
#include <sys/sysctl.h>
#include <unistd.h>
#define __libc_sigprocmask sigprocmask

mode_t
sysctl_getumask(void)
{
        sigset_t sigset, sigoset;
        size_t len;
        mode_t mask;
        u_short smask;

        // First try requesting the umask without temporarily modifying it.
        // Note that this does not work if the sysctl
        // security.bsd.unprivileged_proc_debug is set to 0.
        len = sizeof(smask);
        if (sysctl((int[4]){ CTL_KERN, KERN_PROC, KERN_PROC_UMASK, getpid() },
            4, &smask, &len, NULL, 0) == 0)
                return (smask);

#if 1
	// ... with errno set above
	return 0;
#else
        // Since it's possible that the caller is opening files inside a signal
        // handler, protect them as best we can.
        sigfillset(&sigset);
        (void)__libc_sigprocmask(SIG_BLOCK, &sigset, &sigoset);
        (void)umask(mask = umask(0));
        (void)__libc_sigprocmask(SIG_SETMASK, &sigoset, NULL);
        return (mask);
#endif
}
*/
import "C"

func osGetUmask() (modet, error) {
	m, err := C.sysctl_getumask()
	return modet(m), err
}
