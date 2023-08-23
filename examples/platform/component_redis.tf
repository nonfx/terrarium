# component[Redis Cache]: An in-memory data structure store used as a cache or message broker.
module "tr_component_redis" {
  source = "cloudposse/elasticache-redis/aws"

  for_each = local.tr_component_redis

  availability_zones         = local.azs
  vpc_id                     = module.core_vpc.default_vpc_id
  allowed_security_group_ids = [module.core_vpc.default_security_group_id]
  subnets                    = module.core_vpc.elasticache_subnets
  apply_immediately          = true
  automatic_failover_enabled = false
  engine_version             = each.value.version

  context = module.this.context
}

module "this" {
  source = "cloudposse/label/null"
  # Cloud Posse recommends pinning every module to a specific version
  # version = "x.x.x"
  namespace = "var.namespace"
  stage     = "var.stage"
  name      = "var.name"
}
