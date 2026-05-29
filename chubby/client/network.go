// Assumptions the client makes about addresses of Chubby servers.

package client

// We assume that any Chubby node must have one of these addresses.
// Yes this is gross but we're doing it anyway because of time constraints
var PossibleServerAddrs = map[string]bool{
	"127.0.0.1:5379": true,
	"127.0.0.1:6379": true,
	"127.0.0.1:7379": true,
}
