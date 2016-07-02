package main

import (
	"flag"
	"fmt"
	"github.com/yawning/bulb"
	"golang.org/x/net/proxy"
	"io"
	"net"
)

var localAddr *string = flag.String("l", "localhost:9999", "listening address")
var remoteAddr *string = flag.String("r", "google.com:80", "backend address")
var useTor *bool = flag.Bool("t", true, "Use Tor for backend")
var torSocket *string = flag.String("s", "/var/run/tor/control", "Tor control socket")
var cookieAuth *string = flag.String("c", "", "tor auth password")
var proxyheader *bool = flag.Bool("p", false, "Include a PROXY header")

func Proxy(srvConn, cliConn *net.TCPConn) {
	// channels to wait on the close event for each connection
	serverClosed := make(chan struct{}, 1)
	clientClosed := make(chan struct{}, 1)

	go broker(srvConn, cliConn, clientClosed)
	go broker(cliConn, srvConn, serverClosed)

	// wait for one half of the proxy to exit, then trigger a shutdown of the
	// other half by calling CloseRead(). This will break the read loop in the
	// broker and allow us to fully close the connection cleanly without a
	// "use of closed network connection" error.
	var waitFor chan struct{}
	select {
	case <-clientClosed:
		// the client closed first and any more packets from the server aren't
		// useful, so we can optionally SetLinger(0) here to recycle the port
		// faster.
		srvConn.SetLinger(0)
		srvConn.CloseRead()
		waitFor = serverClosed
	case <-serverClosed:
		cliConn.CloseRead()
		waitFor = clientClosed
	}

	// Wait for the other connection to close.
	// This "waitFor" pattern isn't required, but gives us a way to track the
	// connection and ensure all copies terminate correctly; we can trigger
	// stats on entry and deferred exit of this function.
	<-waitFor
}

// This does the actual data transfer.
// The broker only closes the Read side.
func broker(dst, src net.Conn, srcClosed chan struct{}) {
	// We can handle errors in a finer-grained manner by inlining io.Copy (it's
	// simple, and we drop the ReaderFrom or WriterTo checks for
	// net.Conn->net.Conn transfers, which aren't needed). This would also let
	// us adjust buffersize.
	_, err := io.Copy(dst, src)
	if err != nil {
		srcClosed <- struct{}{}
		return
	}
	src.Close()
	srcClosed <- struct{}{}
}

func proxyConn(dialer *proxy.Dialer, conn *net.TCPConn) {
	rConn, err := (*dialer).Dial("tcp", *remoteAddr)
	if err != nil {
		panic(err)
	}

	if *proxyheader {
		client_ip, client_port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			panic(err)
		}
		my_ip, my_port, _ := net.SplitHostPort(conn.LocalAddr().String())
		rConn.Write([]byte(fmt.Sprintf("PROXY TCP4 %s %s %d %d\r\n", client_ip, my_ip, client_port, my_port)))
	}

	back, okBack := (rConn).(*net.TCPConn)
	if !okBack {
		return
	}

	Proxy(conn, back)
}

func handleConn(in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	var dialer proxy.Dialer
	if *useTor {
		tor, err := bulb.Dial("unix", *torSocket)
		if err != nil {
			panic(err)
		}
		defer tor.Close()

		if err := tor.Authenticate(*cookieAuth); err != nil {
			panic(err)
		}

		dialer, err = tor.Dialer(nil)
		if err != nil {
			panic(err)
		}
	} else {
		dialer = proxy.Direct
	}

	for conn := range in {
		proxyConn(&dialer, conn)
		out <- conn
	}
}

func closeConn(in <-chan *net.TCPConn) {
	for conn := range in {
		conn.Close()
	}
}

func main() {
	flag.Parse()

	fmt.Printf("Listening: %v\nProxying: %v\n\n", *localAddr, *remoteAddr)

	addr, err := net.ResolveTCPAddr("tcp", *localAddr)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	pending, complete := make(chan *net.TCPConn), make(chan *net.TCPConn)

	for i := 0; i < 5; i++ {
		go handleConn(pending, complete)
	}
	go closeConn(complete)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		pending <- conn
	}
}
