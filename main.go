package main

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-isatty"
	"io"
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
	useIsatty()
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
