package resources

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Resource interface {
	Open() io.Reader
	Size() int
	Type() string
}

var resources map[string]Resource

type resource struct {
	data  []byte
	ctype string
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

type handler int

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rsc := resources[req.URL.Path]
	if rsc == nil {
		http.NotFound(w, req)
		return
	}

	header := w.Header()
	header.Set("Content-Type", rsc.Type())
	header.Set("Content-Size", fmt.Sprintf("%d", rsc.Size()))
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
	resources["hello"] = &resource{[]byte("hello there"), "text/plain"}
}
