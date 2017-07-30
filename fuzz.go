// For use with go-fuzz, "github.com/dvyukov/go-fuzz"
//
// +build gofuzz

package mode

//go:generate go-fuzz-build bitbucket.org/dchapes/mode
// Then:
//	go-fuzz -bin=mode-fuzz.zip -workdir=fuzz

// Fuzz is for use with go-fuzz, "github.com/dvyukov/go-fuzz"
func Fuzz(data []byte) int {
	_, err := ParseWithUmask(string(data), 0)
	if err != nil {
		return 0
	}
	return 1
}
