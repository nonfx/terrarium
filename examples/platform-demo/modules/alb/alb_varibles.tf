variable "extract_resource_name" {
  type        = string
  description = "The base name to use for all resources created by this module."
}

variable "environment" {
  type        = string
  description = "The environment in which the infrastructure is being deployed (e.g. dev, prod, etc.)."
}

variable "tf_component_lb" {
  type        = any
  description = "A map of objects that define the load balancers to create."
}

variable "tr_component_ecs_services" {
  type        = any
  description = "A map of objects that define the ECS services to create."
}

variable "vpc_id" {
  type        = string
  description = "The ID of the VPC in which to create the load balancer."
}

variable "public_subnet_ids" {
  type        = any
  description = "A list of IDs of the public subnets in which to create the load balancer."
}

variable "security_group_ids" {
  type        = any
  description = "A list of IDs of the security groups to associate with the load balancer."
}
