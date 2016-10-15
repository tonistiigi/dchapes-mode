// +build ctest

package mode

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"testing/quick"

	//"bitbucket.org/dchapes/mode/cmode"
	"bitbucket.org/dchapes/mode/cmode/_plain_c"
)

func randomUnixPerm(rand *rand.Rand) int {
	return rand.Intn(0777 + 1)
}

func randomFileMode(rand *rand.Rand) os.FileMode {
	highbits := os.FileMode(rand.Intn(07777+1)) << (32 - 12)
	highbits &= os.ModeDir | os.ModeSetuid | os.ModeSetgid | os.ModeSticky // XXX
	perm := randomUnixPerm(rand)
	return highbits | os.FileMode(perm)
}

func TestCQuick(t *testing.T) {
	cnt := 0
	f := func(fm os.FileMode, m mode, umask int) (os.FileMode, error) {
		cnt++
		s, err := ParseWithUmask(string(m), umask)
		if err != nil {
			return 0, err
		}
		//return s.Apply(fm), nil
		result := s.Apply(fm)
		//t.Logf("Go: %14v and %9q %#04o gives %v", fm, m, umask, result)
		//t.Logf("%q → %v", m, s)
		return result, nil
	}

	g := func(fm os.FileMode, m mode, umask int) (os.FileMode, error) {
		s, err := cmode.ParseWithUmask(string(m), umask)
		if err != nil {
			return 0, err
		}
		//return s.Apply(fm), nil
		result := s.Apply(fm)
		//t.Logf(" C: %14v and %9q %#04o gives %v", fm, m, umask, result)
		return result, nil
	}

	cfg := &quick.Config{
		MaxCountScale: 5,
		Values: func(args []reflect.Value, rand *rand.Rand) {
			args[0] = reflect.ValueOf(randomFileMode(rand))
			args[2] = reflect.ValueOf(randomUnixPerm(rand))
			if rand.Intn(20) == 0 {
				p := randomUnixPerm(rand)
				ps := fmt.Sprintf("%#o", p)
				args[1] = reflect.ValueOf(mode(ps))
			} else {
				args[1] = mode("").Generate(rand, 50)
				// C's setmode is buggy in reguarding 'X' in some cases;
				// e.g. 0111, -uX should ignore X and return 0011 but
				// C returns 0111
				/*
					s := args[1].String()
					if strings.Contains(s, "X") {
						s = strings.Replace(s, "X", "x", -1)
						args[1] = reflect.ValueOf(mode(s))
					}
				*/
			}
		},
	}

	/*
		// Example of C's bug
		f(0114, "-Xu", 0037)
		g(0114, "-Xu", 0037)
		t.Log()
		f(0114, "-uX", 0037)
		g(0114, "-uX", 0037)
		t.Log()
	*/

	if err := quick.CheckEqual(f, g, cfg); err != nil {
		t.Error(err)
	}
	t.Log("checked", cnt, "cases against C's setmode/getmode")
}
