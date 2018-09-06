package app

import (
	"github.com/lonord/router-service/base"
	"github.com/lonord/router-service/bridge"
	"github.com/lonord/router-service/client"
	"github.com/lonord/router-service/dnsmasq"
	"github.com/lonord/router-service/forward"
	"github.com/lonord/router-service/netutil"
)

type SubProcess struct {
	Dnsmasq *dnsmasq.DnsmasqProcess
	Bridge  *bridge.Bridge
	Forward *forward.Forward
}

type MainContext struct {
	cfg        *ba.Config
	action     *MainAction
	WebService *WebService
	SubProcess *SubProcess
}

func NewMainContext(cfg *ba.Config) *MainContext {
	sub := &SubProcess{
		Dnsmasq: dnsmasq.NewDnsmasqProcess(ba.DefaultFileReader, cfg),
		Bridge:  bridge.NewBridge(ba.DefaultCmdExecutor, cfg),
		Forward: forward.NewForward(ba.DefaultCmdExecutor, cfg),
	}
	action := NewMainAction(sub, cfg)
	return &MainContext{
		cfg:        cfg,
		SubProcess: sub,
		action:     action,
		WebService: NewWebService(action, cfg),
	}
}

type MainAction struct {
	cfg *ba.Config
	sub *SubProcess
}

func NewMainAction(sub *SubProcess, cfg *ba.Config) *MainAction {
	return &MainAction{
		sub: sub,
		cfg: cfg,
	}
}

func (a *MainAction) CreateNetSpeedReader() (*WrappedNetSpeedReader, error) {
	outerIf := a.cfg.OuterIf
	innerIf := a.cfg.BridgeName
	r := netutil.NewNetSpeedReader(ba.DefaultFileReader, outerIf, innerIf)
	err := r.Init()
	if err != nil {
		return nil, err
	}
	return &WrappedNetSpeedReader{
		reader:  r,
		outerIf: outerIf,
		innerIf: innerIf,
	}, nil
}

func (a *MainAction) GetOnlineClients() ([]client.ClientInfo, error) {
	return client.ReadClients(ba.DefaultCmdExecutor, a.cfg)
}

func (a *MainAction) RestartDnsmasq() error {
	return a.sub.Dnsmasq.Restart()
}

type NetStatus struct {
	InnerNetSpeed netutil.NetSpeed `json:"inner"`
	OuterNetSpeed netutil.NetSpeed `json:"outer"`
}

type WrappedNetSpeedReader struct {
	reader  *netutil.NetSpeedReader
	outerIf string
	innerIf string
}

func (r *WrappedNetSpeedReader) Read() (NetStatus, error) {
	statusList, err := r.reader.Read()
	ns := NetStatus{}
	if err != nil {
		return ns, err
	}
	for _, s := range statusList {
		if s.DevName == r.outerIf {
			ns.OuterNetSpeed = s
		}
		if s.DevName == r.innerIf {
			ns.InnerNetSpeed = s
		}
	}
	return ns, nil
}
