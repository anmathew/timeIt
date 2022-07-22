package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	ver       = "1.03c"
	binName   = "timeIt"
	tag       = binName + "/" + ver
	logFile   = "timeIt.log"
	winFolder = "C:\\Temp\\timeIt\\"
	linFolder = "/tmp/timeIt/"
)

var (
	folder = linFolder
)

func main() {
	if len(os.Args) > 1 {
		if strings.Contains(os.Args[1], "--timeIt.help") {
			showhelp()
			os.Exit(0)
		}
		if strings.Contains(os.Args[1], "--timeIt.ver") {
			fmt.Printf("timeIt: %v\n", ver)
			os.Exit(0)
		}
	}

	binBase := filepath.Base(os.Args[0])

	wg := &sync.WaitGroup{}
	wg.Add(1)
	var ourSum string
	go func() {
		defer wg.Done()
		ourSum = getShasum(os.Args[0])
	}()

	if runtime.GOOS == "windows" {
		folder = winFolder
	}
	cmdToRun := fmt.Sprintf("%s%s", folder, binBase)
	cmdToRun = filepath.FromSlash(cmdToRun)

	cmdSum := getShasum(cmdToRun)
	wg.Wait()

	if ourSum == cmdSum {
		fmt.Fprintf(os.Stderr, "Err 3.1: You cannot evoke %s directly\n", binName)
		os.Exit(1)
	}

	t1 := time.Now()
	cmd := exec.Command(cmdToRun, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	checkErr(cmd.Run())
	t2 := time.Now()
	delta := t2.Sub(t1)
	writeThis := fmt.Sprintf("%v %v took %v to run cmd: %v\n", time.Now().Local().UnixMicro(), tag, delta, cmd)

	// log results to logFile
	appendTo := fmt.Sprintf("%s%s", folder, logFile)
	appendTo = filepath.FromSlash(appendTo)
	f, err := os.OpenFile(appendTo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkErr(err)
	defer f.Close()
	_, err = f.WriteString(writeThis)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}

func showhelp() {
	fmt.Printf("timeIt: runs and times your binary and logs the results\n")
	fmt.Println("See https://github.com/anmathew/timeIt")
	fmt.Println("Usage:")
	fmt.Println("--timeIt.ver   Show current version and exit")
	fmt.Println("--timeIt.help  Show this help section and exit")
}

func getShasum(filename string) (sha string) {
	f, err := os.Open(filename)
	if err != nil {
		checkErr(err)
		return
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		checkErr(err)
		return
	}
	sha = fmt.Sprintf("%x", h.Sum(nil))
	return
}
