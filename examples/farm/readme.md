# Terrarium Farm Documentation

Terrarium Farm is a collection of modules and interfaces designed to enhance the functionality of Terraform and provide smart autocomplete features through the Terrarium VSCode extension. This documentation provides an overview of the Terrarium Farm modules and dependency interfaces.

## Terrarium Farm Modules

The `modules` directory contains a set of Terraform module calls in HCL (HashiCorp Configuration Language) without any attributes. When you run `terraform init`, all the specified modules are downloaded into the `.terraform/modules` directory. These modules are then parsed by the Terrarium CLI to extract information about the Terraform providers, resources, modules, attribute mappings, and more.

Here's an example of a module call in HCL format:

```hcl
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "4.0.2"
}
```

To hide a module from the module indexer in Terrarium (`terrarium harvest modules`) while still using it in resource and resource attribute mapping harvesting (`terrarium harvest mappings`), you can prefix the module name with `tr-hide-`. This allows you to exclude specific modules from indexing while retaining their functionality.

Here's an example of a hidden module call:

```hcl
module "tr-hide-banking-demo" {
  source = "github.com/cldcvr/codepipes-tutorials//tfs/aws-ecr-apprunner-vpc?ref=terrarium-sources"
}
```

## Terrarium Farm Dependency Interfaces

The `dependency-interfaces` directory contains dependency interfaces for Terrarium. These interfaces are harvested into the database and used to provide boilerplate code for platform implementation. They are also utilized in the app dependency specification.

---

We hope this documentation provides you with a clear understanding of the Terrarium Farm modules and dependency interfaces. If you have any further questions or need assistance, please refer to the official Terrarium documentation or reach out to the Terrarium community for support.