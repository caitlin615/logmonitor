package listeners

import (
	"errors"
	"fmt"
	"time"

	"github.com/caitlin615/logmonitor/log"
)

// This ensures adherence to the Listener interface
var _ = Listener(Alert{})

var (
	// ErrInHighTrafficState is the error returned when the listener is has previously reported that there's high traffic
	ErrInHighTrafficState = errors.New("Already in high traffic state")
	// ErrLowTrafficState is the error returned when the threhold for requests per second hasn't been met
	ErrLowTrafficState = errors.New("Low traffic state")
)

// Alert is a Listener that will output summary reports
type Alert struct {
	triggerInterval   time.Duration
	thresholdInterval time.Duration
	logs              log.Lines
	rpsThreshold      int64

	isInHighAlertState bool
}

// NewAlertListener returns an Alert listener with the specified requests per second threshold.
// TODO: Defaults to a 2 minute threshold interval, should this be configurable?
func NewAlertListener(reqPerSecondThreshold int64) Alert {
	return Alert{
		triggerInterval:   10 * time.Second,
		thresholdInterval: 2 * time.Minute,
		rpsThreshold:      reqPerSecondThreshold,
	}
}

// Report returns the summary report during based on what's currently in the logs
func (a *Alert) Report() (string, error) {
	// TODO: Mutexes when reading a.logs
	now := time.Now().UTC()
	start := now.Add(-a.thresholdInterval)
	var count int64

	// Count the number of logs between now and the threshold
	for _, line := range a.logs {
		if start.Before(line.Date.UTC()) {
			count++
		}
	}

	averageReqPerSec := count / int64(a.thresholdInterval.Seconds())
	highTraffic := averageReqPerSec > a.rpsThreshold

	// Previous report determined that we are in a high traffic state
	if a.isInHighAlertState {
		// Now we aren't getting high traffic, so send an alert that it's over
		if !highTraffic {
			a.isInHighAlertState = false
			return fmt.Sprintf("High traffic state ended: %s", now), nil
		}
		return "", ErrInHighTrafficState
	}

	if highTraffic {
		a.isInHighAlertState = true
		return fmt.Sprintf("High traffic generated an alert - hits = %d, triggered at %s", count, now), nil
	}
	return "", ErrLowTrafficState
}

// Start starts the Alert listener
func (a Alert) Start(listenChan log.Channel) OutputChannel {
	recv := make(OutputChannel)
	// start a goroutine that will listen for log entries
	go func() {
		for {
			select {
			case in := <-listenChan:
				// TODO: Use mutexes since this will be accessed across multiple goroutines
				a.logs = append(a.logs, in)
			}
		}
	}()

	// Start a goroutine to send a report into the output channel every X seconds based on the trigger time
	go func() {
		clock := time.Tick(a.triggerInterval)
		for range clock {
			if report, err := a.Report(); err == nil {
				recv <- report
			}
		}
	}()

	return recv
}
