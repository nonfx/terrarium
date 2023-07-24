

## Postgres component outputs

output "tr_component_postgres_host" {
  value = { for k, v in module.tr_component_postgres : k => v.db_instance_address }
}

output "tr_component_postgres_port" {
  value = { for k, v in module.tr_component_postgres : k => v.db_instance_port }
}
