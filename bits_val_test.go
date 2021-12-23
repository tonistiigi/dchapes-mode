//go:build !windows
// +build !windows

package mode

import (
	"syscall"
	"testing"
)

func TestValues(t *testing.T) {
	tests := [...]struct {
		name     string
		value    uint
		expected uint
	}{
		//{"ifDir", ifDir, syscall.S_IFDIR},
		{"isUID", isUID, syscall.S_ISUID},
		{"isGID", isGID, syscall.S_ISGID},
		{"isTXT", isTXT, syscall.S_ISVTX},
		{"iRUser", iRUser, syscall.S_IRUSR},
		{"iWUser", iWUser, syscall.S_IWUSR},
		{"iXUser", iXUser, syscall.S_IXUSR},
	}
	for _, tc := range tests {
		if tc.value != tc.expected {
			t.Errorf("value for %8v is %#5o, want %#5o",
				tc.name, tc.value, tc.expected)
		}
	}
}
