package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

const (
	version = "0.1.0"
	outdir  = "nrsc"
)

var ignoredDirs = map[string]bool{
	".svn": true,
	".hg":  true,
	".git": true,
}

type File struct {
	path string
	info os.FileInfo
}

// iterfiles iterats of directory tree, returns a channel with files to process
func iterfiles(root string) chan *File {
	out := make(chan *File)

	walkfn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if ignoredDirs[strings.ToLower(info.Name())] {
				return filepath.SkipDir
			}
		} else {
			if info.Name()[0] != '.' {
				out <- &File{path, info}
			}
		}
		return nil
	}

	go func() {
		filepath.Walk(root, walkfn)
		close(out)
	}()

	return out
}

// writeResource write resource code to out
func writeResource(prefix int, file *File, out io.Writer) error {
	data, err := ioutil.ReadFile(file.path)
	if err != nil {
		return err
	}
	path, info := file.path, file.info

	key := path[prefix:]
	fmt.Fprintf(out, "\t\"%s\": &resource{\n", key)
	fmt.Fprintf(out, "\t\tsize: %d,\n", info.Size())
	fmt.Fprintf(out, "\t\tmtime: time.Unix(%d, 0),\n", info.ModTime().Unix())
	mtype := mime.TypeByExtension(filepath.Ext(path))
	fmt.Fprintf(out, "\t\tmtype: \"%s\",\n", mtype)
	fmt.Fprintf(out, "\t\tdata: []byte{")
	for _, b := range data {
		fmt.Fprintf(out, "%d, ", b)
	}
	fmt.Fprintf(out, "\t\t},\n\t},\n")

	return nil
}

// die prints error and exists the program with exit status 1
func die(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
	os.Exit(1)
}

// dirExists return true if path exists and is a directory
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}

	return false
}

// writeResources writes the go code for the resources file
func writeResources(root string, out io.Writer) error {
	prefix := len(root)
	if root[len(root)-1] != '/' {
		prefix += 1
	}

	fmt.Fprintf(out, "package nrsc\nimport \"time\"\n")
	fmt.Fprintf(out, "var resources = map[string]Resource {\n")

	for file := range iterfiles(root) {
		if err := writeResource(prefix, file, out); err != nil {
			return fmt.Errorf("can't write %s - %s", file.path, err)
		}
	}

	fmt.Fprintf(out, "\n}")

	return nil
}

func main() {
	var showVersion bool

	flag.BoolVar(&showVersion, "version", false, "show version and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s RESOURCE_DIR\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		fmt.Printf("nrsc %s\n", version)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		die("wrong number of parameters")
	}

	root := flag.Arg(0)

	if !dirExists(root) {
		die("%s is not a directory", root)
	}

	if !dirExists(outdir) {
		if err := os.Mkdir(outdir, 0700); err != nil {
			die("can't create nrsc directory - %s", err)
		}
	}

	ok := false

	defer func() {
		if !ok {
			fmt.Printf("cleaning %s\n", outdir)
			os.RemoveAll(outdir)
		}
	}()

	outfile := fmt.Sprintf("%s/nrsc.go", outdir)
	err := ioutil.WriteFile(outfile, []byte(iface), 0666)
	if err != nil {
		die("can't create %s - %s", outfile, err)
	}

	outfile = fmt.Sprintf("%s/data.go", outdir)
	out, err := os.Create(outfile)
	if err != nil {
		die("can't create %s - %s", outfile, err)
	}

	if err := writeResources(root, out); err != nil {
		die("%s", err)
	}

	ok = true
}
