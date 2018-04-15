# tcplistener
tcplistener to control the tcp listener。

```go
type TCPListener struct {
	tcpListener *net.TCPListener
	stop        chan bool
	locked      bool
	mutex       *sync.Mutex
}
```

use stop chan to control the stop of listener.

use mutex to Lock listenter when received stop signal.



### Test

```shell
go test
2018/04/15 23:39:30 Accept() lock status: false
2018/04/15 23:39:32 -- Accept() err : tcp listener stopped by signal
2018/04/15 23:39:32 -- failed times: 3; fail rate: 0.000375
2018/04/15 23:39:32 Accept() lock status: true

2018/04/15 23:39:32 AcceptTCP() lock status: false
2018/04/15 23:39:42 -- AcceptTCP() err : tcp listener stopped by signal
2018/04/15 23:39:42 -- failed times: 0; fail rate: 0.000000
2018/04/15 23:39:42 AcceptTCP() lock status: true

PASS
ok  	github.com/jungle85gopy/tcplistener	12.499s

```

