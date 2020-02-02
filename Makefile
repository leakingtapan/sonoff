PKG=github.com/leakingtapan/sonoff
IMAGE?=chengpan/sonoff-server
VERSION=v0.1.0-dirty
GIT_COMMIT?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS?=""
#LDFLAGS?="-X ${PKG}/pkg/driver.driverVersion=${VERSION} -X ${PKG}/pkg/driver.gitCommit=${GIT_COMMIT} -X ${PKG}/pkg/driver.buildDate=${BUILD_DATE} -s -w"

COMMANDS=sonoff
BINARIES=$(addprefix bin/,$(COMMANDS))

.EXPORT_ALL_VARIABLES:

.PHONY: build
build:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -ldflags ${LDFLAGS} -o bin/sonoff-server ./cmd/sonoff-server/main.go

.PHONY: server
server:
	go run cmd/sonoff-server/main.go

switch:
	go run cmd/sonoff-switch/main.go

binaries: $(BINARIES)

bin/%: cmd/%
	go build -o $@ ./$<

clean:
	rm -rf bin/
