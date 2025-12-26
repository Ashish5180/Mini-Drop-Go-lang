package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"mini-dropbox/internals/master"
	"mini-dropbox/internals/node"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Start storage nodes
	go node.StartNode(ctx, "8001")
	go node.StartNode(ctx, "8002")

	// Start master node
	go master.StartMaster(ctx, "9000")

	fmt.Println("Mini-Dropbox started: master on 9000, nodes on 8001, 8002")
	fmt.Println("Press Ctrl+C to stop")

	// Block until signal
	<-ctx.Done()
	fmt.Println("Shutting down gracefully...")
	// Allow background goroutines to finish cleanup
	<-make(chan struct{})
}
