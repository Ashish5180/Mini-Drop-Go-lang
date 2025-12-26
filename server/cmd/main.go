package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"mini-dropbox/internals/master"
	"mini-dropbox/internals/node"
)

const shutdownTimeout = 5 * time.Second

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup

	// Start storage nodes
	wg.Add(2)
	go func() {
		defer wg.Done()
		node.StartNode(ctx, "8001")
	}()
	go func() {
		defer wg.Done()
		node.StartNode(ctx, "8002")
	}()

	// Start master node
	wg.Add(1)
	go func() {
		defer wg.Done()
		master.StartMaster(ctx, "9000")
	}()

	fmt.Println("Mini-Dropbox started: master on 9000, nodes on 8001, 8002")
	fmt.Println("Press Ctrl+C to stop")

	// Block until signal
	<-ctx.Done()
	fmt.Println("Shutting down gracefully...")

	// Wait for cleanup with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Shutdown complete")
	case <-time.After(shutdownTimeout):
		fmt.Println("Shutdown timeout reached")
	}
}
