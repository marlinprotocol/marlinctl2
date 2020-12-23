GO=go
GOBUILD=$(GO) build
BINDIR=build
BINCLI=marlinctl
INSTALLLOC=/usr/local/bin/$(BINCLI)
RELEASE=$(shell git describe --tags --abbrev=0)
BUILDCOMMIT=$(shell git rev-parse HEAD)
BUILDLINE=$(shell git rev-parse --abbrev-ref HEAD)
CURRENTTIME=$(shell date -u '+%d-%m-%Y_%H-%M-%S')@UTC

release:
	$(GOBUILD) -ldflags="\
						-X github.com/marlinprotocol/ctl2/cmd.compilationChain=iris \
						-X github.com/marlinprotocol/ctl2/version.applicationVersion=$(RELEASE)-release \
						-X github.com/marlinprotocol/ctl2/version.buildCommit=$(BUILDLINE)@$(BUILDCOMMIT) \
						-X github.com/marlinprotocol/ctl2/version.buildTime=$(CURRENTTIME) \
						-linkmode=external" \
				-o $(BINDIR)/marlinctl
clean:
	rm -rf $(BINDIR)/*

install:
	cp $(BIN) $(INSTALLLOC)

uninstall:
	rm $(INSTALLLOC)