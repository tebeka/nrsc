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

func writeResource(root string, file *File, out io.Writer) error {
	data, err := ioutil.ReadFile(file.path)
	if err != nil {
		return err
	}
	path, info := file.path, file.info

	key := path[len(root)-1:]
	fmt.Fprintf(out, "\t\"%s\": &resource{\n", key)
	fmt.Fprintf(out, "\t\tsize: %d,\n", info.Size())
	fmt.Fprintf(out, "\t\tmtime: time.Unix(%d, 0),\n", info.ModTime().Unix())
	mtype := mime.TypeByExtension(filepath.Ext(path))
	fmt.Fprintf(out, "\t\tmtype: \"%s\",\n", mtype)
	fmt.Fprintf(out, "\t\tdata: {")
	for _, b := range data {
		fmt.Printf("%d, ", b)
		break
	}
	fmt.Println("\t\t},\n\t},")

	return nil
}

func die(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
	os.Exit(1)
}

func main() {
	root := flag.String("root", "", "resource root")
	flag.Parse()

	if len(*root) == 0 {
		die("<root> is required")
	}

	for file := range iterfiles(*root) {
		writeResource(*root, file, os.Stdout)
	}
}
