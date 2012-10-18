`nrsc` - Resource Compiler for Go
=================================
`nrsc` compiles a directory of resource into a Go source file so you can still
deploy a single executable as a web server with all the CSS, image files, JS ...
included.

Invocation
==========
::

    nrsc path_to_resource_dir

This will create a local directory called `nrsc` which you can import in your
code.

API
===
The `nrsc` package has the following interface

`nrsc.Handle(prefix string)`
    This will register with the `net/http` module to handle all paths starting with prefix. 

    When a request is handled, `prefix` is stripped and then a resource is
    located and served.

    Resource that are not found will cause an HTTP 404 response.
    

`nrsc.Get(path string) Resource`
    Will return a resource interface (or `nil` if not found) (see resource interface below).
    This allows you more control on how to serve.


Resource Interface
------------------

`func Open() io.Reader`
    Returns a reader to resource data

`func Size() int64`
    Returns resource size (to be used with `Content-Length` HTTP header).

`func MimeType() string`
    Returns mime type (to be used with `Content-Type` HTTP header).

`func ModTime() time.Time`
    Returns modification time (to be used with `Last-Modified` HTTP header).


Example Code
------------
::

    package main

    import (
            "fmt"
            "net/http"
            "os"

            "./nrsc"
    )

    func indexHandler(w http.ResponseWriter, req *http.Request) {
            fmt.Fprintf(w, "Hello World\n")
    }

    func main() {
            nrsc.Handle("/static/")
            http.HandleFunc("/", indexHandler)
            if err := http.ListenAndServe(":8080", nil); err != nil {
                    fmt.Fprintf(os.Stderr, "error: %s\n", err)
                    os.Exit(1)
            }
    }

Contact
=======
https://bitbucket.org/tebeka/nrsc
    
License
=======
MIT (see `LICENSE.txt`)
