/* "interface" code written to source file. */
package main

const iface = `
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
	Size() int64
	MimeType() string
	ModTime() time.Time
}

type resource struct {
	size  int64
	mtime time.Time
	mtype string
	data  []byte
}

func (rsc *resource) Open() io.Reader {
	return bytes.NewReader(rsc.data)
}

func (rsc *resource) Size() int64 {
	return rsc.size
}

func (rsc *resource) MimeType() string {
	return rsc.mtype
}

func (rsc *resource) ModTime() time.Time {
	return rsc.mtime
}

type handler int

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rsc := resources[req.URL.Path]
	if rsc == nil {
		http.NotFound(w, req)
		return
	}

	if len(rsc.MimeType()) != 0 {
		w.Header().Set("Content-Type", rsc.MimeType())
	}
	w.Header().Set("Content-Size", fmt.Sprintf("%d", rsc.Size()))
	w.Header().Set("Last-Modified", rsc.ModTime().UTC().Format(http.TimeFormat))

	io.Copy(w, rsc.Open())
}

func Get(path string) Resource {
	return resources[path]
}

func Handle(prefix string) {
	var h handler
	http.Handle(prefix, http.StripPrefix(prefix, h))
}
`
