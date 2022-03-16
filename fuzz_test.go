//go:build go1.18
// +build go1.18

package mode

import (
	"math/rand"
	"os"
	"testing"
)

func FuzzParseWithUmask(f *testing.F) {
	f.Add("644")
	f.Add("+t")
	f.Add("-t")
	f.Add("u+rX,go=u-w")
	f.Add(complicatedModeChange())
	{
		// XXX should this use a static or random seed?
		//rand := rand.New(rand.NewSource(1))
		rand := rand.New(rand.NewSource(int64(rand.Uint64())))
		for i := 0; i < 100; i++ {
			m := randMode(rand, 50)
			f.Add(string(m))
		}
	}

	f.Fuzz(func(t *testing.T, str string) {
		set, err := ParseWithUmask(str, 0)
		if err != nil {
			t.Skip(err)
		}
		//_ = set.String()
		//t.Logf("%q → %v", str, set)
		for p := os.FileMode(0); p <= os.ModePerm; p++ {
			_ = set.Apply(p)
		}
	})
}
