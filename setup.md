# Setup

## Environment Variables

For the list of available configurations, refer to [the CLI config package](src/cli/internal/config)

CLI config is to be set in `~/.terrarium/config.yaml` file or by exporting the environment variables.

## CLI Installation & Setup

CLI connects with sqlite Database to store the persistent data.

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

    - Download pre-compiled binary from GitHub Release

      - There are downloadable assets associated with each CLI release in Github. These can be downloaded directly from [the Terrarium Github release page](https://github.com/cldcvr/terrarium/releases)
      - The release assets can also be downloaded via wget or curl:
      ```sh
      wget https://github.com/cldcvr/terrarium/releases/download/v0.1/terrarium_v0.1_darwin_arm64.tar.gz

      curl -LO https://github.com/cldcvr/terrarium/releases/download/v0.1/terrarium-v0.1-darwin-arm64.tar.gz
      ```
2. Seed & Run Database

   ```sh
   terrarium farm update
   ```

3. Setup Configuration

   For the list of available configurations, refer to [the config package](src/cli/internal/config).

## Farm Database Updates

To update the farm data when a new farm repo release is available, use the following command:

```sh
rm ~/.terrarium/farm.sqlite
terrarium farm update
```
