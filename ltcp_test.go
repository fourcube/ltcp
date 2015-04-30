package ltcp_test

import (
	"net"
	"testing"

	"github.com/fourcube/ltcp"
)

func TestListen(t *testing.T) {
	listenAddress := "127.0.0.1:12345"
	done := make(chan struct{})

	err := ltcp.Listen(listenAddress, ltcp.EchoHandler, done)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	close(done)
}

func TestAddressAlreadyInUse(t *testing.T) {
	listenAddress := "127.0.0.1:12345"
	done := make(chan struct{})

	err := ltcp.Listen(listenAddress, ltcp.EchoHandler, done)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = ltcp.Listen(listenAddress, ltcp.EchoHandler, done)
	if err == nil {
		t.Errorf("Expected listen to fail when address is in use, got no error")
	}

	close(done)
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

	close(done)
}

func TestServerActuallyResponds(t *testing.T) {
	listenAddress := "127.0.0.1:12345"
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
	if data != "foo" {
		t.Errorf("Expected to receive '%s' from the echo handler, got '%s'", testPayload, data)
	}

	close(done)
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

	close(done)
}
