variable "common_name_prefix" {
  description = "The common name prefix for resources"
  type        = string
  default     = "demo"
}

variable "environment" {
  description = "The environment (e.g., prod, dev, staging)"
  type        = string
  default     = "dev"
}

variable "common_tags" {
  description = "Common tags for resources"
  type        = map(string)
  default     = {}
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

variable "key_name" {
  description = "Key name for EC2 instance"
  type        = string
  default     = "demo-key"
}

variable "ami_id" {
  description = "AMI ID for EC2 instance"
  type        = string
  default     = "ami-0f2e255ec956ade7f"
}
