package listeners

import (
	"github.com/caitlin615/logmonitor/log"
)

// OutputChannel is a channel that accepts strings.
type OutputChannel chan string

// Listener is an interface for listening to a log Channel and providing a channel
// that it will output to
type Listener interface {
	// Start should start listening on the supplied log channel and return an OutputChannel
	// that the caller will listen and typically output that to somewhere
	Start(log.Channel) (recv OutputChannel)
}
