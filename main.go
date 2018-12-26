package main

import (
	"context"
	"github.com/shoma/rcid/rcid"
	"github.com/shoma/rcid/runner"
	"log"
	"time"
)

func main() {
	rcid.Main()
	//testRunner()
}

func testRunner() {
	args := []string{"images"}
	cmd := runner.NewCommand(
		"/usr/local/bin/docker",
		args,
		"",
		nil,
		60*time.Second,
	)
	cmd.SetStdoutFn(log.Print)
	cmd.SetStderrFn(log.Print)
	cmd.SetCtx(context.Background())
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

}
