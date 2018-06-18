package forward

import (
	"fmt"
	"log"
	"strings"

	"../base"
)

type Forward struct {
	cfg    *ba.Config
	execFn ba.CmdExecutorFn
}

func NewForward(fn ba.CmdExecutorFn, c *ba.Config) *Forward {
	return &Forward{
		execFn: fn,
		cfg:    c,
	}
}

func (f *Forward) SetupForward() error {
	deleteIptablesRule(f.execFn, f.cfg)
	f.execFn("echo 1 > /proc/sys/net/ipv4/ip_forward")
	err := addIptablesRule(f.execFn, f.cfg)
	if err != nil {
		return err
	}
	log.Println("iptables rules setted up")
	return nil
}

func (f *Forward) ClearForward() error {
	err := deleteIptablesRule(f.execFn, f.cfg)
	if err != nil {
		return err
	}
	log.Println("iptables rules cleared")
	return nil
}

func addIptablesRule(execFn ba.CmdExecutorFn, c *ba.Config) error {
	_, err1 := execFn(fmt.Sprint("iptables -t nat -A", generateNatRule(c)))
	_, err2 := execFn(fmt.Sprint("iptables -I", generateForwardSourceRule(c)))
	_, err3 := execFn(fmt.Sprint("iptables -I", generateForwardDestinationRule(c)))
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

func deleteIptablesRule(execFn ba.CmdExecutorFn, c *ba.Config) error {
	_, err1 := execFn(fmt.Sprint("iptables -t nat -D", generateNatRule(c)))
	_, err2 := execFn(fmt.Sprint("iptables -D", generateForwardSourceRule(c)))
	_, err3 := execFn(fmt.Sprint("iptables -D", generateForwardDestinationRule(c)))
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

func generateNatRule(c *ba.Config) string {
	return fmt.Sprintf("POSTROUTING -s %s -o %s -j MASQUERADE", dealWithIPNetmask(c.BridgeAddr), c.OuterIf)
}

func generateForwardSourceRule(c *ba.Config) string {
	return fmt.Sprintf("FORWARD -s %s -j ACCEPT", dealWithIPNetmask(c.BridgeAddr))
}

func generateForwardDestinationRule(c *ba.Config) string {
	return fmt.Sprintf("FORWARD -d %s -j ACCEPT", dealWithIPNetmask(c.BridgeAddr))
}

func dealWithIPNetmask(ip string) string {
	return strings.Join([]string{ba.GetSubnetPrefix(ip), ".0/24"}, "")
}
