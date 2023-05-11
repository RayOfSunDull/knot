.PHONY: build full install bin config python

python:
	cython ./utils/export.py -o ./aux/export.c --embed -3

	gcc ./aux/export.c \
		-Wno-deprecated-declarations \
		-Wl,--copy-dt-needed-entries \
		-o ./bin/export -I/usr/include/python3.11 \
		-L./export-venv/lib/python3.11/site-packages \
		-lpython3
	
	rm ./aux/export.c

build:
	go build -o ./bin/knot ./src

bin: 
	install -v ./bin/knot $(HOME)/bin/knot

config:
	mkdir -p $(HOME)/.config/knot

	cp -TRnv templates $(HOME)/.config/knot/templates

	cp -Tv bin/export $(HOME)/.config/knot/export


install: bin config

full: build python install