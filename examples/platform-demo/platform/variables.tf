variable "vpc_config" {
  description = "A map of objects that define the VPC to create."
  type        = any
  default = {
    "default" = {
      "vpc_cidr_block" = "10.0.0.0/16",
      "number_of_azs"  = 2
    }
  }

}

variable "lb_config" {
  description = "A map of objects that define the LB to create."
  type        = any
  default = {
    "default" = {
      "load_balancer_type" = "application",
      "create_alb"         = true
    }
  }
}

variable "common_name_prefix" {
  description = "The common name prefix to use for all resources."
  type        = string
  default     = "terrarium"
}

variable "environment" {
  description = "The environment to deploy to."
  type        = string
  default     = "testing"
}

variable "ecs_config" {
  description = "A map of objects that define the ECS cluster to create."
  type        = any
  default = {
    "default" : {
      name : "default",
      engine : {
        type : "FARGATE",
        default_weight : 50,
        spot_weight : 50,
      }
    }
  }
}

variable "postres_config" {
  description = "A map of objects that define the Postgres database to create."
  type        = any
  default     = {}
}

variable "redis_config" {
  description = "A map of objects that define the Redis database to create."
  type        = any
  default = {
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
}

variable "domain_name" {
  description = "The DNS zone domain name to create records in."
  type        = string
  default     = "platform.test.codepipes.io"
}

variable "zone_id" {
  description = "The DNS zone ID to create records in."
  type        = string
  default     = "placeholder"
}
