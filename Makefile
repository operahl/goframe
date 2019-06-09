#!/bin/bash
export LANG=zh_CN.UTF-8

ENVARG=GOPATH=$(CURDIR) GO111MODULE=on
LINUXARG=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
BUILDARG=-ldflags " -s -X main.buildTime=`date '+%Y-%m-%dT%H:%M:%S'` -X main.gitHash=`git rev-parse HEAD`"

export GOPATH

dep:
	cd src; ${ENVARG} go get ./...; cd -

ofoodhttp:
	cd src/app/http; ${ENVARG} go build ${BUILDARG} -o ../../../bin/ofoodhttp main.go; cd -

linux_ofoodhttp:
	cd src/app/http; ${ENVARG} ${LINUXARG} go build ${BUILDARG} -o ../../../lbin/ofoodhttp main.go; cd -

clean:
	rm -fr bin/*
	rm -fr lbin/*
	chmod -R 766 pkg/*
	\rm -r pkg/*
	rm src/go.sum
