package main

import "os/exec"

func CheckDepedence() error {
	err := check("brctl")
	if err != nil {
		return err
	}
	err = check("iptables")
	if err != nil {
		return err
	}
	err = check("dnsmasq")
	if err != nil {
		return err
	}
	err = check("ifconfig")
	if err != nil {
		return err
	}
	return nil
}

func check(executable string) error {
	_, err := exec.LookPath(executable)
	return err
}
