.POSIX:
.SUFFIXES:

GO ?= go
RM ?= rm
SCDOC ?= scdoc

REV ?= $(shell git rev-parse --short HEAD)
VERSION ?= $(shell git describe --tags --long)

GOFLAGS     =
GO_LDFLAGS ?= -s -w
GO_LDFLAGS += -X main.date=$(shell date -u -I) -X main.commit=$(REV) -X main.version=$(VERSION)

.PHONY: all
all: porcelain porcelain.1

porcelain:
	$(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)"

porcelain.1: porcelain.1.scd
	$(SCDOC) < $< >$@


.PHONY: validate
validate: GOLANGCI_LINT = golangci-lint
validate:
	$(GOLANGCI_LINT) run

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man

.PHONY: install
install: porcelain porcelain.1
	mkdir -p $(DESTDIR)$(BINDIR)
	mkdir -p $(DESTDIR)$(MANDIR)/man1
	cp -f porcelain $(DESTDIR)$(BINDIR)
	cp -f porcelain.1 $(DESTDIR)$(MANDIR)/man1

.PHONY: clean
clean:
	$(RM) porcelain porcelain.1
