version :=	$(shell git describe --always || echo "0.0.0")
.DEFAULT_GOAL := build

build:
	go build -ldflags "-w -s -X main.version=$(version)"