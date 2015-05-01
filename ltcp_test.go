package ltcp_test

import (
	"log"
	"net"
	"testing"

	"github.com/fourcube/ltcp"
)

func TestListen(t *testing.T) {
	listenAddress := "127.0.0.1:54429"
	done := make(chan struct{})

	err := ltcp.Listen(listenAddress, ltcp.EchoHandler, done)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	done <- ltcp.Shutdown
	<-done
}

func TestShutdown(t *testing.T) {
	listenAddress := "127.0.0.1:54429"
	done := make(chan struct{})

	err := ltcp.Listen(listenAddress, ltcp.EchoHandler, done)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	done <- ltcp.Shutdown
	<-done

	_, err = net.Dial("tcp", listenAddress)
	if err == nil {
		t.Errorf("Expected error when connecting to shut down server, got nil")
	}
}

func TestAddressAlreadyInUse(t *testing.T) {
	listenAddress := "127.0.0.1:54429"
	doneA := make(chan struct{})
	doneB := make(chan struct{})

	err := ltcp.Listen(listenAddress, ltcp.EchoHandler, doneA)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = ltcp.Listen(listenAddress, ltcp.EchoHandler, doneB)
	if err == nil {
		t.Errorf("Expected listen to fail when address is in use, got no error")
	}

	doneA <- ltcp.Shutdown
	<-doneA
	// doneB will be closed because of the error
	<-doneB
}

func TestListenAny(t *testing.T) {
	done := make(chan struct{})

	addr, err := ltcp.ListenAny(ltcp.EchoHandler, done)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if addr == nil {
		t.Errorf("Expected server to listen on some address, got nil")
	}

	done <- ltcp.Shutdown
	<-done
}

func TestServerActuallyResponds(t *testing.T) {
	listenAddress := "127.0.0.1:54429"
	done := make(chan struct{})

	err := ltcp.Listen(listenAddress, ltcp.EchoHandler, done)

	conn, err := net.Dial("tcp", listenAddress)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	testPayload := "foo"
	recvBuf := make([]byte, 32)
	conn.Write([]byte(testPayload))

	n, err := conn.Read(recvBuf)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	data := string(recvBuf[:n])
	if data != testPayload {
		t.Errorf("Expected to receive '%s' from the echo handler, got '%s'", testPayload, data)
	}

	done <- ltcp.Shutdown
	log.Printf("After send done")
	<-done
}

func TestServerActuallyRespondsDuringListenAny(t *testing.T) {
	done := make(chan struct{})
	addr, err := ltcp.ListenAny(ltcp.EchoHandler, done)

	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	testPayload := "foo"
	recvBuf := make([]byte, 32)
	conn.Write([]byte(testPayload))

	n, err := conn.Read(recvBuf)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	data := string(recvBuf[:n])
	if data != "foo" {
		t.Errorf("Expected to receive '%s' from the echo handler, got '%s'", testPayload, data)
	}

	done <- ltcp.Shutdown
	<-done
}
