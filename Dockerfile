ARG TERRAFORM_VERSION=latest

FROM golang:1.20 AS go-base
WORKDIR /usr/src/app
COPY go.work go.work.sum ./
COPY src/api/go.mod src/api/go.sum ./src/api/
COPY src/pkg/go.mod src/pkg/go.sum ./src/pkg/
COPY src/cli/go.mod src/cli/go.sum ./src/cli/
ENV GOPRIVATE=github.com/cldcvr
RUN --mount=type=cache,target=/go/pkg/mod/ \
	cd src/cli && go mod edit -replace github.com/cldcvr/terrarium/src/pkg=../pkg && cd - && \
	cd src/api && go mod edit -replace github.com/cldcvr/terrarium/src/pkg=../pkg && cd - && \
	go mod download && go work sync
COPY src ./src
RUN \
	cd src/cli && go mod edit -replace github.com/cldcvr/terrarium/src/pkg=../pkg && cd - && \
	cd src/api && go mod edit -replace github.com/cldcvr/terrarium/src/pkg=../pkg && cd -

FROM go-base AS api-build
WORKDIR /usr/src/app
RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=/go/pkg/mod/ \
	cd ./src/api && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server

FROM go-base AS cli-build
WORKDIR /usr/src/app
RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=/go/pkg/mod/ \
	cd ./src/cli/terrarium && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/terrarium

FROM alpine AS api-runner
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=api-build /go/bin/server .
ENTRYPOINT ["./server"]

FROM hashicorp/terraform:${TERRAFORM_VERSION} AS farm-harvester
RUN apk update && apk add make && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=cli-build /go/bin/terrarium /bin/
COPY Makefile ./
ENV FARM_DIR=./farm \
	TR_DB_RETRY_ATTEMPTS=20
# workaround to use Makefile
RUN mkdir -p ./src/pkg && \
	mkdir -p ./src/cli
ENTRYPOINT [ "make", "farm-harvest" ]

FROM go-base AS unit-test
WORKDIR /usr/src/app
# following copies the already available go modules from cache
# to the disk. so that it's available on runtime, and is not re-downloaded.
RUN --mount=type=cache,target=/go/pkg/mod/ \
	cp -r /go/pkg/mod /go/pkg/mod.bak
RUN rm -rf /go/pkg/mod && mv /go/pkg/mod.bak /go/pkg/mod
COPY Makefile ./
COPY examples ./examples
ENTRYPOINT [ "make", "test" ]
