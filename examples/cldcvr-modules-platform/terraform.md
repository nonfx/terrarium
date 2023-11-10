# Terraform Platform-Based Code: Essential Functions and Techniques

This document delineates key Terraform functions and techniques pivotal in writing platform-based code, providing clear definitions, practical examples, and guidelines for optimal usage.

## 1. `try` Function

### Definition
The `try` function in Terraform is used to handle errors in expressions. It evaluates all its argument expressions sequentially and returns the result of the first expression that does not produce any errors.

### Example
Consider the `tr_component_postgres` module in a Terraform platform code, where we assign a value to the `username` variable as follows:

```hcl
username = try(each.value.username, "postgres")
```

In this example, Terraform attempts to evaluate `each.value.username`. The value for this expression is derived from a `for_each` construct, which we will explore later. If a value is successfully retrieved, it is assigned to the `username` variable. If, however, an error occurs during the evaluation (e.g., if the value is not found), the `try` function ensures a fallback value of "postgres" is used instead.

### When and Where to Use
This function is particularly useful when dealing with modules that might not have predefined default values. As a platform engineer, ensuring that sensible default values are in place is crucial. This approach minimizes the amount of information required from the end-user while also providing flexibility for them to override these defaults if necessary.

## 2. `for_each` Meta-Argument

### Definition
The `for_each` meta-argument in Terraform is utilized to iterate over a map or set, creating an instance of a module or resource for each item. Each created instance has a unique infrastructure object associated with it.

### Example
In the `tr_component_postgres` module of a Terraform platform code, `for_each` is used to create multiple instances based on a local configuration:

```hcl
module "tr_component_postgres" {
  source  = "github.com/cldcvr/cldcvr-xa-terraform-aws-db-instance?ref=v0.1.0"
  for_each = local.database_configuration
  ...
}
```

Here, `local.database_configuration` is a map with configuration for each database instance to be created:

```hcl
database_configuration = {
  "default" : {
    "engine_version" : "11.20",
    "family"         : "postgres"
  },
  "activity" : {
    "engine_version" : "8.0.36",
    "family"         : "mysql"
  }
}
```

Using `for_each`, Terraform will create two distinct instances of the `tr_component_postgres` module: one for "default" with PostgreSQL 11.20, and another for "activity" with MySQL 8.0.36, assuming the module supports both configurations.

### When and Where to Use
The `for_each` meta-argument is highly recommended for use in a variety of components. It introduces a level of flexibility and dynamism, enabling the concise and efficient creation of multiple resource instances based on input configurations.

## 3. `replace` Function

### Definition
The `replace` function in Terraform is used to perform string substitution. It searches a given string for another substring and replaces it with a third string. The signature of the function is `replace(string, search, replace)`.

### Example
In our platform code, we use the `replace` function in combination with the `try` function to construct resource names while ensuring special characters are removed. Here's how it is applied in the `tr_component_postgres` module:

```hcl
name = replace(try(each.value.name, "${local.name}-${each.key}-db-instance"), "-", "")
```

In this example:
- `try(each.value.name, "${local.name}-${each.key}-db-instance")`: First, we attempt to retrieve the value of `each.value.name`. If it's not available, we fall back to a default naming convention which combines a local variable `local.name`, the key of the current item in iteration `each.key`, and the suffix "-db-instance".
- `replace( ... , "-", "")`: We then pass the resulting string to the `replace` function. Here, we aim to remove any hyphens "-" from the string. This is particularly useful to adhere to naming conventions that might not allow special characters.

### When and Where to Use
The `replace` function proves invaluable when enforcing specific string formats or removing/replacing certain characters in a string. This is particularly crucial in resource naming to ensure consistency and compliance with naming conventions. Some resources may have strict naming constraints, disallowing special characters such as hyphens. Using the `replace` function helps perform a name sanity check, ensuring that the names passed to resources meet the required criteria and prevent potential deployment issues.

In combination with the `try` function, as demonstrated in the example, this approach enhances the robustness of our Terraform configurations, ensuring default values are provided and name formats are standardized. This results in a more predictable and error-resistant provisioning process.

## 4. `for` Expression and `cidrsubnet` Function

### Definition

#### `for` Expression
The `for` expression in Terraform is used to transform lists, sets, or maps based on certain conditions or transformations. It can be used to iterate through each element in a collection and modify it or filter it based on specific criteria.

#### `cidrsubnet` Function
The `cidrsubnet` function is used to calculate a subnet address within a given IP network address prefix. It takes three arguments: the prefix, the number of additional bits to add to the prefix, and the subnet number to calculate.

### Example

In our Terraform configurations, we have the following example to create a list of subnet CIDR blocks:

```hcl
private_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k + coalesce(length(local.azs), 2))]
```

In this example:
- `local.azs` is assumed to be a map or a list of availability zones.
- `local.vpc_cidr` is the CIDR block for the VPC.
- For each availability zone in `local.azs`, we calculate a subnet using the `cidrsubnet` function.
- `cidrsubnet(local.vpc_cidr, 4, k + coalesce(length(local.azs), 2))`: We take the VPC CIDR block, add 4 additional bits to the prefix, and calculate the subnet number based on the index of the availability zone plus the greater of the number of availability zones or 2. This ensures unique subnet calculations for each availability zone.


