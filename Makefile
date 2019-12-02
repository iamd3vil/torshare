BIN := torshare.bin
DIR := dist
HASH := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format=%ci ${HASH})
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
VERSION := ${HASH} (${COMMIT_DATE})

test:
	go test ./...

coverage:
	echo $(shell go test ./... --cover | awk '{if ($$1 != "?") print $$5; else print "0.0";}' | sed 's/\%//g' | awk '{s+=$$1} END {printf "%.2f\n", s}') / $(shell go test ./... --cover | wc -l) \
	|  awk -F "/" '{print ($$1/$$2)}'

build:
	mkdir -p ${DIR}
	go build -o ${DIR}/${BIN} -ldflags="-X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'"

build-all: build

run: build
	./${DIR}/${BIN}