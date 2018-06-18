package dnsmasq

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"../base"
)

type DnsmasqProcess struct {
	cfg          *ba.Config
	proc         *exec.Cmd
	fileReaderFn ba.FileReaderFn
	statusLock   sync.RWMutex
	running      bool
	stopChan     chan bool
}

func NewDnsmasqProcess(reader ba.FileReaderFn, c *ba.Config) *DnsmasqProcess {
	return &DnsmasqProcess{
		cfg:          c,
		fileReaderFn: reader,
		running:      false,
	}
}

func (p *DnsmasqProcess) Start() error {
	if !p.isRunning() {
		p.stopChan = make(chan bool)
		internalArgs := collectInternalArgs(p.fileReaderFn, p.cfg)
		args := make([]string, len(internalArgs)+len(p.cfg.DnsmasqArgs))
		copy(args, internalArgs)
		copy(args[len(internalArgs):], p.cfg.DnsmasqArgs)
		log.Printf("run dnsmasq with args: %v", args)
		cmd := exec.Command("dnsmasq", args...)
		err := ba.ExecPipeCmd(cmd)
		if err != nil {
			return err
		}
		p.proc = cmd
		p.setRunning(true)
		log.Println("dnsmasq process started")
		go func() {
			err := cmd.Wait()
			code := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
			if code != 0 {
				if err != nil {
					log.Printf("dnsmasq process exited %d with error: %v", code, err)
				} else {
					log.Printf("dnsmasq process exited with code %d", code)
				}
			} else {
				log.Printf("dnsmasq process exited")
			}
			p.setRunning(false)
			close(p.stopChan)
		}()
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
		if p.stopChan != nil {
			<-p.stopChan
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
	p.statusLock.RLock()
	defer p.statusLock.RUnlock()
	return p.running
}

func (p *DnsmasqProcess) setRunning(r bool) {
	p.statusLock.Lock()
	defer p.statusLock.Unlock()
	p.running = r
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
	ipPrefix := ba.GetSubnetPrefix(c.BridgeAddr)
	args = append(args, fmt.Sprintf("--dhcp-range=%s.50,%s.250,12h", ipPrefix, ipPrefix))
	return args
}
