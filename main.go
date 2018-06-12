package main

import (
	"log"

	"./app"
	"./base"
	"./bridge"
	"./dnsmasq"
	"./forward"
)

type serviceContext struct {
	webService *app.WebService
	dnsmasq    *dnsmasq.DnsmasqProcess
	bridge     *bridge.Bridge
	forward    *forward.Forward
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	// check
	err := CheckDepedence()
	if err != nil {
		return err
	}
	// config
	cfg, err := ba.ReadConfig(readConfigFile)
	if err != nil {
		return err
	}
	// service ctx
	ctx := &serviceContext{}
	// start
	startService(ctx, cfg)
	return nil
}

func startService(ctx *serviceContext, cfg *ba.Config) error {
	//
	return nil
}

func stopService(ctx *serviceContext) error {
	//
	return nil
}

func readConfigFile() ([]byte, error) {
	//
	return nil, nil
}
