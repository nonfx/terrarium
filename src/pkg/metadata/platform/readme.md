# Terrarium Platform Metadata YAML Format Documentation

The Terrarium Platform Metadata YAML format is used to declare information about the components and their relationships within the Terrarium platform template. It helps organize and describe various components and their dependencies, allowing Terrarium tools to understand and validate the platform's structure.

## Components

The `components` section in the metadata defines the different dependency interfaces that are implemented in the Terrarium platform template. Each component is represented as a YAML object with the following properties:

- `id` (string): A unique identifier for the component. It helps in referencing the component in other parts of the metadata or code.
- `taxonomy` (array of strings): Represents the taxonomy or categories to which the component belongs. This helps in organizing components based on their functionalities or roles.
- `title` (string): A descriptive title for the component, providing a brief overview of its purpose.
- `description` (string): A detailed description of the component's functionality and its significance within the platform.
- `inputs` (JSON Schema): Defines the input parameters required by the component. It follows the JSON Schema format to specify the input properties, their data types, titles, and descriptions.
- `outputs` (JSON Schema): Defines the output properties produced by the component. It also follows the JSON Schema format to specify the output properties, their data types, titles, and descriptions.

## Graph

The `graph` section in the metadata defines the relationships between different terraform blocks in the Terrarium platform template. It represents a directed acyclic graph (DAG) of component dependencies. Each graph node is represented as a YAML object with the following properties:

- `id` (string): A unique identifier for the graph node, typically corresponding to the component ID or resource name.
- `requirements` (array of strings): Specifies the IDs of other graph nodes that the current node depends on. This indicates the dependencies between components and their order of execution during Terraform generation.

## Example

Below is an example of a Terrarium Platform Metadata YAML file:

```yaml
components:
  - id: postgres
    taxonomy: [storage, database, rdbms]
    title: PostgreSQL Database
    description: A relational database management system using SQL.
    inputs:
      properties: # JSON schema format
        version:
          type: number
          title: Engine Version
          description: The version of PostgreSQL to use.
    outputs:
      properties:
        host:
          type: string
          title: Host
          description: The host address of the PostgreSQL server.
        port:
          type: number
          title: Port
          description: The port number on which the PostgreSQL server is listening.
        username:
          type: string
          title: Username
          description: The username for accessing the PostgreSQL database.
        password:
          type: string
          title: Password
          description: The password for accessing the PostgreSQL database.

graph:
    - id: module.tr_component_postgres
      requirements:
        - module.postgres_security_group
        - module.core_vpc
    - id: module.postgres_security_group
      requirements:
        - module.core_vpc
    - id: module.core_vpc
      requirements:
        - resource.random_string.random
    - id: resource.random_string.random
      requirements: []
    # (Other graph nodes are defined here as well)
```

In this example, the metadata defines a "postgres" component with its taxonomy, title, description, inputs, and outputs. Additionally, the graph section establishes the relationships between different terraform blocks using their IDs and their corresponding requirements.

By using the Terrarium Platform Metadata YAML format, DevOps professionals can create well-structured and organized Terrarium platforms with clear component definitions and their dependencies, facilitating better collaboration and understanding among team members.
