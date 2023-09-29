variable "ecs_config" {
  type = any
}

variable "tr_component_ecs_services" {
  type = any
}

variable "environment" {
  type    = string
  default = "dev"
}

variable "vpc_id" {
  type        = any
  description = "The ID of the VPC in which to create the load balancer."
}

variable "public_subnet_ids" {
  type        = any
  description = "A list of IDs of the public subnets in which to create the load balancer."
}

variable "vpc_cidr_block" {
  type = any
  description = "The CIDR block of the VPC in which to create the load balancer."
}

variable "private_subnet_ids" {
  type = any
  description = "A list of IDs of the private subnets in which to create the ECS Service."
}

variable "extract_resource_name" {
  type        = string
  description = "The base name to use for all resources created by this module."
}

variable "lb_config" {
  type        = any
  description = "A map of objects that define the load balancers to create."
}

variable "certificate_arn" {
  type        = string
  description = "The ARN of the certificate to use for HTTPS listeners."
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
