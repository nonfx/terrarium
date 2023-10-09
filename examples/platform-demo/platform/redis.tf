# [Redis]:This is a Terraform module that creates a Redis cluster in AWS.
module "tr_component_redis" {
  source = "../modules/redis"
  for_each = local.tr_component_redis
  tr_component_redis = local.tr_component_redis
  redis_config = var.redis_config

  #Common Variables
  extract_resource_name = local.extract_resource_name
  environment           = var.environment
  common_name_prefix    = var.common_name_prefix

  #Variables from VPC
  redis_availability_zones = module.vpc.availability_zones
  redis_vpc_id             = module.vpc.vpc_id
  redis_subnet             = module.vpc.private_subnet_ids

}

