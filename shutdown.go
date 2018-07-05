package grace

import (
	"context"
	"os"
	"os/signal"
	"syscall"
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
