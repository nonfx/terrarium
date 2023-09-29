module "this" {
  source  = "cloudposse/label/null"
  # Cloud Posse recommends pinning every module to a specific version
  # version = "x.x.x"
  namespace  = var.common_name_prefix
  stage      = var.environment
  name       = "${keys(var.tr_component_redis)[0]}-redis"
}

module "redis" {
  source = "cloudposse/elasticache-redis/aws"
  # Always pin every module to a specific version
  version = "0.52.0"

  for_each = var.tr_component_redis

  description                = "Redis cluster for ${keys(var.tr_component_redis)[0]}"


  availability_zones         = length(var.redis_availability_zones["default"]) >= coalesce(var.redis_config.default.number_of_azs, 2) ? slice(var.redis_availability_zones["default"], 0, coalesce(var.redis_config.default.number_of_azs, 2)) : var.redis_availability_zones["default"]
  zone_id                    = can(var.redis_config.default.redis_zone_id) ? var.redis_config.default.redis_zone_id : ""
  vpc_id                     = var.redis_vpc_id["default"]
  subnets                    = length(var.redis_availability_zones["default"]) == 2 ? slice(var.redis_subnet["default"], 0, 2) : []

  port                       = coalesce(can(var.redis_config.default.port) ? var.redis_config.default.port : null, 6379)

  cluster_size               = var.redis_config.default.cluster_size
  instance_type              =  var.redis_config.default.instance_type
  apply_immediately          = var.environment == "production" ? false : true
  automatic_failover_enabled = var.environment == "production" ? true : false
  cluster_mode_enabled       = var.environment == "production" ? true : false
  engine_version             = each.value.version
  family                     = var.redis_config.default.family
  create_security_group      = true
  at_rest_encryption_enabled = var.environment == "production" ? true : var.redis_config.default.at_rest_encryption_enabled
  transit_encryption_enabled = var.environment == "production" ? true : var.redis_config.default.transit_encryption_enabled
  cluster_mode_replicas_per_node_group = var.environment == "production" ? var.redis_config.default.cluster_mode_replicas_per_node_group : null
  cluster_mode_num_node_groups         = var.environment == "production" ? var.redis_config.default.cluster_mode_num_node_groups : null
  # parameter = [
  #   {
  #     name  = "notify-keyspace-events"
  #     value = "lK"
  #   }
  # ]

  context = module.this.context
}

