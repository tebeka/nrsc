package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

const (
	root = "/tmp/nrsc-test"
	port = 9888
)

func TestText(t *testing.T) {
	expected := map[string]string{
		"Content-Size": "12",
		"Content-Type": "text/plain; charset=utf-8",
	}
	checkPath(t, "ht.txt", expected)
}

func TestSub(t *testing.T) {
	expected := map[string]string{
		"Content-Size": "1150",
		"Content-Type": "image/x-icon",
	}
	checkPath(t, "sub/favicon.ico", expected)
}

func createMain() error {
	filename := fmt.Sprintf("%s/main.go", root)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, code, port)
	return nil
}

func initDir() error {
	// Ignore error value, since it might not be there
	os.RemoveAll(root)

	err := os.Mkdir(root, 0777)
	if err != nil {
		return err
	}

	return createMain()
}

func build() {
	cmd := exec.Command("go", "build")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func get(path string) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:%d/static/%s", port, path)
	return http.Get(url)
}

func startServer(t *testing.T) *exec.Cmd {
	cmd := exec.Command(fmt.Sprintf("%s/nrsc-test", root))
	// Ignore errors, test will fail anyway if server not running
	cmd.Start()

	// Wait for server
	url := fmt.Sprintf("http://localhost:%d", port)
	start := time.Now()
	for time.Since(start) < time.Duration(2*time.Second) {
		_, err := http.Get(url)
		if err == nil {
			return cmd
		}
		time.Sleep(time.Second / 10)
	}

	if cmd.Process != nil {
		cmd.Process.Kill()
	}
	t.Fatalf("can't connect to server")
	return nil
}

func init() {
	build()

	if err := initDir(); err != nil {
		panic(err)
	}

	cwd, _ := os.Getwd()
	path := func(name string) string {
		return fmt.Sprintf("%s/%s", cwd, name)
	}
	os.Chdir(root)
	defer os.Chdir(cwd)

	cmd := exec.Command(path("nrsc"), "-root", path("test-resources"))
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	build()
}

func checkHeaders(t *testing.T, expected map[string]string, headers http.Header) {
	for key := range expected {
		v1 := expected[key]
		v2 := headers.Get(key)
		if v1 != v2 {
			t.Fatalf("bad header %s: %s <-> %s", key, v1, v2)
		}
	}

	key := "Last-Modified"
	value := headers.Get(key)
	if value == "" {
		t.Fatalf("no %s header", key)
	}
}

func checkPath(t *testing.T, path string, expected map[string]string) {
	server := startServer(t)
	if server == nil {
		return
	}
	defer server.Process.Kill()

	resp, err := get(path)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("bad reply - %s", resp.Status)
	}

	checkHeaders(t, expected, resp.Header)
}

const code = `
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
	if err := http.ListenAndServe(":%d", nil); err != nil {
		fmt.Fprintf(os.Stderr, "error: %%s\n", err)
		os.Exit(1)
	}
}
`
