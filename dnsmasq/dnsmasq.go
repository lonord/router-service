package dnsmasq

import (
	"fmt"
	"os/exec"
	"strings"

	"../base"
)

type DnsmasqProcess struct {
	cfg          *ba.Config
	proc         *exec.Cmd
	fileReaderFn ba.FileReaderFn
}

func NewDnsmasqProcess(reader ba.FileReaderFn, c *ba.Config) *DnsmasqProcess {
	return &DnsmasqProcess{
		cfg:          c,
		fileReaderFn: reader,
	}
}

func (p *DnsmasqProcess) Start() error {
	if !p.isRunning() {
		internalArgs := collectInternalArgs(p.fileReaderFn, p.cfg)
		args := make([]string, len(internalArgs)+len(p.cfg.DnsmasqArgs))
		copy(internalArgs, args)
		copy(p.cfg.DnsmasqArgs, args[len(internalArgs):])
		cmd := exec.Command("dnsmasq", args...)
		err := ba.ExecPipeCmd(cmd)
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

func collectInternalArgs(fileReaderFn ba.FileReaderFn, c *ba.Config) []string {
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
	dhcpIPChunks := strings.Split(c.BridgeAddr, ".")
	ipPrefix := strings.Join(dhcpIPChunks[:3], ".")
	args = append(args, fmt.Sprintf("--dhcp-range=%s.50,%s.250,12h", ipPrefix, ipPrefix))
	return args
}
