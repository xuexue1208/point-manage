APP_NAME        := point-manage
APP_VERSION     := $(shell git describe --abbrev=0 --tags)
BUILD_VERSION   := $(shell git log -1 --oneline | base64)
GIT_REVISION    := $(shell git rev-parse --short HEAD)
GIT_BRANCH      := $(shell git name-rev --name-only HEAD)
GO_VERSION      := $(shell go version)
SOURCE          := .
TARGET_DIR      := /usr/point-manage

ifeq (${DEVOS}, windows)
	BUILD_TIME      := $(shell datetime)
else
	BUILD_TIME      := $(shell date "+%FT%T%z")
endif

ifeq (${GOOS}, windows)
	BUILD_TARGET    := ${APP_NAME}.exe
else
	BUILD_TARGET    := ${APP_NAME}
endif

all:
	go build -ldflags                           \
	"                                           \
	-X 'main.AppName=${APP_NAME}' \
	-X 'main.AppVersion=${APP_VERSION}' \
	-X 'main.BuildVersion=${BUILD_VERSION}' \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.GitRevision=${GIT_REVISION}' \
	-X 'main.GitBranch=${GIT_BRANCH}' \
	-X 'main.GoVersion=${GO_VERSION}' \
	-w -s                               \
	"                                           \
	-o ${BUILD_TARGET} ${SOURCE}

ifneq (${GOOS}, windows)
	mkdir -p "output/conf"
	mv ${BUILD_TARGET} output/
endif

clean:
	rm -rf output/

docker:
	cp -r deploy/docker output/
	mv output/${APP_NAME} output/docker/

run:
	go build main.go
	./main

install:
	go build main.go
	/usr/bin/nohub ./main &