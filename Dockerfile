ARG TERRAFORM_VERSION=latest

FROM golang:1.20-alpine AS go-base
WORKDIR /usr/src/app
COPY go.work go.work.sum ./
COPY src/api/go.mod src/api/go.sum ./src/api/
COPY src/pkg/go.mod src/pkg/go.sum ./src/pkg/
COPY src/seeder/go.mod src/seeder/go.sum ./src/seeder/

RUN --mount=type=cache,target=/go/pkg/mod/ \
	go mod download && go work sync

COPY src ./src

FROM go-base AS api-build
WORKDIR /usr/src/app
RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=/go/pkg/mod/ \
	cd ./src/api/cmd && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server

FROM go-base AS seed-build
WORKDIR /usr/src/app
RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=/go/pkg/mod/ <<EOT
	cd ./src/seeder/resources && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/seed_resources
	cd ../modules && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/seed_modules
	cd ../mappings && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/seed_mappings
EOT

FROM alpine AS api-runner
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=api-build /go/bin/server .
ENTRYPOINT ["./server"]

FROM hashicorp/terraform:${TERRAFORM_VERSION} AS seed-runner
RUN apk update && apk add make && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=seed-build /go/bin/seed_resources ./.bin/
COPY --from=seed-build /go/bin/seed_modules ./.bin/
COPY --from=seed-build /go/bin/seed_mappings ./.bin/
COPY Makefile ./

# hack make target to not trigger build since the build is already ready
RUN <<EOT
	mkdir -p ./src/pkg
	mkdir -p ./src/seed
	touch ./.bin/seed_resources
	touch ./.bin/seed_modules
	touch ./.bin/seed_mappings
EOT

ENTRYPOINT [ "make", "seed" ]

FROM golang:1.20 AS unit-test
WORKDIR /usr/src/app
COPY --from=go-base /go /go
COPY . .
ENTRYPOINT [ "make", "test" ]
