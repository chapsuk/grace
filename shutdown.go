package grace

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// ShutdownContext returns child context from passed context which will be canceled
// on incoming signals: SIGINT, SIGTERM, SIGHUP.
// Ends immediately by os.Exit(1) after second signal
func ShutdownContext(c context.Context) context.Context {
	ctx, cancel := context.WithCancel(c)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		select {
		case <-ctx.Done():
			return
		case <-ch:
			cancel()
			<-ch
			os.Exit(1)
		}
	}()
	return ctx
}
