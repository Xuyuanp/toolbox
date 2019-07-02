package http

import (
	"context"
	stdhttp "net/http"
	"time"
)

// Graceful server
func Graceful(ctx context.Context, srv *stdhttp.Server, wait time.Duration) error {
	chErr := make(chan error, 1)
	go func() {
		<-ctx.Done()
		err := shutdown(srv, wait)
		select {
		case chErr <- err:
		default:
		}
	}()

	err := srv.ListenAndServe()
	if err == stdhttp.ErrServerClosed {
		return <-chErr
	}
	return err
}

func shutdown(srv *stdhttp.Server, wait time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	return srv.Shutdown(ctx)
}
