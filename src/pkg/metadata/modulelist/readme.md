# Farm Module List File Format

Module list can be used to provide Terraform modules to be loaded by CLI `harvest` commands.

## Overview

The client `harvest` commands can read Terraform files directly from a given directory path. The user is required to pre-initialize the working directory by running commands such as `terraform init`. This can lead to a version conflict if some of the loaded modules require incompatible provider versions.

To avoid this caveat the user may instead pass in a file containing a list of modules to be loaded. The `harvest` commands then process each module in a separate workspace hence avoiding any version conflicts between different modules.

In case of the module list file the CLI performs workspace initialization for each module on-the-fly. In order to improve performance it relies on the Terraform built-in plugin cache. The cache location can be controlled by `terraform.plugin_cache_dir` configuration variable pre-set to `~/.terraform.d/plugin-cache`.

## Module List File Structure

Each entry has a required `source` and an optional `version` attribute that maps to Terraform module call `source` and `version` attributes.
In addition to these each may also declare a boolean `export` attribute. Only modules that set it to `true` will be processed by the `module-harvest` command - i.e. will be included in the user's module library.
Finally each entry must declare a unique `name` identifier. This value will be used as exported module's name in the library.

## Example Module List File

```yaml
farm:
  - name: vpc # REQUIRED: unique user-defined module name (will be used as name for exported modules)
    source: "terraform-aws-modules/vpc/aws" # REQUIRED: Terraform module source
    version: "4.0.2" # OPTIONAL: Terraform module version
    export: true # OPTIONAL: if true the module can be imported to the module library, otherwise it will be used only in discovery of resource attributes and mappings
  - name: voting-demo
    source: "github.com/cldcvr/codepipes-tutorials//voting/infra/aws/eks?ref=terrarium-sources"
    export: false # this modules will only be used to discover resource attributes and mappings between resources

```
