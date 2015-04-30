# ltcp

[![GoDoc](https://godoc.org/github.com/fourcube/ltcp?status.svg)](http://godoc.org/github.com/fourcube/ltcp) [![Build Status](https://travis-ci.org/fourcube/ltcp.svg?branch=master)](https://travis-ci.org/fourcube/ltcp)

ltcp is a simple library providing an easy way to spawn tcp handlers. It might be useful for your tests.

Handlers are functions of type `func (net.Conn)`.

Sample code:

```go
package main

import (
	"github.com/fourcube/ltcp"
 	"net"	
)

func main() {
	// Listens on some random, available port on all interfaces
	// handles all connections with the ltcp.EchoHandler
	//
	// You can close the 'done' channel to stop the listener.
	//
	done := make(chan struct{})
	addr, err :=	ltcp.ListenAny(ltcp.EchoHandler, done)

	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		// ...
	}

	// ... do whatever with the connection

	conn.Close()
	close(done)
}
```


