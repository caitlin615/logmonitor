FROM golang:1.10

RUN touch /var/log/access.log # since the program will read this by default

WORKDIR /go/src/github.com/caitlin615/logmonitor

ADD . /go/src/github.com/caitlin615/logmonitor

ENTRYPOINT ["go", "run", "main.go"]
