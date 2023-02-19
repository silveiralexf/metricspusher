#!/usr/bin/make -f

.PHONY: build

ifndef GIT_COMMIT
GIT_COMMIT=$(shell git rev-parse --short HEAD)
endif

FROM_PUSHGATEWAY = "prom-pushgateway"
FROM_SDK = "golang:latest"

GO_DIR ?= $(shell pwd)
GO_PKG ?= $(shell go list -e -f "{{ .ImportPath }}")

GOOS ?= $(shell go env GOOS || echo linux)
GOARCH ?= $(shell go env GOARCH || echo amd64)
CGO_ENABLED ?= 0

ARTIFACT_NAME = "metricspusher"
REPO_NAME = "silveiralexf"
RELEASE_VERSION = "1.0.1"

build:
	@DOCKER_BUILDKIT=1 \
	  docker build -t ${REPO_NAME}/${ARTIFACT_NAME}:${RELEASE_VERSION} \
	  	--platform=linux/amd64 \
	  	--target runtime \
		-f "${DOCKERFILE_PATH}" \
		--progress=plain \
		--network host \
		--build-arg FROM_SDK=${FROM_SDK} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		--build-arg RELEASE_VERSION=${RELEASE_VERSION} .

push:
	@DOCKER_BUILDKIT=1 \
	docker tag ${REPO_NAME}/${ARTIFACT_NAME}:${RELEASE_VERSION} ${REPO_NAME}/${ARTIFACT_NAME}:latest \
	&& docker push ${REPO_NAME}/${ARTIFACT_NAME}:${RELEASE_VERSION} \
	&& docker push ${REPO_NAME}/${ARTIFACT_NAME}:latest


local-build:
	@CGO_ENABLED=0 \
	GOOS=${GOOS} \
	GOARCH=${GOARCH} \
	go build -o  ./bin/${ARTIFACT_NAME} --trimpath . \
	&& cp ./bin/${ARTIFACT_NAME} ${GOPATH}/bin/${ARTIFACT_NAME} \
	&& echo "successfully build on ./bin/${ARTIFACT_NAME}"

setup:
	@docker run -d --name ${FROM_PUSHGATEWAY} \
				  -p 9091:9091 \
				  prom/pushgateway

destroy:
	@docker ps -a | grep ${FROM_PUSHGATEWAY} && docker stop ${FROM_PUSHGATEWAY} && docker rm ${FROM_PUSHGATEWAY}