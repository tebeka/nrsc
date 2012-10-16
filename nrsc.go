package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ignoredDirs = map[string]bool{
	".svn": true,
	".hg":  true,
	".git": true,
}

type Path struct {
	path string
	info os.FileInfo
}


func iterfiles(root string) chan *Path {
	out := make(chan *Path)

	walkfn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if ignoredDirs[strings.ToLower(info.Name())] {
				return filepath.SkipDir
			}
		} else {
			if info.Name()[0] != '.' {
				out <- &Path{path, info}
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

func main() {
	root := flag.String("root", "", "resource root")
	flag.Parse()

	if len(*root) == 0 {
		fmt.Fprint(os.Stderr, "error: <root> is required\n")
		os.Exit(1)
	}

	for path := range iterfiles(*root) {
		fmt.Println(path.path)
	}

}
