#!/bin/bash

VERSION=1.0

cd $(dirname $0)

if [ ! -e "dist" ]; then
	mkdir dist
fi

gobuild() {
	if [ -e "dist/tmp" ]; then
		rm -rf dist/tmp
		mkdir dist/tmp
	fi
	go build -o dist/tmp/usr/local/bin/router-service .
	mkdir -p dist/tmp/etc/router-service
	cp config.yml dist/tmp/etc/router-service/config.yml
}

build_deb() {
	fpm -s dir -t deb -p $1 -n router-service -v ${VERSION} -d "bridge-utils (> 0)" -d "iptables (> 0)" -d "dnsmasq (> 0)" -d "net-tools (> 0)" -d "net-tools (> 0)" -C dist/tmp .
}

# build deb x64
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
gobuild
rm -rf dist/linux/x64
mkdir -p dist/linux/x64
build_deb dist/linux/x64/router-service_v${VERSION}_amd64.deb

# build deb armhf
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm
gobuild
rm -rf dist/linux/armhf
mkdir -p dist/linux/armhf
build_deb dist/linux/armhf/router-service_v${VERSION}_armhf.deb
