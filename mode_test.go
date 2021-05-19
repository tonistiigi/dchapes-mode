package mode

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"testing/quick"
)

func TestParseOctal(t *testing.T) {
	for i := 0; i <= 07777; i++ {
		s := fmt.Sprintf("%#o", i)
		_, err := Parse(s)
		if err != nil {
			t.Errorf("cmode.Parse(%q): %v", s, err)
		}
	}
}

func TestParseQuick(t *testing.T) {
	f := func(m mode) bool {
		//t.Log("mode:", m)
		_, err := ParseWithUmask(string(m), 0)
		return err == nil
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestParseError(t *testing.T) {
	tests := [...]struct {
		input string
		err   error
	}{
		{"", ErrSyntax},
		{",", ErrSyntax},
		{"q+r", ErrSyntax},
		{"017777", ErrSyntax},
		{"01000000000000000007777", nil},
		{"9", nil},
		{"z", ErrSyntax},
		{"u", ErrSyntax},
		{"g", ErrSyntax},
		{"o", ErrSyntax},
		{"go", ErrSyntax},
	}
	for _, tc := range tests {
		//t.Log(tc)
		_, err := ParseWithUmask(tc.input, 0)
		if err == nil {
			t.Errorf("Parse(%q) succeeded, expected %T:%[2]q", tc.input, tc.err)
		} else if tc.err != nil && err != tc.err {
			t.Errorf("Parse(%q) failed with %T:%[2]q, expected %T:%[3]q",
				tc.input, err, tc.err)
			//} else if tc.err == nil {
			//	t.Logf("Parse(%q) failed with %T:%[2]q", tc.input, err)
		}
	}
}

func TestApply(t *testing.T) {
	//savedmask := syscall.Umask(0)
	//defer syscall.Umask(savedmask)

	tests := [...]struct {
		omode os.FileMode
		nmode os.FileMode
		umask int
		str   string
	}{
		{0777, 0000, 0000, "="},
		{0777, 0777, 0000, "+"},
		{0421, 0421, 0000, "+"},
		{0777, 0777, 0000, "-"},
		{0421, 0421, 0000, "-"},
		{00421, os.ModeSticky | 0421, 0000, "+t"},
		{os.ModeSticky | 0421, 0421, 0000, "-t"},

		{0000, 0400, 0000, "u+r"},
		{0000, 0200, 0000, "u+w"},
		{0000, 0100, 0000, "u+x"},
		{0000, 0000, 0000, "u+X"},
		{0001, 0101, 0000, "u+X"},
		{os.ModeDir | 000, os.ModeDir | 0100, 0000, "u+X"},
		{00000, os.ModeSetuid | 0000, 0000, "u+s"},

		{0000, 0040, 0000, "g+r"},
		{0000, 0020, 0000, "g+w"},
		{0000, 0010, 0000, "g+x"},
		{0000, 0000, 0000, "g+X"},
		{0100, 0110, 0000, "g+X"},
		{os.ModeDir | 000, os.ModeDir | 0010, 0000, "g+X"},
		{0000, os.ModeSetgid | 0000, 0000, "g+s"},

		{0000, 0004, 0000, "o+r"},
		{0000, 0002, 0000, "o+w"},
		{0000, 0001, 0000, "o+x"},
		{0000, 0000, 0000, "o+X"},
		{0100, 0101, 0000, "o+X"},
		{os.ModeDir | 000, os.ModeDir | 0001, 0000, "o+X"},
		{0000, 0000, 0000, "o+s"},

		{0000, 0421, 0000, "0421"},
		{0000, os.ModeSticky | os.ModeSetuid | os.ModeSetgid | 0421, 0000, "07421"}, // XXX

		{0000, 0444, 0000, "+r"},
		{0000, 0222, 0000, "+w"},
		{0000, 0111, 0000, "+x"},
		{0000, 0000, 0000, "+X"},
		{0010, 0111, 0000, "+X"},
		{os.ModeDir | 000, os.ModeDir | 0111, 0000, "+X"},
		{0000, os.ModeSetuid | os.ModeSetgid | 0000, 0000, "+s"},
		{0000, 0440, 0034, "+r"},
		{0000, 0220, 0052, "+w"},
		{0000, 0110, 0061, "+x"},

		{0777, 0333, 0000, "-r"},
		{0777, 0555, 0000, "-w"},
		{0777, 0666, 0000, "-x"},
		{0777, 0777, 0000, "-X"},
		{os.ModeDir | 0777, os.ModeDir | 0777, 0000, "-X"},
		{os.ModeSetuid | os.ModeSetgid | 0777, 0777, 0000, "-s"},
		{0777, 0337, 0037, "-r"},
		{0777, 0557, 0057, "-w"},
		{0777, 0667, 0067, "-x"},

		{0777, 0377, 0000, "u-r"},
		{0777, 0577, 0000, "u-w"},
		{0777, 0677, 0000, "u-x"},
		{0777, 0777, 0000, "u-X"},
		{os.ModeDir | 0777, os.ModeDir | 0777, 0000, "u-X"},
		{os.ModeSetuid | 0777, 0777, 0000, "u-s"},

		{0000, 0600, 0000, "u=rw"},
		{0777, 0677, 0000, "u=rw"},
		{0777, 0077, 0000, "u="},
		{0000, 0666, 0000, "=rw"},
		{0777, 0666, 0000, "=rw"},

		{0042, 0442, 0000, "u=g"},
		{0742, 0442, 0000, "u=g"},
		{0722, 0022, 0000, "u=g-w"},
		{0742, 0442, 0000, "u=g-w"},
		{0762, 0462, 0000, "u=g-w"},
		{0772, 0572, 0000, "u=g-w"},

		{0021, 0121, 0000, "u+o"},
		{0021, 0221, 0000, "u+g"},
		{0021, 0321, 0000, "u+go"},

		{0011, 0111, 0037, "u+uo"},
		{0111, 0011, 0037, "-u"},
		{0111, 0011, 0037, "-Xu"}, // differs from C's 0111
		{0111, 0011, 0037, "-uX"}, // differs from C's 0111

		{0000, 0764, 0000, "a=r,ug=rw,u+x"},
		{0777, 0764, 0000, "a=r,ug=rw,u+x"},

		{0000, os.ModeSticky | 0624, 0000, "o+r,g+w,a+t,g-s,u+og"},
		{os.ModeSetgid | 0100, os.ModeSticky | 0724, 0000, "o+r,g+w,a+t,g-s,u+og"},
		{
			os.ModeDir | os.ModeSetuid | os.ModeSetgid | 0000,
			os.ModeDir | os.ModeSetuid | os.ModeSticky | 0624,
			0000, "o+r,g+w,a+t,g-s,u+og",
		},

		// From the chmod(1) examples:
		{0777, 0644, 0000, "644"},
		{0777, 0755, 0000, "go-w"},
		{0666, 0664, 0002, "=rw,+X"},
		{0766, 0775, 0002, "=rw,+X"}, // should this be 0764?
		{0766, 0765, 0012, "=rw,+X"},
		{0000, 0755, 0000, "u=rwx,go=rx"},
		{0000, 0755, 0000, "u=rwx,go=u-w"},
		{0777, 0700, 0000, "go="},
		{0770, 0755, 0000, "go=u-w"},
		{0670, 0644, 0000, "go=u-w"},
		{0570, 0555, 0000, "go=u-w"},
		{0370, 0311, 0000, "go=u-w"},
	}
	var lastSet Set
	var lastStr string
	//lastMask := savedmask
	var lastMask int
	for _, tc := range tests {
		//t.Logf("%#12o=%[1]v → %#12o=%[2]v via %q\tumask=%#o",
		//	tc.omode, tc.nmode, tc.str, tc.umask)
		set := lastSet
		if tc.str != lastStr || tc.umask != lastMask {
			//_ = syscall.Umask(tc.umask)
			var err error
			set, err = ParseWithUmask(tc.str, tc.umask)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.str, err)
			}
			//t.Logf("%q → %v", tc.str, set)
			lastSet = set
			lastStr = tc.str
			lastMask = tc.umask
		}
		if g, e := set.Apply(tc.omode), tc.nmode; g != e {
			var umaskStr string
			if tc.umask != 0 {
				umaskStr = fmt.Sprintf(" with umask=%#o", tc.umask)
			}
			t.Errorf("%q%s applied to %v gave %v, expected %v",
				tc.str, umaskStr, tc.omode, g, e)
			t.Errorf("%q%s applied to %#o gave %#o, expected %#o",
				tc.str, umaskStr, tc.omode, g, e)
		}
	}
}

