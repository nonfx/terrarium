# Terraform Providers Schema Processor

This program processes the Terraform providers schema and stores it in a database. It extracts information about providers, resource types, and attributes from the schema file and creates corresponding database entries.

## Prerequisites

Before running the program, make sure you have the following prerequisites:

1. Terraform: Ensure that Terraform is installed on your system.
2. Initialized Terraform Project: Set up a Terraform project directory and run `terraform init` to initialize the project. This step is necessary to download the required providers.

## Getting Started

Follow these steps to run the program:

1. Generate the Terraform providers schema file: Run the following command in your Terraform project directory to generate the schema file in JSON format:

    ```sh
    terraform providers schema -json > .terraform/providers/schema.json
    ```

    This command will generate a file named `schema.json` containing the providers schema.

2. Run the program: Execute the main program file, providing the path to the `tf_resources.json` file as a command-line argument. The program will connect to the database, load the providers schema, and create corresponding entries in the database.

    ```sh
    terrarium harvest resources
    ```

    Note: Make sure you have the terrarium binary installed (`make install`) and the database connection details configured appropriately.

3. Monitor the program execution: The program will log progress messages and errors to the console. Check the console output for information about the providers, resource types, and attributes that were created in the database.

## Additional Notes

- The program assumes that you have the required providers declared in your Terraform project. If any providers are missing, make sure to include them in your project configuration and run `terraform init` to download them.

- Ensure that the database connection details (host, port, credentials, etc.) are correctly configured in the `~/.terrarium.yaml` or via environment variables, depending on your setup.
