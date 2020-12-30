GO=go
GOBUILD=$(GO) build
BINDIR=build
BINCLI=marlinctl
INSTALLLOC=/usr/local/bin/$(BINCLI)
RELEASE=$(shell git describe --tags --abbrev=0)
BUILDCOMMIT=$(shell git rev-parse HEAD)
BUILDLINE=$(shell git rev-parse --abbrev-ref HEAD)
CURRENTTIME=$(shell date -u '+%d-%m-%Y_%H-%M-%S')@UTC
CFGTEMPLATE=$(shell cat version/marlincfg_texttemplate.yaml | sed ':a;N;$!ba;s/\n/\\n/g')

release:
	$(GOBUILD) -ldflags="\
	-X github.com/marlinprotocol/ctl2/version.applicationVersion=$(RELEASE)-release \
	-X github.com/marlinprotocol/ctl2/version.buildCommit=$(BUILDLINE)@$(BUILDCOMMIT) \
	-X github.com/marlinprotocol/ctl2/version.buildTime=$(CURRENTTIME) \
	-linkmode=external" \
	-o $(BINDIR)/marlinctl
clean:
	rm -rf $(BINDIR)/*

install:
	cp $(BINDIR)/marlinctl $(INSTALLLOC)

uninstall:
	rm $(INSTALLLOC)