var marks = [...]struct{ name, modeChange string }{
	{"Octal", "04741"},
	{"Simple", "u+rw"},
	{"Complicated", complicatedModeChange()},
}

func complicatedModeChange() string {
	// This generates a fairly ridiculous string; it'd probably
	// be better to come up with a complicated but sensible string
	// (i.e. something that might plausibly be seen in real use).
	rng := rand.New(rand.NewSource(0))
	v := mode("").Generate(rng, 50)
	return string(v.Interface().(mode))
}

func BenchmarkParse(b *testing.B) {
	for _, bm := range marks {
		b.Log(bm)
		_, err := ParseWithUmask(bm.modeChange, 0)
		if err != nil {
			b.Errorf("Parse(%q): %v", bm.modeChange, err)
			continue
		}
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = ParseWithUmask(bm.modeChange, 0)
			}
		})
	}
}

func BenchmarkApply(b *testing.B) {
	for _, bm := range marks {
		set, err := ParseWithUmask(bm.modeChange, 0)
		if err != nil {
			b.Errorf("Parse(%q): %v", bm.modeChange, err)
			continue
		}
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = set.Apply(os.ModeDir | os.ModeSetuid | 0421)
			}
		})
	}
}

func BenchmarkGetUmask(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = getumask()
	}
}

// The String() methods are only useful for debugging so
// benchmarking them is kinda silly; mostly done for completness.

func BenchmarkString(b *testing.B) {
	for _, bm := range marks {
		set, err := ParseWithUmask(bm.modeChange, 0)
		if err != nil {
			b.Errorf("Parse(%q): %v", bm.modeChange, err)
			continue
		}
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = set.String()
			}
		})
	}
}
