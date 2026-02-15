.PHONY: all build test lint arch-lint install clean

all: build test lint arch-lint
.DEFAULT_GOAL := all

build:
	@task build

test:
	@task test

lint:
	@task lint

arch-lint:
	@task lint:arch

install:
	@task install

clean:
	@task clean
