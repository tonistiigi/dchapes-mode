package mode

import (
	"os"
	"testing"
)

func TestBits(t *testing.T) {
	cnt := 0
	// Check an individual os.FileMode <=> modet conversion.
	check := func(tfm os.FileMode, tm modet) {
		cnt++
		//t.Logf("%#9o <=> %#5o", tfm, tm)
		m := fileModeToBits(tfm)
		if m != tm {
			t.Errorf("fileModeToBits(%v) gave %#4o, want %#4o",
				tfm, m, tm)
		}
		fm := bitsToFileMode(tfm|fmBits, tm)
		if fm != tfm {
			t.Errorf("bitsToFileMode(%v,%#4o) gave %v, want %v",
				tfm, tm, fm, tfm)
			//t.Errorf("bitsToFileMode(%v,%#4o) gave %#9o, want %#9o",
			//	tfm, tm, fm, tfm)
		}
	}

	// All the possible extra bits that we move/manipulate during the conversion.
	tests := [...]struct {
		fm os.FileMode
		m  modet
	}{
		{0, 0},
		{os.ModeSetuid, isUID},
		{os.ModeSetgid, isGID},
		{os.ModeSticky, isTXT},
		{os.ModeSetuid | os.ModeSetgid, isUID | isGID},
		{os.ModeSetuid | os.ModeSticky, isUID | isTXT},
		{os.ModeSetgid | os.ModeSticky, isGID | isTXT},
		{os.ModeSetuid | os.ModeSetgid | os.ModeSticky, isUID | isGID | isTXT},
	}

	// Individual extra os.FileMode bits we shouldn't touch.
	exbits := []os.FileMode{
		os.ModeDir,
		os.ModeAppend,
		os.ModeExclusive,
		os.ModeTemporary,
		os.ModeSymlink,
		os.ModeDevice,
		os.ModeNamedPipe,
		os.ModeSocket,
		os.ModeCharDevice,
	}

	for _, tc := range tests {
		t.Logf("%#9o <=> %#5o", tc.fm|os.ModePerm, tc.m|0777)

		// With each set of extra bits we move/manimuplate,
		// exhaustively check every possible Unix
		// permision combination (without any
		// of the "extra" bits we don't touch).
		for bits := 0; bits <= 0777; bits++ {
			check(tc.fm|os.FileMode(bits), tc.m|modet(bits))
		}

		// Just test combinations of one and two extra
		// bits with a single Unix permission bit set.
		for _, ex1 := range exbits {
			for _, ex2 := range exbits {
				check(tc.fm|ex1|ex2|0421, tc.m|0421)
			}
		}
	}
	t.Log("checked", cnt, "conversions")
}

// Kinda silly to benchmark these; mostly done for completness.

func BenchmarkFileModeToBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fileModeToBits(os.ModeDir | os.ModeSetuid | 0421)
	}
}

func BenchmarkBitsToFileMode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = bitsToFileMode(os.ModeDir, isUID|0421)
	}
}
