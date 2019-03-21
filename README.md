# HTTP log monitoring console program

[![Go Report Card](https://goreportcard.com/badge/github.com/caitlin615/logmonitor)](https://goreportcard.com/report/github.com/caitlin615/logmonitor)
[![Coverage Status](https://coveralls.io/repos/github/caitlin615/logmonitor/badge.svg?branch=master)](https://coveralls.io/github/caitlin615/logmonitor?branch=master)
[![GoDoc](https://godoc.org/github.com/caitlin615/logmonitor?status.svg)](https://godoc.org/github.com/caitlin615/logmonitor)


Consume an actively written-to [w3c-formatted HTTP access log](https://en.wikipedia.org/wiki/Common_Log_Format).
It should default to reading `/var/log/access.log` and be overridable.
Example log lines:

```
127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 1234
127.0.0.1 - jill [09/May/2018:16:00:41 +0000] "GET /api/user HTTP/1.0" 200 1234
127.0.0.1 - frank [09/May/2018:16:00:42 +0000] "GET /api/user HTTP/1.0" 200 1234
127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "GET /api/user HTTP/1.0" 200 1234
```

* Display stats every 10s about the traffic during those 10s:
  * the sections of the web site with the most hits, as well as interesting summary statistics on the traffic as a whole.
  * A section is defined as being what's before the second '/' in the path.
  * For example, the section for "http://my.site.com/pages/create” is "http://my.site.com/pages".
* Make sure a user can keep the app running and monitor the log file continuously
* Whenever total traffic for the past 2 minutes exceeds a certain number on average, add a message saying that `High traffic generated an alert - hits = {value}, triggered at {time}`.
  * The default threshold should be 10 requests per second and should be overridable.
* Whenever the total traffic drops again below that value on average for the past 2 minutes, print or displays another message detailing when the alert recovered.

## Run via Docker

### Build the Docker container
```
docker build -t caitlin615:logmonitor .
```

### Run with defaults

```
docker run --rm -it caitlin615:logmonitor
```

### Run with custom `*.log` file

```
docker run --rm -it caitlin615:logmonitor -filename myAccessFile.log
```

### Run with custom high traffic alert threshold (requests per second)

```
docker run --rm -it -e ALERT_REQ_PER_SECOND_THRESHOLD="30" caitlin615:logmonitor
```

## Generating random log entries

Where `myAccessFile.log` is the filename where the script should write the logs to.

```
go run cmd/generate_logs/main.go -filename myAccessFile.log
```

Or with Docker, run this command to override the entrypoint in the Dockerfile

```
docker run --rm -it --entrypoint go caitlin615:logmonitor run cmd/generate_logs/main.go -filename myAccessFile.log
```

### Things I didn't get to but would like to have done
- [ ] Mutexes when handling list of logs within goroutines
- [ ] Error channels
