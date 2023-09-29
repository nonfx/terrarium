variable "postgres_subnet_ids" {
  description = "A list of subnet IDs to launch the PostgreSQL instance in"
  type        = any
}

variable "postgres_vpc_id" {
  description = "The ID of the VPC to launch the PostgreSQL instance in"
  type        = any
}

variable "postgres_app_config" {
  description = "A map of objects that define the PostgreSQL instance to create"
  type        = any
}

variable "postgres_cidr_blocks" {
  description = "The CIDR blocks to allow access to the PostgreSQL instance"
  type        = any
}

variable "postgres_database_subnet_group" {
  description = "A map of objects that define the Postgres database subnet group to create."
  type        = any
}
