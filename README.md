Onion Proxy
-----------

This is a small go module acting as a relay between for TCP connections that
should be forwarded to a listening Onion service on the Tor network.

It optionally speaks the [PROXY](http://www.haproxy.org/download/1.5/doc/proxy-protocol.txt)
protocol to forward client addresses to the backing service.

Configuration
=====

An instance of the proxy is parameterized by the following options:

* ```-l localhost:9999``` Where the proxy listens.
* ```-t=false``` Disable use of tor for backend socket resolution.
* ```-r sixteencharacter.onion:80``` Where the proxy forwards.
* ```-s /var/run/tor/control``` Where the tor control channel is.
* ```-c passwordauth``` The tor auth password, if using password auth.
* ```-p``` include a PROXY header on forwarded streams.

Usage
=====

Install `onionproxy` via the go command line:

```go
go install github.com/willscott/onionproxy
```

Then add it as a daemon in your startup script, run it locally, or otherwise
invoke it.

```bash
onionproxy -l 0.0.0.0:80 -r sixteencharacter.onion:80 &
onionproxy -l 0.0.0.0:25 -r sixteencharacter.onion:25 -p &
```
