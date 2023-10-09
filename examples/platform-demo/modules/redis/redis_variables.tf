variable "common_name_prefix" {
  type        = string
  description = "The common name prefix to use for all resources"
}

variable "environment" {
  type        = string
  description = "The environment to deploy the resources in"
}

variable "extract_resource_name" {
  type        = string
  description = "The name of the resource"
}

variable "tr_component_redis" {
  description = "A map of objects that define the Redis cluster to create."
  type        = any
}

variable "redis_availability_zones" {
  type        = any
  description = "The availability zones to deploy the Redis cluster in"
}

variable "redis_vpc_id" {
  type        = any
  description = "The ID of the VPC to deploy the Redis cluster in"
}


variable "redis_subnet" {
  type        = any
  description = "The subnets to deploy the Redis cluster in"
}

variable "redis_config" {
  description = "A map of objects that define the Redis database to create."
  type        = any
}
