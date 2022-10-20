package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"github.com/mattn/go-isatty"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type flusher interface {
	Flush() error
}

func main() {
	//filePath()
	//buffering()
	//useIsatty()
	multiWrite()
}

func multiWrite() {
	tmp, _ := ioutil.TempFile(os.TempDir(), "tmp")
	defer tmp.Close()

	hash := sha256.New()

	w := io.MultiWriter(tmp, hash)

	written, _ := io.Copy(w, os.Stdin)

	fmt.Printf("Wrote %d bytes to %s \nSHA256: %x \n",
		written,
		tmp.Name(),
		hash.Sum(nil),
	)
}

func useIsatty() {
	var output io.Writer
	if isatty.IsTerminal(os.Stdout.Fd()) {
		output = os.Stdout
	} else {
		output = bufio.NewWriter(os.Stdout)
	}

	for i := 0; i < 100; i++ {
		fmt.Fprintln(output, strings.Repeat("x", 100))
	}
	if _o, ok := output.(flusher); ok {
		_o.Flush()
	}
}

func buffering() {
	b := bufio.NewWriter(os.Stdout)
	for i := 0; i < 100; i++ {
		fmt.Fprintln(b, strings.Repeat("x", 100))
	}
	b.Flush()
}

func filePath() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if ok, err := path.Match("/data/*.html", r.URL.Path); err != nil || !ok {
			http.NotFound(w, r)
			return
		}

		name := filepath.Join(cwd, "data", filepath.Base(r.URL.Path))
		f, err := os.Open(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()
		io.Copy(w, f)
	})
	http.ListenAndServe(":8080", nil)
}
