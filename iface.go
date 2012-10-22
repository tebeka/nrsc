/* "interface" code written to source file. */
package main

const iface = `
package nrsc

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Resource interface {
	Open() io.Reader
	Size() int64
	ModTime() time.Time
}

type resource struct {
	size  int64
	mtime time.Time
	data  []byte
}

func (rsc *resource) Open() io.Reader {
	return bytes.NewReader(rsc.data)
}

func (rsc *resource) Size() int64 {
	return rsc.size
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

	mtype := mime.TypeByExtension(filepath.Ext(req.URL.Path))
	if len(mtype) != 0 {
		w.Header().Set("Content-Type", mtype)
	}
	w.Header().Set("Content-Size", fmt.Sprintf("%d", rsc.Size()))
	w.Header().Set("Last-Modified", rsc.ModTime().UTC().Format(http.TimeFormat))

	io.Copy(w, rsc.Open())
}

// Get returns the named resource (nil if not found)
func Get(path string) Resource {
	return resources[path]
}

// Handle register HTTP handler under prefix
func Handle(prefix string) {
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	var h handler
	http.Handle(prefix, http.StripPrefix(prefix, h))
}

// LoadTemplates loads named templates from resources.
// If the argument "t" is nil, it is created from the first resource.
func LoadTemplates(t *template.Template, filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("no files named in call to LoadTemplates")
	}

	for _, filename := range filenames {
		rsc := Get(filename)
		if rsc == nil {
			return nil, fmt.Errorf("can't find %s", filename)
		}

		data, err := ioutil.ReadAll(rsc.Open())
		if err != nil {
			return nil, err
		}

		var tmpl *template.Template
		name := filepath.Base(filename)
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(string(data))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
`
