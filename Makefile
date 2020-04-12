PKG=github.com/leakingtapan/sonoff
IMAGE?=chengpan/sonoff
VERSION=dev
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GO_BUILD_EXTR_ENV?=
LDFLAGS?=""
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux
GO111MODULE=on
OS=$(shell uname | tr '[:upper:]' '[:lower:]')
#LDFLAGS?="-s -w -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT} -X"

COMMANDS=sonoff
BINARIES=$(addprefix bin/,$(COMMANDS))

ifneq ($(GO_BUILD_EXTR_ENV),)
GO_BUILD_ENV += ${GO_BUILD_EXTR_ENV}
endif

.EXPORT_ALL_VARIABLES:

build: bin/sonoff
	mkdir -p bin

server:
	go run cmd/sonoff/main.go server

switch:
	go run cmd/sonoff/main.go switch

binaries: $(BINARIES)

bin/%: cmd/%
	${GO_BUILD_ENV} go build -ldflags ${LDFLAGS} -o $@ ./$<

clean:
	rm -rf bin/

image:
	docker build -t chengpan/sonoff:armv7 .
	docker build -t chengpan/sonoff:amd64 .
	docker push chengpan/sonoff:armv7
	docker push chengpan/sonoff:amd64
	docker manifest create chengpan/sonoff:latest chengpan/sonoff:armv7 chengpan/sonoff:amd64
	docker manifest annotate chengpan/sonoff:latest chengpan/sonoff:armv7 --os linux --arch arm
	docker manifest annotate chengpan/sonoff:latest chengpan/sonoff:amd64 --os linux --arch amd64
	docker manifest push chengpan/sonoff:latest

bin/swagger:
	curl -o ./bin/swagger -L https://github.com/go-swagger/go-swagger/releases/download/v0.23.0/swagger_${OS}_amd64
	chmod +x ./bin/swagger

swagger-gen:
	find ./apis ! -name swagger.yaml -type f -delete
	./bin/swagger generate model -f ./apis/swagger.yaml --model-package=apis/models
	./bin/swagger generate client -f ./apis/swagger.yaml --model-package=apis/models --client-package=apis/client
