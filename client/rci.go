package main

import (
	"context"
	"github.com/shoma/rcid/pb"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	addr    = pflag.String("addr", "127.0.0.1:5555", "The address to connect rcid via gRPC")
	path    = pflag.String("path", "", "The path to execute command")
	args    = pflag.StringArray("args", []string{}, "Set args of execute command")
	env     = pflag.StringArray("env", []string{}, "Set environment variables. eg. VAR1=value1")
	workdir = pflag.String("workdir", "", "Set working directory to the command")
	timeout = pflag.Int("timeout", 60, "Set client side timeout sec to connect the host")
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer func() {
		err := conn.Close()
		log.Fatalf("failed to Close the connection: %+v", err)
	}()

	client := pb.NewRunnerClient(conn)

	pflag.Parse()
	run(client)

}

func parseEnvArg() map[string]string {
	parsed := make(map[string]string)
	if len(*env) == 0 {
		return parsed
	}
	for _, e := range *env {
		arr := strings.Split(e, "=")
		if arr[0] == "" {
			continue
		}
		parsed[arr[0]] = strings.Join(arr[1:], "")
	}
	return parsed
}

func run(client pb.RunnerClient) {
	to := time.Duration(*timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), to)
	hostname, _ := os.Hostname()
	attr := make(map[string]string)
	attr["Host"] = hostname
	ctx.Value(attr)
	defer cancel()

	req := pb.CommandRequest{
		Path:    *path,
		Arg:     *args,
		Workdir: *workdir,
		Env:     parseEnvArg(),
		Timeout: int32(*timeout + 20),
	}

	stream, err := client.Run(ctx, &req)
	if err != nil {
		log.Fatalf("revice error msg: %+v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Run(_) = _, %+v", client, err)
		}
		if res.Stderr != "" {
			log.Println(res.Stderr)
		}

		if res.Stdout != "" {
			log.Println(res.Stdout)
		}
	}
}
