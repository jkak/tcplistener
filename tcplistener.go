package tcplistener

import (
	"errors"
	"net"
	"sync"
	"time"
)

// TCPListener for wrap tcp listener. and control the stop and mutex
// locked is false by default, set to true when Stop() call mutex.Lock()
type TCPListener struct {
	tcpListener *net.TCPListener
	stop        chan bool
	locked      bool
	mutex       *sync.Mutex
}

// NewTCPListener new a ptr of TCPListener for given Listener
func NewTCPListener(lisn net.Listener) (tcpLisn *TCPListener, err error) {
	listener, ok := lisn.(*net.TCPListener)
	if !ok {
		return nil, errors.New("assert error for base listener")
	}

	tcpLisn = &TCPListener{}
	tcpLisn.tcpListener = listener
	tcpLisn.stop = make(chan bool)
	tcpLisn.locked = false
	tcpLisn.mutex = new(sync.Mutex)
	return
}

// Accept accept a new net listener
func (tl *TCPListener) Accept() (conn net.Conn, err error) {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()

	for {
		// process the close signal
		select {
		case <-tl.stop:
			tl.tcpListener.Close()
			return nil, errors.New("tcp listener stopped by signal")
		default:
		}
		// set timeout before accept conn
		tl.tcpListener.SetDeadline(time.Now().Add(time.Second))

		conn, err = tl.tcpListener.Accept()
		if err != nil {
			errNet, ok := err.(net.Error)
			if ok && errNet.Timeout() && errNet.Temporary() {
				continue // continue while timeout
			}
		}
		return
	}
}

// AcceptTCP accept a new tcp listener
func (tl *TCPListener) AcceptTCP() (tcpc *net.TCPConn, err error) {
	tl.mutex.Lock()
	defer tl.mutex.Unlock()

	for {
		// process the close signal
		select {
		case <-tl.stop:
			tl.tcpListener.Close()
			return nil, errors.New("tcp listener stopped by signal")
		default:
		}

		// set timeout before accept conn
		tl.tcpListener.SetDeadline(time.Now().Add(time.Second))

		tcpc, err = tl.tcpListener.AcceptTCP()
		if err != nil {
			errNet, ok := err.(net.Error)
			if ok && errNet.Timeout() && errNet.Temporary() {
				continue // continue while timeout
			}
		}
		return
	}
}

// Start to enable mutex
func (tl *TCPListener) Start() {
	go func() {
		// tl.locked for fix bug of directly Unlock here.
		// bug: "panic: sync: unlock of unlocked mutex"
		if tl.locked {
			tl.mutex.Unlock()
			tl.locked = false
		}
	}()
}

// Stop to disable mutex
func (tl *TCPListener) Stop() {
	tl.stop <- true
	go func() {
		if !tl.locked {
			tl.mutex.Lock()
			tl.locked = true
		}
	}()
}
