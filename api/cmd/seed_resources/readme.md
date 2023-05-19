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
    terraform providers schema -json > tf_resources.json
    mv tf_resources.json path_to_terrarium/cache_data/tf_resources.json
    ```

    This command will generate a file named `tf_resources.json` containing the providers schema.

2. Run the program: Execute the main program file, providing the path to the `tf_resources.json` file as a command-line argument. The program will connect to the database, load the providers schema, and create corresponding entries in the database.

    ```sh
    go run ./api/cmd/seed_resources
    ```

    Note: Make sure you have the necessary Go dependencies installed (`go mod download`) and the database connection details configured appropriately.

3. Monitor the program execution: The program will log progress messages and errors to the console. Check the console output for information about the providers, resource types, and attributes that were created in the database.

## Additional Notes

- The program assumes that you have the required providers declared in your Terraform project. If any providers are missing, make sure to include them in your project configuration and run `terraform init` to download them.

- Ensure that the database connection details (host, port, credentials, etc.) are correctly configured in the program's code or via environment variables, depending on your setup.

- If any errors occur during the execution of the program, it will log error messages and panic. Review the error messages to identify the cause of the issue.

- It's recommended to review the code and customize it according to your specific database and schema requirements before running it in a production environment.
