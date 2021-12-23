//go:build gofuzzbeta || go1.18
// +build gofuzzbeta go1.18

package mode

import (
	"testing"
)

func FuzzParseWithUmask(f *testing.F) {
	f.Add("644")
	f.Add("+t")
	f.Add("-t")
	f.Add("u+rX,go=u-w")
	f.Add(complicatedModeChange())
	f.Fuzz(func(t *testing.T, str string) {
		_, err := ParseWithUmask(str, 0)
		if err != nil {
			t.Skip()
		}
		//t.Logf("%q → %v", str, set)
	})
}
