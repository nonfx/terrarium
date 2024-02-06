#  --- Input

variable "name" {
  type = string
  description = "Name of the VPC"
  default = ""
}

variable "cidr" {
  type = string
  description = "The IPv4 CIDR block for the VPC."
  default = ""
}

variable "azs" {
  type = set(string)
  description = "A list of availability zones names or ids in the region"
  default = []
}

variable "private_subnets" {
  type = set(string)
  description = "A list of private subnets inside the VPC"
  default = []

}

variable "public_subnets" {
  type = set(string)
  description = "A list of public subnets inside the VPC"
  default = []
}

variable "database_subnets" {
  type = set(string)
  description = "A list of database subnets inside the VPC"
  default = []
}

variable "elasticache_subnets" {
  type = set(string)
  description = "A list of elasticache subnets inside the VPC"
  default = []
}

# --- Config

resource "random_id" "vpc" {
  byte_length = 8
}

resource "random_pet" "default_security_group" {
  length = 2
}

resource "random_id" "subnet" {
  byte_length = 8
  count = length(var.private_subnets) + length(var.public_subnets) + length(var.database_subnets) + length(var.elasticache_subnets)
}

# --- Output

output "vpc_id" {
  description = "The ID of the VPC"
  value = "vpc-${random_id.vpc.hex}"
}

output "vpc_cidr_block" {
  description = "The CIDR block of the VPC"
  value = var.cidr
}

output "default_security_group_id" {
  description = "The ID of the security group created by default on VPC creation"
  value = "sg-${random_pet.default_security_group.id}"
}

output "private_subnets" {
  description = "List of IDs of private subnets"
  value = [for i in range(length(var.private_subnets)) : "subnet-${random_id.subnet[i].hex}"]
}

output "public_subnets" {
  description = "List of IDs of public subnets"
  value = [for i in range(length(var.public_subnets)) : "subnet-${random_id.subnet[length(var.private_subnets) + i].hex}"]
}

output "database_subnets" {
  description = "List of IDs of database subnets"
  value = [for i in range(length(var.database_subnets)) : "subnet-${random_id.subnet[length(var.private_subnets) + length(var.public_subnets) + i].hex}"]
}

output "elasticache_subnets" {
  description = "List of IDs of elasticache subnets"
  value = [for i in range(length(var.elasticache_subnets)) : "subnet-${random_id.subnet[length(var.private_subnets) + length(var.public_subnets) + length(var.database_subnets) + i].hex}"]
}
