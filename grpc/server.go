package grpc

import (
	"github.com/shoma/rcid/grpc/runner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/reflection"
	"math"
	"os"
)

const (
	maxStreams   = math.MaxUint32
	maxSendBytes = math.MaxInt32
)

func New() (*grpc.Server, error) {
	var opts []grpc.ServerOption
	opts = append(opts, grpc.MaxSendMsgSize(maxSendBytes))
	opts = append(opts, grpc.MaxConcurrentStreams(maxStreams))
	opts = append(opts, getStreamInterceptor())

	server := grpc.NewServer(opts...)

	runner.Register(server)
	service.RegisterChannelzServiceToServer(server)

	if os.Getenv("USE_GRPC_REFLECTION") == "1" {
		reflection.Register(server)
	}

	return server, nil
}
