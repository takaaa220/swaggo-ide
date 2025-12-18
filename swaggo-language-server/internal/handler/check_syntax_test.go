package handler

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/logger"
)

func TestMain(m *testing.M) {
	logger.Setup(os.Stderr, logger.LogDebug)
	os.Exit(m.Run())
}

func TestCheckSyntax_ShutdownDuringTimer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	h := &LSPHandler{
		checkSyntaxDebounce: 100 * time.Millisecond,
		checkSyntaxReq:      make(chan checkSyntaxRequest),
		cancel:              cancel,
	}

	done := make(chan struct{})
	go func() {
		h.checkSyntax(ctx)
		close(done)
	}()

	// Set timer
	h.requestCheckSyntax("file:///test.go", "package main")

	// Cancel before timer fires
	cancel()

	// Wait for checkSyntax goroutine to finish (confirms no panic)
	select {
	case <-done:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("checkSyntax did not shutdown in time")
	}
}

func TestCheckSyntax_TimerFiresAfterShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	h := &LSPHandler{
		checkSyntaxDebounce: 10 * time.Millisecond,
		checkSyntaxReq:      make(chan checkSyntaxRequest),
		cancel:              cancel,
	}

	done := make(chan struct{})
	go func() {
		h.checkSyntax(ctx)
		close(done)
	}()

	// Set timer with very short debounce
	h.requestCheckSyntax("file:///test.go", "package main")

	// Wait for timer to fire and try to send
	time.Sleep(15 * time.Millisecond)

	// Cancel after timer has fired but possibly before send completes
	cancel()

	// Wait for checkSyntax goroutine to finish (confirms no panic)
	select {
	case <-done:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("checkSyntax did not shutdown in time")
	}
}

func TestCheckSyntax_MultipleRequestsThenShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	h := &LSPHandler{
		checkSyntaxDebounce: 5 * time.Millisecond,
		checkSyntaxReq:      make(chan checkSyntaxRequest),
		cancel:              cancel,
	}

	done := make(chan struct{})
	go func() {
		h.checkSyntax(ctx)
		close(done)
	}()

	// Send multiple requests rapidly
	for i := 0; i < 10; i++ {
		h.requestCheckSyntax("file:///test.go", "package main")
		time.Sleep(1 * time.Millisecond)
	}

	// Cancel while timers might be firing
	cancel()

	// Wait for checkSyntax goroutine to finish (confirms no panic)
	select {
	case <-done:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("checkSyntax did not shutdown in time")
	}
}
