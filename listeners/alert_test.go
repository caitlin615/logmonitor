package listeners

import (
	"strings"
	"testing"
	"time"

	"github.com/caitlin615/logmonitor/log"
)

func addLines(num int) log.Lines {
	ll := make(log.Lines, num)
	for i := 0; i < num; i++ {
		line := log.RandomLine()
		line.Date = line.Date.Add(time.Duration(-i) * time.Millisecond)
		ll = append(ll, line)
	}
	return ll
}

func TestAlertReport(t *testing.T) {

	alert := NewAlertListener(5)
	alert.thresholdInterval = 1 * time.Minute
	alert.logs = append(alert.logs, addLines(50)...)

	_, err := alert.Report()
	if err != ErrLowTrafficState {
		t.Errorf("expected low alert error, got: %v", err)
	}

	alert.logs = append(alert.logs, addLines(500)...)
	report, err := alert.Report()
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(report, "High traffic generated an alert") {
		t.Errorf("expected high alert report, got: %s", report)
	}

	alert.logs = append(alert.logs, addLines(500)...)
	_, err = alert.Report()
	if err != ErrInHighTrafficState {
		t.Errorf("expected already in high traffic state error, got: %v", err)
	}

	// Clear them out and add a few
	// This is as if there was no traffic over the previous interval, then some requests
	// came in and we're back below the threshold
	alert.logs.Clear()
	alert.logs = append(alert.logs, addLines(10)...)

	report, err = alert.Report()
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(report, "High traffic state ended") {
		t.Errorf("expected high traffic ended report, got: %s", report)
	}
}
