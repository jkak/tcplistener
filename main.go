package tcplistener

import (
	"net"
	"sync"
)

// TCPListener for wrap tcp listener. and control the stop and mutex
type TCPListener struct {
	*net.TCPListener
	stop  chan bool
	mutex sync.Mutex
}

// NewTCPListener new a ptr of TCPListener
func NewTCPListener(ls net.Listener) (tcpl *TCPListener, err error) {
	return
}

// Accept accept a new net listener after timeout
func (tl *TCPListener) Accept() (conn net.Conn, err error) {
	return
}

// AcceptTCP accept a new tcp listener after timeout
func (tl *TCPListener) AcceptTCP() (tcpc *net.TCPConn, err error) {
	return
}

// Start to enable mutex
func (tl *TCPListener) Start() {
}

// Stop to disable mutex
func (tl *TCPListener) Stop() {
}
