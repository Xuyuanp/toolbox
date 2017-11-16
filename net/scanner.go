package net

import (
	"crypto/tls"
	"errors"
	"net"
)

var _ ListenerScanner = &listenerScanner{}

// errors returned by scanner
var (
	ErrNilListener = errors.New("imhttp: nil listener")
)

// ListenerScanner provides a convenient interface for accepting connection from a net listener.
type ListenerScanner interface {
	Scan() bool
	Conn() net.Conn
	Err() error
	Close() error
}

type listenerScanner struct {
	ln   net.Listener
	err  error
	conn net.Conn
}

// NewScanner creates a ListenerScanner from the provided listener.
// Don't call Close method if you want to close the listener by yourself.
func NewScanner(ln net.Listener) ListenerScanner {
	return &listenerScanner{ln: ln, err: checkNilListener(ln)}
}

func checkNilListener(ln net.Listener) error {
	if ln == nil {
		return ErrNilListener
	}
	return nil
}

// Scan advances the scanner to the next connection, which will then be available through the Conn method.
// It returns false when the listener is nil or the listener.Accept method returns a non-nil error.
func (s *listenerScanner) Scan() (ok bool) {
	if s.err != nil || s.ln == nil {
		return false
	}
	s.conn, s.err = s.ln.Accept()
	return s.err == nil
}

// Conn returns the most recent connection accepted by listener after a successful call to Scan.
func (s *listenerScanner) Conn() net.Conn {
	return s.conn
}

// Err returns the first non-nill error that encountered by scanner.
func (s *listenerScanner) Err() error {
	return s.err
}

// Close closes the internal listener.
func (s *listenerScanner) Close() error {
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}

// ScanNet returns a ListenerScanner accepts normal connection.
func ScanNet(network, addr string) ListenerScanner {
	scanner := &listenerScanner{}
	scanner.ln, scanner.err = net.Listen(network, addr)
	return scanner
}

// ScanTLS returns a ListenerScanner accepts tls connection.
func ScanTLS(network, addr string, config *tls.Config) ListenerScanner {
	scanner := &listenerScanner{}
	scanner.ln, scanner.err = tls.Listen(network, addr, config)
	return scanner
}
