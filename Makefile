.PHONY: all install

all: install

build:
	bash build.sh

install:
	go build -v -o $(shell go env GOPATH)/bin/mac_tools main.go

clean:
	rm -rf output
	rm -rf mac_tools.tgz

test:
	mkdir -p output
	mac_tools -from '@cur' -to output

tar: clean build
	tar -zcvf mac_tools.tgz ./*
