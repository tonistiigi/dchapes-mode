//go:build !darwin && !dragonfly && !freebsd && !netbsd && !openbsd
// +build !darwin,!dragonfly,!freebsd,!netbsd,!openbsd

package mode

import "errors"

func osGetUmask() (modet, error) { return 0, errors.New("not supported") }
