package log

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// RandomLine returns a randomly generated Line.
func RandomLine() Line {
	hosts := []string{"http://example.com", "", "https://example.com", "http://www.example.com", "https://www.example.com"}
	paths := []string{"report", "club", "api", "user"}
	users := []string{"frank", "jack", "eddie", "sadie", "leo", "lola", "lily"}
	methods := []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"}
	statusCodes := []int{100, 101, 102, 200, 201, 202, 203, 204, 205, 206, 207, 208, 226, 300, 301, 302, 303, 304, 305, 307, 308, 400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418, 421, 422, 423, 424, 426, 428, 429, 431, 444, 451, 499, 500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511, 599}
	size := rand.Intn(5000)

	url := hosts[rand.Intn(len(hosts))]
	pathPerms := rand.Intn(len(paths)) + 1 // always make sure there's at least one path
	for i := 0; i < pathPerms; i++ {
		url += "/" + paths[rand.Intn(len(paths))]
	}

	ip := make([]string, 4)
	for i := range ip {
		ip[i] = strconv.Itoa(rand.Intn(255))
	}

	return Line{
		IPAddress: strings.Join(ip, "."),
		Identity:  "-",
		UserID:    users[rand.Intn(len(users))],
		Date:      time.Now().UTC(),
		Request: LineRequest{
			Method:   methods[rand.Intn(len(methods))],
			URL:      url,
			Protocol: "HTTP/1.0",
		},
		StatusCode: statusCodes[rand.Intn(len(statusCodes))],
		Size:       size,
	}
}
