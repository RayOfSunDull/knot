.PHONY: build full install bin config


build:
	go build -o ./bin/knot ./src

bin: 
	install -v ./bin/knot $(HOME)/bin/knot

config:
	mkdir -p $(HOME)/.config/knot

	cp -TRnv templates $(HOME)/.config/knot/templates


install: bin config

full: build install