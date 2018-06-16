package client

import (
	"testing"

	"../base"
)

func TestReadClients(t *testing.T) {
	cfg := &ba.Config{
		BridgeName: "br-lan",
	}
	clients, err := ReadClients(func(s string) (string, error) {
		return arpResult, nil
	}, cfg)
	if err != nil {
		t.Error(err)
	}
	if len(clients) != 11 {
		t.Error("length dismatch")
	}
	if clients[0].HostName != "andro07" {
		t.Error("HostName dismatch")
	}
	if clients[0].IPAddr != "192.168.4.55" {
		t.Error("IPAddr dismatch")
	}
	if clients[0].MACAddr != "e4:90:7e:08:43:aa" {
		t.Error("MACAddr dismatch")
	}
	if clients[10].HostName != "?" {
		t.Error("HostName dismatch")
	}
	if clients[10].IPAddr != "192.168.4.92" {
		t.Error("IPAddr dismatch")
	}
	if clients[10].MACAddr != "60:6b:ff:25:29:aa" {
		t.Error("MACAddr dismatch")
	}
}

const arpResult = `andro07 (192.168.4.55) at e4:90:7e:08:43:aa [ether] on br-lan
lonord-iPhone (192.168.4.132) at f4:31:c3:63:ab:aa [ether] on br-lan
? (192.168.4.109) at <incomplete> on br-lan
andro94d34 (192.168.4.93) at 80:38:96:85:29:aa [ether] on br-lan
lmbp (192.168.4.127) at 00:0e:c6:a6:0c:aa [ether] on br-lan
iPhone-2 (192.168.4.64) at 20:3c:ae:3c:6f:aa [ether] on br-lan
cuijihone (192.168.4.71) at 9c:e3:3f:47:5a:aa [ether] on br-lan
lumi-gate61 (192.168.4.59) at 78:11:dc:e1:d4:aa [ether] on br-lan
? (192.168.4.90) at <incomplete> on br-lan
pi-station (192.168.4.105) at b8:27:eb:60:3a:aa [ether] on br-lan
? (192.168.4.74) at <incomplete> on br-lan
zhimi-airer-v3370 (192.168.4.81) at 78:11:dc:49:aa:14 [ether] on br-lan
? (192.168.4.54) at 34:5b:bb:8a:aa:e9 [ether] on br-lan
? (192.168.4.92) at 60:6b:ff:25:29:aa [ether] on br-lan
chiPhone (192.168.4.91) at <incomplete> on br-lan
`
