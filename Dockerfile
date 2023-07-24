ARG TERRAFORM_VERSION=latest

FROM golang:1.20 AS go-base
WORKDIR /usr/src/app
COPY go.work go.work.sum ./
COPY src/api/go.mod src/api/go.sum ./src/api/
COPY src/pkg/go.mod src/pkg/go.sum ./src/pkg/
COPY src/cli/go.mod src/cli/go.sum ./src/cli/
ENV GOPRIVATE=github.com/cldcvr
RUN --mount=type=cache,target=/go/pkg/mod/ \
	go mod download && go work sync
COPY src ./src

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
ENV FARM_DIR=./farm
# workaround to use Makefile
RUN mkdir -p ./src/pkg && \
	mkdir -p ./src/cli
ENTRYPOINT [ "make", "farm-harvest" ]

FROM golang:1.20 AS unit-test
WORKDIR /usr/src/app
COPY --from=go-base /go /go
COPY . .
ENTRYPOINT [ "make", "test" ]
