package context

import (
	"context"
	"os"
	"os/signal"
)

func WithSignals(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
    return signal.NotifyContext(parent, signals...)
}
