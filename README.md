Onion Proxy
-----------

This is a small go module acting as a relay between for TCP connections that
should be forwarded to a listening Onion service on the Tor network.

It optionally speaks the (PROXY)[http://www.haproxy.org/download/1.5/doc/proxy-protocol.txt]
protocol to forward client addresses to the backing service.

Usage
=====

An instance of the proxy is parameterized by the following options:

* ```-l localhost:9999``` Where the proxy listens.
* ```-r thirteenchars.onion:80``` Where the proxy forwards.
* ```-s /var/run/tor/control``` Where the tor control channel is.
* ```-c passwordauth``` The tor auth password, if using password auth.
