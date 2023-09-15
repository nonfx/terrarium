# Terrarium Platform Template & Framework Documentation

The Terrarium project introduces a suite of tools aimed to assist DevOps professionals in creating reusable Terraform templates. Central to these tools is the Terrarium Platform Template and Framework. This document will guide you through the process of creating a Terrarium Platform Template using the framework.

## Terrarium Platform Template

The Terrarium Platform Template, henceforth referred to as the 'Platform', is a Terraform template written in HashiCorp Language (HCL) using the Terrarium Platform Framework. This template helps implement Terrarium dependency interfaces, ensuring that your code remains modular and reusable.

### Terrarium Dependency Interface

Terrarium dependency interfaces act as a contract between application developers and DevOps professionals, enabling the specification of application dependencies. They streamline the process of injecting Infrastructure as Code (IaC) dependencies, significantly simplifying app development. For Dependency interface format, refer the component heading in platform metadata documentation [here](../../src/pkg/metadata/platform), and for the app dependency format, [click here](../../src/pkg/metadata/app).

## Terrarium Platform Framework

The Terrarium Platform Framework provides a structured way of implementing Terrarium Dependency Interfaces (Components). It specifies how dependency inputs are provided to a component and how outputs are structured.

### Terrarium Platform Component

A Terrarium Platform Component is a Terraform module call intended to implement a specific dependency interface. Multiple components can reuse other Terraform blocks within the Terrarium Platform Template, such as module calls, resources, data, etc. Each component must follow the naming convention `module.tr_component_<interface name>`.

#### Inputs

In the framework, dependency interface inputs are provided via Terraform local variables. These variables are named using the convention `local.tr_component_<interface name>`. The variable contains an object that houses the app dependency instance name as the key and an object of dependency input values as the value. As a platform author, you can set default values in this object, which would be replaced at the time of Terraform generation.

All component local variables must be defined in the file `tr_locals.tf` so that the Terrarium tools will be able to regenerate the component input values based on dependencies being asked for.

#### Outputs

In the framework, dependency interface outputs are provided via Terraform outputs. The output name follows the convention `tr_component_<interface name>_<output>`. The value of the output is an object, which is keyed by the app dependency instance name.

### Terrarium Platform Metadata

The platform metadata contains detailed information about the Terrarium dependency interfaces implemented within the platform. This metadata is contained within the `terrarium.yaml` file, which is saved alongside the platform HCL code.

The Terrarium tools (cli & vs-code) provide commands that parse the Terrarium Platform Template, show lint errors, and generate the `terrarium.yaml` metadata file. The platform author should review this file to add any missing descriptions or other information to the interface attributes. The metadata file format specification can be found [here](../../src/pkg/metadata/platform).

Using the platform metadata and the app dependency data, the Terrarium tools can determine whether the required app dependencies are implemented within a given Terrarium platform template.

### Generating Terraform Template

Terrarium tools provide the capability of generating Terraform code specific to the app requirement by parsing the Terrarium Platform Template and picking only the necessary IaC code for given app dependencies.

## Example

The directory [examples/platform](.) contains an example Terrarium Platform Template that is also used in the unit tests.

Here is a quick example of a Terrarium Platform Template:

```tf
locals {
  tr_component_postgres = {
    default = {
      version = 11
    },
  }
}

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

module "postgres_security_group" {
  source = "terraform-aws-modules/security-group/aws"

  name   = "postgres_sg"
  vpc_id = module.core_vpc.vpc_id

  # ingress
  ingress_with_cidr_blocks = [
    {
      rule        = "postgresql-tcp"
      cidr_blocks = module.core_vpc.vpc_cidr_block
    }
  ]
  egress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    },
  ]
}

module "core_vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "core_vpc"
  cidr = "10.0.0.0/16"

  azs              = ["eu-west-1a", "eu-west-1b", "eu-west-1c"]
  private_subnets  = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets   = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  database_subnets = local.database_enabled ? ["10.0.21.0/24", "10.0.22.0/24", "10.0.23.0/24"] : []

  enable_nat_gateway = true
  enable_vpn_gateway = true
}

output "tr_component_postgres_host" {
  value = {for k, v in module.tr_component_postgres: k => v.db_instance_address}
}

output "tr_component_postgres_port" {
  value = {for k, v in module.tr_component_postgres: k => v.db_instance_port}
}
```

Here is an example metadata file for the above template:

`terrarium.yaml`

```yaml
components:
- id: postgres
  title: PostgreSQL Database
  description: A relational database management system using SQL.
  inputs:
    properties:
      version:
        title: Engine Version
        description: Postgres engine version
        type: number
        min: 9
        max: 12
        default: 11
      db_name:
        title: Database Name
        description: The name of the database. By default, the dependency name is used.
        type: string
```

## Command

Run following commands in the platform directory.

To lint platform code:

```sh
terrarium platform lint
```

To generate working terraform code based on App dependencies:

```sh
terrarium generate -c dev -a ../apps/voting-be -a ../apps/voting-fe -a ../apps/voting-worker
```

The `terrarium generate` command generates the terraform code, a `tr_gen_profile.auto.tfvars` profile and `*.env.mustache` files for each app in the destination folder (`./.terrarium`).

These files looks something like this:

`app_voting_be.env.mustache`

```sh
BA_LEDGERDB_HOST="{{ tr_component_postgres_host.value.ledgerdb }}"
BA_LEDGERDB_PASSWORD="{{ tr_component_postgres_password.value.ledgerdb }}"
BA_LEDGERDB_PORT="{{ tr_component_postgres_port.value.ledgerdb }}"
BA_LEDGERDB_USERNAME="{{ tr_component_postgres_username.value.ledgerdb }}"
```

As you can see, the env vars personalised for the app are templated referring to a value in terraform state file.
after provisioning infrastructure with terraform, one can render the above template by providing the terraform state file outputs to it. like this:

```sh
terraform output -json | mustache app_banking_app.env.mustache
```

---

By adhering to the conventions and principles set out in this document, DevOps professionals can streamline their development processes and facilitate better collaboration with application developers.
