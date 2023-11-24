module "tr_component_s3_bucket_for_logs" {
  source = "terraform-aws-modules/s3-bucket/aws"


  for_each = local.tr_component_s3_bucket_for_logs
  bucket = each.value.bucket
  acl    = each.value.acl

  # Allow deletion of non-empty bucket
  force_destroy = each.value.force_destroy

  control_object_ownership = each.value.control_object_ownership
  object_ownership         = each.value.object_ownership

  attach_elb_log_delivery_policy = each.value.attach_elb_log_delivery_policy  # Required for ALB logs
  attach_lb_log_delivery_policy  = each.value.attach_lb_log_delivery_policy  # Required for ALB/NLB logs
}