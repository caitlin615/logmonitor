package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/caitlin615/logmonitor/listeners"
	"github.com/caitlin615/logmonitor/log"
)

var logFilename = flag.String("filename", "/var/log/access.log", "Log filename to read from")

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	alertReqPerSecondThreshold := mustParseInt(getEnvDefault("ALERT_REQ_PER_SECOND_THRESHOLD", "10"))

	fmt.Println("Starting...")

	// Open the file and create a new buffer for reading the contents
	logReader, roFile, err := NewLogReader(*logFilename)
	if err != nil {
		panic(err)
	}

	// Create a log listening channel and start listening to the log reader buffer
	listenChan := make(log.Channel)
	go listenChan.Tail(logReader)

	// TODO: Would be nice to have an error channel that these listeners can write to
	// if they encounter an error and we can decide here to panic or continue
	summary := listeners.NewSummaryListener()
	summaryRecv := summary.Start(listenChan)

	alert := listeners.NewAlertListener(int64(alertReqPerSecondThreshold))
	alertRecv := alert.Start(listenChan)

	// goroutine will listen for receives from both the alert and summary output channels
	// and print whatever they get
	go func() {
		for {
			select {
			case in := <-summaryRecv:
				fmt.Println(in)
			case in := <-alertRecv:
				fmt.Println(in)
			}
		}
	}()

	// Handle ctrl-c because the interrupt signal skips the deferred calls,
	// make sure everything is closed down before exiting
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	for _ = range c {
		fmt.Println("Interrupt received, shutting down cleanly")
		roFile.Close()
		os.Exit(0)
	}
}

// NewLogReader opens the file as read-only, seeks to the end, and returns a bufio reader
func NewLogReader(filename string) (reader *bufio.Reader, file *os.File, err error) {
	file, err = os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	// Seek to end of file
	_, err = file.Seek(0, 2)
	if err != nil {
		return
	}
	reader = bufio.NewReader(file)
	return
}

// the following are helper functions for parsing environment variables
func getEnvDefault(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok && len(v) > 0 {
		return v
	}
	return defaultValue
}

func mustParseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return i
}

func mustParseBool(value string) bool {
	b, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}
	return b
}
