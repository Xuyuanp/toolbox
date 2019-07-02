package http

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	stdhttp "net/http"
	"testing"
	"time"
)

func getFreePort(network string) (string, error) {
	ln, err := net.Listen(network, "127.0.0.1:")
	if err != nil {
		return "", err
	}
	defer ln.Close()
	return ln.Addr().String(), nil
}

func TestGraceful(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	srv := &stdhttp.Server{}
	go func() {
		cancel()
	}()
	err := Graceful(ctx, srv, time.Second)
	if want, got := (error)(nil), err; got != want {
		t.Fatalf("want error: %v, but got: %v", want, got)
	}
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	addr, err := getFreePort("tcp")
	if err != nil {
		t.Fatalf("get free port failed: %v", err)
	}
	var handler stdhttp.HandlerFunc = func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		time.Sleep(3 * time.Second)
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
		fmt.Fprintln(w, "hello")
	}
	srv := &stdhttp.Server{Addr: addr, Handler: handler}

	go func() {
		go func() {
			time.Sleep(time.Millisecond * 200)
			cancel()
		}()
		if _, err := stdhttp.Get("http://" + addr); err != nil {
			t.Fatalf("preform request failed: %v", err)
		}
	}()

	err = Graceful(ctx, srv, time.Millisecond*100)
	if want, got := context.DeadlineExceeded, err; got != want {
		t.Fatalf("want error: %v, but got: %v", want, got)
	}
}

func TestError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	addr, err := getFreePort("tcp")
	if err != nil {
		t.Fatalf("get free port failed: %v", err)
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("listen addr failed: %v", err)
	}
	defer ln.Close()

	srv := &stdhttp.Server{Addr: addr}
	go func() {
		cancel()
	}()
	err = Graceful(ctx, srv, time.Second)
	if err == nil {
		t.Fatalf("want error, but got nil")
	}
}
