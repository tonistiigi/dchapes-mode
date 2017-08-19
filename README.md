Mode
========

[Mode](https://bitbucket.org/dchapes/mode)
is a [Go](http://golang.org/) package that provides
a native Go implementation of BSD's
[`setmode`](https://www.freebsd.org/cgi/man.cgi?query=setmode&sektion=3)
and `getmode` which can be used to modify the mode bits of
an [`os.FileMode`](https://golang.org/pkg/os#FileMode) value
based on a sumbolic value as described by the
Unix [`chmod`](https://www.freebsd.org/cgi/man.cgi?query=chmod&sektion=1) command.

[![GoDoc](https://godoc.org/bitbucket.org/dchapes/mode?status.png)](https://godoc.org/bitbucket.org/dchapes/mode)
[ ![Codeship Status for dchapes/mode](https://app.codeship.com/projects/5d61b6e0-671d-0135-6ac8-72f6a397e706/status)](https://app.codeship.com/projects/241007)

Online package documentation is available via
[https://godoc.org/bitbucket.org/dchapes/mode](https://godoc.org/bitbucket.org/dchapes/mode).

To install:

		go get bitbucket.org/dchapes/mode

or `go build` any Go code that imports it:

		import "bitbucket.org/dchapes/mode"
