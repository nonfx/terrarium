# Setup

## Env

Following env vars are required:

```sh
TR_DB_PASSWORD= # Choose a custom password. This variable is used by docker-compose to set password in postgres local server and is used in API to connect to the database
```

## Run tests

```sh
make test
```

## Start containers

```sh
make docker-run
```

## Stop containers

```sh
make docker-stop

OR

make docker-stop-clean # deletes database
```

## Update terraform resources

1. Follow pre-requisite in [api/cmd/seed_resources/readme.md](api/cmd/seed_resources/readme.md)
2. Run `make seed`
