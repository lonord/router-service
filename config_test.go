package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	c := MainContext{}
	err := ReadConfig(getContentByte, &c)
	expect(t, err, nil)
	expect(t, c.Cfg.OuterIf, "eth0")
	expect(t, c.Cfg.BridgeName, "br-lan")
	expect(t, c.Cfg.BridgeAddr, "192.168.4.1")
	expect(t, c.Cfg.LanNames[0], "eth1")
	expect(t, len(c.Cfg.DnsmasqArgs), 0)
}

func TestReadConfig2(t *testing.T) {
	c := MainContext{}
	err := ReadConfig(getContentByte2, &c)
	expect(t, err, nil)
	expect(t, len(c.Cfg.DnsmasqArgs), 3)
	expect(t, c.Cfg.DnsmasqArgs[0], "--resolv-file=/var/run/dnsmasq/resolv.conf")
	expect(t, c.Cfg.DnsmasqArgs[1], "--log-dhcp")
	expect(t, c.Cfg.DnsmasqArgs[2], "--port=5353")
}

func getContentByte() ([]byte, error) {
	return []byte(configContent), nil
}

func getContentByte2() ([]byte, error) {
	return []byte(configContent2), nil
}

const configContent = `
###############################################################################
#                                                                             #
#                        configure file of soft router                        #
#                                                                             #
###############################################################################

# interface connected to internet, default "eth0"
outerIf: 'eth0'

# bridge name, default "br-lan"
brName: 'br-lan'

# Ethernets to add to bridge
lan:
  - 'eth1'

# address of bridge, default "192.168.4.1"
brAddr: '192.168.4.1'

# dnsmasq additional arguments, optional
# dnsmasqArgs:
#  - '--resolv-file=/var/run/dnsmasq/resolv.conf'
#  - '--log-dhcp'
`

const configContent2 = `
###############################################################################
#                                                                             #
#                        configure file of soft router                        #
#                                                                             #
###############################################################################

# interface connected to internet, default "eth0"
outerIf: 'eth0'

# bridge name, default "br-lan"
brName: 'br-lan'

# Ethernets to add to bridge
lan:
  - 'eth1'

# address of bridge, default "192.168.4.1"
brAddr: '192.168.4.1'

# dnsmasq additional arguments, optional
dnsmasqArgs:
  - '--resolv-file=/var/run/dnsmasq/resolv.conf'
  - '--log-dhcp'
  - '--port=5353'
`
