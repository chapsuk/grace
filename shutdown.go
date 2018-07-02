package grace

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/dig"
)

// ShutdownContext returns child context from passed context which will be canceled
// on incoming signals: SIGINT, SIGTERM, SIGHUP
func ShutdownContext(c context.Context) context.Context {
	ctx, cancel := context.WithCancel(c)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-ch
		cancel()
	}()
	return ctx
}

// ContextResult wrap grace context for dig usage
type ContextResult struct {
	dig.Out
	Context context.Context `name:"grace_context"`
}

// NewShutdownContext returns wrapped ShutdownContext func result for dig usage
func NewShutdownContext() ContextResult {
	return ContextResult{
		Context: ShutdownContext(context.Background()),
	}
}

// ContextParams contains shutdown context for dig usage
type ContextParams struct {
	dig.In
	Context context.Context `name:"grace_context"`
}
