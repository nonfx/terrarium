# Terrarium Platform Framework

The Terrarium Platform is an integral part of the Terrarium Project, designed to streamline the generation of Terraform code using a pre-defined set of patterns and application dependencies. It allows DevOps to establish a generic, reusable platform for developers to dynamically generate Terraform code based on their app's dependencies.

The Terrarium Project encompasses tools such as a CLI, a language server, a VS Code extension, and a UI for Terraform code generation. By utilizing the Terrarium Platform and app dependencies in YAML format, Terraform code can be generated automatically.

## Dependency Interface

The Dependency Interface functions as a contract between the app developer and infrastructure developer. It's defined by metadata such as inputs and outputs, described in JSON schema, along with other descriptions. Terrarium comes prepackaged with these Dependency Interfaces, which are categorized based on taxonomy levels 7 or 8. Apps utilize these dependencies by providing inputs and utilizing outputs. The Terrarium Platform implements these Dependency Interfaces by processing inputs and generating outputs. The taxonomical hierarchy outlines metadata of the dependency aiding both providers and users in aspects like discovery and more.

## Application Dependency

An application instantiates dependencies by naming them with an identifier, providing the required inputs, and specifying expected outputs. The Terrarium Platform implements these dependencies as components that utilize dependency inputs and export dependency outputs as Terraform output variables.

A YAML-formatted definition of an application dependency might look like this:

```yaml
app:
  name: Banking App

  service:
    type: web_service
    inputs:
      port: 3000

  dependencies:
    user_db:
      type: postgres
      inputs:
        version: 11
      outputs:
        host: BE_PG_HOST
        port: BE_PG_PORT
    ledger_db:
      type: postgres
      inputs:
        version: 12
      outputs:
        host: WO_PG_HOST
        port: WO_PG_PORT
    user_cache:
      type: redis
      inputs:
        version: 6
      outputs:
        host: WO_REDIS_HOST
        port: WO_REDIS_PORT
```

## Platform

The Terrarium Platform framework consists of a base Terraform template, a set of components, and metadata. Terrarium comes with recommended open-source platforms pre-defined that cover common use cases and enable cloud usage. DevOps may use the framework to define new platforms or make changes to the default platforms.

### Component

A component is an implementation of a dependency interface. It serves as a way to fulfill dependencies by utilizing a Terraform module. The module is crafted with predictable names, such as module.tr_component_<component_name>. To cater to the required dependencies of the application, the component module should leverage the Terraform for_each argument, allowing it to render zero or more instances of these modules.

Dependency inputs are set by the Terrarium generator as a map in Terraform local variables and are available for the component module call. The output variables are defined as Terraform outputs using the component module output.

An example Component code may look like this:

```tf
module tr_component_postgres {
  source   = "terraform-aws-modules/rds/aws"
  for_each = local.tr_component_postgres

  identifier = each.key
  engine_version = each.value.version
  db_name = each.value.db_name
  engine = "postgres"

  vpc_security_group_ids = [ module.postgres_security_group.security_group_id ]
  vpc_id = module.vpc.id
}

output "tr_component_postgres_host" {
  value = {for k, v in module.tr_component_postgres: k => v.db_instance_address}
}

output "tr_component_postgres_port" {
  value = {for k, v in module.tr_component_postgres: k => v.db_instance_port}
}
```

An example generated local variable set may look like this:

```tf
locals {
  tr_component_postgres = {
    backend_db = {
      version = 11
    },
    worker_db = {
      version = 12
    }
  }

  tr_component_postgres_enabled = true
  tr_taxon_database_enabled = true
  tr_taxon_sql_enabled = true
}
```

For each taxon that is enabled, Terrarium generates a local variable. If one or more dependencies utilize the `postgres` component, then taxons such as `local.tr_taxon_database_enabled`, `local.tr_taxon_sql_enabled`, etc., will be set to true.

### Base Template

The base template is a general Terraform code containing module calls, resource calls, etc., without any predefined prefix or filename constraints. Terrarium detects the blocks from the base template that are required to provision each component. To provision application dependencies, Terrarium generates Terraform code by pulling in the components, required blocks, and recursively repeating the process until all dependencies are met. For example, if the `postgres` component depends on the `core_vpc` and `postgres_security_group` modules, and the `redis` component depends on the `core_vpc` and `redis_security_group` modules, and the application calls for two Postgres databases and one Redis, then the final Terraform code would contain the `core_vpc`, `postgres_security_group`, and `redis_security_group` in the base template, in addition to the components.

Note: The intended dependencies across modules, resources, and data must be traceable using input-output attributes or the "depends_on" attribute. Only the required components and their traceable dependencies will end up in the final Terraform code. The rest of the Terraform code will be cleaned up.

### Metadata

The Terrarium Platform metadata file, `terrarium.yaml`, contains information about all the dependencies implemented in the platform. This file is auto-generated by Terrarium and can be edited by the platform author. The metadata file allows Terrarium tools to detect dependencies implemented within the platform repository without re-parsing the Terraform code. It also provides a way for the author to override the dependency interface metadata coming from the Terrarium core. example:

```yaml
dependency_interface:
  postgres:
    taxonomy: [database, sql]
    title: PostgreSQL Database
    description: A relational database management system using SQL.
    inputs:
      properties: # JSON schema format
        version:
          type: string
          description: The version of PostgreSQL to use.
        db_name:
          type: string
          description: The name of the database.
    outputs:
      properties:
        host:
          type: string
          description: The host address of the PostgreSQL server.
        port:
          type: number
          description: The port number on which the PostgreSQL server is listening.
        username:
          type: string
          description: The username for accessing the PostgreSQL database.
        password:
          type: string
          description: The password for accessing the PostgreSQL database.
```
