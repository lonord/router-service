package dnsmasq

import (
	"testing"

	"../context"
)

func TestCollectInternalArgs(t *testing.T) {
	c := &context.MainContext{
		Cfg: context.Config{
			BridgeAddr: "192.168.8.1",
		},
	}
	args := collectInternalArgs(func(path string) (string, error) {
		return dsContent, nil
	}, c)
	if len(args) != 5 {
		t.Error("arg length dismatch")
	}
	if args[3] != "--trust-anchor=.,19036,8,2,49AAC11D7B6F6446702E54A1607371607A1A41855200FD2CE1CDDE32F24E8FB5" {
		t.Errorf("ds content dismatch [%s]", args[3])
	}
	if args[4] != "--dhcp-range=192.168.8.50,192.168.8.250,12h" {
		t.Errorf("dhcp range content dismatch [%s]", args[3])
	}
}

const dsContent = `. IN DS 19036 8 2 49AAC11D7B6F6446702E54A1607371607A1A41855200FD2CE1CDDE32F24E8FB5
`
