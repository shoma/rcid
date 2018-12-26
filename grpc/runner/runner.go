package runner

import (
	"context"
	"fmt"
	"github.com/shoma/rcid/pb"
	"github.com/shoma/rcid/runner"
	"google.golang.org/grpc"
	"log"
	"time"
)

type runnerServer struct {
	s *grpc.Server
}

func Register(s *grpc.Server) {
	pb.RegisterRunnerServer(s, &runnerServer{
		s: s,
	})
}

func (r *runnerServer) Run(req *pb.CommandRequest, stream pb.Runner_RunServer) error {
	log.Printf("Receive a request: %v\n", req.String())
	g := newGRPCRunner(req, stream)
	g.run()

	return nil
}

type gRPCRunner struct {
	req     *pb.CommandRequest
	stream  pb.Runner_RunServer
	timeOut time.Duration
	cmd     runner.Command
}

func (g *gRPCRunner) run() {
	g.cmd.SetStderrFn(g.sendStderr)
	g.cmd.SetStdoutFn(g.sendStdout)
	err := g.cmd.Run()
	if err != nil {
		log.Print(err)
		g.sendError(err)
	}
}

func (g *gRPCRunner) send(r pb.CommandResult) {
	if g.stream.Context().Err() == context.Canceled {
		//err = status.New(codes.Canceled, "Client cancelled, abandoning.").Err()
		log.Printf("Client cancelled, abandoning. %v\n", g.req.String())
		return
	}

	err := g.stream.Send(&r)
	if err != nil {
		// TODO returns error
		log.Printf("Client cancelled, abandoning. %v\n", err)
	}
}

func (g *gRPCRunner) sendError(err error) {
	r := pb.CommandResult{
		Error:  true,
		Stdout: "",
		Stderr: err.Error(),
	}
	g.send(r)
}

func (g *gRPCRunner) sendStdout(v ...interface{}) {
	log.Println(v)
	r := pb.CommandResult{
		Error:  false,
		Stdout: fmt.Sprint(v...),
		Stderr: "",
	}

	g.send(r)
}

func (g *gRPCRunner) sendStderr(v ...interface{}) {
	log.Println(v)
	r := pb.CommandResult{
		Error:  true,
		Stdout: "",
		Stderr: fmt.Sprint(v...),
	}

	g.send(r)
}

func newGRPCRunner(req *pb.CommandRequest, stream pb.Runner_RunServer) *gRPCRunner {
	g := &gRPCRunner{req: req, stream: stream}
	t := req.GetTimeout()
	timeout := time.Second * time.Duration(t)
	g.timeOut = timeout

	cmd := runner.NewCommand(
		req.GetPath(),
		req.GetArg(),
		req.GetWorkdir(),
		req.GetEnv(),
		timeout,
	)
	cmd.SetCtx(g.stream.Context())

	g.cmd = *cmd

	return g
}
