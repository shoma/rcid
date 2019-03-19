package rcid

import (
	"context"
	"errors"
	"github.com/shoma/rcid/grpc"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Serve() error {
	ctx := context.Background()

	server, err := grpc.New()
	if err != nil {
		return errors.New(err.Error())
	}

	listener, err := net.Listen("tcp", ":5555")
	if err != nil {
		return errors.New(err.Error())
	}
	log.Println("Start to listen gRPC port :5555")

	ctx, cancel := context.WithCancel(ctx)
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error { return server.Serve(listener) })

	defer cancel()

	// Waiting for SIGTERM or Interrupt signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		log.Println("received SIGTERM, exiting server gracefully")
	case <-ctx.Done():
	}

	server.GracefulStop()

	return nil
}

func Main()int {
	err := Serve()
	if err != nil {
		log.Printf("%s\n", err.Error())
		return 1
	}
	return 0
}
