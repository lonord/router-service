package ba

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type CmdExecutorFn func(cmd string) (string, error)

func DefaultCmdExecutor(cmd string) (string, error) {
	command := exec.Command("/bin/bash", "-c", cmd)
	content, err := command.Output()
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ExecPipeCmd(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go pipeReader(bufio.NewReader(stdout))
	go pipeReader(bufio.NewReader(stderr))
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func pipeReader(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Printf(line)
	}
}
