package dnsmasq

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"../cmdutil"
	"../context"
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
	ctx          *context.MainContext
	proc         *exec.Cmd
	fileReaderFn FileReaderFn
}

func NewDnsmasqProcess(reader FileReaderFn, c *context.MainContext) *DnsmasqProcess {
	return &DnsmasqProcess{
		ctx:          c,
		fileReaderFn: reader,
	}
}

func (p *DnsmasqProcess) Start() error {
	if !p.isRunning() {
		internalArgs := collectInternalArgs(p.fileReaderFn, p.ctx)
		args := make([]string, len(internalArgs)+len(p.ctx.Cfg.DnsmasqArgs))
		copy(internalArgs, args)
		copy(p.ctx.Cfg.DnsmasqArgs, args[len(internalArgs):])
		cmd := exec.Command("dnsmasq", args...)
		err := cmdutil.ExecPipeCmd(cmd)
		if err != nil {
			return err
		}
		p.proc = cmd
	}
	return nil
}

func (p *DnsmasqProcess) Stop() error {
	if p.isRunning() {
		err := p.proc.Process.Kill()
		if err != nil {
			return err
		}
		p.proc = nil
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

func collectInternalArgs(fileReaderFn FileReaderFn, c *context.MainContext) []string {
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
	dhcpIPChunks := strings.Split(c.Cfg.BridgeAddr, ".")
	ipPrefix := strings.Join(dhcpIPChunks[:3], ".")
	args = append(args, fmt.Sprintf("--dhcp-range=%s.50,%s.250,12h", ipPrefix, ipPrefix))
	return args
}
