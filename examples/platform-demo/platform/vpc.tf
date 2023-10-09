module "vpc" {
  source = "../modules/vpc"

  tr_component_vpc = var.vpc_config

  #Common Variables
  environment           = var.environment
  common_name_prefix    = var.common_name_prefix
  extract_resource_name = local.extract_resource_name

}
