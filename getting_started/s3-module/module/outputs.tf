output "tr_component_s3_bucket_id" {
    value = {for k, v in module.tr_component_s3_bucket : k => v.s3_bucket_id }
}
output "tr_component_s3_bucket_arn" {
    value = {for k, v in module.tr_component_s3_bucket : k => v.s3_bucket_arn }
}