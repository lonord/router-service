package main

import "os/exec"

type CmdExecutorFn func(cmd string) (string, error)

func DefaultCmdExecutor(cmd string) (string, error) {
	command := exec.Command("/bin/bash", "-c", cmd)
	content, err := command.Output()
	if err != nil {
		return "", err
	}
	return string(content), nil
}
