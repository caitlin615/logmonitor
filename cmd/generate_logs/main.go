package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/caitlin615/logmonitor/log"
)

var logFilename = flag.String("filename", "/var/log/access.log", "Log filename to read from")

func main() {
	flag.Parse()

	fmt.Println("Starting random log writer...")
	file, err := os.OpenFile(*logFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	clock := time.Tick(10 * time.Millisecond)
	for _ = range clock {
		// Randomly skip
		r := rand.Int()
		if r%2 == 0 {
			continue
		}
		line := log.RandomLine()
		_, err := file.WriteString(line.String() + "\n")
		if err != nil {
			panic(err)
		}
	}
}
