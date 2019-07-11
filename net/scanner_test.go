package net

import (
	"errors"
	"net"
	stdnet "net"
	"testing"
)

var _ stdnet.Listener = (*mockListener)(nil)

type mockListener struct {
	onAccept func() (stdnet.Conn, error)
	onAddr   func() stdnet.Addr
	onClose  func() error
}

func (ln *mockListener) Accept() (stdnet.Conn, error) {
	return ln.onAccept()
}

func (ln *mockListener) Addr() stdnet.Addr {
	return ln.onAddr()
}

func (ln *mockListener) Close() error {
	return ln.onClose()
}

func TestScanner(t *testing.T) {
	t.Run("nil listener", func(t *testing.T) {
		var nilLn stdnet.Listener
		scanner := NewScanner(nilLn)
		if want, got := ErrNilListener, scanner.Err(); want != got {
			t.Fatalf("scanner.Err -> want: %v, but got: %v", want, got)
		}
		if want, got := false, scanner.Scan(); want != got {
			t.Fatalf("scanner.Scan -> want: %v, but got: %v", want, got)
		}
	})

	ln := &mockListener{}
	scanner := NewScanner(ln)
	if want, got := (error)(nil), scanner.Err(); want != got {
		t.Fatalf("scanner.Err -> want: %v, but got: %v", want, got)
	}

	if want, got := ln, scanner.Listener(); want != got {
		t.Fatalf("scanner.Listener -> want: %v, but got: %v", want, got)
	}

	mockConn := &stdnet.TCPConn{}

	ln.onAccept = func() (stdnet.Conn, error) {
		return mockConn, nil
	}

	if want, got := true, scanner.Scan(); want != got {
		t.Fatalf("scanner.Scan -> want: %v, but got: %v", want, got)
	}

	if want, got := mockConn, scanner.Conn(); want != got {
		t.Fatalf("scanner.Conn -> want: %v, but got: %v", want, got)
	}

	if want, got := (error)(nil), scanner.Err(); want != got {
		t.Fatalf("scanner.Err -> want: %v, but got: %v", want, got)
	}

	accErr := errors.New("testing")
	ln.onAccept = func() (stdnet.Conn, error) {
		return nil, accErr
	}
	if want, got := false, scanner.Scan(); want != got {
		t.Fatalf("scanner.Scan -> want: %v, but got: %v", want, got)
	}
	if want, got := accErr, scanner.Err(); want != got {
		t.Fatalf("scanner.Err -> want: %v, but got: %v", want, got)
	}
}

func TestScannerCloser(t *testing.T) {
	scanner := ScanNet("tcp", "127.0.0.1:")
	defer scanner.Close()

	if want, got := (error)(nil), scanner.Err(); want != got {
		t.Fatalf("scanner.Err -> want: %v, but got: %v", want, got)
	}

	ln := scanner.Listener()
	addr := ln.Addr()
	type asyncAddr struct {
		err  error
		addr net.Addr
	}
	chAddr := make(chan *asyncAddr, 1)
	go func() {
		aa := &asyncAddr{}
		conn, err := stdnet.Dial(addr.Network(), addr.String())
		if err != nil {
			aa.err = err
		}
		defer conn.Close()
		chAddr <- aa
		aa.addr = conn.LocalAddr()
		chAddr <- aa
	}()

	if want, got := true, scanner.Scan(); want != got {
		t.Fatalf("scanner.Scan -> want: %v, but got: %v", want, got)
	}
	aa := <-chAddr
	if aa.err != nil {
		t.Fatalf("dial failed: %v", aa.err)
	}

	conn := scanner.Conn()
	if conn == nil {
		t.Fatalf("accept failed")
	}
	defer conn.Close()

	cAddr := aa.addr
	sAddr := conn.RemoteAddr()
	if cAddr.String() != sAddr.String() {
		t.Fatalf("addr not matched")
	}
}
