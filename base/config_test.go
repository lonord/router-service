package ba

import (
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	cfg, err := ReadConfig(getContentByte)
	expect(t, err, nil)
	expect(t, cfg.OuterIf, "eth0")
	expect(t, cfg.BridgeName, "br-lan")
	expect(t, cfg.BridgeAddr, "192.168.4.1")
	expect(t, cfg.LanNames[0], "eth1")
	expect(t, len(cfg.DnsmasqArgs), 0)
	expect(t, cfg.RPCHost, "0.0.0.0")
	expect(t, cfg.RPCPort, 2018)
}

func TestReadConfig2(t *testing.T) {
	cfg, err := ReadConfig(getContentByte2)
	expect(t, err, nil)
	expect(t, len(cfg.DnsmasqArgs), 3)
	expect(t, cfg.DnsmasqArgs[0], "--resolv-file=/var/run/dnsmasq/resolv.conf")
	expect(t, cfg.DnsmasqArgs[1], "--log-dhcp")
	expect(t, cfg.DnsmasqArgs[2], "--port=5353")
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(b, a) {
		t.Errorf("Expected %#v (type %v) - Got %#v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
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

# rpc hostname, default 0.0.0.0
# rpcHost: '0.0.0.0'

#rpc port, default 2018
# rpcPort: 2018
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

# rpc hostname, default 0.0.0.0
# rpcHost: '0.0.0.0'

#rpc port, default 2018
# rpcPort: 2018
`