### Explanation for advance for like `for k, v in local.tr_component_service_web : k => merge(v, var.ec2_config["default"])`

This expression is a `for` loop and it operates over a map. Here’s a breakdown:

- `local.tr_component_service_web`: This refers to a local variable that is expected to be a map. Each element in the map has a key (`k`) and a value (`v`).

- `for k, v in local.tr_component_service_web`: This part of the expression is iterating over each key-value pair in the `local.tr_component_service_web` map. `k` represents the key, and `v` represents the value of each pair.

- `k => merge(v, var.ec2_config["default"])`: For each iteration, this expression constructs a new key-value pair. The key is `k`, which remains unchanged. The value is the result of the `merge(v, var.ec2_config["default"])` function.

  - `merge(v, var.ec2_config["default"])`: This `merge` function call is merging two maps. The first map is `v`, which is the value from the current iteration over the `local.tr_component_service_web` map. The second map is `var.ec2_config["default"]`, which is a map defined elsewhere in your Terraform code, presumably under a variable named `ec2_config` with a key of "default".

  The `merge` function takes these two maps and combines them into a single map. If there are any overlapping keys between `v` and `var.ec2_config["default"]`, the values from `v` will overwrite the values from `var.ec2_config["default"]`.

The result of this entire `for` expression is a new map. Each key in this new map corresponds to a key in the original `local.tr_component_service_web` map. Each value is a merged map, combining specific configurations from the original value with default configurations from `var.ec2_config["default"]`.

This construct is particularly useful when you have a set of configurations (in this case, service configurations) and you want to ensure that each configuration has a set of default values, but also allow for specific overrides or additions on a per-service basis.

### When and Where to Use
The combination of `for` and `cidrsubnet` is particularly useful when you need to generate a list of subnet CIDR blocks dynamically based on the number of availability zones or another list/map of resources.

The `for` expression provides a concise way to iterate through a collection and apply a transformation or calculation to each element.

The `cidrsubnet` function is crucial when dealing with network configurations in Terraform, as it allows for precise subnet calculations within a given IP address range.

By using these functions together, you can create scalable and dynamic Terraform configurations that adapt based on the input variables or existing resources, leading to more maintainable and flexible infrastructure as code.

## 5. Working with Maps and the `merge` Function

### Definition
The `merge` function in Terraform is used to combine two or more maps into a single map. If there are any common keys between the maps, the values from the rightmost map in the `merge` function arguments will overwrite those in the leftmost.

### Example and Explanation
Consider the following example in a Terraform module where we iterate over a map and use the `merge` function:

```hcl
locals {
  tr_component_service_web = {
    "serviceA" = { memory = "512Mi", cpu = "250m" },
    "serviceB" = { memory = "1Gi", cpu = "500m" },
    // other services…
  }
}

resource "kubernetes_deployment" "this" {
  for_each = local.tr_component_service_web

  metadata {
    name = each.key
  }

  spec {
    template {
      spec {
        container {
          resources {
            limits = merge(each.value, var.ec2_config["default"])
          }
        }
      }
    }
  }
}
```

In the `kubernetes_deployment` resource block, we use the `for_each` meta-argument to iterate over the `local.tr_component_service_web` map. For each key-value pair in the map (referred to as `k` and `v` in the pseudo-code in the question), a `kubernetes_deployment` resource is created. The name of the deployment is set to the key of the current map item (`each.key`), and the resource limits are set using the `merge` function.

The `merge` function is used here to combine `each.value` (which contains resource limits like memory and CPU for each service) and a default EC2 configuration (`var.ec2_config["default"]`). If `each.value` and `var.ec2_config["default"]` have any common keys, the values from `each.value` will overwrite those from `var.ec2_config["default"]`, ensuring that service-specific configurations take precedence over default configurations.

### When and Where to Use
The `merge` function is particularly useful when you have a set of default configurations that you want to apply to a number of resources, but you also need the ability to override these defaults for specific instances. In this example, it ensures that each Kubernetes deployment has the necessary resource limits, while also allowing for service-specific overrides.

## 6. `anytrue` Function

### Definition
The `anytrue` function in Terraform returns true if any of the given booleans are true.

### Example

```hcl
result = anytrue([false, false, true, false])
```

In this example, `result` would be `true` because there is at least one `true` value in the list.

### When and Where to Use
This function is useful when you need to perform a operation where you want to check if any one of the options is true to set something. For example, in the `tr_component_postgres` module, we use the `anytrue` function to determine whether or not to create a read replica:

```hcl
read_replica = anytrue([for k, v in local.database_configuration : v.read_replica])
```

## 7. `coalesce` Function

### Definition
The `coalesce` function in Terraform returns the first non-null value in a list of arguments.

### Example

```hcl
result = coalesce(null, "default", "specific")
```

In this example, `result` would be `"default"` because it is the first non-null value in the list of arguments.

### When and Where to Use
This function is particularly useful when you have a list of potential values and you want to select the first one that is not null. It can help provide default values in cases where a specific value might not be available. This is demonstrated in the `for` and `cidrsubnet` example, where the `coalesce` function ensures that there is always a number to add to the index `k` in subnet calculations, preventing potential null value errors.

