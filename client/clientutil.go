package client

import (
	"fmt"
	"strings"

	"github.com/lonord/router-service/base"
)

type ClientInfo struct {
	MACAddr  string `json:"mac"`
	IPAddr   string `json:"ip"`
	HostName string `json:"host"`
}

func ReadClients(execFn ba.CmdExecutorFn, cfg *ba.Config) ([]ClientInfo, error) {
	output, err := execFn(fmt.Sprintf("arp -a -i %s", cfg.BridgeName))
	if err != nil {
		return nil, err
	}
	clients := []ClientInfo{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		chunks := strings.Split(line, " ")
		if len(chunks) < 4 || chunks[2] != "at" {
			continue
		}
		if !strings.HasPrefix(chunks[1], "(") || !strings.HasSuffix(chunks[1], ")") {
			continue
		}
		if !strings.Contains(chunks[3], ":") {
			continue
		}
		ip := chunks[1][1 : len(chunks[1])-1]
		if !strings.HasPrefix(ip, ba.GetSubnetPrefix(cfg.BridgeAddr)) {
			continue
		}
		client := ClientInfo{
			HostName: chunks[0],
			IPAddr:   ip,
			MACAddr:  chunks[3],
		}
		clients = append(clients, client)
	}
	return clients, nil
}
