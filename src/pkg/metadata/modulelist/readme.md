# Module List File Format

The Module List File Format is used to specify Terraform modules for the `harvest` CLI commands. This format allows for more efficient module management and initialization.

## Overview

The client `harvest` commands can read Terraform files directly from a given directory path. The user is required to pre-initialize the working directory by running commands such as `terraform init`. This can lead to a version conflict if some of the loaded modules require incompatible provider versions.

To avoid this caveat the user may instead pass in a file containing a list of modules to be loaded. The `harvest` commands then process each module in a separate workspace hence avoiding any version conflicts between different modules.

In case of the module list file the CLI performs workspace initialization for each module on-the-fly. In order to improve performance it relies on the Terraform built-in plugin cache. The cache location can be controlled by `terraform.plugin_cache_dir` configuration variable pre-set to `~/.terraform.d/plugin-cache`.

## Module List File Structure

Each module entry in the list must include:

- `name`: A unique identifier for the module (required).
- `source`: The Terraform module source (required).
- `version`: The Terraform module version (optional).
- `export`: A boolean indicating if the module should be processed by the module harvest command (optional).
- `group`: A string to group compatible modules for efficient initialization (optional).

Modules with the `export` attribute set to `true` are included in the user's module library. Those without it or set to `false` are used only for resource attribute discovery and mapping.

The `group` attribute allows users to specify a group name for compatible modules. Modules in the same group are initialized together, which can significantly speed up the `terraform init` process. This is particularly useful for large sets of modules with shared dependencies.

## Example Module List File

```yaml
farm:
  - name: vpc
    source: "terraform-aws-modules/vpc/aws"
    version: "4.0.2"
    export: true
    group: "group1" # Grouping modules for efficient initialization
  - name: voting-demo
    source: "github.com/cldcvr/codepipes-tutorials//voting/infra/aws/eks?ref=terrarium-sources"
    export: false
    group: "group1" # Same group as 'vpc' for shared initialization
  - name: rds
    source: "terraform-aws-modules/rds/aws"
    export: true
    group: "group2" # Different group for separate initialization
```

In this example, `vpc` and `voting-demo` are part of the `group1` group and will be initialized together, while `database` is in a separate `group2` group.
