package tcplistener

import (
	"errors"
	"net"
	"sync"
	"time"
)

// TCPListener for wrap tcp listener. and control the stop and mutex
type TCPListener struct {
	tcpListener *net.TCPListener
	stop        chan bool
	mutex       sync.Mutex
}

// NewTCPListener new a ptr of TCPListener for given Listener
func NewTCPListener(lisn net.Listener) (tcpLisn *TCPListener, err error) {
	listener, ok := lisn.(*net.TCPListener)
	if !ok {
		return nil, errors.New("assert error for base listener")
	}
	tcpLisn.tcpListener = listener
	tcpLisn.stop = make(chan bool)
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
		tl.mutex.Unlock()
	}()
}

// Stop to disable mutex
func (tl *TCPListener) Stop() {
	tl.stop <- true
	go func() {
		tl.mutex.Lock()
	}()
}
