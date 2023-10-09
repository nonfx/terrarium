# Terraform Code Comparison: Application vs Platform Code

This README.md aims to highlight the key differences between two sets of Terraform code used for setting up infrastructure for different purposes: application-specific and platform code. By the end of this document, we will understand the differences and the rationale behind the decisions made for the platform code.

## Application Code Overview
The application-specific code is relatively simple and straightforward. It focuses on setting up a few essential AWS services like EC2, S3, and RDS. However it focuses on a specific usecase and does not offer flexibility.

Example:

```hcl
module "ec2_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  name    = "single-instance"
  ...
}
```

## Platform Code Overview
The platform code is more complex and feature-rich, offering a set of functionalities catered towards a complete platform setup. This includes services like VPC, multiple EC2 instances, RDS, and Elasticache for Redis.

Example:

```hcl
module "tr_component_service_web" {
  source = "terraform-aws-modules/ec2-instance/aws"

  for_each = merge(local.tr_component_service_web, var.ec2_config)

  ...
}
```

## Key Differences

### 1. Modularity

**Application Code**
- Less modular; resources are often directly declared in the `main.tf`.

**Platform Code**
- Highly modular; uses Terraform modules to manage VPC, EC2 instances, RDS, etc.

### 2. Flexibility & Complexity

**Application Code**
- Less flexible and less complex; generally hardcoded or has fewer variables.

**Platform Code**
- Highly flexible; uses Terraform's `for_each`, `locals`, and other advanced features for dynamic configuration.

Example:

```hcl
# Platform Code
for_each = merge(local.tr_component_service_web, var.ec2_config)
```

### 3. Security Group Configuration

**Application Code**
- Limited security group rules.

**Platform Code**
- Detailed security group rules for ingress and egress traffic.

Example:

```hcl
# Platform Code
security_group = {
  ingress = [
    {
      type        = "ingress"
      cidr_blocks = ["0.0.0.0/0"]
      ...
    }
  ],
  ...
}
```

### 4. Tags and Naming Conventions

**Application Code**
- Basic or no tagging.

**Platform Code**
- Detailed tagging and custom naming conventions used for easy identification and management.

Example:

```hcl
# Platform Code
tags = local.tags
```

### 5. Data Handling

**Application Code**
- Minimal usage of data sources and locals.

**Platform Code**
- Extensive usage of `locals` and data sources like `aws_availability_zones` for dynamic data handling.

Example:

```hcl
# Platform Code
locals {
  name   = "${var.common_name_prefix}-${var.environment}-demo"
  region = data.aws_region.current.name
}
```

## Why These Decisions?

- **Modularity** allows for easier management and scalability.
- **Flexibility** enables the platform to adapt to varying requirements without code changes.
- **Advanced security configurations** ensure a more secure and compliant environment.
- **Tags and custom naming conventions** assist in resource tracking and billing.
- **Dynamic data handling** allows the platform to be agnostic to environment specifics, thereby being more reusable.

Feel free to dive into each Terraform file to understand the specifics. The comments within the code offer more insights into each block's functionality.
