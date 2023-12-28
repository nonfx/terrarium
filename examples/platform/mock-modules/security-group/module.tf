# --- Input

variable "name" {
  type        = string
  description = "Name of the security group"
  default = ""
}

variable "vpc_id" {
  type        = string
  description = "The VPC ID where the security group is created"
  default = ""
}

variable "ingress_with_cidr_blocks" {
  type        = list(any)
  description = "A list of ingress rules with CIDR blocks"
  default     = []
}

variable "egress_with_cidr_blocks" {
  type        = list(any)
  description = "A list of egress rules with CIDR blocks"
  default     = []
}

# --- Config

resource "random_id" "security_group" {
  byte_length = 8
}

# --- Output

output "security_group_id" {
  description = "The ID of the security group"
  value       = "sg-${random_id.security_group.hex}"
}
