## Terrarium Components

#### The Terrarium Farm
The Terrarium Farm is the library shared across all the components. The library is composed of Git repos and a Postgres DB. It can reside locally on a single machine or be shared using a remote Postgres DB. The components discussed below interact with the Farm to accomplish particular functions.

TODO:
- discuss how each artifact is put into the farm
- how does the taxonomy get established?


#### Curation and Publishing

These functions are generally the purview of the "Devops" persona. The goal is to create and maintain the various objects that are used downstream other personas (such as "Developer"). This is accomplished either by directly adding objects to the Terrarium Farm or by executing a process that results in objects being added to the Farm. The first objects supported are based on Terraform - i.e. infrastructure objects.

There are 3 types of objects related to infrastructure:
1. Terraform modules - These can be imported via the "harvest" command and then used to build up the other objects or new Terraform templates using the Terriaum code generator.
2. Platform (? is this the correct term) - A base infrastructure pattern created using modules (or otherwise?) that is imported via ???. A platform is used base layer that the dependencies used in applications plugs into. The platform code (Terraform) lives in the Farm repo and is imported into Terrarium using ???
3. Dependencies - Typically platform independent objects that are resolved to concrete code at generation time. For example, a Postgres DB would be a dependency that when plugged into a platform defined for AWS might resolve using an RDS module. Dependencies are specified using ????...

TODO:
- discuss how taxonomy is assigned and/or used for each type
- walk through how each type is created and added to Farm
- usage of Farm viewer UI?
- usage of IDE plugin?

#### Composer

This is generally used by the "Developer" persona to select application dependencies resulting in a manifest to be stored in a Git repo with their source code (i.e. terrarium.yaml). From this, the Terraform template can be generated for their desired platform with all the dependencies.

TODO:
- reference UI for this and how that flow helps to get to terrarium.yaml
- discuss integration with Code Pipes. How one goes from a terrarium.yaml to build/deploy via CP.
- format of terrarium.yaml
- usage of CLI command "terrarium <???>" when in repo with terrarium.yaml








