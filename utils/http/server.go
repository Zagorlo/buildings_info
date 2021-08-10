package http

import (
	"io"
	"io/ioutil"
)

func DrainClose(r io.ReadCloser) {
	_, _ = io.Copy(ioutil.Discard, r)
	_ = r.Close()
}
