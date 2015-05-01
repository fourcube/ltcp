// Package ltcp provides a simple function for launching tcp servers e.g. for
// testing.
package ltcp

import (
	"io"
	"log"
	"net"
)

// Shutdown is sent to a listener to signal that it should shutdown
var Shutdown = struct{}{}

// ConnectionHandler is a alias for the func(net.Conn) type. ConnectionHandlers can be passed
// to the Listen* functions. Every client connection is processed by a ConnectionHandler.
type ConnectionHandler func(net.Conn)

// Listen on 'addr' serving all client connections with 'handler'
//
// Stop the server by sending ltcp.Shutdown over the 'done' channel.
func Listen(addr string, handler ConnectionHandler, done chan struct{}) (err error) {
	ln, err := net.Listen("tcp", addr)

	if err != nil {
		log.Printf("Error listening on %s, %v ", addr, err)
		close(done)
		return
	}

	go listen(ln, handler, done)
	return
}

// ListenAny picks a random port and listens on all interfaces serving connections with 'handler'
//
// Stop the server by sending ltcp.Shutdown over the 'done' channel.
func ListenAny(handler ConnectionHandler, done chan struct{}) (addr *net.TCPAddr, err error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Printf("Couldn't listen on any ip address, %v", err)
		close(done)
		return
	}

	addr, err = net.ResolveTCPAddr(ln.Addr().Network(), ln.Addr().String())
	if err != nil {
		log.Printf("Couldn't resolve %v, %v", ln.Addr(), err)
		close(done)
		return
	}

	go listen(ln, handler, done)
	return
}

func listen(ln net.Listener, handler ConnectionHandler, done chan struct{}) {
	// This implementation does not block, but instead runs it's
	// accept -> handle loop inside a goroutine
	go func() {
		_, ok := <-done
		ln.Close()

		if ok {
			close(done)
		}
	}()

	go func() {
		for {
			conn, err := ln.Accept()

			// We will see an error when the connection is closed but accept is still blocking
			// for a connection
			if err != nil {
				log.Printf("Error during accept, %v ", err)
				return
			}

			// Do all the work inside the supplied goroutine so we can quickly accept
			// other connections
			go handler(conn)
		}
	}()
}

// EchoHandler simply returns everything that is received to the client himself
func EchoHandler(client net.Conn) {
	var err error
	defer client.Close()

	for err == nil {
		// ...simply return everything the client sends to itself
		_, err = io.Copy(client, client)

		if err != nil && err != io.EOF {
			log.Printf("Error during echo %v", err)
		}
	}
}
