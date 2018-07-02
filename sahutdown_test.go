package grace_test

import (
	"testing"
	"time"

	"github.com/chapsuk/grace"
	"go.uber.org/dig"
)

func TestDigUsage(t *testing.T) {

	c := dig.New()

	if err := c.Provide(grace.NewShutdownContext); err != nil {
		t.Fatalf("provide grace context error: %v", err)
	}

	err := c.Invoke(func(cparam grace.ContextParams) {
		select {
		case <-cparam.Context.Done():
			t.Fatalf("context canceled without signal")
		case <-time.Tick(100 * time.Millisecond):
		}
	})

	if err != nil {
		t.Fatalf("invoke grace context error: %v", err)
	}
}
