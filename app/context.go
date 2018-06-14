package app

import (
	"../base"
	"../bridge"
	"../dnsmasq"
	"../forward"
	"../netutil"
)

type SubProcess struct {
	Dnsmasq *dnsmasq.DnsmasqProcess
	Bridge  *bridge.Bridge
	Forward *forward.Forward
}

type MainContext struct {
	Cfg        *ba.Config
	Action     *MainAction
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
		Cfg:        cfg,
		SubProcess: sub,
		Action:     action,
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

func (a *MainAction) GetDnsmasqClients() ([]dnsmasq.DnsmasqLease, error) {
	return a.sub.Dnsmasq.ReadLeases()
}

type NetStatus struct {
	InnerNetSpeed netutil.NetSpeed
	OuterNetSpeed netutil.NetSpeed
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
