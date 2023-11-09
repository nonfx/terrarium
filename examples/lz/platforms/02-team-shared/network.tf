# create a subnet for each environment
module "subnet_dev" {
  source = "../../modules/subnet"

  # create a subnet for each env
  for_each = local.tr_component_environment

  env_type = each.key
  network = data.terraform_remote_state.common.vpc_network_name
}

output "subnet" {
  value = { for k, v in module.subnet_dev : k => v.name }
}
