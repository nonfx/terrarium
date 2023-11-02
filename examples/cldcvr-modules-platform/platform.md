# Terraform Code Conversion to Terrarium Platform-Based Infrastructure as Code (IaC)

## Overview

This guide is designed to assist you in converting existing Terraform code into platform-based Infrastructure as Code (IaC) suitable for use with Terrarium. If you have queries related to Platform and Platform Engineering, please refer to our FAQ.md for detailed answers. For insights on common functions utilized during this conversion process, consult Terraform.md for comprehensive explanations and examples.

## Pre-Conversion Knowledge Requirements

The following prerequisites are expected of the engineer:

1. Proficiency in Terraform, including its modules, functions, and their integrations. Prior experience in deploying infrastructure using IaC is crucial.
2. Possession of either pre-existing modular or non-modular Terraform IaC, with the intention to convert it to a platform-based format.
3. Access to a cloud account for testing purposes, although it is not a prerequisite for writing the code.

## Preparation Steps

Ensure the following tools and resources are available and properly configured before starting the conversion process:

1. Terraform
2. Terrarium (Refer to our installation guide for more details.)
3. A code editor of your choice (The guide uses VSCode as an example, but feel free to use any editor you are comfortable with.)

## Scope of Learning

This guide does not cover the basics of writing Terraform code. For beginners, we recommend exploring our internal knowledge portal or the numerous online courses available. The focus here is to elevate your IaC skills through platform-based code conversion using Terrarium.

## Getting Started with Conversion

Now that you are acquainted with Platform and Platform Engineering concepts, and have your Terraform-based IaC ready, let's proceed to convert it into Platform-Based Code. This progression will enhance your IaC skills and prepare you for the next phase in DevOps evolution.

First, examine a sample of Terraform code that hasn’t undergone platform-based modification, located at terrarium/examples/cldcvr-modules-platform/tf-without-platform/main.tf. A detailed breakdown of the resources created by this code can be found in terrarium/examples/cldcvr-modules-platform/tf-without-platform/terraform.md.

Let's delve into each resource and observe the transformations applied:

### 1. VPC (Virtual Private Cloud)

The initial module is the VPC. The responsibility lies with the Platform Engineer to decide the number of VPCs permitted on a Platform. This decision should align with the organization's requirements, whether it involves multiple VPCs assigned to specific resources or a single VPC encompassing all resources. For simplicity, we will proceed with a single VPC, requiring no modifications.

**Original Code:**

```hcl
module "vpc" {
  source = "github.com/cldcvr/cldcvr-xa-terraform-aws-vpc?ref=v0.1.0"
  ...
}
```

### 2. EC2 Instance (Service_web/Compute)

In our pursuit to establish a platform with EC2 as the compute model, serving as our web server, notable differences emerge in the code before and after platform adoption.

**Pre-Platform Code:**

```hcl
module "ec2_instance" {
  source = "terraform-aws-modules/ec2-instance/aws"
  name   = "single-instance"
  ami    = var.ami_id

  instance_type          = var.instance_type
  key_name               = var.key_name
  monitoring             = true
  vpc_security_group_ids = [module.ec2_sg.id]
  subnet_id              = element(module.vpc.public_subnets, 0)

  tags = local.Tags
}
```

**Post-Platform Code:**

```hcl
module "tr_component_service_web" {
  source = "terraform-aws-modules/ec2-instance/aws"

  for_each = {
    for k, v in local.tr_component_service_web : k => merge(v, var.ec2_config["default"])
  }

  name                        = "${local.name}-ec2-instance"
  ami                         = each.value.ami
  instance_type               = each.value.instance_type
  monitoring                  = each.value.monitoring_enabled
  create_iam_instance_profile = true
  associate_public_ip_address = length(local.tr_component_service_web) > 0 ? true : false
  subnet_id                   = module.vpc.public_subnets[0]
  vpc_security_group_ids      = [module.ec2_sg[each.key].id]
  tags                        = local.Tags
}
```

**Key Changes:**

1. **Module Name:** Transitioned from `ec2_instance` to `tr_component_service_web`, aligning with Terrarium conventions. This naming convention indicates that the compute object is public-facing (denoted by `service_web`). All Terrarium platform components should begin with `tr_component` as a prefix in their module name. If you want the component to be private-facing, use `tr_component_service_private` as the prefix.

