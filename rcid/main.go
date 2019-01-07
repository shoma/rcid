package rcid

import (
	"context"
	"github.com/shoma/rcid/grpc"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Serve() {
	ctx := context.Background()

	server, err := grpc.New()
	if err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}

	listener, err := net.Listen("tcp", ":5555")
	log.Println("Start to listen gRPC port :5555")

	ctx, cancel := context.WithCancel(ctx)
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error { return server.Serve(listener) })

	// Waiting for SIGTERM or Interrupt signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	select {
	case <-sigCh:
		log.Println("received SIGTERM, exiting server gracefully")
	case <-ctx.Done():
	}

	server.GracefulStop()
	cancel()
}

func Main() {
	Serve()
}
