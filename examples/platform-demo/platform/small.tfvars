lb_config = {
  "default" = {
    "load_balancer_type" = "application",
    "create_alb"         = true
  }
}

postres_config     = {}
common_name_prefix = "testing"
ecs_config = {
  "default" : {
    name : "default",
    engine : {
      type : "FARGATE",
      default_weight : 50,
      spot_weight : 50,
    }
  }
}
environment = "testing"
redis_config = {
  "default" : {
    "cluster_size" : 1,
    "cluster_mode_replicas_per_node_group" : 1,
    "cluster_mode_num_node_groups" : 1,
    "instance_type" : "cache.t2.micro",
    "family" : "redis5.0",
    "at_rest_encryption_enabled" : true,
    "transit_encryption_enabled" : true,
    "port" : null,
    "number_of_azs" = 1
  }
}
domain_name = "platform.test.codepipes.io"
zone_id     = "placeholder"
