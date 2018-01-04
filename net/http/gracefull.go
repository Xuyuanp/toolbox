package http

import (
	"context"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Gracefull server
func Gracefull(srv *stdhttp.Server, wait time.Duration) error {
	chErr := make(chan error, 1)
	chSigs := make(chan os.Signal)
	signal.Notify(chSigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(chSigs)

	go func() {
		chErr <- srv.ListenAndServe()
	}()

	ctx := context.Background()
	if wait > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, wait)
		defer cancel()
	}
	select {
	case <-chSigs:
		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
	case err := <-chErr:
		return err
	}
	err := <-chErr
	if err != stdhttp.ErrServerClosed {
		return err
	}
	return nil
}
