.POSIX:
.SUFFIXES:

REV     ?= $(shell git rev-parse --short HEAD)
VERSION ?= $(shell git describe --tags --long)

GO      ?= go
GOFLAGS ?=

porcelain: GO_LDFLAGS ?= -s -w
porcelain: GO_LDFLAGS += -X main.date=$(shell date -u -I)
porcelain: GO_LDFLAGS += -X main.commit=$(REV)
porcelain: GO_LDFLAGS += -X main.version=$(VERSION)
porcelain:
	$(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)" ./cmd/$@

porcelain.1: SCDOC ?= scdoc
porcelain.1: porcelain.1.scd
	$(SCDOC) < $< >$@

.PHONY: validate
validate:
	golangci-lint run -j$(shell nproc) || $(GO) vet $(GOFLAGS)

.PHONY: check
check: GO_TESTFLAGS ?= -cover
check:
	$(GO) test $(GOFLAGS) $(GO_TESTFLAGS) ./...

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man

.PHONY: install
install: porcelain porcelain.1
	test -d $(DESTDIR)$(BINDIR) || install -Dm 00755 -d $(DESTDIR)$(BINDIR)
	install -m 00755 porcelain $(DESTDIR)$(BINDIR)/.
	test -d $(DESTDIR)$(MANDIR)/man1 || install -Dm 00755 -d $(DESTDIR)$(MANDIR)/man1
	install -m 00755 porcelain.1 $(DESTDIR)$(MANDIR)/man1/.

.PHONY: release
release: GORELEASER ?= goreleaser
release: porcelain.1
	$(GORELEASER) release

.PHONY: clean
clean:
	rm -fv dist/. porcelain porcelain.1
