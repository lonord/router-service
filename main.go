package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lonord/router-service/app"
	"github.com/lonord/router-service/base"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	configFilePath := flag.String("c", "/etc/router-service/config.yml", "config file path")
	flag.Parse()
	// check
	err := CheckDepedence()
	if err != nil {
		return err
	}
	// config
	cfg, err := ba.ReadConfig(readConfigFile(*configFilePath))
	if err != nil {
		return err
	}
	// service ctx
	ctx := app.NewMainContext(cfg)
	// start
	startService(ctx)
	// listen signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)
	for {
		select {
		case sig := <-signalChan:
			switch sig {
			case syscall.SIGHUP:
				fallthrough
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGQUIT:
				fallthrough
			case syscall.SIGTERM:
				log.Println("got signal ", sig.String())
				stopService(ctx)
				return nil
			}
		}
	}
}

func startService(ctx *app.MainContext) {
	err := ctx.SubProcess.Dnsmasq.Start()
	if err != nil {
		log.Println(err)
	}
	err = ctx.SubProcess.Bridge.SetupBridge()
	if err != nil {
		log.Println(err)
	}
	err = ctx.SubProcess.Forward.SetupForward()
	if err != nil {
		log.Println(err)
	}
	ctx.WebService.Start()
	log.Println("service started")
}

func stopService(ctx *app.MainContext) {
	log.Println("service is going to shutdown")
	ctx.WebService.Stop()
	err := ctx.SubProcess.Forward.ClearForward()
	if err != nil {
		log.Println(err)
	}
	err = ctx.SubProcess.Bridge.ClearBridge()
	if err != nil {
		log.Println(err)
	}
	err = ctx.SubProcess.Dnsmasq.Stop()
	if err != nil {
		log.Println(err)
	}
	log.Println("bye")
}

func readConfigFile(p string) ba.ConfigReaderFn {
	return func() ([]byte, error) {
		return ioutil.ReadFile(p)
	}
}