2. **Leveraging `for_each` for Multiple Compute Instances** In scenarios where multiple compute services each require a unique EC2 instance, the `for_each` construct becomes invaluable. It enables the creation of multiple module instances, each configured according to specific requirements. In our context, this translates to generating distinct EC2 instances, each governed by its unique configuration settings.
The `for_each` loop is combined with a nested loop that merges default settings with service-specific configurations. This amalgamation of configurations is sourced from both local variables and explicit variable declarations. Local variables typically contain configurations specified by developers, while the `variable` declarations house configurations added by DevOps engineers.

    **Code Implementation:**
    ```hcl
    for_each = {
      for k, v in local.tr_component_service_web : k => merge(v, var.ec2_config["default"])
    }
    ```

    **Local Variable Storage (tr_gen_locals.tf):**
    ```hcl
    tr_component_service_web = {
      "default" : {
        port : 80
      }
    }
    ```

3. **Variable Declaration (variables.tf):**
    ```hcl
    variable "ec2_config" {
      description = "The configuration for the EC2 instance"
      default = {
        "default" : {
          ami : "ami-0e18308d78c527c8a"
          instance_type : "t2.micro"
          monitoring_enabled : true
        }
      }
    }
    ```

    In the provided code snippets, the port configuration is exclusively specified in the local variable, signifying that this setting is to be determined by the developer. On the other hand, configurations such as `ami`, `instance_type`, and monitoring preferences are specified by the Platform engineers in the variable declaration. The `merge` function is then employed to seamlessly integrate these configurations.

    **Understanding Variable Values:**

    Upon merging `local.tr_component_service_web` with `var.ec2_config`, the resulting modified variable is structured as follows:

    ```hcl
    modified_variable = {
      "default" : {
        ami : "ami-0e18308d78c527c8a"
        instance_type : "t2.micro"
        monitoring_enabled : true
        port: 80
      }
    }
    ```

    Here, `each.key` is set to "default", and `each.value` references the associated configuration objects. As a result, configurations such as the port can be accessed via `each.value.port`.

4. **Decision-Making with Manipulations:**

    The `associate_public_ip_address` setting is determined based on the length of `local.tr_component_service_web`. If it is set, the instances are made public; otherwise, they remain private. We can make many such simple manipulations in order to reduce the number of options available to developers, thereby simplifying the process of creating infrastructure.

    **Code Example:**
    ```hcl
    associate_public_ip_address = length(local.tr_component_service_web) > 0 ? true : false
    ```

    We make similiar modifications in other modules as well, to convert them into platform code.

## Local Files and Variable Handling in Terrarium

When converting Terraform code to work with the Terrarium platform, handling and structuring variables is a crucial aspect. This process entails segregating configurations based on their origin: options intended for the developer and options meant for the platform engineer.

### Developer Options: `tr_gen_locals.tf`

The `tr_gen_locals.tf` file is where configurations expected from the developer are specified. Here, default values and specific settings for various Terraform modules are laid out. An example of this can be seen in the configuration for the `tr_component_postgres` module:

```hcl
locals {
  tr_component_postgres = {
    "default" : {
      "engine_version" : "11.20",
      "family" : "postgres"
    }
  }
  ...
}
```

In the example above, the developer is expected to provide values for the `engine_version` and `family` attributes of the `tr_component_postgres` module. It is crucial to maintain the map of the object as `default` to ensure seamless integration and value replacement by Terrarium.

### Platform Engineer Options: Direct Code Insertion and `variables.tf`

Platform engineers have the option to hard-code values directly into the Terraform modules or fetch them from the `variables.tf` file. While the former approach offers simplicity, it deviates from best practices as it reduces code generality and reusability. Nonetheless, it is a viable option, as demonstrated in the `tr_component_postgres` module:

```hcl
module "tr_component_postgres" {
  source = "github.com/cldcvr/cldcvr-xa-terraform-aws-db-instance?ref=v0.1.0"
  ...
  engine_version = try(each.value.engine_version, "14")
  ...
}
```

In the code snippet above, the `engine_version` is set directly within the module. Although this method is straightforward, leveraging the `variables.tf` file is recommended for a cleaner and more maintainable codebase.

### Ensuring Generic and Maintainable Code

To uphold code quality and maintainability, it is imperative to abstract configurations away from the modules and into appropriate variables or local files. This practice not only enhances code readability but also simplifies future modifications and updates.

By meticulously organizing configurations and adhering to best practices, we pave the way for a robust and flexible IaC setup, fully leveraging the capabilities of the Terrarium platform.
