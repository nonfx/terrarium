# Landing Zone implementation using Terrarium

## Make commands

1. Generate the terraform code

    ```sh
    make lz             # generate or re generate the landing zone terraform code
    make lz_cleanup lz  # or force delete previously generated code and generate fresh (not required in most cases)
    ```

2. Lint platform

    ```sh
    make lint_all               # run terrarium lint on all the platforms
    make lint_cleanup lint_all  # delete previously generated platform `terrarium.yaml` files and then run the lint command on each platforms (not required in most cases)
    ```

3. Dry Run

    ```sh
    make dryrun # Print the files and directories identified for the platforms, platforms, requirements and target dirs for generated code
    ```

4. Help

    ```sh
    make help
    ```

Also See:

- [Requirements Doc](./requirements/readme.md)
