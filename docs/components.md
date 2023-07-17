## Terrarium Components

#### The Terrarium Farm
The Terrarium Farm is the library shared across all the components. The library is composed of Git repos and a Postgres DB. It can reside locally on a single machine or be shared using a remote Postgres DB. The components discussed below interact with the Farm to accomplish particular functions.

######Example usage:
You have an open source Terraform module that is being used often and you want to see it suggested automatically in Terrarium (say via the VS Code plugin). To do this, you would create pull request in the terrarium-farm Git repo which, upon approval, would add it to the shared PostgreSQL database.

Other examples of use-cases for interacting with the farm:
- incrementing the version of an existing module
- add a new version of an existing module
- mapping improvements
- taxonomy addition
- dependency interface addition/update

######Process:
1. Verify the new item isn't already in the farm; verify the license & minimum number of stars, last update time. For a module, verify that it is compatible with the cloud provider version being used.
2. Raise a pull request and pass the automated checks.
3. Get approval from farm maintainers.

TODO:
- discuss how each artifact is put into the farm
- how does the taxonomy get established?


#### Curation and Publishing

These functions are generally the purview of the "Devops" persona. The goal is to create and maintain the various objects that are used downstream other personas (such as "Developer"). This is accomplished either by directly adding objects to the Terrarium Farm or by executing a process that results in objects being added to the Farm. The first objects supported are based on Terraform - i.e. infrastructure objects.

There are 3 types of objects related to infrastructure:
1. Terraform modules - These can be imported via the "harvest" command and then used to build up the other objects or new Terraform templates using the Terriaum code generator.
2. Terrarium Platform Template - A base infrastructure pattern written in HCL that follows the [Terrarium Platform Framework](https://github.com/cldcvr/terrarium/blob/main/platform/definition/readme.md) that is imported using the Terrarium CLI/API. A platform is used as the base layer that the dependencies used in applications plugs into. The platform code (HCL) is referenced in the Farm repo but may live in it's own repo.
3. Dependencies - Typically platform independent objects that are resolved to concrete code at generation time. For example, a Postgres DB would be a dependency that when plugged into a platform defined for AWS might resolve using an RDS module. Dependencies are specified using YAML. See [Dependency specification format](TBD)

TODO:
- discuss how taxonomy is assigned and/or used for each type
- walk through how each type is created and added to Farm
- usage of IDE plugin?

#### Composer

This is generally used by the "Developer" persona to select application dependencies resulting in a manifest to be stored in a Git repo with their source code (i.e. terrarium.yaml). From this, the Terraform template can be generated for their desired platform with all the dependencies.

TODO:
- reference UI for this and how that flow helps to get to terrarium.yaml
- discuss integration with Code Pipes. How one goes from a terrarium.yaml to build/deploy via CP.
- format of terrarium.yaml
- usage of CLI command "terrarium <???>" when in repo with terrarium.yaml


#### Terrarium CLI
```
terrarium:
  farm:
    - resources
    - mappings
    - modules
  platform: # (run in platform dir)
    - lint
    - compile # generate terraform code with given dependencies
    - init,plan,apply # proxy terraform commands on platform
    - implement # add scaffolding code for implementing a dependency in platform
  terraform: # helper to be used by devops (run in terraform dir)
    - module # find module from farm
    - attribute # find attribute from farm
    - suggestion # return suggestions based on user code
    - language-server # protocol to implement IDE independent auto-complete.
  dependency: # helpers to be used by devs. (run in app dir)
    - list # show all available dependencies
    - add-to-app # add a dependency to app yaml
    - lint # show any linting errors in the app yaml
    - show # render app yaml in a user-friendly layout
    - find-platforms # show platforms available in terrarium that satisfies all app dependencies.
```



