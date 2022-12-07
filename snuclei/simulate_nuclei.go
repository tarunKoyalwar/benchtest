package main

import (
	"crypto/sha256"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"

	"github.com/remeh/sizedwaitgroup"
)

var targetUrl string = "http://192.168.29.119:9000/"
var sizedwg bool = false
var errorCount int = 0
var m *sync.Mutex = &sync.Mutex{}
var totaltasks int = 0
var savetraceAt string
var fixedworkers bool = false
var concurrency int
var wg *sync.WaitGroup = &sync.WaitGroup{}

var sometext string = `Lorem ipsum dolor sit amet,
consectetur adipiscing elit, sed do eiusmod tempor 
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, 
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo 
consequat. Duis aute irure dolor in reprehenderit in voluptate velit 
esse cillum dolore eu fugiat nulla pariatur.Excepteur sint occaecat 
cupidatat non proident, sunt in culpa qui
officia deserunt mollit anim id est laborum`

func main() {

	c := flag.Int("cores", 4, "Value will be used by GOMAXPROCS")
	flag.StringVar(&targetUrl, "u", "http://192.168.29.119:9000/", "Target URL for benchmarking")
	flag.IntVar(&totaltasks, "count", 1000, "Total Tasks count")
	flag.StringVar(&savetraceAt, "trace", "", "Save trace file at location")
	flag.IntVar(&concurrency, "c", 2000, "Experimental concurrency")
	flag.BoolVar(&fixedworkers, "fixed", false, "Fixed/Bounded Goroutines")
	flag.BoolVar(&sizedwg, "sizedwg", false, "Uses sized waitgroup same as nuclei")

	flag.Parse()

	t := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	http.DefaultClient.Transport = &t

	runtime.GOMAXPROCS(*c)

	if savetraceAt != "" {
		f, er := os.OpenFile(savetraceAt, os.O_CREATE|os.O_WRONLY, 0644)
		if er != nil {
			panic(er)
		}
		trace.Start(f)
		defer trace.Stop()
		defer f.Close()
	}

	if sizedwg {
		fmt.Println("[+]Simulating SizedWaitgroup Goroutines...")
		SizedWaitGroups()
		fmt.Printf("[-]Go routines after scheduling all tasks %v\n", runtime.NumGoroutine())
		fmt.Printf("[-]Number of errors %v\n", errorCount)
		os.Exit(0)
	}

	if !fixedworkers {
		fmt.Println("[+]Simulating Unlimited Goroutines...")
		UnlimitedGoroutines()
	} else {
		fmt.Println("[+]Simulating Bounded Goroutines...")
		WorkerGoroutines()
	}

	fmt.Printf("[-]Go routines after scheduling all tasks %v\n", runtime.NumGoroutine())
	wg.Wait()
	fmt.Printf("[-]Number of errors %v\n", errorCount)

}

func UnlimitedGoroutines() {
	for i := 0; i < totaltasks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nucleitask()
		}()
	}
}

func SizedWaitGroups() {
	sized := sizedwaitgroup.New(concurrency)
	for i := 0; i < totaltasks; i++ {
		sized.Add()
		go func() {
			defer sized.Done()
			nucleitask()
		}()
	}
	sized.Wait()
}

func WorkerGoroutines() {
	var ch chan func() = make(chan func(), concurrency)

	worker := func(chx <-chan func()) {
		defer wg.Done()
		// recieve only channel
		for {
			w, ok := <-ch
			if !ok {
				//closed
				break
			}
			// execute assigned work
			w()
		}
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(ch)
	}

	// assign tasks
	for i := 0; i < totaltasks; i++ {
		ch <- nucleitask
	}

	close(ch)
}

func nucleitask() {

	resp, err := http.Get(targetUrl)
	if err != nil {
		m.Lock()
		errorCount += 1
		m.Unlock()
	}
	if resp != nil {
		io.Copy(io.Discard, resp.Body)
		defer resp.Body.Close()
	}

	cpuintensive()
}

func cpuintensive() {
	sha := sha256.New()

	for i := 0; i < 1000; i++ {
		sha.Sum([]byte(sometext))
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 1000; i++ {
		sha.Sum([]byte(sometext))
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 1000; i++ {
		sha.Sum([]byte(sometext))
	}

}
