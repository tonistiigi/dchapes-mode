// For use with go-fuzz, "github.com/dvyukov/go-fuzz"
// See fuzz_test.go for Go1.18+ version.
//
//go:build gofuzz
// +build gofuzz

package mode

//go:generate go-fuzz-build hg.sr.ht/~dchapes/mode
// Then:
//	go-fuzz -bin=mode-fuzz.zip -workdir=fuzz
// or
//	limits -S -w 1g -v 4g nice go-fuzz -bin=mode-fuzz.zip -workdir=fuzz -dumpcover

// Fuzz is for use with go-fuzz, "github.com/dvyukov/go-fuzz"
func Fuzz(data []byte) int {
	_, err := ParseWithUmask(string(data), 0)
	if err != nil {
		return 0
	}
	return 1
}
