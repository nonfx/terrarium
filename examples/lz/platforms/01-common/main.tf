# create a common network for all requirements
module "vpc" {
  source = "../../modules/vpc"

  # this demonstrates how we can aggregate all requirements across all envs and provision tf resources accordingly
  name = format("network_for_%s", local.count_everything)
}

locals {
  count_everything = sum([length(local.tr_component_environment), length(local.tr_component_postgres), length(local.tr_component_bucket)])
}

output "vpc_network_name" {
  value = module.vpc.vpc_name
}
