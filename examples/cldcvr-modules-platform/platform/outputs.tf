#VPC Output
output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

#EC2 Output
output "public_ip" {
  description = "The public IP address of the EC2 instance"
  value       = { for k, v in module.tr_component_service_web : k => v.public_ip }
}

output "public_dns" {
  description = "The public DNS name of the EC2 instance"
  value       = { for k, v in module.tr_component_service_web : k => v.public_dns }
}

#DB Output
output "postgres_host" {
  description = "The hostname of the DB instance"
  value       = { for k, v in module.tr_component_postgres : k => v.rds_endpoint }
}

#Redis Output
output "redis_host" {
  description = "The hostname of the Redis instance"
  value       = { for k, v in module.tr_component_redis : k => v.endpoint }
}
