# Terrarium Platform Template & Framework

The Terrarium project introduces a suite of tools aimed at assisting DevOps professionals in creating reusable Terraform templates. Central to these tools is the Terrarium Platform Template and Framework. This document will guide you through the process of creating a Terrarium Platform Template using the framework.

## Terrarium Platform Template

The Terrarium Platform Template, henceforth referred to as the 'Platform', is a Terraform template written in HashiCorp Language (HCL) using the Terrarium Platform Framework. This template implements [Terrarium dependency interfaces](../../src/pkg/metadata/dependency/readme.md), ensuring that your code remains modular and reusable. Which is then used for generating IaC code specific to different project requirements specified in the [App Manifest](../../src/pkg/metadata/app/readme.md).

### No-op Platform

This directory [examples/platform](.) contains an example No-op Terrarium Platform Template that is also used in the unit tests.
This example implements a mock of various Terraform modules for storage and compute types such that it can be used to test the Terrarium lint & generate commands and can also run Terraform plan and apply successfully. It'll generate random mock values for output instead of actually provisioning resources on the cloud.

## Terrarium Platform Framework

The Terrarium Platform Framework provides a structured way of implementing Terrarium Dependencies (as components). It specifies how dependency inputs are provided to a component and how outputs are structured.

### Terrarium Platform Component

A Terrarium Platform Component is a Terraform module call intended to implement a specific dependency interface. Multiple components can reuse other Terraform blocks within the Terrarium Platform Template, such as module calls, resources, data, etc. Each component implementation must follow this naming convention: `module.tr_component_<interface name>`.

**Example:**

```tf
# A relational database management system using SQL.
# @title: PostgreSQL Database
module "tr_component_postgres" {
  source = "terraform-aws-modules/rds/aws"

  for_each = local.tr_component_postgres

  identifier     = each.key
  engine_version = coalesce(each.value.version, 11)
  db_name        = coalesce(each.value.db_name, each.key)
  engine         = "postgres"
  family         = format("postgres%s", each.value.version)

  instance_class = coalesce(lookup(var.db_instance_class, each.key, null), var.all_db_instance_class)

  vpc_security_group_ids = [module.postgres_security_group.security_group_id]
  subnet_ids             = module.core_vpc.database_subnets
}
```

#### Inputs

In the framework, dependency interface inputs (coming from apps) are provided to the component by using Terraform local variables. These variables are named using the convention `local.tr_component_<interface name>`. The variable contains an object that houses the app dependency identifiers as the key and an object of dependency input values as the value. As a platform author, you can set default values in this object, which would be replaced at the time of Terraform generation.

```tf
locals {
  tr_component_postgres = {

    # This 'default' block is only for documentation and will be replaced with actual
    # app dependency identifiers like 'my_db' instead of 'default' at the time of running 'terrarium generate'.
    default = {

      # Engine Version of the PostgreSQL engine to use
      # @enum: 11, 12, 13
      version = 11

    },

  }
}
```

#### Outputs

In the framework, dependency interface outputs are provided via Terraform outputs. The output name follows the convention `tr_component_<interface name>_<output>`. The value of the output is an object, which is keyed by the app dependency identifiers.

**Example:**

```tf
output "tr_component_postgres_host" {
  description = "The host address of the PostgreSQL server."
  value       = { for k, v in module.tr_component_postgres : k => v.db_instance_address }
}

output "tr_component_postgres_port" {
  description = "The port number on which the PostgreSQL server is listening."
  value       = { for k, v in module.tr_component_postgres : k => v.db_instance_port }
}
```

#### Comment Annotation

Comment annotations are supported in Terrarium to add metadata for the component and its inputs.

Basically, a simple comment line above the component or input field declaration defines its description. Additional information like @title, @enum, etc., JSON Schema declaration format tags can be set using the '@' symbol.

Component metadata - are set on the module call of the component. Example:

```tf
# A relational database management system using SQL.
# @title: PostgreSQL Database
module "tr_component_postgres" {
```

Component inputs metadata - are set at the component inputs declaration in the 'default' object. Example:

```tf
locals {
  tr_component_postgres = {
    default = {

      # Version of the PostgreSQL engine to use
      # @enum: 11, 12, 13
      version = 11

    },
  }
}
```

## Platform Lint

The Terrarium tools (CLI & VS Code) provide a lint command to parse the Terrarium Platform Template, show any framework lint errors, and generate the platform metadata file (`terrarium.yaml`). The generated platform metadata should be pushed to git as part of the platform development workflow. Read more about Platform Metadata [here](../../src/pkg/metadata/platform/readme.md).

NOTE: The Terrarium lint command **does not** return terraform errors if any. Hence, it is recommended to check for terraform errors separately before running the lint command.

### Lint Command

```sh
# recommended pre-run
terraform init && terraform plan

# help
terrarium platform lint -h

# example
terrarium platform lint
```

## Terrarium Generate

Once a T8 Platform template is developed, it can be used to generate IaC for different project-specific requirements. The following step explains how to use an existing platform template and generate project-specific IaC using Terrarium CLI.

### Generate Command

```sh
# help
terrarium generate -h

# example
terrarium generate -a ../apps/voting-worker
```

### Demo

![Demo GIF](./vhs/platform_usage.gif)

#### Example Platform

No-op platform at [examples/platform](./terrarium.yaml) is used in this example.

#### Example App Project

Sample app manifests that are used in this example project. Note these apps are interdependent.

- [voting-worker](../apps/voting-worker/terrarium.yaml)
- [voting-be](../apps/voting-be/terrarium.yaml)
- [voting-fe](../apps/voting-fe/terrarium.yaml)

#### Generate Flow

1. Terrarium Generate

    ```sh
    terrarium generate -c dev -a ../apps/voting-be -a ../apps/voting-fe -a ../apps/voting-worker
    ```

    The `terrarium generate` command generates the Terraform code. Along with app-based environment variables templates in `*.env.mustache` files for each app in the destination folder (`./.terrarium`).

    Example app_voting_be.env.mustache file:

    ```sh
    VO_MSG_QUEUE_HOST="{{ tr_component_redis_host.value.msg_queue }}"
    VO_MSG_QUEUE_PASSWORD="{{ tr_component_redis_password.value.msg_queue }}"
    VO_MSG_QUEUE_PORT="{{ tr_component_redis_port.value.msg_queue }}"
    VO_VOTING_BE_HOST="{{ tr_component_server_web_host.value.voting_be }}"
    ```

    As you can see, the env vars are personalized for each app and are templated using [mustache](https://mustache.github.io/) referring to a value in Terraform state outputs.

2. Provision infrastructure

    Run Terraform Plan & Apply. Since this is a no-op example template, it only generates random outputs without provisioning any cloud resources.

    ```sh
    cd .terrarium
    terraform init && terraform plan && terraform apply
    ```

3. Render App Env

    The App Env variable templates can be rendered into actual values easily by using Terraform outputs and mustache tools.

    ```sh
    terraform output -json | mustache app_voting_worker.env.mustache
    terraform output -json | mustache app_voting_be.env.mustache
    terraform output -json | mustache app_voting_fe.env.mustache
    ```
