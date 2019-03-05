package grace_test

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/chapsuk/grace"
	"github.com/chapsuk/keymon"
)

const (
	success = 1
	failed  = 2
)

func TestShutdownContext(t *testing.T) {
	osExit := make(chan int)
	keymon.Patch(os.Exit, func(code int) { osExit <- code })
	defer keymon.Unpatch(os.Exit)

	assertSignals(t, grace.ShutdownContext(context.Background()), osExit, syscall.SIGTERM)
	assertSignals(t, grace.ShutdownContext(context.Background()), osExit, syscall.SIGHUP)
	assertSignals(t, grace.ShutdownContext(context.Background()), osExit, syscall.SIGINT)
	assertParent(t, grace.ShutdownContext)
}

func TestStopContext(t *testing.T) {
	osExit := make(chan int)
	keymon.Patch(os.Exit, func(code int) { osExit <- code })
	defer keymon.Unpatch(os.Exit)

	assertSignals(t, grace.StopContext(context.Background()), osExit, syscall.SIGTERM)
	assertSignals(t, grace.StopContext(context.Background()), osExit, syscall.SIGINT)
	assertParent(t, grace.StopContext)
}

func TestReloadChannel(t *testing.T) {
	reloadChan := grace.ReloadChannel(context.Background())

	for i := 0; i < 5; i++ {
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		select {
		case <-reloadChan:
			t.Logf("#%d reload signal as expected", i)
		case <-time.After(5 * time.Millisecond):
			t.Fatalf("#%d reload signal timeout", i)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	rchan := grace.ReloadChannel(ctx)

	cancel()
	select {
	case <-rchan:
		t.Fatal("Received the reload signal when parent context canceled")
	case <-time.After(5 * time.Millisecond):
		t.Log("As expected the reload signal not published when context canceled")
	}
}

func assertParent(t *testing.T, createFunc func(context.Context) context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	cctx := grace.ShutdownContext(ctx)

	cancel()
	select {
	case <-cctx.Done():
		t.Log("Returns child context as expected")
	case <-time.After(5 * time.Millisecond):
		t.Fatal("Returns not child context")
	}
}

func assertSignals(t *testing.T, ctx context.Context, osExit chan int, sig syscall.Signal) {
	testResult := make(chan int)
	go func() {
		select {
		case <-ctx.Done():
			testResult <- success
		case <-time.After(time.Second):
			testResult <- failed
		}
	}()

	syscall.Kill(syscall.Getpid(), sig)
	switch <-testResult {
	case failed:
		t.Fatalf("Wait context done timeout, signal '%s'", sig)
	case success:
		t.Logf("Context canceled as expected for '%s' signal", sig)
	default:
		t.Fatal("Unexpected test result")
	}

	if osExit != nil {
		syscall.Kill(syscall.Getpid(), sig)
		select {
		case exitCode := <-osExit:
			if exitCode != 1 {
				t.Fatalf("Incorrect exit code for second '%s' signal,  expected: 1 actual: %d",
					sig, exitCode)
			}
			t.Logf("Received second '%s' signal, exit code 1 as expected", sig)
		case <-time.After(time.Second):
			t.Fatalf("Waiting os.Exit timeout, signal '%s'", sig)
		}
	}
}
