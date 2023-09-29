output "alb_names" {
  value = [for k, v in module.alb : v.name]
  description = "A list of the names of the ALBs that were created."
}

output "alb_dns_names" {
  value = [for k, v in module.alb : v.dns_name]
  description = "A list of the DNS names of the ALBs that were created."
}

output "alb_arns" {
  value = [for k, v in module.alb : v.arn]
  description = "A list of the ARNs of the ALBs that were created."
}

output "alb_security_group_ids" {
  value = [for k, v in module.alb : v.security_group_id]
  description = "A list of the security group IDs of the ALBs that were created."
}

output "alb_target_group_arns" {
  value = module.alb.alb_target_group_arns
  description = "A list of the ARNs of the target groups associated with the ALBs."
}

output "alb_log_bucket_names" {
  value = [for k, v in module.s3_bucket : v.bucket]
  description = "A list of the names of the S3 buckets used for ALB access logs."
}
