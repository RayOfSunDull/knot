build:
	go build -o ./bin/knot ./src

.PHONY: install bin config


install: bin config

bin: build
	install -v ./bin/knot $(HOME)/bin/knot

config:
	mkdir -p $(HOME)/.config/knot

	cp -TRnv templates $(HOME)/.config/knot/templates