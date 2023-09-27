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

      - The release assets can also be downloaded via wget or curl.

      - The below steps configure terrarium cli version 0.5. Similar steps can be followed to configure other versions of CLI.

        ```sh
        wget https://github.com/cldcvr/terrarium/releases/download/v0.5/terrarium-v0.5-macos-arm64.tar.gz
        ```

        ```sh
        curl -LO https://github.com/cldcvr/terrarium/releases/download/v0.5/terrarium-v0.5-macos-arm64.tar.gz
        ```
      - Once the the file is downloaded, the tar utility can be used to extract the binary.

        ```sh
        tar -xzvf terrarium-v0.5-macos-arm64.tar.gz
        ```

      - To make the binary runnable from anywhere, it should be moved to a directory that is included in system's PATH.

        Common directories for user binaries include `/usr/local/bin` or `~/bin` (if it exists and is in PATH).

        ```
        mv terrarium /usr/local/bin/
        ```

2. Seed & Run Database

   ```sh
   rm ~/.terrarium/farm.sqlite
   terrarium farm update
   ```


## Farm Database Updates

To update the farm data when a new farm repo release is available, use the following command:

```sh
rm ~/.terrarium/farm.sqlite
terrarium farm update
```

## Environment Variables

For the list of available configurations, refer to [the CLI config package](src/cli/internal/config/defaults.yaml)

CLI config is to be set in `~/.terrarium/config.yaml` file or by exporting the environment variables.

