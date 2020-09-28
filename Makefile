PKG=github.com/leakingtapan/sonoff
IMAGE?=chengpan/sonoff
VERSION=dev
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GO_BUILD_EXTR_ENV?=
LDFLAGS?=""
GO_BUILD_ENV=CGO_ENABLED=0
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
	docker buildx build -t cheng/sonoff:latest --platform=linux/arm64,linux/amd64 --output="type=image,push=false" .
	
bin/swagger:
	curl -o ./bin/swagger -L https://github.com/go-swagger/go-swagger/releases/download/v0.23.0/swagger_${OS}_amd64
	chmod +x ./bin/swagger

swagger-gen:
	find ./apis ! -name swagger.yaml -type f -delete
	./bin/swagger generate model -f ./apis/swagger.yaml --model-package=apis/models
	./bin/swagger generate client -f ./apis/swagger.yaml --model-package=apis/models --client-package=apis/client
