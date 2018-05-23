package log

import "testing"
import "math/rand"

func init() {
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	rand.Seed(99)
}
func TestNewLogLine(t *testing.T) {
	line := RandomLine()

	if line.IPAddress != "171.229.234.222" {
		t.Error(line.String())
	}
	if line.Request.URL != "http://www.example.com/api/api/user" {
		t.Errorf("bad request url: %s", line.Request.URL)
	}
	section, err := line.Request.Section()
	if err != nil {
		t.Error(err)
	}
	if section != "http://www.example.com/api" {
		t.Errorf("bad request section: %s", section)
	}
}

func TestLinesClearBefore(t *testing.T) {
	lines := Lines{}
	for i := 0; i < 6; i++ {
		lines = append(lines, RandomLine())
	}

	someLine := lines[2]
	lines = lines.ClearBefore(someLine.Date)

	firstLine := lines[0]
	if firstLine.String() != someLine.String() {
		t.Errorf("ClearBefore failed:\nexpected: %s\ngot: %s", someLine.String(), firstLine.String())
	}
}

func TestLinesMostActiveUser(t *testing.T) {
	lines := Lines{}
	for i := 0; i < 20; i++ {
		lines = append(lines, RandomLine())
	}

	user, hits := lines.MostActiveUser()
	if user != "sadie" && hits != 5 {
		t.Errorf("incorrect MostActiveUser: %s %d", user, hits)
	}
}
