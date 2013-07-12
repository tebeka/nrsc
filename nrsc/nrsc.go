package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

const (
	Version = "0.4.0"
)

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func exists(path string, dir bool) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	if dir {
		return fi.IsDir()
	} else {
		return !fi.IsDir()
	}
}

func tmpZip() string {
	return fmt.Sprintf("%s/nrsc-%d.zip", os.TempDir(), time.Now().Unix())
}

func mkzip(root, zip string, zipArgs []string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err = os.Chdir(root); err != nil {
		return err
	}
	defer os.Chdir(pwd)

	args := []string{"-r", zip, "."}
	args = append(args, zipArgs...)

	cmd := exec.Command("zip", args...)
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

func appendZip(exe, zipFile string) error {
	zipfo, err := os.Open(zipFile)
	if err != nil {
		return err
	}
	defer zipfo.Close()

	exefo, err := os.OpenFile(exe, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}
	defer exefo.Close()

	_, err = exefo.Seek(0, os.SEEK_END)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipfo, exefo)
	return err
}

func fixZipOffset(exe string) error {
	return exec.Command("zip", "-q", "-A", exe).Run()
}

func main() {
	flag.Usage = func() {
        fmt.Println("usage: nrsc EXECTABLE RESOURCE_DIR [ZIP OPTIONS]")
	}

	version := flag.Bool("version", false, "show version and exit")
	flag.Parse()

	if *version {
		fmt.Printf("nrsc version %s\n", Version)
	}

	if flag.NArg() < 2 {
		die("error: Wrong number of arguments\n")
	}

	exe := flag.Arg(0)
	if !exists(exe, false) {
		die("error: `%s` is not a file", exe)
	}

	root := flag.Arg(1)
	if !exists(root, true) {
		die("error: `%s` is not a directory", root)
	}

	zip := tmpZip()
	defer os.Remove(zip)

	if err := mkzip(root, zip, flag.Args()[2:]); err != nil {
		die("error: can't create zip - %s", err)
	}

	if err := appendZip(exe, zip); err != nil {
		die("error: can't append zip to %s - %s", exe, err)
	}

	if err := fixZipOffset(exe); err != nil {
		die("error: can't fix zip offset in %s - %s", exe, err)
	}
}
