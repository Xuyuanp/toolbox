package net

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
)

var _ ScannerCloser = &scanner{}

// errors returned by scanner
var (
	ErrNilListener = errors.New("nil listener")
)

// Scanner provides a convenient interface for accepting connection from a net listener.
type Scanner interface {
	Scan() bool
	Conn() net.Conn
	Listener() net.Listener
	Err() error
}

// ScannerCloser is a closable Scanner
type ScannerCloser interface {
	Scanner
	io.Closer
}

type scanner struct {
	ln   net.Listener
	err  error
	conn net.Conn
}

// NewScanner creates a Scanner from the provided listener.
// Don't call Close method if you want to close the listener by yourself.
func NewScanner(ln net.Listener) Scanner {
	return &scanner{ln: ln, err: checkNilListener(ln)}
}

func checkNilListener(ln net.Listener) error {
	if ln == nil {
		return ErrNilListener
	}
	return nil
}

// Scan advances the scanner to the next connection, which will then be available through the Conn method.
// It returns false when the listener is nil or the listener.Accept method returns a non-nil error.
func (s *scanner) Scan() (ok bool) {
	if s.err != nil || s.ln == nil {
		return false
	}
	s.conn, s.err = s.ln.Accept()
	return s.err == nil
}

// Listener returns the internal listener
func (s *scanner) Listener() net.Listener {
	return s.ln
}

// Conn returns the most recent connection accepted by listener after a successful call to Scan.
func (s *scanner) Conn() net.Conn {
	return s.conn
}

// Err returns the first non-nill error that encountered by scanner.
func (s *scanner) Err() error {
	return s.err
}

// Close closes the internal listener.
func (s *scanner) Close() error {
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}

// ScanNet returns a Scanner accepts normal connection.
func ScanNet(network, addr string) ScannerCloser {
	scanner := &scanner{}
	scanner.ln, scanner.err = net.Listen(network, addr)
	return scanner
}

// ScanTLS returns a Scanner accepts tls connection.
func ScanTLS(network, addr string, config *tls.Config) ScannerCloser {
	scanner := &scanner{}
	scanner.ln, scanner.err = tls.Listen(network, addr, config)
	return scanner
}
