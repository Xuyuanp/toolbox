package context

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWithSignals(t *testing.T) {
	testSig := syscall.SIGUSR1
	ctx, cancel := WithSignals(context.Background(), testSig)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatalf("context shouldn't be done right now")
	default:
	}

	if err := syscall.Kill(os.Getpid(), testSig); err != nil {
		t.Fatalf("kill current process failed: %v", err)
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatalf("timeout")
	}

	if want, got := context.Canceled, ctx.Err(); want != got {
		t.Fatalf("want error: %v, but got: %v", want, got)
	}
}

func TestCancel(t *testing.T) {
	testSig := syscall.SIGUSR1
	ctx, cancel := WithSignals(context.Background(), testSig)

	select {
	case <-ctx.Done():
		t.Fatalf("context shouldn't be done right now")
	default:
	}
	go func() {
		cancel()
	}()

	<-ctx.Done()
	if want, got := context.Canceled, ctx.Err(); want != got {
		t.Fatalf("want error: %v, but got: %v", want, got)
	}
}
