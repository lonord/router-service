package main

import (
	"fmt"
	"strings"
)

func SetupForward(execFn CmdExecutorFn, c *MainContext) error {
	deleteIptablesRule(execFn, c)
	execFn("echo 1 > /proc/sys/net/ipv4/ip_forward")
	err := addIptablesRule(execFn, c)
	if err != nil {
		return err
	}
	return nil
}

func ClearForward(execFn CmdExecutorFn, c *MainContext) error {
	err := deleteIptablesRule(execFn, c)
	if err != nil {
		return err
	}
	return nil
}

func addIptablesRule(execFn CmdExecutorFn, c *MainContext) error {
	_, err1 := execFn(fmt.Sprint("iptables -t nat -A ", generateNatRule(c)))
	_, err2 := execFn(fmt.Sprint("iptables -I ", generateForwardSourceRule(c)))
	_, err3 := execFn(fmt.Sprint("iptables -I ", generateForwardDestinationRule(c)))
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}
	return nil
}

func deleteIptablesRule(execFn CmdExecutorFn, c *MainContext) error {
	_, err1 := execFn(fmt.Sprint("iptables -t nat -D ", generateNatRule(c)))
	_, err2 := execFn(fmt.Sprint("iptables -D ", generateForwardSourceRule(c)))
	_, err3 := execFn(fmt.Sprint("iptables -D ", generateForwardDestinationRule(c)))
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}
	return nil
}

func generateNatRule(c *MainContext) string {
	return fmt.Sprintf("POSTROUTING -s %s -o %s -j MASQUERADE", dealWithIPNetmask(c.Cfg.BridgeAddr), c.Cfg.OuterIf)
}

func generateForwardSourceRule(c *MainContext) string {
	return fmt.Sprintf("FORWARD -s %s -j ACCEPT", dealWithIPNetmask(c.Cfg.BridgeAddr))
}

func generateForwardDestinationRule(c *MainContext) string {
	return fmt.Sprintf("FORWARD -d %s -j ACCEPT", dealWithIPNetmask(c.Cfg.BridgeAddr))
}

func dealWithIPNetmask(ip string) string {
	ipChunks := strings.Split(ip, ".")
	return strings.Join([]string{ipChunks[0], ipChunks[1], ipChunks[2], "0/24"}, ".")
}
