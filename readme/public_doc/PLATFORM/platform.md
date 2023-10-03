---
title: "Terrarium Platform Template & Framework Documentation"
slug: "terrarium-platform-template-framework-doc"
excerpt: "A guide to creating Terrarium Platform Templates using the Terrarium Platform Framework."
hidden: false
category: 64fad762908a8c005066ee0a
---

# Terrarium Platform Template & Framework Documentation

The Terrarium project introduces a suite of tools aimed to assist DevOps professionals in creating reusable Terraform templates. Central to these tools is the Terrarium Platform Template and Framework. This document will guide you through the process of creating a Terrarium Platform Template using the framework.

## Terrarium Platform Template

The Terrarium Platform Template, henceforth referred to as the 'Platform,' is a Terraform template written in HashiCorp Language (HCL) using the Terrarium Platform Framework. This template helps implement Terrarium dependency interfaces, ensuring that your code remains modular and reusable.

### Terrarium Dependency Interface

Terrarium dependency interfaces act as a contract between application developers and DevOps professionals, enabling the specification of application dependencies. They streamline the process of injecting Infrastructure as Code (IaC) dependencies, significantly simplifying app development. For Dependency interface format, refer to the component heading in platform metadata documentation [here](../../src/pkg/metadata/platform), and for the app dependency format, [click here](../../src/pkg/metadata/app).

## Terrarium Platform Framework

The Terrarium Platform Framework provides a structured way of implementing Terrarium Dependency Interfaces (Components). It specifies how dependency inputs are provided to a component and how outputs are structured.

### Terrarium Platform Component

A Terrarium Platform Component is a Terraform module call intended to implement a specific dependency interface. Multiple components can reuse other Terraform blocks within the Terrarium Platform Template, such as module calls, resources, data, etc. Each component must follow the naming convention `module.tr_component_<interface name>`.

#### Inputs

In the framework, dependency interface inputs (coming from apps) are provided via Terraform local variables. These variables are named using the convention `local.tr_component_<interface name>`. The variable contains an object that houses the app dependency instance name as the key and an object of dependency input values as the value. As a platform author, you can set default values in this object, which would be replaced at the time of Terraform generation.

#### Outputs

In the framework, dependency interface outputs are provided via Terraform outputs. The output name follows the convention `tr_component_<interface name>_<output>`. The value of the output is an object, which is keyed by the app dependency instance name.

### Terrarium Platform Metadata

The platform metadata contains detailed information about the Terrarium dependency interfaces implemented within the platform. This metadata is contained within the `terrarium.yaml` file, which is saved alongside the platform HCL code.

The Terrarium tools (cli & vs-code) provide commands that parse the Terrarium Platform Template, show lint errors, and generate the `terrarium.yaml` metadata file. The metadata file format specification can be found [here](../../src/pkg/metadata/platform/readme.md).

Using the platform metadata and the app manifest, the Terrarium tools can determine whether the required app dependencies are implemented within a given Terrarium platform template.

### Generating Terraform Template

Terrarium tools provide the capability of generating Terraform code specific to the app requirement by parsing the Terrarium Platform Template and picking only the necessary IaC code for given app dependencies.

## Example

The directory [examples/platform](.) contains an example Terrarium Platform Template that is also used in the unit tests.

Here is a quick example of a Terrarium Platform Template:

```tf
... (the code section)
```

Here is an example metadata file for the above template:

`terrarium.yaml`

```yaml
... (the metadata section)
```

### Command

Run the following commands in the platform directory.

To generate working Terraform code based on App dependencies:

```sh
terrarium generate -a ../apps/voting-be -a ../apps/voting-fe -a ../apps/voting-worker
```

To lint platform code:

```sh
terrarium platform lint
```

By adhering to the conventions and principles set out in this document, DevOps professionals can streamline their development processes and facilitate better collaboration with application developers.
