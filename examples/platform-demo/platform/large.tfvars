lb_config = {
  "default" = {
    "load_balancer_type" = "application",
    "create_alb"         = true
  }
}

postres_config     = {}
common_name_prefix = "prod"
ecs_config = {
  "default" : {
    name : "vote",
    engine : {
      type : "FARGATE",
      default_weight : 100,
      spot_weight : 0,
    }
  }
}
environment = "testing"
redis_config = {
  "prod" : {
    "cluster_size" : 3,
    "cluster_mode_replicas_per_node_group" : 2,
    "cluster_mode_num_node_groups" : 2,
    "instance_type" : "cache.m5.large",
    "family" : "redis6.x",
    "at_rest_encryption_enabled" : true,
    "transit_encryption_enabled" : true,
    "port" : 6379,
    "number_of_azs" : 2
  }
}
domain_name = "platform.test.codepipes.io"
zone_id     = "Z04638322TOMMCGJDXCAO"
