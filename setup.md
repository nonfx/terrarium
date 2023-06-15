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

Or to re-build the containers with latest code and run:

```sh
make docker-build docker-run
```

## Stop containers

```sh
make docker-stop

OR

make docker-stop-clean # deletes database
```

## Seed database

1. Setup `~/.netrc` with git credentials if using private repository as terraform dependency.
2. Run `make docker-seed` to run terraform init and all seed commands.
3. If you want to make a fresh seed, then truncate the db tables manually, and then run `make clean_tf docker-seed`. To delete the .terraform directory and lock files and start fresh.
4. In order to take backup of the newly seeded database, run `make db-dump`. The backup will be saved in `data/cc_terrarium.sql` directory.
