package tcplistener

import (
	"log"
	"net"
	"testing"
	"time"
)

var (
	port = ":8080"
)

func TestTCPListnerAccept(t *testing.T) {
	var name = "Accept()"

	base, err := net.Listen("tcp", port)
	if err != nil {
		t.Fatalf("Listen base err:%s", err)
	}
	tcpLisn, err := NewTCPListener(base)
	if err != nil {
		t.Fatalf("new Listen err:%s", err)
	}
	tcpLisn.Start()
	go func() {
		for {
			// listener accept request here.
			newConn, err := tcpLisn.Accept()
			if err != nil {
				// I want see the stop signal status
				log.Printf("-- %s err : %s", name, err)
				if tcpLisn.locked {
					break
				}
			}
			// sleep and Close()
			if newConn != nil {
				newConn.Close()
			}
		}
	}()

	// request one time to see the locked status.
	var dia net.Dialer
	c, err := dia.Dial("tcp", port)
	if err != nil {
		t.Logf("%s err : %s", name, err)
	}
	network := c.LocalAddr().Network()
	laddr := *c.LocalAddr().(*net.TCPAddr)
	defer c.Close()
	log.Printf("%s lock status: %v\n", name, tcpLisn.locked)

	n := 8000
	if testing.Short() {
		n = 100
	}
	failed := 0
	for i := 0; i < n; i++ {
		var dia net.Dialer
		dia.Timeout = time.Millisecond * 20
		c, err := dia.Dial(network, port)
		if err != nil {
			t.Logf("current dial should be failed")
			failed++
			continue
		}
		addr := c.LocalAddr().(*net.TCPAddr)
		// log.Println("port:", addr.Port, laddr.Port, addr.IP.Equal(laddr.IP))
		if addr.Port == laddr.Port || !addr.IP.Equal(laddr.IP) {
			failed++
			t.Logf("current dial %v should be failed", addr)
		}
		c.Close()
	}
	tcpLisn.Stop()
	time.Sleep(time.Millisecond * 10)
	log.Printf("-- failed times: %d; fail rate: %f\n", failed, float64(failed)/float64(n))
	log.Printf("%s lock status: %v\n\n", name, tcpLisn.locked)
}
