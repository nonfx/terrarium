module "tr_component_sqs" {
  source  = "terraform-aws-modules/sqs/aws"
  for_each = local.tr_component_sqs
  name = each.value.name

  kms_master_key_id                 = var.kms_master_key_id
  kms_data_key_reuse_period_seconds = var.kms_data_key_reuse_period_seconds

  tags = var.tags
}