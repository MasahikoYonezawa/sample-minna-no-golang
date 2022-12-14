package main

import (
	"bufio"
	"context"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/mattn/go-isatty"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

type flusher interface {
	Flush() error
}

var globalWG sync.WaitGroup

func main() {
	//filePath()
	//buffering()
	//useIsatty()
	//multiWrite()
	//mathRand()
	//cryptRand()
	//useHumanize()
	//tr(os.Stdin, os.Stdout, os.Stderr)
	//stopRoutine()
	//stopRoutineWitContext()
	signalHandling()
}

func signalHandling() {
	defer fmt.Println("done")
	trapSignals := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, trapSignals...)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigCh
		fmt.Println("Got signal", sig)

		cancel()
	}()

	doMain(ctx)
}

func doMain(ctx context.Context) {
	defer fmt.Println("done doMain")
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}
		fmt.Println("do something")
	}
}

func stopRoutineWitContext() {
	ctx, cancel := context.WithCancel(context.Background())
	queue := make(chan string)
	for i := 0; i < 2; i++ {
		globalWG.Add(1)
		go fetchURLWithContext(ctx, queue)
	}

	queue <- "https://www.example.com"
	queue <- "https://www.example.net"
	queue <- "https://www.example.net/foo"
	queue <- "https://www.example.net/bar"

	cancel()
	globalWG.Wait()
}

func fetchURLWithContext(ctx context.Context, queue chan string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker exit")
			globalWG.Done()
			return
		case url2 := <-queue:
			fmt.Println("fetching", url2)
		}
	}
}

func stopRoutine() {
	queue := make(chan string)
	for i := 0; i < 2; i++ {
		globalWG.Add(1)
		go fetchURL(queue)
	}

	queue <- "https://www.example.com"
	queue <- "https://www.example.net"
	queue <- "https://www.example.net/foo"
	queue <- "https://www.example.net/bar"

	close(queue)
	globalWG.Wait()
}

func fetchURL(queue chan string) {
	for {
		url, more := <-queue
		if more {
			fmt.Println("fetching", url)
		} else {
			fmt.Println("worker exit")
			globalWG.Done()
			return
		}
	}
}

func tr(src io.Reader, dst io.Writer, errDst io.Writer) error {
	cmd := exec.Command("tr", "a-z", "A-Z")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		_, err = io.Copy(stdin, src)
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.EPIPE {
			fmt.Println("ignore EPIPE")
		} else if err != nil {
			log.Println("failed to write to STDIN", err)
		}
		stdin.Close()
		wg.Done()
	}()

	go func() {
		io.Copy(dst, stdout)
		stdout.Close()
		wg.Done()
	}()

	go func() {
		io.Copy(errDst, stderr)
		stderr.Close()
		wg.Done()
	}()

	wg.Wait()

	return cmd.Wait()
}

func useHumanize() {
	name := os.Args[1]
	s, _ := os.Stat(name)
	fmt.Printf(
		"%s: %s \n",
		name,
		humanize.Bytes(uint64(s.Size())),
	)
}

func cryptRand() {
	var s int64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &s); err != nil {
		s = time.Now().UnixNano()

	}
	rand.Seed(s)
	n := rand.Intn(100)
	fmt.Println(n)
}

func mathRand() {
	rand.Seed(42)
	n := rand.Intn(100)
	fmt.Println(n)
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
