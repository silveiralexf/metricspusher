# syntax=docker/dockerfile:1

ARG FROM_SDK
ARG FROM_RUNTIME
ARG GIT_COMMIT
ARG RELEASE_VERSION

FROM ${FROM_SDK} AS builder
ARG GIT_COMMIT
ARG RELEASE_VERSION
ENV GOPATH="/go"
WORKDIR  ${GOPATH}/src/github.com/silveiralexf/metricsbuilder
COPY . .
ENV GIT_COMMIT=$GIT_COMMIT
ENV RELEASE_VERSION=$RELEASE_VERSION
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
        -ldflags="-X github.com/silveiralexf/metricsbuilder/cmd.releaseNumber=v${RELEASE_VERSION}-${GIT_COMMIT}" \
        -o  ${GOPATH}/bin/metricspusher --trimpath . 

FROM scratch as runtime
WORKDIR /bin
COPY --from=builder /go/bin/metricspusher /bin/metricspusher
ENTRYPOINT [ "/bin/metricspusher" ]

