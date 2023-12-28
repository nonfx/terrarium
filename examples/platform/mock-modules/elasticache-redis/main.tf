# --- Input

variable "availability_zones" {
  type        = list(string)
  description = "A list of availability zones in which to place the cluster"
  default     = ["us-west-1a", "us-west-1b"]
}

variable "vpc_id" {
  type        = string
  description = "The VPC ID where the cluster is created"
  default     = "vpc-12345678"
}

variable "allowed_security_group_ids" {
  type        = list(string)
  description = "A list of Security Group IDs that are allowed to access the cluster"
  default     = ["sg-12345678"]
}

variable "subnets" {
  type        = list(string)
  description = "A list of subnet IDs for the cluster"
  default     = ["subnet-12345678", "subnet-23456789"]
}

variable "apply_immediately" {
  type        = bool
  description = "Specifies whether any modifications are applied immediately, or during the next maintenance window"
  default     = true
}

variable "automatic_failover_enabled" {
  type        = bool
  description = "Specifies whether a read-only replica will be automatically promoted to primary if the existing primary fails"
  default     = false
}

variable "engine_version" {
  type        = string
  description = "The version number of the cache engine to use"
  default     = "6.x"
}

# --- Config

resource "random_pet" "redis_instance" {
  length = 2
}

resource "random_password" "redis_password" {
  length  = 16
  special = true
}

# --- Output

output "host" {
  description = "The DNS name of the cache instance"
  value       = "${random_pet.redis_instance.id}.mock-elasticache.com"
}

output "port" {
  description = "The port number on which each of the cache nodes will accept connections"
  value       = "6379" # Default Redis port
}

output "password" {
  description = "The password used to access a password-protected Redis server"
  value       = random_password.redis_password.result
  sensitive = true
}
