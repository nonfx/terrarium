module "tr_component_s3_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"


  for_each = local.tr_component_s3_bucket
  bucket = each.value.bucket
  acl    = var.acl

  control_object_ownership = var.control_object_ownership
  object_ownership         = var.object_ownership

  versioning = var.versioning
  }