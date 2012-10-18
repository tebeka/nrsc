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

var ignoredDirs = map[string]bool{
	".svn": true,
	".hg":  true,
	".git": true,
}

type File struct {
	path string
	info os.FileInfo
}

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

func die(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
	os.Exit(1)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}

	return false
}

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
	var root string
	outdir := "nrsc"

	flag.StringVar(&root, "root", "", "root direcotry")
	flag.Parse()

	if len(root) == 0 {
		die("<root> is required")
	}

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
