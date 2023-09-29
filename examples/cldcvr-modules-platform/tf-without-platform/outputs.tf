# VPC Outputs
output "vpc_id" {
  value = module.vpc.vpc_id
  description = "The ID of the VPC."
}

output "public_subnets" {
  value = module.vpc.public_subnets
  description = "The IDs of the public subnets within the VPC."
}

output "private_subnets" {
  value = module.vpc.private_subnets
  description = "The IDs of the private subnets within the VPC."
}

# EC2 Instance Outputs
output "ec2_instance_id" {
  value = module.ec2_instance.id
  description = "The ID of the EC2 instance."
}

# EC2 Security Group Outputs
output "ec2_security_group_id" {
  value = module.ec2_sg.security_group_id
  description = "The ID of the EC2 security group."
}

# Database Outputs
output "database_id" {
  value = module.aws_sql_database.database_id
  description = "The ID of the database."
}

output "database_endpoint" {
  value = module.aws_sql_database.endpoint
  description = "The endpoint of the database."
}

# Database Security Group Outputs
output "db_security_group_id" {
  value = module.db_sg.security_group_id
  description = "The ID of the database security group."
}

# Redis Outputs
output "redis_id" {
  value = module.redis.id
  description = "The ID of the Redis instance."
}

output "redis_endpoint" {
  value = module.redis.endpoint
  description = "The endpoint of the Redis instance."
}

# Redis Security Group Outputs
output "redis_security_group_id" {
  value = module.redis_sg.security_group_id
  description = "The ID of the Redis security group."
}
