package runner

import (
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"
)

type Command struct {
	path     string
	arg      []string
	workdir  string
	env      map[string]string
	ctx      context.Context
	timeout  time.Duration
	stdoutFn func(v ...interface{})
	stderrFn func(v ...interface{})
	errExit  bool
	lastLog  string
	exited   bool
}

func (c *Command) SetCtx(ctx context.Context) {
	c.ctx = ctx
}

func (c *Command) SetStderrFn(stderrFn func(v ...interface{})) {
	c.stderrFn = stderrFn
}

func (c *Command) SetStdoutFn(stdoutFn func(v ...interface{})) {
	c.stdoutFn = stdoutFn
}

func NewCommand(path string, arg []string, workdir string, env map[string]string, timeout time.Duration) *Command {
	return &Command{path: path, arg: arg, workdir: workdir, env: env, timeout: timeout}
}

func (c *Command) envSlice() []string {
	var rtn []string
	for k, v := range c.env {
		rtn = append(rtn, fmt.Sprintf("%s=%s", k, v))
	}

	return rtn
}

func (c *Command) Run() error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.path, c.arg...)

	cmd.Dir = c.workdir
	cmd.Env = c.envSlice()

	stdout, err := cmd.StdoutPipe()
	defer func() {
		err := stdout.Close()
		log.Printf("failed to close the stdout: %+v", err)
	}()

	if err != nil {
		return errors.Wrap(err, "failed to pipe stdout")

	}

	stderr, err := cmd.StderrPipe()
	defer func(){
		err := stderr.Close()
		log.Printf("failed to close the stderr: %+v", err)
	}()

	if err != nil {
		return errors.Wrap(err, "failed to pipe stderr")
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	wg.Add(2)
	go func() { defer wg.Done(); readOutput(stdout, c.stdoutFn) }()
	go func() { defer wg.Done(); readOutput(stderr, c.stderrFn) }()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return errors.Wrap(err, "failed to wait the cmd")
	}

	c.errExit = true
	c.exited = false
	if cmd.ProcessState.Exited() {
		c.lastLog = cmd.ProcessState.String()
		c.exited = true
		if cmd.ProcessState.Success() {
			c.errExit = false
		}
	}

	return nil
}

func (c *Command) GetLastLog() string {
	return c.lastLog
}

func (c Command) IsErrorExit() bool {
	return c.errExit
}

func (c Command) Exited() bool {
	return c.exited
}

func readOutput(r io.Reader, callback func(v ...interface{})) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		callback(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
