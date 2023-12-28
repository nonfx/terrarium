# An in-memory data structure store used as a cache or message broker.
# @title: Redis Cache
module "tr_component_redis" {
  source = "./mock-modules/elasticache-redis"

  for_each = local.tr_component_redis

  availability_zones         = local.azs
  vpc_id                     = module.core_vpc.vpc_id
  allowed_security_group_ids = [module.core_vpc.default_security_group_id]
  subnets                    = module.core_vpc.elasticache_subnets
  apply_immediately          = true
  automatic_failover_enabled = false
  engine_version             = each.value.version
}
