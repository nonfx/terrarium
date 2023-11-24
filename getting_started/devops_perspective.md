## Getting started with Devops perspective

> For Devops we have terraform modules already created. You can check them in below path

```
examples/platform-demo/platform
```

- As a Devops we know that we are eager to know how terrarium is getting synced with terraform 

- So here are the details 

    - Source code repo for terrarium

        ```
        ```

    - Code is written on Golang




### How to write a basic platform ??

# Let say we have a requirement for s3 bucket which store logs for ALB.

# Steps

Step 1: Pick up a module for anywhere for eg: i have picked it up from terraform registry
> source : https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest

Step 2: Create a main.tf file and paste it like below

```
module "s3_bucket_for_logs" {
  source = "terraform-aws-modules/s3-bucket/aws"

  bucket = "my-s3-bucket-for-logs"
  acl    = "log-delivery-write"

  # Allow deletion of non-empty bucket
  force_destroy = true

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  attach_elb_log_delivery_policy = true  # Required for ALB logs
  attach_lb_log_delivery_policy  = true  # Required for ALB/NLB logs
}
```

Step 3: Now change it according to your platform requirement. Suppose you want to take bucket name from developer and you want all the remaining values to be same so your main.tf will be like below.

```
locals{
    tr_component_s3_bucket_for_logs = {
        "defaults" : {
         "bucket" : "terrarium-test-bucket"
        }
    }
}

module "tr_component_s3_bucket_for_logs" {
  source = "terraform-aws-modules/s3-bucket/aws"

  for_each = local.tr_component_s3_bucket_for_logs
  bucket = each.value.bucket
  acl    = "log-delivery-write"

  # Allow deletion of non-empty bucket
  force_destroy = true

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  attach_elb_log_delivery_policy = true  # Required for ALB logs
  attach_lb_log_delivery_policy  = true  # Required for ALB/NLB logs
}

output "tr_component_s3_bucket_id" {
    value = {for k, v in module.tr_component_s3_bucket_for_logs : k => v.s3_bucket_id }
}
```
> [!IMPORTANT]
> For arguments you see in the above code  tr_components_ works with modules , locals and outputs for now so if you want to take arguments from developer just define it in locals  as shown  above 

### Generate terrarium.yaml file from terraform code

```
terrarium platform lint
```

### Generated terrarium yaml file 

```
components:
    - id: s3_bucket_for_logs
      title: S3 Bucket For Logs
      inputs: {}
      outputs:
        type: object
graph:
    - id: local.tr_component_s3_bucket_for_logs
      requirements: []
    - id: module.tr_component_s3_bucket_for_logs
      requirements:
        - local.tr_component_s3_bucket_for_logs
    - id: output.tr_component_s3_bucket_arn
      requirements:
        - module.tr_component_s3_bucket_for_logs
    - id: output.tr_component_s3_bucket_id
      requirements:
        - module.tr_component_s3_bucket_for_logs

```