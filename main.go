package main

import (
	"log"

	"./app"
	"./base"
)

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
	ctx := app.NewMainContext(cfg)
	// start
	startService(ctx)
	return nil
}

func startService(ctx *app.MainContext) error {
	//
	return nil
}

func stopService(ctx *app.MainContext) error {
	//
	return nil
}

func readConfigFile() ([]byte, error) {
	//
	return nil, nil
}
