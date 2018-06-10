package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

type FileReaderFn func(path string) (string, error)

func DefaultFileReader(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

type DnsmasqProcess struct {
	ctx          *MainContext
	proc         *exec.Cmd
	fileReaderFn FileReaderFn
}

func NewDnsmasqProcess(reader FileReaderFn, c *MainContext) *DnsmasqProcess {
	return &DnsmasqProcess{
		ctx:          c,
		fileReaderFn: reader,
	}
}

func (p *DnsmasqProcess) Start() error {
	if !p.isRunning() {
		//
	}
	return nil
}

func (p *DnsmasqProcess) Stop() error {
	if p.isRunning() {
		err := p.proc.Process.Kill()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *DnsmasqProcess) Restart() error {
	err := p.Stop()
	if err != nil {
		return err
	}
	err = p.Start()
	if err != nil {
		return err
	}
	return nil
}

func (p *DnsmasqProcess) isRunning() bool {
	return p.proc != nil && !p.proc.ProcessState.Exited()
}

func collectInternalArgs(fileReaderFn FileReaderFn) []string {
	args := []string{
		"--keep-in-foreground",
		"--conf-dir=/etc/dnsmasq.d,.dpkg-dist,.dpkg-old,.dpkg-new",
		"--local-service",
	}
	trustAnchor, err := fileReaderFn("/usr/share/dns/root.ds")
	if err == nil {
		lines := strings.Split(trustAnchor, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			a := strings.Replace(line, ". IN DS ", "--trust-anchor=.,", 1)
			a = strings.Replace(a, " ", ",", -1)
			args = append(args, a)
		}
	}
	return args
}

func execCmd(cmd *exec.Cmd) (*exec.Cmd, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go pipeReader(bufio.NewReader(stdout))
	go pipeReader(bufio.NewReader(stderr))
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
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
