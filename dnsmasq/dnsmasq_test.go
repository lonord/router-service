package dnsmasq

import (
	"testing"

	"../base"
)

func TestCollectInternalArgs(t *testing.T) {
	c := &ba.Config{
		BridgeAddr: "192.168.8.1",
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

func TestConvertLeases(t *testing.T) {
	leases := convertLeases(leasesContent)
	if len(leases) != 6 {
		t.Errorf("leases length dismatch")
	}
	if leases[0].ExpireTimeStamp != 1524448523 {
		t.Errorf("ExpireTimeStamp dismatch")
	}
	if leases[0].MACAddr != "00:0e:c6:a6:0c:99" {
		t.Errorf("MACAddr dismatch")
	}
	if leases[0].IPAddr != "192.168.4.127" {
		t.Errorf("IPAddr dismatch")
	}
	if leases[0].HostName != "lmbp" {
		t.Errorf("HostName dismatch")
	}
	if leases[0].ClientID != "01:00:0e:c6:a6:0c:99" {
		t.Errorf("ClientID dismatch")
	}
	if leases[1].ExpireTimeStamp != 1528766444 {
		t.Errorf("ExpireTimeStamp dismatch")
	}
	if leases[1].MACAddr != "dc:a9:04:86:fd:99" {
		t.Errorf("MACAddr dismatch")
	}
	if leases[1].IPAddr != "192.168.4.90" {
		t.Errorf("IPAddr dismatch")
	}
	if leases[1].HostName != "*" {
		t.Errorf("HostName dismatch")
	}
	if leases[1].ClientID != "01:dc:a9:04:86:fd:99" {
		t.Errorf("ClientID dismatch")
	}
	if leases[4].ExpireTimeStamp != 1528766261 {
		t.Errorf("ExpireTimeStamp dismatch")
	}
	if leases[4].MACAddr != "60:6b:ff:25:29:99" {
		t.Errorf("MACAddr dismatch")
	}
	if leases[4].IPAddr != "192.168.4.92" {
		t.Errorf("IPAddr dismatch")
	}
	if leases[4].HostName != "*" {
		t.Errorf("HostName dismatch")
	}
	if leases[4].ClientID != "*" {
		t.Errorf("ClientID dismatch")
	}
}

const dsContent = `. IN DS 19036 8 2 49AAC11D7B6F6446702E54A1607371607A1A41855200FD2CE1CDDE32F24E8FB5
`

const leasesContent = `1524448523 00:0e:c6:a6:0c:99 192.168.4.127 lmbp 01:00:0e:c6:a6:0c:99
1528766444 dc:a9:04:86:fd:99 192.168.4.90 * 01:dc:a9:04:86:fd:99
1528751440 9c:f3:87:bd:d7:99 192.168.4.109 cuir 01:9c:f3:87:bd:d7:99
1528770147 78:11:dc:e1:d4:99 192.168.4.59 lumi-gate1 *
1528766261 60:6b:ff:25:29:99 192.168.4.92 * *
1528762195 34:5b:bb:8a:39:99 192.168.4.54 * 01:34:5b:bb:8a:39:99
`
