package listeners

import (
	"fmt"
	"time"

	"github.com/caitlin615/logmonitor/log"
)

// This ensures adherence to the Listener interface
var _ = Listener(Summary{})

// Summary is a Listener that will output summary reports
type Summary struct {
	triggerInterval time.Duration
	logs            log.Lines
}

// NewSummaryListener returns an Summary listener that will report every 10 seconds
// TODO: Defaults to a 10 seconds trigger interval, should this be configurable?
func NewSummaryListener() Summary {
	return Summary{triggerInterval: 10 * time.Second}
}

// Add appends the line to the log storage
func (s *Summary) Add(line log.Line) {
	s.logs = append(s.logs, line)
}

// Report returns the summary report during based on what's currently in the logs
func (s *Summary) Report() (report SummaryReport, err error) {
	section, hits := s.logs.SectionWithMostHits()
	if len(s.logs) == 0 {
		err = fmt.Errorf("no requests available to summarize")
		return
	}

	report.Section = SummaryReportItem{section, hits}
	mau, mauCount := s.logs.MostActiveUser()
	report.MostActiveUser = SummaryReportItem{mau, mauCount}
	report.Error4XX = SummaryReportItem{"Requests with error code 4XX", s.logs.ErrorCode4XX()}
	report.Error5XX = SummaryReportItem{"Requests with error code 5XX", s.logs.ErrorCode5XX()}
	defer s.logs.Clear()
	return
}

// Start starts the Summary listener
func (s Summary) Start(listenChan log.Channel) OutputChannel {
	recv := make(OutputChannel)
	// start a goroutine that will listen for log entries
	go func() {
		for {
			select {
			case in := <-listenChan:
				// TODO: Use mutexes since this will be accessed across multiple goroutines
				s.Add(in)
			}
		}
	}()

	// Start a goroutine to send a report into the output channel every X seconds based on the trigger time
	go func() {
		clock := time.Tick(s.triggerInterval)
		for _ = range clock {
			if report, err := s.Report(); err == nil {
				recv <- report.String()
			}
		}
	}()

	return recv
}

// SummaryReport is the data structure that holds all the information for the report
type SummaryReport struct {
	Section        SummaryReportItem
	MostActiveUser SummaryReportItem
	Error4XX       SummaryReportItem
	Error5XX       SummaryReportItem
}

func (sr SummaryReport) String() string {
	return fmt.Sprintf(`Section with the most hits: %s (%d),
* Most Active User: %s (%d)
* %s: %d
* %s: %d
`, sr.Section.Key, sr.Section.Value,
		sr.MostActiveUser.Key, sr.MostActiveUser.Value,
		sr.Error4XX.Key, sr.Error4XX.Value,
		sr.Error5XX.Key, sr.Error5XX.Value)
}

// SummaryReportItem ...
type SummaryReportItem struct {
	Key   string
	Value int
}
