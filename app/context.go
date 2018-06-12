package app

import (
	"../base"
	"../dnsmasq"
	"../netutil"
)

type MainContext struct {
	Cfg    ba.Config
	Action *MainAction
}

type NetStatus struct {
	InnerNetSpeed netutil.NetSpeed
	OuterNetSpeed netutil.NetSpeed
}

type MainAction interface {
	CreateNetSpeedReader() (*netutil.NetSpeedReader, error)
	GetDnsmasqClients() ([]dnsmasq.DnsmasqLease, error)
}
