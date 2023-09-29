# General variables
variable "environment" {
  description = "The name for identifying the type of environment"
  type        = string
}

variable "common_name_prefix" {
  description = "The prefix used to name all resources created."
  type        = string
}

variable "number" {
  description = "The count of the resource"
  default     = 001
}

variable "tags" {
  type        = map(string)
  description = "Any tags that should be present on the Virtual Network resources"
  default     = {}
}

variable "extract_resource_name" {
  type        = string
  description = "The name of the resource"
}

#VPC Variables
variable "vpc_cidr_block" {
  description = "CIDR Block to be used by VPC"
  default     = "10.0.0.0/16"
  type        = string
}

variable "number_of_azs" {
  description = "Number of Availability Zones to use for the VPC"
  default     = 2
  type        = number
}

variable "vpc_enable_dns_hostnames" {
  description = "A boolean flag to enable/disable DNS support in the VPC. Defaults to true."
  default     = true
  type        = bool
}

variable "tr_component_vpc" {
  description = "A map of objects that define the VPC to create."
  type        = any
}
