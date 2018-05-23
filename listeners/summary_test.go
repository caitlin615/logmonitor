package listeners

import (
	// "fmt"
	"math/rand"
	"testing"

	"github.com/caitlin615/logmonitor/log"
)

func init() {
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	rand.Seed(99)
}

func TestSummaryReport(t *testing.T) {
	summary := NewSummaryListener()
	for i := 0; i < 100; i++ {
		summary.Add(log.RandomLine())
	}
	report, err := summary.Report()
	if err != nil {
		t.Error(err)
	}
	expected := SummaryReport{
		Section:        SummaryReportItem{"/user", 9},
		MostActiveUser: SummaryReportItem{"eddie", 20},
		Error4XX:       SummaryReportItem{"", 45},
		Error5XX:       SummaryReportItem{"", 22},
	}
	if report.Section.Key != expected.Section.Key && report.Section.Value != expected.Section.Value {
		t.Errorf("bad Section summary: got: %v", report.Section)
	}

	if report.MostActiveUser.Key != expected.MostActiveUser.Key && report.MostActiveUser.Value != expected.MostActiveUser.Value {
		t.Errorf("bad MostActiveUser summary: got: %v", report.MostActiveUser)
	}
	if report.Error4XX.Value != expected.Error4XX.Value {
		t.Errorf("bad Error4XX summary: got: %v", report.Error4XX)
	}
	if report.Error5XX.Value != expected.Error5XX.Value {
		t.Errorf("bad Error5XX summary: got: %v", report.Error5XX)
	}
}
