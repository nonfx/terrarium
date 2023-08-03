# Dependency YAML Configuration Documentation

The Dependency YAML configuration provides a way to define dependency components and their properties. This document serves as a guide to understanding and utilizing the Dependency YAML configuration.

## Introduction

The Dependency YAML configuration allows DevOps users to define various components, their input and output properties, and their relationships. Each component is identified by a unique `id` and can have inputs and outputs specified using a JSON schema format. Once the changes made are merged the pipeline will trigger the cli call to update the same to the dependecy table.

## Depencency Components

Dependency components are the building blocks of the Dependency YAML configuration. They represent different functionalities or services within the system. Each component is defined by an `id`, inputs, and outputs.

## Example

The directory [farm-example/dependency-interfaces/storage/database/rdbms](.) contains an example Dependency YAML configuration Template.

The following example demonstrates a simplified version of the Terrarium YAML configuration:

```yaml
- id: rdbms
  inputs:
    properties:
      version:
        type: number
      engine:
        type: string
        title: Engine
  outputs:
    properties:
      host:
        type: string
      port:
        type: number



The `rdbms` component represents a Relational Database Management System. It includes the following inputs and outputs:

The `rdbms` component has the following input properties:

- `version`: The version of the RDBMS (Relational Database Management System) to use. This property accepts a number data type.
- `engine`: The engine of the RDBMS. This property accepts a string data type.


The `rdbms` component provides the following output properties:

- `host`: The host address of the RDBMS server.
- `port`: The port number on which the RDBMS server is listening.
```

## Component Extensions

Terrarium YAML configuration allows user to extend components to create variations with specific input properties. Extensions are useful when user wants to define components with slightly different attributes while inheriting most of the properties from a base component.

The `postgres` component can be extended to create specific versions of PostgreSQL databases. Here are some examples:

```
- id: postgres@11
  extends:
    id: rdbms
    inputs:
      version: 11
      engine: postgres

- id: postgres@12
  extends:
    id: rdbms
    inputs:
      version: 12
      engine: postgres

- id: postgres@13
  extends:
    id: rdbms
    inputs:
      version: 13
      engine: postgres

The postgres@11, postgres@12, and postgres@13 components extend the postgres component, defining specific properties for PostgreSQL versions 11, 12, and 13, respectively.

```
By adhering to the conventions and principles set out in this document, DevOps professionals can streamline their development processes and facilitate better collaboration with application developers.
