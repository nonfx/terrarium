# Setup

## Prerequisites

Before you begin, ensure you have the following installed:

- Go version >= 1.20
- Docker & Docker Compose CLI
- GitHub cli (gh)

## Environment Variables

The following environment variables are required:

```sh
TR_DB_PASSWORD= # Choose a custom password. This variable is used by docker-compose to set the password in the local Postgres server and is used in the API to connect to the database.
```

For the list of available configurations, refer to [the CLI config package](src/cli/internal/config) & [the API config package](src/api/internal/config).

API config is set in the `.env` file in the current folder. While the CLI config is to be set in `~/.terrarium.yaml` file or by exporting the environment variables.

## CLI Installation & Setup

CLI connects with PostgreSQL Database to store the persistent data.

1. Install CLI

    There are multiple ways to install the terrarium CLI:

    - Clone the repo and compile:

      ```sh
      git clone git@github.com:cldcvr/terrarium.git
      cd terrarium
      make install
      ```

    - Install using go package manager

      ```sh
      go install github.com/cldcvr/terrarium/src/cli/terrarium@latest
      ```

    - Download pre-compiled binary from GitHub Release (coming soon)

2. Seed & Run Database

   ```sh
   echo "TR_DB_PASSWORD=<choose a password>" > .env
   make farm-release-pull docker-run
   ```

3. Setup Configuration

   Configure the DB password chosen above in the CLI as well. You can set the environment variable using:

   ```sh
   export TR_DB_PASSWORD=<your_password>
   ```

   Alternatively, you can write it into the configuration file in the home directory:

   ```sh
   echo "db:\n  password: <your_password>" > ~/.terrarium.yaml
   ```

   For the list of available configurations, refer to [the config package](src/cli/internal/config).

## Farm Database Updates

To update the farm data when a new farm repo release is available, use the following command:

```sh
make farm-release-pull db-update
```

## Run Tests

```sh
make test
```

## Re-build and Run API Containers

```sh
make docker-build docker-run
```

### Stop Containers

To stop containers, you have two options:

```sh
make docker-stop
```

or

```sh
make docker-stop-clean # This will also delete the database.
```
