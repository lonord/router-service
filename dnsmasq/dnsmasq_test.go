package dnsmasq

import (
	"testing"

	"github.com/lonord/router-service/base"
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
	if args[3] != "--trust-anchor=.,20326,8,2,E06D44B80B8F1D39A95C0B0D7C65D08458E880409BBC683457104237C7F8EC8D" {
		t.Errorf("ds content dismatch [%s]", args[3])
	}
	if args[4] != "--dhcp-range=192.168.8.50,192.168.8.250,12h" {
		t.Errorf("dhcp range content dismatch [%s]", args[3])
	}
}

const dsContent = `. IN DS 20326 8 2 E06D44B80B8F1D39A95C0B0D7C65D08458E880409BBC683457104237C7F8EC8D
`

const leasesContent = `1524448523 00:0e:c6:a6:0c:99 192.168.4.127 lmbp 01:00:0e:c6:a6:0c:99
1528766444 dc:a9:04:86:fd:99 192.168.4.90 * 01:dc:a9:04:86:fd:99
1528751440 9c:f3:87:bd:d7:99 192.168.4.109 cuir 01:9c:f3:87:bd:d7:99
1528770147 78:11:dc:e1:d4:99 192.168.4.59 lumi-gate1 *
1528766261 60:6b:ff:25:29:99 192.168.4.92 * *
1528762195 34:5b:bb:8a:39:99 192.168.4.54 * 01:34:5b:bb:8a:39:99
`
