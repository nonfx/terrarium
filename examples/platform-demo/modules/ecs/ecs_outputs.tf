output "ecs_cluster_names" {
  value = [for cluster in local.tr_component_ecs_combined : cluster.name]
}

output "ecs_service_names" {
  value = [for service_key, service_value in var.tr_component_ecs_services : service_value.name]
}

output "alb_outputs" {
  value = module.alb
}
