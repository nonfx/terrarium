module "tr_component_bucket" {
  source = "../../modules/gcs"

  for_each = local.tr_component_bucket

  name = each.value.name
}

output "tr_component_bucket_uri" {
  value = { for k, v in module.tr_component_bucket : k => v.uri}
}
