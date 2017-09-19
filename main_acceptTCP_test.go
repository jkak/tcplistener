package tcplistener

import (
	"log"
	"net"
	"testing"
	"time"
)

var (
	portT = ":8081"
	nameT = "AcceptTCP:"
)

func TestTCPListnerAcceptTCP(t *testing.T) {
	base, err := net.Listen("tcp", portT)
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
			newConn, err := tcpLisn.AcceptTCP()
			if err != nil {
				// I want see the stop signal
				log.Printf("-- %s err:%s", nameT, err)
				if tcpLisn.locked {
					break
				}
			}
			if newConn != nil {
				newConn.Close()
			}
		}
	}()

	var dia net.Dialer
	c, err := dia.Dial("tcp", portT)
	if err != nil {
		t.Logf("accept() err:%s", err)
	}
	network := c.LocalAddr().Network()
	laddr := *c.LocalAddr().(*net.TCPAddr)
	defer c.Close()
	log.Printf("%s lock status: %v\n", nameT, tcpLisn.locked)

	n := 8000
	if testing.Short() {
		n = 100
	}
	failedT := 0
	for i := 0; i < n; i++ {
		var dia net.Dialer
		dia.Timeout = time.Millisecond * 20
		c, err := dia.Dial(network, portT)
		if err != nil {
			// t.Error("Dial should fail")
			failedT++
			continue
		}
		addr := c.LocalAddr().(*net.TCPAddr)
		// log.Println("port:", addr.Port, laddr.Port, addr.IP.Equal(laddr.IP))
		if addr.Port == laddr.Port || !addr.IP.Equal(laddr.IP) {
			failedT++
			// t.Errorf("Dial %v should fail", addr)
		}
		c.Close()
	}
	tcpLisn.Stop()
	time.Sleep(time.Millisecond * 10)
	log.Printf("== failed : %d, rate: %f\n", failedT, float64(failedT)/float64(n))
	log.Printf("%s lock status: %v\n\n", nameT, tcpLisn.locked)
}
