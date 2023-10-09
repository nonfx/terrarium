output "redis_cluster_address" {
  value = { for k, v in module.redis : k => v.host }
}

output "redis_cluster_endpoint" {
  value =  { for k, v in module.redis : k => v.endpoint }
}


output "redis_cluster_port" {
  value = { for k, v in module.redis : k => v.port }
}

output "redis_cluster_security_group_id" {
  value = { for k, v in module.redis : k => v.security_group_id }
}

output "redis_cluster_id" {
  value = { for k, v in module.redis : k => v.id }
}

output "redis_security_group_name" {
  value       = { for k, v in module.redis : k => v.security_group_name }
  description = "The name of the created security group"
}

output "redis_cluster_reader_endpoint_address" {
  value       = { for k, v in module.redis : k => v.reader_endpoint_address }
  description = "The address of the endpoint for the reader node in the replication group, if the cluster mode is disabled."
}

output "redis_member_clusters" {
  value       = { for k, v in module.redis : k => v.member_clusters }
  description = "Redis cluster members"
}


output "redis_cluster_arn" {
  value       = { for k, v in module.redis : k => v.arn }
  description = "Elasticache Replication Group ARN"
}

output "redis_engine_version_actual" {
  value       = { for k, v in module.redis : k => v.engine_version_actual }
  description = "The running version of the cache engine"
}

output "redis_cluster_enabled" {
  value       = { for k, v in module.redis : k => v.cluster_enabled }
  description = "Indicates if cluster mode is enabled"
}
