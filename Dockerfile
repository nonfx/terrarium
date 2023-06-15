FROM golang:1.19-alpine AS go-base
	WORKDIR /usr/src/app
	COPY go.work go.work.sum ./
	COPY api/go.mod api/go.sum ./api/
	COPY api/pkg/terraform-config-inspect/go.mod api/pkg/terraform-config-inspect/go.sum ./api/pkg/terraform-config-inspect/
	RUN go mod download && go work sync
	COPY api ./api

FROM go-base AS api-build
	WORKDIR /usr/src/app
	RUN cd api/cmd/server && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server

FROM go-base AS seed-build
	WORKDIR /usr/src/app
	RUN cd api/cmd/seed_resources && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/seed_resources
	RUN cd api/cmd/seed_modules && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/seed_modules
	RUN cd api/cmd/seed_mappings && CGO_ENABLED=0 GOOS=linux go build -o /go/bin/seed_mappings

FROM alpine AS api-runner
	RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
	WORKDIR /app
	COPY --from=api-build /go/bin/server .
	ENTRYPOINT ["./server"]

FROM hashicorp/terraform:1.4 AS seed-runner
	RUN apk update && apk add make && rm -rf /var/cache/apk/*
	WORKDIR /app
	COPY --from=seed-build /go/bin/seed_resources ./.bin/
	COPY --from=seed-build /go/bin/seed_modules ./.bin/
	COPY --from=seed-build /go/bin/seed_mappings ./.bin/
	COPY Makefile ./
	RUN touch ./.bin/seed_resources && touch ./.bin/seed_modules && touch ./.bin/seed_mappings
	ENTRYPOINT [ "make", "seed" ]

FROM golang:1.19 AS unit-test
	WORKDIR /usr/src/app
	COPY --from=go-base /go /go
	COPY . .
	ENTRYPOINT [ "make", "test" ]
