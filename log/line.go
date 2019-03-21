package log

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/caitlin615/logmonitor/counter"
)

// Channel is a channel that accepts LogLines
type Channel chan Line

// Tail will listen on the reader and send all Lines into the Channel
// This should be run within a goroutine
func (lc Channel) Tail(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n') // TODO: Use reader.Readline()
		if err != nil && err != io.EOF {
			continue
		}
		line = strings.Trim(line, "\n")
		if len(line) > 0 {
			if logLine, err := NewLine(line); err == nil {
				// FIXME: Not sure why this needs to be in a goroutine, but otherwise it blocks
				// the for loop from continuing
				go logLine.Send(lc)
			}
		}
	}
}

const dateFormat = "02/Jan/2006:15:04:05 -0700"
const missingData = "-" // according to https://en.wikipedia.org/wiki/Common_Log_Format > A "-" in a field indicates missing data.

// LineRequest is the portion of the log line which contains request line from the client.
// It includes method of the request, the resource requested, and the HTTP protocol.
type LineRequest struct {
	Method   string
	URL      string
	Protocol string
}

func (lr *LineRequest) String() string {
	return fmt.Sprintf("%s %s %s", lr.Method, lr.URL, lr.Protocol)
}

// Line is a data structure that holds parsed information about a
// w3c-formatted HTTP access log (https://en.wikipedia.org/wiki/Common_Log_Format)
type Line struct {
	IPAddress  string
	Identity   string
	UserID     string
	Date       time.Time
	Request    LineRequest
	StatusCode int
	Size       int
}

// FIXME: Add tests specifically for this regex
var re = regexp.MustCompile(`([^ ]*) ([^ ]*) ([^ ]*) (?:-|\[([^\]]*)\]) \"(.*)\" (-|[0-9]{3}) ([0-9]*)`)

// ErrInvalidLine is the error if the line supplied did not match the regex
var ErrInvalidLine = errors.New("Invalid Line")

// NewLine returns a Line from a raw string
func NewLine(raw string) (Line, error) {
	line := Line{}
	parsed := re.FindStringSubmatch(raw)
	if len(parsed) != 8 {
		return line, ErrInvalidLine
	}
	req := NewLineRequest(parsed[5])
	line.IPAddress = parsed[1]
	line.Identity = parsed[2]
	line.UserID = parsed[3]
	line.Request = *req

	if statusCode, err := strconv.Atoi(parsed[6]); err == nil {
		line.StatusCode = statusCode
	}

	if size, err := strconv.Atoi(parsed[7]); err == nil {
		line.Size = size
	}

	if reqTime, err := time.ParseInLocation(dateFormat, parsed[4], time.UTC); err == nil {
		line.Date = reqTime
	}

	return line, nil
}

func (l *Line) String() string {
	// date will be 0 if it wasn't populated, so use the "-" for missing data unless
	// the date was populated and can be formatted
	date := missingData
	if !l.Date.IsZero() {
		date = l.Date.Format(dateFormat)
	}

	return fmt.Sprintf("%s %s %s [%s] \"%s\" %d %d",
		l.IPAddress,
		l.Identity,
		l.UserID,
		date,
		l.Request.String(),
		l.StatusCode,
		l.Size,
	)
}

// Send sends the Line into the Channel
func (l *Line) Send(c Channel) {
	c <- *l
}

// Section returns what's before the second '/' in the URL's path.
// For example, the section for "http://my.site.com/pages/create” is "http://my.site.com/pages".
func (lr *LineRequest) Section() (string, error) {
	u, err := url.Parse(lr.URL)
	if err != nil {
		return "", err
	}

	paths := strings.Split(u.EscapedPath(), "/")
	if len(paths) < 1 {
		// No path provided in the url, so we can't determine the section
		return "", fmt.Errorf("No path provided in the url, so the section cannot be determined. URL = %v", u)
	}
	u.Path = "/" + paths[1]
	return u.String(), nil
}

// NewLineRequest returns a LineRequest from a raw string
func NewLineRequest(raw string) *LineRequest {
	parsed := strings.Split(raw, " ")
	if len(parsed) < 3 {
		return &LineRequest{
			Method:   missingData,
			URL:      missingData,
			Protocol: missingData,
		}
	}
	return &LineRequest{
		Method:   parsed[0],
		URL:      parsed[1],
		Protocol: parsed[2],
	}
}

// Lines is a typalias for a list so methods can be used against it
type Lines []Line

// MostActiveUser returns the user with the most requests and the number of requests
func (ll *Lines) MostActiveUser() (string, int) {
	m := counter.New()
	m.SortByFunc = counter.SortDesc
	for _, line := range *ll {
		m.Increment(line.UserID)
	}
	sort.Sort(m)
	top := m.Dict[0]
	return top.Key, top.Value
}

// ErrorCode4XX returns the number of logs that have an error code 4xx
func (ll *Lines) ErrorCode4XX() int {
	count := 0
	for _, line := range *ll {
		if line.StatusCode >= 400 && line.StatusCode < 500 {
			count++
		}
	}
	return count
}

// ErrorCode5XX returns the number of logs that have an error code 5xx
func (ll *Lines) ErrorCode5XX() int {
	count := 0
	for _, line := range *ll {
		if line.StatusCode >= 500 {
			count++
		}
	}
	return count
}

// SectionWithMostHits returns the request section that has the largest occurrence
// and the number of time it appears (hits)
func (ll *Lines) SectionWithMostHits() (string, int) {
	// TODO: Mutexes when reading s.logs
	m := counter.New()
	for _, line := range *ll {
		if section, err := line.Request.Section(); err == nil {
			m.Increment(section)
		}
	}

	if m.Len() == 0 {
		return "", -1
	}

	m.SortByFunc = counter.SortDesc
	sort.Sort(m)
	top := m.Dict[0]
	return top.Key, top.Value
}

// Clear ...
func (ll *Lines) Clear() {
	empty := Lines{}
	*ll = empty
}

// ClearBefore returns a new Lines with all entries with a date before time t removed
// TODO: clear in place
func (ll Lines) ClearBefore(t time.Time) Lines {
	sort.Sort(ll)
	start := 0
	for i, item := range ll {
		if item.Date.After(t) {
			start = i - 1
			break
		}
	}
	return ll[start:]
}

// Len is part of sort.Interface.
func (ll Lines) Len() int {
	return len(ll)
}

// Swap is part of sort.Interface.
func (ll Lines) Swap(i, j int) {
	ll[i], ll[j] = ll[j], ll[i]
}

// Less is part of sort.Interface. Sort by Date
func (ll Lines) Less(i, j int) bool {
	return ll[i].Date.Before(ll[j].Date)
}
