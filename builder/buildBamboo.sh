#!/bin/bash
apk update && apk add curl git bash go haproxy supervisor net-tools && rm -rf /var/cache/apk/*
export GOROOT=/usr/lib/go
export GOPATH=/gopath
export GOBIN=/gopath/bin
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

cd /gopath/src/github.com/QubitProducts/bamboo
go get github.com/tools/godep && \
go get -t github.com/smartystreets/goconvey && \
go build && \
mkdir /var/bamboo && \
cp  /gopath/src/github.com/QubitProducts/bamboo/bamboo /var/bamboo/bamboo && \
mkdir -p /run/haproxy && \
mkdir -p /var/log/supervisor &&\
cd /

rm -rf /tmp/* /var/tmp/*
rm -f /etc/ssh/ssh_host_*
rm -rf /gopath
rm -rf /usr/lib/go
