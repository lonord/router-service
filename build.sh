#!/bin/bash

VERSION=1.1

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
	gobuild
	rm -rf dist/linux/$1
	mkdir -p dist/linux/$1
	fpm -s dir -t deb -a $1 -p dist/linux/$1/router-service_v${VERSION}_$1.deb -n router-service -v ${VERSION} -d "bridge-utils (> 0)" -d "iptables (> 0)" -d "dnsmasq (> 0)" -d "net-tools (> 0)" -C dist/tmp .
}

# build deb x64
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
build_deb amd64

# build deb armhf
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm
build_deb armhf
