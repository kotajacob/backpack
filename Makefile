# backpack
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr/local
GO ?= go
GOFLAGS ?=
RM ?= rm -f

all: backpack

backpack:
	$(GO) build $(GOFLAGS)

clean:
	$(RM) backpack

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f backpack $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/backpack

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/backpack

.DEFAULT_GOAL := all

.PHONY: all backpack clean install uninstall
