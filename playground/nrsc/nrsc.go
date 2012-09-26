package nrsc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Resource interface {
	Open() io.Reader
	Size() int
	Type() string
	Time() time.Time
}

var resources map[string]Resource

type resource struct {
	data  []byte
	ctype string
	time  time.Time
}

func (rsc *resource) Open() io.Reader {
	return bytes.NewReader(rsc.data)
}

func (rsc *resource) Size() int {
	return len(rsc.data)
}

func (rsc *resource) Type() string {
	return rsc.ctype
}

func (rsc *resource) Time() time.Time {
	return rsc.time
}

type handler int

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rsc := resources[req.URL.Path]
	if rsc == nil {
		http.NotFound(w, req)
		return
	}

	w.Header().Set("Content-Type", rsc.Type())
	w.Header().Set("Content-Size", fmt.Sprintf("%d", rsc.Size()))
	w.Header().Set("Last-Modified", rsc.Time().UTC().Format(http.TimeFormat))

	io.Copy(w, rsc.Open())
}

func Get(path string) Resource {
	return resources[path]
}

func Handle(prefix string) {
	var h handler
	http.Handle(prefix, http.StripPrefix(prefix, h))
}

func init() {
	resources = make(map[string]Resource)
	resources["hello"] = &resource{
		[]byte{
			0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x74,
			0x68, 0x65, 0x72, 0x65, 0xa,
		},
		"text/plain",
		time.Unix(0, 0),
	}
}
