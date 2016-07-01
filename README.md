Onion Proxy
-----------

This is a small go module acting as a relay between for TCP connections that
should be forwarded to a listening Onion service on the Tor network.

It optionally speaks the (PROXY)[http://www.haproxy.org/download/1.5/doc/proxy-protocol.txt]
protocol to forward client addresses to the backing service.

Usage
=====

An instance of the proxy is parameterized by the following options:

* ```Port``` and ```Host``` specify the interface and port on the local machine.
* ```OnionService``` and ```OnionPort``` specify the backend.
* ```SocksPort``` specifies where the local tor socks port is.
* ```ActiveConnections``` and ```ConnectionTimeout``` specify how traffic is relayed.
* ```ProxyHeader``` specifies whether a PROXY header is prefixed to streams.
