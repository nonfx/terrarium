---
title: "App Manifest"
slug: "app-manifest"
excerpt: "Simplify application dependency management with structured terrarium.yaml files and Terrarium tools, enabling seamless integration with infrastructure-as-code workflows."
hidden: false
category: 64fad762908a8c005066ee0a
---

# App Manifest

The **App Manifest Format** is a specification that allows application developers to define and manage the dependencies of their applications in a structured manner. This format is used in conjunction with the Terrarium tools, which handle the resolution and provisioning of these dependencies.

## Overview

An application often relies on various services, databases, and other components to function properly. The **App Dependency Format** provides a standardized way to declare and configure these dependencies in a `terrarium.yaml` file. By following this format, developers can easily specify the required dependencies, their inputs, and even map their outputs to environment variables.

The Terrarium tools leverage a [Terrarium platform template](../../../../platform/definition/readme.md) to resolve these dependencies. The platform template utilizes the dependency inputs and exports dependency outputs as Terraform output variables. This allows for seamless integration with infrastructure-as-code workflows.

## `terrarium.yaml` Structure

The `terrarium.yaml` file is where the application's infrastructure dependencies are defined. It follows a specific structure based on the `App` struct, which is used to parse the YAML file. The following elements are available in the `terrarium.yaml` file:

### App

The `App` section represents the main application configuration. It contains the following fields:

- `ID`: A required identifier for the app in the project. It must start with an alphabetic character, can only contain alphanumeric characters, and must not exceed 20 characters in length.
- `Name`: A human-friendly name for the application.
- `EnvPrefix`: The prefix used for environment variables in this app. If not set, it defaults to an empty string.
- `Compute`: Denotes a specific dependency that best classifies the app itself. It can only be of the type `compute/*`. The ID of this dependency is automatically set to the app ID, and it is used to set up the deployment pipeline in Code Pipes and allow other apps to use this app as a dependency.
- `Dependencies`: Lists the required services, databases, and other components that the application relies on. This field is an array of `Dependency` objects.

### Dependency

The `dependencies` section represents a single dependency of the application. It includes the following fields:

- `ID`: A required identifier for the dependency in the project. It must start with an alphabet character, can only contain alphanumeric characters, and must not exceed 20 characters in length.
- `Use`: Indicates the specific dependency interface ID that is used as an app dependency. It may include version as a short-hand expression instead of adding version to the inputs block.
- `EnvPrefix`: Used to prefix the output environment variables related to this dependency. If not set, it defaults to the dependency ID in uppercase.
- `Inputs`: Represents customization options for the selected dependency interface. It is a key-value map where the keys represent the input names, and the values represent the corresponding input values.
- `Outputs`: Maps dependency outputs to environment variables. Each entry in this map consists of an environment variable name as key and dependency output name in the value. The format for the environment variable name gets the prefix later automatically.
- `NoProvision`: Indicates whether the dependency is provisioned in another app. If set to `true`, it means this dependency is shared, and its inputs are set in another app while its outputs are made available in the current app.

## Example `terrarium.yaml` File

```yaml
id: banking_app
name: Banking App
env_prefix: BA

compute:
  use: server_web
  inputs:
    port: 3000

dependencies:
  - id: user_db
    use: postgres@11
    env_prefix: USER
  - id: ledger_db
    use: postgres
    env_prefix: LEDGER
    inputs:
      db_name: ledger
      version: 11
    outputs:
      PG_CON: "host={{host}} user={{username}} password={{password}} dbname={{dbname}} port={{port}} sslmode={{sslmode}}"
  - id: user_cache
    use: redis
  - id: auth_app
    no_provision: true
    use: server_web
    outputs:
      URL: "{{endpoint}}"
```

In the example above, we have defined an application with the ID `banking_app` and the name `Banking App`. The environment variables in this app are prefixed with `BA`. The compute base used for the app is classified as a `server_web` and is configured with a `port` input set to `3000`.

The application has several dependencies:

- The `user_db` dependency is of type `postgres@11` and has the ID `user_db`. It should generate standard postgres environment variables prefixed with `BA_USER_` (e.g., `BA_USER_PGHOST`). No custom inputs or outputs are specified for this dependency.
- The `ledger_db` dependency is of type `postgres` and has the ID `ledger_db`. Its environment variables are prefixed with `LEDGER`. It takes inputs for `db_name` and `version`. Additionally, it maps the `BA_LEDGER_PG_CON` output to the environment variable by resolving the Mustache template with dependency outputs.
- The `user_cache` dependency is of type `redis` and has the ID `user_cache`. Its environment variables are prefixed with the default prefix for the dependency ID in uppercase.
- The `auth_app` dependency is of type `server_web` and has the ID `auth_app`. It is marked as `no_provision`, indicating that it is provisioned in another app. It exports the `BA_AUTH_APP_URL` output, which is mapped to the environment variable by resolving the Mustache template with dependency outputs.

By following this format and providing the necessary information in the `terrarium.yaml` file, developers can effectively manage and configure their application's dependencies using the Terrarium tools.
